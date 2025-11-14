package graphBetaAndroidManagedDeviceAppConfigurationPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AndroidManagedDeviceAppConfigurationPolicyResourceModel) (graphmodels.AndroidManagedStoreAppConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewAndroidManagedStoreAppConfiguration()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	// Validate and set targeted mobile apps
	if !data.TargetedMobileApps.IsNull() && !data.TargetedMobileApps.IsUnknown() {
		var appIds []string
		if diags := data.TargetedMobileApps.ElementsAs(ctx, &appIds, false); diags.HasError() {
			return nil, fmt.Errorf("failed to extract targeted mobile app IDs: %v", diags)
		}

		// Validate app IDs against Intune
		if err := validateAndroidMobileAppIds(ctx, client, appIds); err != nil {
			return nil, fmt.Errorf("validation failed for targeted_mobile_apps: %w", err)
		}
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.TargetedMobileApps, requestBody.SetTargetedMobileApps); err != nil {
		return nil, fmt.Errorf("failed to set targeted mobile apps: %s", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	convert.FrameworkToGraphString(data.PackageId, requestBody.SetPackageId)
	convert.FrameworkToGraphBase64String(data.PayloadJson, requestBody.SetPayloadJson)

	if !data.ProfileApplicability.IsNull() && !data.ProfileApplicability.IsUnknown() {
		if err := convert.FrameworkToGraphEnum(
			data.ProfileApplicability,
			graphmodels.ParseAndroidProfileApplicability,
			requestBody.SetProfileApplicability,
		); err != nil {
			return nil, fmt.Errorf("failed to set profile applicability: %w", err)
		}
	}

	convert.FrameworkToGraphBool(data.ConnectedAppsEnabled, requestBody.SetConnectedAppsEnabled)

	if !data.PermissionActions.IsNull() && !data.PermissionActions.IsUnknown() {
		permissionsElements := data.PermissionActions.Elements()
		graphPermissions := make([]graphmodels.AndroidPermissionActionable, 0, len(permissionsElements))

		for _, permElement := range permissionsElements {
			if permObj, ok := permElement.(types.Object); ok {
				attrs := permObj.Attributes()

				permission := graphmodels.NewAndroidPermissionAction()

				if permAttr, exists := attrs["permission"]; exists {
					if permStr, ok := permAttr.(types.String); ok && !permStr.IsNull() {
						permission.SetPermission(permStr.ValueStringPointer())
					}
				}

				if actionAttr, exists := attrs["action"]; exists {
					if actionStr, ok := actionAttr.(types.String); ok && !actionStr.IsNull() {
						if err := convert.FrameworkToGraphEnum(
							actionStr,
							graphmodels.ParseAndroidPermissionActionType,
							permission.SetAction,
						); err != nil {
							return nil, fmt.Errorf("failed to set permission action: %w", err)
						}
					}
				}

				graphPermissions = append(graphPermissions, permission)
			}
		}

		requestBody.SetPermissionActions(graphPermissions)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
