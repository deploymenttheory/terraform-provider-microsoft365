package graphBetaIOSManagedDeviceAppConfigurationPolicy

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
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *IOSManagedDeviceAppConfigurationPolicyResourceModel) (graphmodels.IosMobileAppConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewIosMobileAppConfiguration()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if !data.TargetedMobileApps.IsNull() && !data.TargetedMobileApps.IsUnknown() {
		var appIds []string
		if diags := data.TargetedMobileApps.ElementsAs(ctx, &appIds, false); diags.HasError() {
			return nil, fmt.Errorf("failed to extract targeted mobile app IDs: %v", diags)
		}

		if err := validateIOSMobileAppIds(ctx, client, appIds); err != nil {
			return nil, fmt.Errorf("validation failed for targeted_mobile_apps: %w", err)
		}
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.TargetedMobileApps, requestBody.SetTargetedMobileApps); err != nil {
		return nil, fmt.Errorf("failed to set targeted mobile apps: %s", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	convert.FrameworkToGraphBytes(data.EncodedSettingXml, requestBody.SetEncodedSettingXml)

	if !data.Settings.IsNull() && !data.Settings.IsUnknown() {
		settingsElements := data.Settings.Elements()
		graphSettings := make([]graphmodels.AppConfigurationSettingItemable, 0, len(settingsElements))

		for _, settingElement := range settingsElements {
			if settingObj, ok := settingElement.(types.Object); ok {
				attrs := settingObj.Attributes()

				setting := graphmodels.NewAppConfigurationSettingItem()

				if keyAttr, exists := attrs["app_config_key"]; exists {
					if keyStr, ok := keyAttr.(types.String); ok && !keyStr.IsNull() {
						setting.SetAppConfigKey(keyStr.ValueStringPointer())
					}
				}

				if keyTypeAttr, exists := attrs["app_config_key_type"]; exists {
					if keyTypeStr, ok := keyTypeAttr.(types.String); ok && !keyTypeStr.IsNull() {
						if err := convert.FrameworkToGraphEnum(
							keyTypeStr,
							graphmodels.ParseMdmAppConfigKeyType,
							setting.SetAppConfigKeyType,
						); err != nil {
							return nil, fmt.Errorf("failed to set app config key type: %w", err)
						}
					}
				}

				if keyValueAttr, exists := attrs["app_config_key_value"]; exists {
					if keyValueStr, ok := keyValueAttr.(types.String); ok && !keyValueStr.IsNull() {
						setting.SetAppConfigKeyValue(keyValueStr.ValueStringPointer())
					}
				}

				graphSettings = append(graphSettings, setting)
			}
		}

		requestBody.SetSettings(graphSettings)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
