package graphBetaTargetedManagedAppConfigurations

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	deviceappmanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructBaseResourceUpdate constructs the PATCH request body for base resource properties
func constructBaseResourceUpdate(ctx context.Context, data *TargetedManagedAppConfigurationResourceModel) (graphmodels.TargetedManagedAppConfigurationable, error) {
	tflog.Debug(ctx, "Constructing base resource update request")

	requestBody := graphmodels.NewTargetedManagedAppConfiguration()
	odataType := "#microsoft.graph.targetedManagedAppConfiguration"
	requestBody.SetOdataType(&odataType)

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	return requestBody, nil
}

// constructTargetAppsUpdate constructs the POST request body for /targetApps endpoint
func constructTargetAppsUpdate(ctx context.Context, data *TargetedManagedAppConfigurationResourceModel) (*deviceappmanagement.TargetedManagedAppConfigurationsItemTargetAppsPostRequestBody, error) {
	tflog.Debug(ctx, "Constructing target apps update request")

	appsRequest := deviceappmanagement.NewTargetedManagedAppConfigurationsItemTargetAppsPostRequestBody()

	// Set app group type
	if !data.AppGroupType.IsNull() && !data.AppGroupType.IsUnknown() {
		if err := convert.FrameworkToGraphEnum(
			data.AppGroupType,
			graphmodels.ParseTargetedManagedAppGroupType,
			appsRequest.SetAppGroupType,
		); err != nil {
			return nil, fmt.Errorf("failed to set app group type: %s", err)
		}
	}

	// Set apps
	if !data.Apps.IsNull() && !data.Apps.IsUnknown() {
		var appModels []ManagedMobileAppResourceModel
		diags := data.Apps.ElementsAs(ctx, &appModels, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert apps set: %v", diags)
		}

		graphApps := make([]graphmodels.ManagedMobileAppable, 0, len(appModels))
		for _, appModel := range appModels {
			app := graphmodels.NewManagedMobileApp()

			if appModel.MobileAppIdentifier != nil {
				switch appModel.MobileAppIdentifier.Type.ValueString() {
				case "android_mobile_app":
					identifier := graphmodels.NewAndroidMobileAppIdentifier()
					androidOdataType := "#microsoft.graph.androidMobileAppIdentifier"
					identifier.SetOdataType(&androidOdataType)
					convert.FrameworkToGraphString(appModel.MobileAppIdentifier.PackageId, identifier.SetPackageId)
					app.SetMobileAppIdentifier(identifier)
				case "ios_mobile_app":
					identifier := graphmodels.NewIosMobileAppIdentifier()
					iosOdataType := "#microsoft.graph.iosMobileAppIdentifier"
					identifier.SetOdataType(&iosOdataType)
					convert.FrameworkToGraphString(appModel.MobileAppIdentifier.BundleId, identifier.SetBundleId)
					app.SetMobileAppIdentifier(identifier)
				case "windows_app":
					identifier := graphmodels.NewWindowsAppIdentifier()
					windowsOdataType := "#microsoft.graph.windowsAppIdentifier"
					identifier.SetOdataType(&windowsOdataType)
					convert.FrameworkToGraphString(appModel.MobileAppIdentifier.WindowsAppId, identifier.SetWindowsAppId)
					app.SetMobileAppIdentifier(identifier)
				}
			}

			convert.FrameworkToGraphString(appModel.Version, app.SetVersion)
			graphApps = append(graphApps, app)
		}
		appsRequest.SetApps(graphApps)
	}

	return appsRequest, nil
}

// constructCustomSettingsUpdate constructs the PATCH request body for custom settings
func constructCustomSettingsUpdate(ctx context.Context, data *TargetedManagedAppConfigurationResourceModel) (graphmodels.TargetedManagedAppConfigurationable, error) {
	tflog.Debug(ctx, "Constructing custom settings update request")

	requestBody := graphmodels.NewTargetedManagedAppConfiguration()
	odataType := "#microsoft.graph.targetedManagedAppConfiguration"
	requestBody.SetOdataType(&odataType)

	if len(data.CustomSettings) > 0 {
		graphCustomSettings := make([]graphmodels.KeyValuePairable, 0, len(data.CustomSettings))
		for _, settingModel := range data.CustomSettings {
			setting := graphmodels.NewKeyValuePair()
			convert.FrameworkToGraphString(settingModel.Name, setting.SetName)
			convert.FrameworkToGraphString(settingModel.Value, setting.SetValue)
			graphCustomSettings = append(graphCustomSettings, setting)
		}
		requestBody.SetCustomSettings(graphCustomSettings)
	}

	return requestBody, nil
}

// constructSettingsUpdate constructs the POST request body for /changeSettings endpoint
func constructSettingsUpdate(ctx context.Context, data *TargetedManagedAppConfigurationResourceModel) (*deviceappmanagement.TargetedManagedAppConfigurationsItemChangeSettingsPostRequestBody, error) {
	tflog.Debug(ctx, "Constructing settings update request")

	settingsRequest := deviceappmanagement.NewTargetedManagedAppConfigurationsItemChangeSettingsPostRequestBody()

	if data.SettingsCatalog != nil && len(data.SettingsCatalog.Settings) > 0 {
		settings := ConstructSettingsCatalogSettings(ctx, *data.SettingsCatalog)
		settingsRequest.SetSettings(settings)
	}

	return settingsRequest, nil
}
