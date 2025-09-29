// Main entry point to construct the intune settings catalog profile resource for the Terraform provider.
package graphBetaReuseablePolicySettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	sharedConstructor "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *sharedmodels.ReuseablePolicySettingsResourceModel) (graphmodels.DeviceManagementReusablePolicySettingable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementReusablePolicySetting()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	settings := sharedConstructor.ConstructSettingsCatalogSettings(ctx, data.Settings)
	if len(settings) > 0 && settings[0].GetSettingInstance() != nil {
		settingInstance := settings[0].GetSettingInstance()
		requestBody.SetSettingInstance(settingInstance)

		// Sets the required setting definition ID at the root level of the settings catalog req.
		// This logic may need to change when other examples are identifed.
		if settingDefId := settingInstance.GetSettingDefinitionId(); settingDefId != nil {
			requestBody.SetSettingDefinitionId(settingDefId)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
