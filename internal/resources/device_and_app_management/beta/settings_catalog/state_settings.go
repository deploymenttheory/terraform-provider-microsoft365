package graphBetaSettingsCatalog

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, remoteSettings []graphmodels.DeviceManagementConfigurationSettingable) {
	if remoteSettings == nil {
		tflog.Debug(ctx, "Remote settings are nil")
		return
	}

	tflog.Debug(ctx, "Starting to map settings state to Terraform state")

	// Convert to JSON
	jsonData, err := json.Marshal(&DeviceConfigV2GraphServiceModel)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal settings data to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Store in Terraform state
	data.Settings = types.StringValue(string(jsonData))

	tflog.Debug(ctx, "Finished mapping settings state to Terraform state")
}
