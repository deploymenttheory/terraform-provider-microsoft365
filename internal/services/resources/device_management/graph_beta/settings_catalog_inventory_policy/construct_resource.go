package graphBetaSettingsCatalogInventoryPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	configPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *InventoryPolicyResourceModel) (graphmodels.DeviceManagementConfigurationPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementConfigurationPolicy()

	convert.FrameworkToGraphString(data.Name, requestBody.SetName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	platformStr := data.Platforms.ValueString()
	var platform graphmodels.DeviceManagementConfigurationPlatforms
	switch platformStr {
	case "none":
		platform = graphmodels.NONE_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
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
	case "unknownFutureValue":
		platform = graphmodels.UNKNOWNFUTUREVALUE_DEVICEMANAGEMENTCONFIGURATIONPLATFORMS
	}
	requestBody.SetPlatforms(&platform)

	// "extensibility" is not recognized by the SDK enum, so set it via AdditionalData
	techStr := data.Technologies.ValueString()
	additionalData := requestBody.GetAdditionalData()
	if additionalData == nil {
		additionalData = make(map[string]any)
	}
	additionalData["technologies"] = techStr
	requestBody.SetAdditionalData(additionalData)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if data.ConfigurationPolicy != nil {
		settings := configPolicy.ConstructSettingsCatalogSettings(ctx, *data.ConfigurationPolicy)
		requestBody.SetSettings(settings)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
