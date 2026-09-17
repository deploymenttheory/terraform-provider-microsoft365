package graphBetaLinuxPlatformScript

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUnitLinuxPlatformScriptSettingsState(t *testing.T) {
	for _, tc := range []struct {
		name      string
		fixture   string
		context   string
		frequency string
		retries   string
		script    string
	}{
		{
			name: "minimal", fixture: "get_settings_minimal.json",
			context: "user", frequency: "15minutes", retries: "1",
			script: "#!/bin/bash\nprintf \"Linux platform script test\\n\"\n",
		},
		{
			name: "maximal", fixture: "get_settings_maximal.json",
			context: "root", frequency: "1day", retries: "3",
			script: "#!/bin/bash\n# Updated script: café\nprintf \"Linux platform script updated\\n\"\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := os.ReadFile("tests/responses/validate_read/" + tc.fixture)
			require.NoError(t, err)
			// Remote values must replace stale state as well as populate imports.
			state := LinuxPlatformScriptResourceModel{
				ScriptContent: types.StringValue("old script"), ExecutionContext: types.StringValue("old context"),
				ExecutionFrequency: types.StringValue("old frequency"), ExecutionRetries: types.StringValue("old retries"),
			}
			require.NoError(t, MapRemoteSettingsStateToTerraform(context.Background(), &state, raw))
			require.Equal(t, tc.context, state.ExecutionContext.ValueString())
			require.Equal(t, tc.frequency, state.ExecutionFrequency.ValueString())
			require.Equal(t, tc.retries, state.ExecutionRetries.ValueString())
			require.Equal(t, tc.script, state.ScriptContent.ValueString())
		})
	}
}

func TestUnitLinuxPlatformScriptSettingsStateInvalidResponse(t *testing.T) {
	raw, err := os.ReadFile("tests/responses/validate_read/get_settings_minimal.json")
	require.NoError(t, err)
	var invalidScript map[string]any
	require.NoError(t, json.Unmarshal(raw, &invalidScript))
	invalidScript["value"].([]any)[3].(map[string]any)["settingInstance"].(map[string]any)["simpleSettingValue"].(map[string]any)["value"] = "invalid-base64!"
	invalidScriptJSON, err := json.Marshal(invalidScript)
	require.NoError(t, err)

	for name, response := range map[string][]byte{
		"malformed_json": []byte("{"),
		"empty_settings": []byte("{\"value\":[]}"),
		"graph_error":    []byte("{\"error\":{\"code\":\"BadRequest\",\"message\":\"Invalid settings request\"}}"),
		"invalid_script": invalidScriptJSON,
	} {
		t.Run(name, func(t *testing.T) {
			state := LinuxPlatformScriptResourceModel{
				ScriptContent: types.StringValue("old script"), ExecutionContext: types.StringValue("user"),
				ExecutionFrequency: types.StringValue("15minutes"), ExecutionRetries: types.StringValue("1"),
			}
			before := state
			require.Error(t, MapRemoteSettingsStateToTerraform(context.Background(), &state, response))
			require.Equal(t, before, state, "an invalid settings response must preserve existing values")
		})
	}
}
