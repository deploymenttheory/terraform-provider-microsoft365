package graphBetaApplication

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
)

// validateRequest validates the request body during creation or update operations.
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *ApplicationResourceModel, currentID string, isCreate bool) error {
	tflog.Debug(ctx, "Starting validation of application request")

	// Only check for duplicate names if prevent_duplicate_names is true
	if !data.PreventDuplicateNames.IsNull() && data.PreventDuplicateNames.ValueBool() {
		if err := checkForDuplicateApplicationName(ctx, client, data.DisplayName.ValueString(), currentID); err != nil {
			return err
		}
	}

	// Validate app roles configuration for create operations
	// if isCreate {
	// 	if err := validateAppRolesIsEnabledIsTrue(ctx, data); err != nil {
	// 		return err
	// 	}
	// }

	tflog.Debug(ctx, "Successfully validated application request")
	return nil
}

// validateAppRolesIsEnabledIsTrue validates that all app roles have is_enabled set to true during create operations.
// The is_enabled=false value is only valid during update operations when disabling a role before deletion.
func validateAppRolesIsEnabledIsTrue(ctx context.Context, data *ApplicationResourceModel) error {
	if data.AppRoles.IsNull() || data.AppRoles.IsUnknown() {
		return nil
	}

	tflog.Debug(ctx, "Validating app roles for create operation")

	var appRoles []ApplicationAppRole
	data.AppRoles.ElementsAs(ctx, &appRoles, false)

	for _, role := range appRoles {
		if !role.IsEnabled.IsNull() && !role.IsEnabled.IsUnknown() && !role.IsEnabled.ValueBool() {
			displayName := "unknown"
			if !role.DisplayName.IsNull() && !role.DisplayName.IsUnknown() {
				displayName = role.DisplayName.ValueString()
			}
			return fmt.Errorf("app role '%s' has is_enabled set to false. During resource creation, all app roles must have is_enabled set to true (or omitted to use the default). The is_enabled=false value is only valid during update operations when disabling a role before deletion", displayName)
		}
	}

	tflog.Debug(ctx, "App roles validation passed for create operation")
	return nil
}

// checkForDuplicateApplicationName checks if an application with the given display name already exists
// It queries the Microsoft Graph API with a filter on displayName
// If a duplicate is found and it's not the current resource (during update), it returns an error
func checkForDuplicateApplicationName(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName string, currentID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Checking for duplicate application name: %s", displayName))

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
