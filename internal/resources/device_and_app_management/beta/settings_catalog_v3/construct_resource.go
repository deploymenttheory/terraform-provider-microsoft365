package graphBetaSettingsCatalog

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the settings catalog profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *SettingsCatalogProfileResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, "Constructing Settings Catalog resource")
	construct.DebugPrintStruct(ctx, "Constructed Settings Catalog Resource from model", data)

	profile := graphmodels.NewDeviceManagementConfigurationPolicy()

	Name := data.Name.ValueString()
	description := data.Description.ValueString()
	profile.SetName(&Name)
	profile.SetDescription(&description)

	platformStr := data.Platforms.ValueString()
	var platform graphmodels.DeviceManagementConfigurationPlatforms
	switch platformStr {
	case "android":
		platform = graphmodels.ANDROID_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "androidEnterprise":
		platform = graphmodels.ANDROIDENTERPRISE_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "aosp":
		platform = graphmodels.AOSP_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "iOS":
		platform = graphmodels.IOS_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "linux":
		platform = graphmodels.LINUX_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "macOS":
		platform = graphmodels.MACOS_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "windows10":
		platform = graphmodels.WINDOWS10_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	case "windows10X":
		platform = graphmodels.WINDOWS10X_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	}
	profile.SetPlatforms(&platform)

	var technologiesStr []string
	for _, tech := range data.Technologies {
		technologiesStr = append(technologiesStr, tech.ValueString())
	}
	parsedTechnologies, _ := graphmodels.ParseDeviceManagementConfigurationTechnologies(strings.Join(technologiesStr, ","))
	profile.SetTechnologies(parsedTechnologies.(*graphmodels.DeviceManagementConfigurationTechnologies))

	if len(data.RoleScopeTagIds) > 0 {
		var tagIds []string
		for _, tag := range data.RoleScopeTagIds {
			tagIds = append(tagIds, tag.ValueString())
		}
		profile.SetRoleScopeTagIds(tagIds)
	} else {
		profile.SetRoleScopeTagIds([]string{"0"})
	}

	// Construct settings and set them to profile
	settings := constructSettingsCatalogSettings(ctx, data.Settings)
	profile.SetSettings(settings)

	tflog.Debug(ctx, "Finished constructing Windows Settings Catalog resource")
	return profile, nil
}

func constructSettingsCatalogSettings(ctx context.Context, settingsJSON types.String) []graphmodels.DeviceManagementConfigurationSettingable {
	tflog.Debug(ctx, "Constructing settings catalog settings")

	// Parse the settings structure
	var settingsData struct {
		SettingsDetails []struct {
			SettingInstance map[string]interface{} `json:"settingInstance"`
		} `json:"settingsDetails"`
	}

	if err := json.Unmarshal([]byte(settingsJSON.ValueString()), &settingsData); err != nil {
		tflog.Error(ctx, "Failed to unmarshal settings JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}

	// Create settings collection
	settingsCollection := make([]graphmodels.DeviceManagementConfigurationSettingable, 0, len(settingsData.SettingsDetails))

	// Get parse node factory registry
	parseNodeRegistry := serialization.NewParseNodeFactoryRegistry()

	// Lock the registry before modifying
	parseNodeRegistry.Lock()
	// Add the JSON factory to the registry
	parseNodeFactory := jsonserialization.NewJsonParseNodeFactory()
	parseNodeRegistry.ContentTypeAssociatedFactories["application/json"] = parseNodeFactory
	parseNodeRegistry.Unlock()

	// Process each setting instance
	for _, detail := range settingsData.SettingsDetails {
		// Convert setting instance to JSON string
		settingJSON, err := json.Marshal(detail.SettingInstance)
		if err != nil {
			tflog.Error(ctx, "Failed to marshal setting instance", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		// Get root parse node
		parseNode, err := parseNodeRegistry.GetRootParseNode("application/json", settingJSON)
		if err != nil {
			tflog.Error(ctx, "Failed to get root parse node", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		// Create setting from parse node
		parsable, err := graphmodels.CreateDeviceManagementConfigurationSettingFromDiscriminatorValue(parseNode)
		if err != nil {
			tflog.Error(ctx, "Failed to create setting from discriminator value", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		// Convert to correct type
		setting, ok := parsable.(graphmodels.DeviceManagementConfigurationSettingable)
		if !ok {
			tflog.Error(ctx, "Failed to convert parsable to DeviceManagementConfigurationSettingable")
			continue
		}

		settingsCollection = append(settingsCollection, setting)
	}

	return settingsCollection
}
