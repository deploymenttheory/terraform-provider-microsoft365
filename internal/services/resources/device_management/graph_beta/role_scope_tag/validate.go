package graphBetaRoleScopeTag

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates that the display name is unique among existing role scope tags
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName string, excludeResourceID *string) error {
	tflog.Debug(ctx, "Starting display name validation", map[string]interface{}{
		"displayName": displayName,
	})

	// Get all existing role scope tags
	existingRoleScopeTags, err := client.
		DeviceManagement().
		RoleScopeTags().
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to retrieve existing role scope tags for validation: %s", err.Error())
	}

	if existingRoleScopeTags == nil || existingRoleScopeTags.GetValue() == nil {
		tflog.Debug(ctx, "No existing role scope tags found, display name is unique")
		return nil
	}

	// Check if any existing role scope tag has the same display name (excluding current resource if specified)
	for _, existingTag := range existingRoleScopeTags.GetValue() {
		if existingTag.GetDisplayName() != nil && *existingTag.GetDisplayName() == displayName {
			// Skip validation if this is the current resource being updated
			if excludeResourceID != nil && existingTag.GetId() != nil && *existingTag.GetId() == *excludeResourceID {
				tflog.Debug(ctx, "Skipping validation for current resource", map[string]interface{}{
					"displayName": displayName,
					"resourceId":  *excludeResourceID,
				})
				continue
			}
			return fmt.Errorf("role scope tag with display name '%s' already exists. Display names must be unique", displayName)
		}
	}

	tflog.Debug(ctx, "Display name validation passed - name is unique", map[string]interface{}{
		"displayName": displayName,
	})

	return nil
}
