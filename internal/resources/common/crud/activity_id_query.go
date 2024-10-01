package crud

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/auditlogs"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func ExtractActivityID(ctx context.Context, errorMessage string) string {
	tflog.Debug(ctx, "Attempting to extract Activity ID from error message", map[string]interface{}{
		"error_message": errorMessage,
	})

	// First, try to find the Activity ID using a regex pattern
	re := regexp.MustCompile(`Activity ID:\s*([0-9a-fA-F-]+)`)
	matches := re.FindStringSubmatch(errorMessage)
	if len(matches) > 1 {
		activityID := matches[1]
		tflog.Debug(ctx, "Activity ID extracted using regex", map[string]interface{}{
			"activity_id": activityID,
		})
		return activityID
	}
	tflog.Debug(ctx, "Regex extraction failed, attempting string split method")

	// If regex fails, try a simpler string split approach
	parts := strings.Split(errorMessage, "Activity ID:")
	if len(parts) > 1 {
		// Take the part after "Activity ID:" and trim any leading/trailing whitespace
		activityIDPart := strings.TrimSpace(parts[1])
		tflog.Debug(ctx, "String after 'Activity ID:'", map[string]interface{}{
			"activity_id_part": activityIDPart,
		})

		// If there's any whitespace or dash in the remaining part, split again and take the first part
		if strings.Contains(activityIDPart, " ") || strings.Contains(activityIDPart, "-") {
			activityID := strings.Split(activityIDPart, " ")[0]
			tflog.Debug(ctx, "Activity ID extracted after splitting", map[string]interface{}{
				"activity_id": activityID,
			})
			return activityID
		}

		tflog.Debug(ctx, "Activity ID extracted without further splitting", map[string]interface{}{
			"activity_id": activityIDPart,
		})
		return activityIDPart
	}

	// If no Activity ID is found, return an empty string
	tflog.Debug(ctx, "No Activity ID found in the error message")
	return ""
}

// Current QueryActivityDetails function (simplified)
func QueryActivityDetails(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, activityID string) (bool, error) {
	// Try with correlationId first
	filter := fmt.Sprintf("correlationId eq '%s'", activityID)
	isUnavailable, err := queryAuditLogs(ctx, client, filter, activityID)
	if err == nil {
		return isUnavailable, nil
	}

	// Log the error from the first attempt
	tflog.Debug(ctx, "First attempt to query audit logs failed", map[string]interface{}{
		"filter": filter,
		"error":  err.Error(),
	})

	// If no results, try with activityDisplayName
	filter = fmt.Sprintf("activityDisplayName eq 'Delete assignment filter'")
	isUnavailable, err = queryAuditLogs(ctx, client, filter, activityID)
	if err == nil {
		return isUnavailable, nil
	}

	// Log the error from the second attempt
	tflog.Debug(ctx, "Second attempt to query audit logs failed", map[string]interface{}{
		"filter": filter,
		"error":  err.Error(),
	})

	// If still no results, we can't determine if the resource is unavailable
	tflog.Warn(ctx, "Unable to determine resource availability", map[string]interface{}{
		"activity_id": activityID,
		"error":       err.Error(),
	})

	// Return an error instead of assuming the resource is available
	return false, fmt.Errorf("unable to determine resource availability: %w", err)
}

func queryAuditLogs(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, filter string, activityID string) (bool, error) {
	tflog.Debug(ctx, "Querying audit logs", map[string]interface{}{
		"filter": filter,
	})

	requestParameters := &auditlogs.DirectoryAuditsRequestBuilderGetQueryParameters{
		Filter: &filter,
		Select: []string{"activityDisplayName", "result", "operationType", "targetResources"},
		Top:    Int32Ptr(1), // Limit to 1 result for efficiency
	}

	configuration := &auditlogs.DirectoryAuditsRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	result, err := client.AuditLogs().DirectoryAudits().Get(ctx, configuration)
	if err != nil {
		tflog.Debug(ctx, "Error querying audit logs", map[string]interface{}{
			"error": err.Error(),
		})
		return false, fmt.Errorf("error querying audit logs: %v", err)
	}

	if result == nil || len(result.GetValue()) == 0 {
		tflog.Debug(ctx, "No audit log entries found", map[string]interface{}{
			"activity_id": activityID,
		})
		return false, fmt.Errorf("no audit log entries found for Activity ID: %s", activityID)
	}

	entry := result.GetValue()[0]
	activityDisplayName := entry.GetActivityDisplayName()
	operationResult := entry.GetResult()
	operationType := entry.GetOperationType()
	targetResources := entry.GetTargetResources()

	tflog.Debug(ctx, "Audit log entry details", map[string]interface{}{
		"activity_display_name": *activityDisplayName,
		"operation_result":      string(*operationResult),
		"operation_type":        *operationType,
	})

	isDeleteOperation := strings.Contains(strings.ToLower(*operationType), "delete") ||
		strings.Contains(strings.ToLower(*operationType), "remove")
	isSuccessfulOperation := *operationResult == models.SUCCESS_OPERATIONRESULT

	if isDeleteOperation && isSuccessfulOperation {
		tflog.Debug(ctx, "Resource identified as unavailable due to successful delete operation")
		return true, nil
	}

	for _, resource := range targetResources {
		if resource.GetModifiedProperties() != nil {
			for _, prop := range resource.GetModifiedProperties() {
				if strings.EqualFold(*prop.GetDisplayName(), "IsDeleted") && strings.EqualFold(*prop.GetNewValue(), "True") {
					tflog.Debug(ctx, "Resource identified as unavailable due to IsDeleted property")
					return true, nil
				}
				if strings.EqualFold(*prop.GetDisplayName(), "Status") && strings.EqualFold(*prop.GetNewValue(), "Deleted") {
					tflog.Debug(ctx, "Resource identified as unavailable due to Status property")
					return true, nil
				}
			}
		}
	}

	if *operationResult == models.FAILURE_OPERATIONRESULT && strings.Contains(strings.ToLower(*activityDisplayName), "get") {
		tflog.Debug(ctx, "Resource might be unavailable due to failed GET operation")
		return true, nil
	}

	tflog.Debug(ctx, "Resource appears to be available")
	return false, nil
}

func Int32Ptr(i int32) *int32 {
	return &i
}
