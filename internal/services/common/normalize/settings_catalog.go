package normalize

import (
	"fmt"
)

// PreserveSecretSettings recursively searches through the server response for settings catalog and if a match for
// #microsoft.graph.deviceManagementConfigurationSecretSettingValue is found, it preserves the original secret value and state from the config settings.
func PreserveSecretSettings(config, resp any) error {
	if config == nil || resp == nil {
		return nil
	}

	switch configV := config.(type) {
	case map[string]any:
		respV, ok := resp.(map[string]any)
		if !ok {
			return fmt.Errorf("expected map[string]any in response, got %T", resp)
		}

		// If this is a secret setting, preserve the config values
		if odataType, ok := configV["@odata.type"].(string); ok &&
			odataType == "#microsoft.graph.deviceManagementConfigurationSecretSettingValue" {
			if value, ok := configV["value"]; ok {
				respV["value"] = value
			}
			if valueState, ok := configV["valueState"]; ok {
				respV["valueState"] = valueState
			}
			return nil
		}

		for k, v := range configV {
			if respChild, ok := respV[k]; ok {
				PreserveSecretSettings(v, respChild)
			}
		}

	case []any:
		respV, ok := resp.([]any)
		if !ok {
			return fmt.Errorf("expected []any in response, got %T", resp)
		}

		for i := range configV {
			if i < len(respV) {
				PreserveSecretSettings(configV[i], respV[i])
			}
		}
	}

	return nil
}
