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
// func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, resp []byte) {
// 	if err := json.Unmarshal(resp, &DeviceConfigV2GraphServiceModel); err != nil {
// 		tflog.Error(ctx, "Failed to unmarshal settings", map[string]interface{}{"error": err.Error()})
// 		return
// 	}

// 	normalizedJSON, err := normalize.JSONAlphabetically(string(resp))
// 	if err != nil {
// 		tflog.Error(ctx, "Failed to normalize JSON", map[string]interface{}{"error": err.Error()})
// 		return
// 	}

// 	data.Settings = types.StringValue(normalizedJSON)
// }

func MapRemoteSettingsStateToTerraform(ctx context.Context, data *SettingsCatalogProfileResourceModel, resp []byte) {
	// Parse the raw response into a temporary map
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(resp, &rawResponse); err != nil {
		tflog.Error(ctx, "Failed to unmarshal settings response", map[string]interface{}{"error": err.Error()})
		return
	}

	// Extract the "value" property, if it exists
	value, ok := rawResponse["value"]
	if !ok {
		tflog.Error(ctx, "Response does not contain 'value' field", nil)
		return
	}

	// Marshal the extracted "value" back to JSON for storage
	valueJSON, err := json.Marshal(value)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal 'value' field", map[string]interface{}{"error": err.Error()})
		return
	}

	// Normalize the JSON
	normalizedJSON, err := normalize.JSONAlphabetically(string(valueJSON))
	if err != nil {
		tflog.Error(ctx, "Failed to normalize JSON", map[string]interface{}{"error": err.Error()})
		return
	}

	// Store the normalized JSON string in Terraform state
	data.Settings = types.StringValue(normalizedJSON)
}
