package graphBetaSettingsCatalog

import (
	"context"
	"encoding/json"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/normalize"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteSettingsStateToTerraform maps the remote settings catalog settings state to the Terraform state
// using the shared DeviceConfigV2GraphServiceMode used also for requests.
// Results are normalized and stored alphabetically in the Terraform state.
func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, resp []byte) {
	if err := json.Unmarshal(resp, &DeviceConfigV2GraphServiceModel); err != nil {
		tflog.Error(ctx, "Failed to unmarshal settings", map[string]interface{}{"error": err.Error()})
		return
	}

	normalizedJSON, err := normalize.JSONAlphabetically(string(resp))
	if err != nil {
		tflog.Error(ctx, "Failed to normalize JSON", map[string]interface{}{"error": err.Error()})
		return
	}

	data.Settings = types.StringValue(normalizedJSON)
}
