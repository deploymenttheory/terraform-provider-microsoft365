package graphBetaGroupPolicyUploadedDefinitionFiles

import (
	"context"
	"fmt"
	"regexp"
	"time"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// checkADMXUploadStatus polls and monitors the status of an ADMX file upload operation until completion.
//
// This function handles the asynchronous nature of ADMX file uploads in Intune, which can take time to process.
// It periodically checks the status of the upload operation and waits for it to complete, with three possible outcomes:
//
// 1. "available" - Upload succeeded and the ADMX file is ready for use
// 2. "uploadFailed" - Upload failed, in which case it automatically cleans up by deleting the failed resource
// 3. "uploadInProgress" - Upload is still processing, function will continue polling until timeout
//
// The function also retrieves detailed error information when an upload fails by examining
// the groupPolicyOperations collection, which contains specific error messages that can help
// diagnose issues such as missing dependency files or format problems.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeout
//   - id: The ID of the group policy uploaded definition file to check
//
// Returns:
//   - string: The final status of the upload operation ("available", "uploadFailed", etc.)
//   - string: Detailed error message if the upload failed, empty string otherwise
//   - error: Any error that occurred during the status check process itself
func (r *GroupPolicyUploadedDefinitionFileResource) checkADMXUploadStatus(ctx context.Context, id string) (string, string, error) {
	maxRetries := 60
	retryInterval := 5 * time.Second

	for i := range maxRetries {

		respResource, err := r.client.
			DeviceManagement().
			GroupPolicyUploadedDefinitionFiles().
			ByGroupPolicyUploadedDefinitionFileId(id).
			Get(ctx, &devicemanagement.GroupPolicyUploadedDefinitionFilesGroupPolicyUploadedDefinitionFileItemRequestBuilderGetRequestConfiguration{
				QueryParameters: &devicemanagement.GroupPolicyUploadedDefinitionFilesGroupPolicyUploadedDefinitionFileItemRequestBuilderGetQueryParameters{
					Expand: []string{"groupPolicyOperations"},
				},
			})

		if err != nil {
			return "", "", fmt.Errorf("failed to get upload status: %v", err)
		}

		if status := respResource.GetStatus(); status != nil {
			statusStr := status.String()

			statusDetails := ""
			if operations := respResource.GetGroupPolicyOperations(); len(operations) > 0 {
				for _, op := range operations {
					if op.GetOperationType() != nil && op.GetOperationType().String() == "upload" {
						if op.GetStatusDetails() != nil {
							statusDetails = *op.GetStatusDetails()
							break
						}
					}
				}
			}

			switch statusStr {
			case "available":
				return statusStr, statusDetails, nil
			case "uploadFailed":
				tflog.Debug(ctx, fmt.Sprintf("ADMX upload failed with details: %s", statusDetails))

				deleteErr := r.client.
					DeviceManagement().
					GroupPolicyUploadedDefinitionFiles().
					ByGroupPolicyUploadedDefinitionFileId(id).
					Delete(ctx, nil)

				if deleteErr != nil {
					tflog.Warn(ctx, fmt.Sprintf("Failed to delete failed upload: %s", deleteErr.Error()))
				}

				return statusStr, statusDetails, nil
			case "uploadInProgress":
				// Still in progress, wait and retry
				tflog.Debug(ctx, fmt.Sprintf("ADMX upload in progress (attempt %d/%d), waiting %v before checking again",
					i+1, maxRetries, retryInterval))
				time.Sleep(retryInterval)
				continue
			default:
				return statusStr, statusDetails, fmt.Errorf("unknown ADMX upload status: %s", statusStr)
			}
		}

		time.Sleep(retryInterval)
	}

	return "", "", fmt.Errorf("timed out waiting for ADMX upload to complete after %d attempts", maxRetries)
}

// checkADMXRemovalStatus polls and monitors the deletion of an ADMX file until it is fully removed.
//
// This function handles the asynchronous nature of ADMX file deletion in Intune, which can take time to process.
// It periodically checks the status of the removal operation and waits for the resource to be completely deleted.
//
// The removal process has the following states:
// 1. "removalInProgress" - Deletion is in progress, function will continue polling
// 2. 404 Not Found - Resource has been completely deleted (success condition)
//
// The function also examines the groupPolicyOperations collection to get detailed information
// about the removal operation status, which can help diagnose any issues during deletion.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeout
//   - id: The ID of the group policy uploaded definition file to monitor
//
// Returns:
//   - error: Any error that occurred during the removal monitoring process, or nil if successful
func (r *GroupPolicyUploadedDefinitionFileResource) checkADMXRemovalStatus(ctx context.Context, id string) error {
	maxRetries := 60
	retryInterval := 5 * time.Second

	for i := range maxRetries {

		respResource, err := r.client.
			DeviceManagement().
			GroupPolicyUploadedDefinitionFiles().
			ByGroupPolicyUploadedDefinitionFileId(id).
			Get(ctx, &devicemanagement.GroupPolicyUploadedDefinitionFilesGroupPolicyUploadedDefinitionFileItemRequestBuilderGetRequestConfiguration{
				QueryParameters: &devicemanagement.GroupPolicyUploadedDefinitionFilesGroupPolicyUploadedDefinitionFileItemRequestBuilderGetQueryParameters{
					Expand: []string{"groupPolicyOperations"},
				},
			})

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 {
				tflog.Debug(ctx, "ADMX removal completed - resource no longer exists (404)")
				return nil
			}
			return fmt.Errorf("failed to check removal status: %v", err)
		}

		if status := respResource.GetStatus(); status != nil {
			statusStr := status.String()

			if operations := respResource.GetGroupPolicyOperations(); len(operations) > 0 {
				for _, op := range operations {
					if op.GetOperationType() != nil && op.GetOperationType().String() == "remove" {
						opStatus := "unknown"
						if op.GetOperationStatus() != nil {
							opStatus = op.GetOperationStatus().String()
						}
						tflog.Debug(ctx, fmt.Sprintf("Remove operation status: %s", opStatus))
					}
				}
			}

			switch statusStr {
			case "removalInProgress":
				tflog.Debug(ctx, fmt.Sprintf("ADMX removal in progress (attempt %d/%d), waiting %v before checking again",
					i+1, maxRetries, retryInterval))
				time.Sleep(retryInterval)
				continue
			case "removalFailed":
				return fmt.Errorf("ADMX removal failed")
			default:
				tflog.Warn(ctx, fmt.Sprintf("Unexpected status during removal: %s", statusStr))
				time.Sleep(retryInterval)
				continue
			}
		}

		time.Sleep(retryInterval)
	}

	return fmt.Errorf("timed out waiting for ADMX removal to complete after %d attempts", maxRetries)
}

// deleteExistingDefinitionFileByContent attempts to find and delete any existing definition file
// with the same target namespace as the one being uploaded. This is determined by parsing the
// ADMX content to extract the target namespace, then querying all existing definition files
// to find a match.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeout
//   - content: The ADMX file content to parse for target namespace
//
// Returns:
//   - error: Any error that occurred during the process, or nil if successful or no match found
func (r *GroupPolicyUploadedDefinitionFileResource) deleteExistingDefinitionFileByContent(ctx context.Context, content string) error {
	tflog.Debug(ctx, "Checking for existing definition files with same target namespace")

	allFiles, err := r.client.
		DeviceManagement().
		GroupPolicyUploadedDefinitionFiles().
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to list existing definition files: %v", err)
	}

	if allFiles == nil || allFiles.GetValue() == nil {
		tflog.Debug(ctx, "No existing definition files found")
		return nil
	}

	targetNamespace := extractTargetNamespaceFromADMX(content)
	if targetNamespace == "" {
		tflog.Debug(ctx, "Could not extract target namespace from ADMX content, skipping pre-deletion check")
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Extracted target namespace from ADMX: %s", targetNamespace))

	for _, file := range allFiles.GetValue() {
		if file.GetTargetNamespace() != nil && *file.GetTargetNamespace() == targetNamespace {
			fileID := file.GetId()
			if fileID == nil {
				continue
			}

			tflog.Info(ctx, fmt.Sprintf("Found existing definition file with matching namespace %s (ID: %s), deleting it", targetNamespace, *fileID))

			deleteErr := r.client.
				DeviceManagement().
				GroupPolicyUploadedDefinitionFiles().
				ByGroupPolicyUploadedDefinitionFileId(*fileID).
				Remove().
				Post(ctx, nil)

			if deleteErr != nil {
				return fmt.Errorf("failed to delete existing definition file %s: %v", *fileID, deleteErr)
			}

			removalErr := r.checkADMXRemovalStatus(ctx, *fileID)
			if removalErr != nil {
				return fmt.Errorf("failed to verify removal of existing definition file %s: %v", *fileID, removalErr)
			}

			tflog.Info(ctx, fmt.Sprintf("Successfully deleted existing definition file with ID: %s", *fileID))
			return nil
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("No existing definition file found with target namespace: %s", targetNamespace))
	return nil
}

// extractTargetNamespaceFromADMX parses ADMX XML content to extract the target namespace.
// ADMX files use XML format with a policyNamespaces element that contains the target namespace.
//
// Example ADMX structure:
//
//	<policyDefinitions xmlns:xsd="..." xmlns:xsi="..." revision="1.0" schemaVersion="1.0">
//	  <policyNamespaces>
//	    <target prefix="mozilla" namespace="Mozilla.Policies.Firefox"/>
//	  </policyNamespaces>
//	</policyDefinitions>
//
// Parameters:
//   - content: The ADMX file content as a string
//
// Returns:
//   - string: The extracted target namespace, or empty string if not found
func extractTargetNamespaceFromADMX(content string) string {
	namespaceRegex := regexp.MustCompile(`<target\s+[^>]*namespace="([^"]+)"`)
	matches := namespaceRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
