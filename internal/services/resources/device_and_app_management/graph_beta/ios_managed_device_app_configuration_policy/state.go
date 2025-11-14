package graphBetaIOSManagedDeviceAppConfigurationPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of an IOSManagedDeviceAppConfigurationPolicyResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *IOSManagedDeviceAppConfigurationPolicyResourceModel, remoteResource graphmodels.ManagedDeviceMobileAppConfigurationable) {
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

	if iosConfig, ok := remoteResource.(graphmodels.IosMobileAppConfigurationable); ok {

		data.EncodedSettingXml = convert.GraphToFrameworkBytes(iosConfig.GetEncodedSettingXml())

		if settings := iosConfig.GetSettings(); len(settings) > 0 {
			settingsElements := make([]attr.Value, 0, len(settings))

			for _, setting := range settings {
				settingAttrs := make(map[string]attr.Value)

				settingAttrs["app_config_key"] = convert.GraphToFrameworkString(setting.GetAppConfigKey())
				settingAttrs["app_config_key_type"] = convert.GraphToFrameworkEnum(setting.GetAppConfigKeyType())
				settingAttrs["app_config_key_value"] = convert.GraphToFrameworkString(setting.GetAppConfigKeyValue())

				settingObj, _ := types.ObjectValue(
					map[string]attr.Type{
						"app_config_key":       types.StringType,
						"app_config_key_type":  types.StringType,
						"app_config_key_value": types.StringType,
					},
					settingAttrs,
				)

				settingsElements = append(settingsElements, settingObj)
			}

			settingsSet, _ := types.SetValue(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"app_config_key":       types.StringType,
						"app_config_key_type":  types.StringType,
						"app_config_key_value": types.StringType,
					},
				},
				settingsElements,
			)

			data.Settings = settingsSet
		} else {
			if data.Settings.IsNull() || data.Settings.IsUnknown() {
				data.Settings = types.SetNull(types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"app_config_key":       types.StringType,
						"app_config_key_type":  types.StringType,
						"app_config_key_value": types.StringType,
					},
				})
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
