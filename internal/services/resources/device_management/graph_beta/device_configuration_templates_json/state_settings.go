package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_management"
)

// MapRemoteSettingsStateToTerraform restores writable settings, including values and
// relationships exposed by separate API calls. It never filters to prior configured keys.
func (r *DeviceConfigurationTemplatesJsonResource) MapRemoteSettingsStateToTerraform(
	ctx context.Context,
	data *DeviceConfigurationTemplatesJsonResourceModel,
	response []byte,
) error {
	var settings map[string]any
	if err := json.Unmarshal(response, &settings); err != nil {
		return fmt.Errorf("decode template settings: %w", err)
	}
	if settings == nil {
		return fmt.Errorf("%w: expected a JSON object", errInvalidProfileResponse)
	}
	if err := r.readRelationships(ctx, data.ID.ValueString(), data.Settings, settings); err != nil {
		return err
	}
	if settings["@odata.type"] == "#microsoft.graph.windows10CustomConfiguration" {
		if err := r.readOmaSettingValues(ctx, data.ID.ValueString(), settings); err != nil {
			return err
		}
	}

	value, err := sharedstater.StateDeviceConfigurationSettings(data.Settings, settings)
	if err != nil {
		return fmt.Errorf("map template settings: %w", err)
	}
	data.Settings = value
	return nil
}

func (r *DeviceConfigurationTemplatesJsonResource) readOmaSettingValues(
	ctx context.Context,
	id string,
	settings map[string]any,
) error {
	entries, _ := settings["omaSettings"].([]any)
	for _, entry := range entries {
		setting, ok := entry.(map[string]any)
		if !ok {
			return fmt.Errorf("%w: omaSettings must contain objects", errInvalidProfileResponse)
		}
		if encrypted, _ := setting["isEncrypted"].(bool); encrypted {
			reference, _ := setting["secretReferenceValueId"].(string)
			if reference == "" {
				return fmt.Errorf(
					"%w: encrypted OMA setting has no secret reference",
					errInvalidProfileResponse,
				)
			}
			plainText, err := r.client.DeviceManagement().
				DeviceConfigurations().
				ByDeviceConfigurationId(id).
				GetOmaSettingPlainTextValueWithSecretReferenceValueId(&reference).
				GetAsGetOmaSettingPlainTextValueWithSecretReferenceValueIdGetResponse(ctx, nil)
			if err != nil {
				return fmt.Errorf("read encrypted OMA setting value: %w", err)
			}
			if plainText == nil || plainText.GetValue() == nil {
				return fmt.Errorf(
					"%w: plaintext endpoint returned no value",
					errInvalidProfileResponse,
				)
			}
			if setting["@odata.type"] == "#microsoft.graph.omaSettingStringXml" {
				setting["value"] = base64.StdEncoding.EncodeToString([]byte(*plainText.GetValue()))
			} else {
				setting["value"] = *plainText.GetValue()
			}
		}
		for _, field := range []string{"isEncrypted", "secretReferenceValueId", "isReadOnly"} {
			delete(setting, field)
		}
	}
	return nil
}
