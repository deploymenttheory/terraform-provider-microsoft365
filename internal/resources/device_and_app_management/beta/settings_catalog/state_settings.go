package graphBetaSettingsCatalog

import (
	"context"
	"encoding/json"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/normalize"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteSettingsStateToTerraform maps the remote settings catalog settings state to the Terraform state
// using the shared DeviceConfigV2GraphServiceMode used also for requests.
// Results are normalized and stored alphabetically in the Terraform state.
func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, remoteSettings []graphmodels.DeviceManagementConfigurationSettingable) {
	if remoteSettings == nil {
		tflog.Debug(ctx, "Remote settings are nil")
		return
	}

	tflog.Debug(ctx, "Starting to map settings state to Terraform state")

	jsonData, err := json.Marshal(&DeviceConfigV2GraphServiceModel)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal settings data to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	normalizedJSON, err := normalize.JSONAlphabetically(string(jsonData))
	if err != nil {
		tflog.Error(ctx, "Failed to normalize JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	data.Settings = types.StringValue(normalizedJSON)

	tflog.Debug(ctx, "Finished mapping settings state to Terraform state")
}
