package graphBetaTargetedManagedAppConfigurations

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *TargetedManagedAppConfigurationResourceModel) (graphmodels.TargetedManagedAppConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewTargetedManagedAppConfiguration()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

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

	// Handle app group type - set default to "allApps" if not specified
	if !data.AppGroupType.IsNull() && !data.AppGroupType.IsUnknown() {
		if err := convert.FrameworkToGraphEnum(
			data.AppGroupType,
			graphmodels.ParseTargetedManagedAppGroupType,
			requestBody.SetAppGroupType,
		); err != nil {
			return nil, fmt.Errorf("failed to set app group type: %w", err)
		}
	}

	// Handle targeted app management levels - set default to "unspecified" if not specified
	if !data.TargetedAppManagementLevels.IsNull() && !data.TargetedAppManagementLevels.IsUnknown() {
		if err := convert.FrameworkToGraphEnum(
			data.TargetedAppManagementLevels,
			graphmodels.ParseAppManagementLevel,
			requestBody.SetTargetedAppManagementLevels,
		); err != nil {
			return nil, fmt.Errorf("failed to set targeted app management levels: %w", err)
		}
	} else {
		// Set default value
		targetedAppManagementLevels := graphmodels.AppManagementLevel(graphmodels.UNSPECIFIED_APPMANAGEMENTLEVEL)
		requestBody.SetTargetedAppManagementLevels(&targetedAppManagementLevels)
	}

	// Handle apps - always set apps field, even if empty
	var appModels []ManagedMobileAppResourceModel
	if !data.Apps.IsNull() && !data.Apps.IsUnknown() {
		diags := data.Apps.ElementsAs(ctx, &appModels, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert apps set: %v", diags)
		}
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

	requestBody.SetApps(graphApps)

	// Handle settings catalog configuration. Defer to the settings catalog resource constructor.
	// given it's complexity.
	if data.SettingsCatalog != nil && len(data.SettingsCatalog.Settings) > 0 {
		tflog.Debug(ctx, "Adding settings catalog settings to create request")
		settings := ConstructSettingsCatalogSettings(ctx, *data.SettingsCatalog)
		requestBody.SetSettings(settings)
	}

	if !data.Assignments.IsNull() && !data.Assignments.IsUnknown() {
		assignments, err := constructAssignments(ctx, data.Assignments)
		if err != nil {
			return nil, fmt.Errorf("failed to construct assignments: %s", err)
		}
		requestBody.SetAssignments(assignments)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
