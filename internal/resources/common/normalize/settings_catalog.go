package normalize

import (
	"fmt"
	"reflect"
)

// PreserveSecretSettings recursively searches through settings catalog HCL JSON structure for secret settings
// and preserves the value and valueState from the config settings. This is performed recursively throughout the JSON
// settings catalog and It returns an error if any unexpected data types or mismatches are encountered.
func PreserveSecretSettings(config, resp interface{}) error {
	switch configV := config.(type) {
	case map[string]interface{}:
		respV, ok := resp.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map[string]interface{} in response, got %s", reflect.TypeOf(resp))
		}

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
				if err := PreserveSecretSettings(v, respChild); err != nil {
					return fmt.Errorf("error in key %q: %w", k, err)
				}
			}
		}

	case []interface{}:
		respV, ok := resp.([]interface{})
		if !ok {
			return fmt.Errorf("expected []interface{} in response, got %s", reflect.TypeOf(resp))
		}
		for i := range configV {
			if i < len(respV) {
				if err := PreserveSecretSettings(configV[i], respV[i]); err != nil {
					return fmt.Errorf("error in array index %d: %w", i, err)
				}
			}
		}

	default:
		return fmt.Errorf("unsupported type: %s", reflect.TypeOf(config))
	}

	return nil
}
