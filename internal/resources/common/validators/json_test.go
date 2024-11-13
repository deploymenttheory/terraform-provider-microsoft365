package validators

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestJSONSchemaValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         types.String
		expectedError string
	}{
		{
			name:          "valid_simple_json",
			value:         types.StringValue(`{"key": "value"}`),
			expectedError: "",
		},
		{
			name:          "valid_nested_json",
			value:         types.StringValue(`{"level1": {"level2": {"level3": "value"}}}`),
			expectedError: "",
		},
		{
			name:          "invalid_json_syntax",
			value:         types.StringValue(`{"key": "value"`),
			expectedError: "Invalid JSON String",
		},
		{
			name: "too_deep_nesting",
			value: types.StringValue(`{
				"1": {"2": {"3": {"4": {"5": {"6": {"7": {"8": {"9": {"10": {"11": {"12": {"13": {"14": {"15": {"16": {"17": {"18": {"19": {"20": {"21": "value"}}}}}}}}}}}}}}}}}}}}}`),
			expectedError: "JSON Nesting Too Deep",
		},
		{
			name:          "null_value",
			value:         types.StringNull(),
			expectedError: "",
		},
		{
			name:          "unknown_value",
			value:         types.StringUnknown(),
			expectedError: "",
		},
		{
			name:          "deep_array_nesting",
			value:         types.StringValue(`{"array":[[[[[[[[[[[[[[[[[[[[["too deep"]]]]]]]]]]]]]]]]]]]]]}`),
			expectedError: "JSON Nesting Too Deep",
		},
		{
			name: "mixed_nesting",
			value: types.StringValue(`{
				"obj": {"array": [{"nested": {"more": [{"deep": "value"}]}}]}
			}`),
			expectedError: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.value,
			}

			response := &validator.StringResponse{}

			JSONSchemaValidator().ValidateString(context.Background(), request, response)

			if test.expectedError == "" {
				assert.False(t, response.Diagnostics.HasError(), "expected no error")
			} else {
				assert.True(t, response.Diagnostics.HasError(), "expected error")
				assert.Contains(t, response.Diagnostics.Errors()[0].Summary(), test.expectedError)
			}
		})
	}
}

func TestJSONDepth(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		expectedDepth int
	}{
		{
			name:          "flat_object",
			json:          `{"key": "value"}`,
			expectedDepth: 1,
		},
		{
			name:          "nested_object",
			json:          `{"level1": {"level2": "value"}}`,
			expectedDepth: 2,
		},
		{
			name:          "flat_array",
			json:          `["value1", "value2"]`,
			expectedDepth: 1,
		},
		{
			name:          "nested_array",
			json:          `[["value1"], ["value2"]]`,
			expectedDepth: 2,
		},
		{
			name:          "mixed_nesting",
			json:          `{"array": [{"nested": "value"}]}`,
			expectedDepth: 3,
		},
		{
			name:          "single_value",
			json:          `"value"`,
			expectedDepth: 0,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var data interface{}
			err := json.Unmarshal([]byte(test.json), &data)
			assert.NoError(t, err)

			depth := getJSONDepth(data)
			assert.Equal(t, test.expectedDepth, depth, "Depth mismatch for %s", test.name)
		})
	}
}
