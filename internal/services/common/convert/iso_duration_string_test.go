package convert

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestFrameworkToGraphISODurationStringSerialization(t *testing.T) {
	for _, raw := range []string{"PT0S", "P0D", "PT30S", "PT1H30M15S", "PT12H", "P1D", "P7D", "P1W", "P1Y2M3DT4H5M6S"} {
		t.Run(raw, func(t *testing.T) {
			body := models.NewDeviceShellScript()
			err := FrameworkToGraphISODurationString(types.StringValue(raw), func(value *string) { body.GetAdditionalData()["executionFrequency"] = *value })
			require.NoError(t, err)
			require.Nil(t, body.GetExecutionFrequency())
			writer := jsonserialization.NewJsonSerializationWriter()
			require.NoError(t, writer.WriteObjectValue("", body))
			encoded, err := writer.GetSerializedContent()
			require.NoError(t, err)
			var payload map[string]any
			require.NoError(t, json.Unmarshal(encoded, &payload))
			require.Equal(t, raw, payload["executionFrequency"])
			require.Equal(t, 1, strings.Count(string(encoded), `"executionFrequency"`))
		})
	}
}

func TestFrameworkToGraphISODurationStringUnset(t *testing.T) {
	for _, value := range []types.String{types.StringNull(), types.StringUnknown()} {
		called := false
		require.NoError(t, FrameworkToGraphISODurationString(value, func(*string) { called = true }))
		require.False(t, called)
	}
}

func TestFrameworkToGraphISODurationStringInvalid(t *testing.T) {
	for _, raw := range []string{"", "P", "PT", "P1DT", "InvalidDuration", "PT0Ssuffix", "P1W1D"} {
		t.Run(raw, func(t *testing.T) {
			called := false
			require.Error(t, FrameworkToGraphISODurationString(types.StringValue(raw), func(*string) { called = true }))
			require.False(t, called)
		})
	}
}
