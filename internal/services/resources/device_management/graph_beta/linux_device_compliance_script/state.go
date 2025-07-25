package graphBetaLinuxDeviceComplianceScript

import (
	"context"
	"encoding/base64"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote DeviceManagementReusablePolicySetting resource state to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *LinuxDeviceComplianceScriptResourceModel, remoteResource graphmodels.DeviceManagementReusablePolicySettingable) {
	if remoteResource == nil {
		return
	}

	// Store the original script content in case we need to preserve it
	originalScriptContent := data.DetectionScriptContent

	// Map basic fields
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.SettingDefinitionId = convert.GraphToFrameworkString(remoteResource.GetSettingDefinitionId())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	// Set version
	if version := remoteResource.GetVersion(); version != nil {
		data.Version = types.Int32Value(*version)
	} else {
		data.Version = types.Int32Null()
	}

	// Extract detection script content from setting instance
	// Expected structure:
	// "settingInstance": {
	//     "@odata.type": "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
	//     "settingDefinitionId": "linux_customcompliance_discoveryscript_reusablesetting",
	//     "settingInstanceTemplateReference": null,
	//     "simpleSettingValue": {
	//         "@odata.type": "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
	//         "settingValueTemplateReference": null,
	//         "value": "base64encodedstring"
	//     }
	// }
	if settingInstance := remoteResource.GetSettingInstance(); settingInstance != nil {
		tflog.Debug(ctx, "Processing setting instance", map[string]any{
			"hasSettingInstance": true,
		})

		// Verify this is a DeviceManagementConfigurationSimpleSettingInstance
		if simpleInstance, ok := settingInstance.(graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable); ok {
			tflog.Debug(ctx, "Successfully cast to simple setting instance")

			// Verify the setting definition ID matches
			if instanceSettingDefId := simpleInstance.GetSettingDefinitionId(); instanceSettingDefId != nil {
				tflog.Debug(ctx, "Setting instance definition ID", map[string]any{
					"settingDefinitionId": *instanceSettingDefId,
				})
			}

			// Get the simpleSettingValue
			if simpleSettingValue := simpleInstance.GetSimpleSettingValue(); simpleSettingValue != nil {
				tflog.Debug(ctx, "Found simple setting value")

				// Verify this is a DeviceManagementConfigurationStringSettingValue
				if stringValue, ok := simpleSettingValue.(graphmodels.DeviceManagementConfigurationStringSettingValueable); ok {
					tflog.Debug(ctx, "Successfully cast to string setting value")

					// Get the base64 encoded value
					if encodedValue := stringValue.GetValue(); encodedValue != nil && *encodedValue != "" {
						tflog.Debug(ctx, "Found encoded value", map[string]any{
							"encodedLength": len(*encodedValue),
						})

						// Decode base64 content
						if decodedBytes, err := base64.StdEncoding.DecodeString(*encodedValue); err == nil {
							decodedContent := string(decodedBytes)
							data.DetectionScriptContent = types.StringValue(decodedContent)
							tflog.Debug(ctx, "Successfully decoded detection script content", map[string]any{
								"decodedLength": len(decodedContent),
							})
						} else {
							tflog.Error(ctx, "Failed to decode base64 script content", map[string]any{
								"error":        err.Error(),
								"encodedValue": *encodedValue,
							})
							data.DetectionScriptContent = types.StringNull()
						}
					} else {
						tflog.Warn(ctx, "Encoded value is null or empty")
						data.DetectionScriptContent = types.StringNull()
					}
				} else {
					tflog.Error(ctx, "Setting value is not a DeviceManagementConfigurationStringSettingValue")
					data.DetectionScriptContent = types.StringNull()
				}
			} else {
				tflog.Error(ctx, "No simpleSettingValue found in setting instance")
				data.DetectionScriptContent = types.StringNull()
			}
		} else {
			tflog.Error(ctx, "Setting instance is not a DeviceManagementConfigurationSimpleSettingInstance")
			data.DetectionScriptContent = types.StringNull()
		}
	} else {
		tflog.Warn(ctx, "No setting instance found in remote resource - preserving original script content")
		// If the API doesn't return the settingInstance (which can happen immediately after creation),
		// preserve the original script content to avoid inconsistent state errors
		if !originalScriptContent.IsNull() && !originalScriptContent.IsUnknown() {
			data.DetectionScriptContent = originalScriptContent
			tflog.Debug(ctx, "Preserved original script content", map[string]any{
				"contentLength": len(originalScriptContent.ValueString()),
			})
		} else {
			data.DetectionScriptContent = types.StringNull()
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state")
}
