package graphBetaAndroidManagedDeviceAppConfigurationPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of an AndroidManagedDeviceAppConfigurationPolicyResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *AndroidManagedDeviceAppConfigurationPolicyResourceModel, remoteResource graphmodels.ManagedDeviceMobileAppConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.TargetedMobileApps = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetTargetedMobileApps())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.Version = convert.GraphToFrameworkInt32(remoteResource.GetVersion())

	if androidConfig, ok := remoteResource.(graphmodels.AndroidManagedStoreAppConfigurationable); ok {

		data.PackageId = convert.GraphToFrameworkString(androidConfig.GetPackageId())
		data.PayloadJson = convert.GraphToFrameworkBase64String(ctx, androidConfig.GetPayloadJson())
		data.ProfileApplicability = convert.GraphToFrameworkEnum(androidConfig.GetProfileApplicability())
		data.ConnectedAppsEnabled = convert.GraphToFrameworkBool(androidConfig.GetConnectedAppsEnabled())
		data.AppSupportsOemConfig = convert.GraphToFrameworkBool(androidConfig.GetAppSupportsOemConfig())

		if permissionActions := androidConfig.GetPermissionActions(); len(permissionActions) > 0 {
			permissionElements := make([]attr.Value, 0, len(permissionActions))

			for _, permission := range permissionActions {
				permissionAttrs := make(map[string]attr.Value)

				permissionAttrs["permission"] = convert.GraphToFrameworkString(permission.GetPermission())
				permissionAttrs["action"] = convert.GraphToFrameworkEnum(permission.GetAction())

				permissionObj, _ := types.ObjectValue(
					map[string]attr.Type{
						"permission": types.StringType,
						"action":     types.StringType,
					},
					permissionAttrs,
				)

				permissionElements = append(permissionElements, permissionObj)
			}

			permissionsSet, _ := types.SetValue(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"permission": types.StringType,
						"action":     types.StringType,
					},
				},
				permissionElements,
			)

			data.PermissionActions = permissionsSet
		} else {
			if data.PermissionActions.IsNull() || data.PermissionActions.IsUnknown() {
				data.PermissionActions = types.SetNull(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"permission": types.StringType,
						"action":     types.StringType,
					},
				})
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
