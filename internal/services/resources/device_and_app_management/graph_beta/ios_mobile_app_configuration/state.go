package graphBetaIOSMobileAppConfiguration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of an IOSMobileAppConfigurationResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *IOSMobileAppConfigurationResourceModel, remoteResource graphmodels.ManagedDeviceMobileAppConfigurationable) {
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

	if version := remoteResource.GetVersion(); version != nil {
		data.Version = types.Int64Value(int64(*version))
	} else {
		data.Version = types.Int64Null()
	}

	// Handle iOS specific properties if this is an iOS mobile app configuration
	if iosConfig, ok := remoteResource.(graphmodels.IosMobileAppConfigurationable); ok {
		// Handle encoded setting XML
		if encodedXml := iosConfig.GetEncodedSettingXml(); encodedXml != nil {
			data.EncodedSettingXml = types.StringValue(string(encodedXml))
		} else {
			data.EncodedSettingXml = types.StringNull()
		}

		// Handle settings
		if settings := iosConfig.GetSettings(); settings != nil {
			settingsElements := make([]attr.Value, 0, len(settings))

			for _, setting := range settings {
				settingAttrs := make(map[string]attr.Value)

				if key := setting.GetAppConfigKey(); key != nil {
					settingAttrs["app_config_key"] = types.StringValue(*key)
				} else {
					settingAttrs["app_config_key"] = types.StringNull()
				}

				settingAttrs["app_config_key_type"] = convert.GraphToFrameworkEnum(setting.GetAppConfigKeyType())

				if keyValue := setting.GetAppConfigKeyValue(); keyValue != nil {
					settingAttrs["app_config_key_value"] = types.StringValue(*keyValue)
				} else {
					settingAttrs["app_config_key_value"] = types.StringNull()
				}

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
			data.Settings = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"app_config_key":       types.StringType,
					"app_config_key_type":  types.StringType,
					"app_config_key_value": types.StringType,
				},
			})
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
