package graphBetaApplication

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
)

// validateRequest validates the application request
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *ApplicationResourceModel, currentID string) error {
	tflog.Debug(ctx, "Starting validation of application request")

	// Only check for duplicate names if prevent_duplicate_names is true
	if !data.PreventDuplicateNames.IsNull() && data.PreventDuplicateNames.ValueBool() {
		if err := checkForDuplicateApplicationName(ctx, client, data.DisplayName.ValueString(), currentID); err != nil {
			return err
		}
	}

	tflog.Debug(ctx, "Successfully validated application request")
	return nil
}

// checkForDuplicateApplicationName checks if an application with the given display name already exists
// It queries the Microsoft Graph API with a filter on displayName
// If a duplicate is found and it's not the current resource (during update), it returns an error
func checkForDuplicateApplicationName(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName string, currentID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Checking for duplicate application name: %s", displayName))

	// Build the filter query: $filter=displayName eq 'name'
	filter := fmt.Sprintf("displayName eq '%s'", escapeODataString(displayName))

	requestConfig := &applications.ApplicationsRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ApplicationsRequestBuilderGetQueryParameters{
			Filter: &filter,
			Select: []string{"id", "displayName"},
		},
	}

	result, err := client.
		Applications().
		Get(ctx, requestConfig)

	if err != nil {
		return fmt.Errorf("failed to query for existing applications: %w", err)
	}

	if result == nil {
		tflog.Debug(ctx, "No existing applications found with this display name")
		return nil
	}

	existingApps := result.GetValue()
	if len(existingApps) == 0 {
		tflog.Debug(ctx, "No existing applications found with this display name")
		return nil
	}

	// Check if any of the found applications are different from the current one (for update scenarios)
	for _, app := range existingApps {
		if app.GetId() == nil {
			continue
		}

		existingID := *app.GetId()
		existingName := ""
		if app.GetDisplayName() != nil {
			existingName = *app.GetDisplayName()
		}

		// If this is an update and the existing app is the current resource, it's fine
		if currentID != "" && existingID == currentID {
			tflog.Debug(ctx, fmt.Sprintf("Found application is the current resource (ID: %s), skipping", existingID))
			continue
		}

		// Otherwise, we found a duplicate
		tflog.Warn(ctx, fmt.Sprintf("Found duplicate application with name '%s' and ID: %s", existingName, existingID))
		return fmt.Errorf("an application with the display name '%s' already exists (ID: %s). To import this existing application, use: terraform import %s.name %s", displayName, existingID, ResourceName, existingID)
	}

	tflog.Debug(ctx, "No duplicate applications found")
	return nil
}

// escapeODataString escapes single quotes in OData filter strings
// Single quotes in OData must be escaped by doubling them
func escapeODataString(s string) string {
	result := ""
	for _, char := range s {
		if char == '\'' {
			result += "''"
		} else {
			result += string(char)
		}
	}
	return result
}
