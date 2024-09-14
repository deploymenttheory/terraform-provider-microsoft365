package construct

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// ComplexNestedStruct represents a complex nested structure for testing
type ComplexNestedStruct struct {
	StringField  types.String   `tfsdk:"string_field"`
	IntField     types.Int64    `tfsdk:"int_field"`
	FloatField   types.Float64  `tfsdk:"float_field"`
	BoolField    types.Bool     `tfsdk:"bool_field"`
	ListField    types.List     `tfsdk:"list_field"`
	SetField     types.Set      `tfsdk:"set_field"`
	MapField     types.Map      `tfsdk:"map_field"`
	NestedStruct *NestedStruct  `tfsdk:"nested_struct"`
	SliceField   []types.String `tfsdk:"slice_field"`
	IgnoredField string         `tfsdk:"-"`
}

type NestedStruct struct {
	NestedString types.String `tfsdk:"nested_string"`
	NestedInt    types.Int64  `tfsdk:"nested_int"`
}

func TestComplexStructReflection(t *testing.T) {
	// Create a complex nested struct
	complexStruct := ComplexNestedStruct{
		StringField: types.StringValue("test string"),
		IntField:    types.Int64Value(42),
		FloatField:  types.Float64Value(3.14),
		BoolField:   types.BoolValue(true),
		ListField: types.ListValueMust(
			types.StringType,
			[]attr.Value{
				types.StringValue("item1"),
				types.StringValue("item2"),
			},
		),
		SetField: types.SetValueMust(
			types.Int64Type,
			[]attr.Value{
				types.Int64Value(1),
				types.Int64Value(2),
			},
		),
		MapField: types.MapValueMust(
			types.StringType,
			map[string]attr.Value{
				"key1": types.StringValue("value1"),
				"key2": types.StringValue("value2"),
			},
		),
		NestedStruct: &NestedStruct{
			NestedString: types.StringValue("nested string"),
			NestedInt:    types.Int64Value(99),
		},
		SliceField: []types.String{
			types.StringValue("slice1"),
			types.StringValue("slice2"),
		},
		IgnoredField: "this should be ignored",
	}

	// Call the function to convert the struct to a map
	result := structToMap(reflect.ValueOf(complexStruct))

	// Convert the result to JSON for easier assertion
	jsonResult, err := json.Marshal(result)
	assert.NoError(t, err)

	// Define the expected JSON structure
	expectedJSON := `{
		"string_field": "test string",
		"int_field": 42,
		"float_field": 3.14,
		"bool_field": true,
		"list_field": ["item1", "item2"],
		"set_field": [1, 2],
		"map_field": {
			"key1": "value1",
			"key2": "value2"
		},
		"nested_struct": {
			"nested_string": "nested string",
			"nested_int": 99
		},
		"slice_field": ["slice1", "slice2"]
	}`

	// Compare the result with the expected JSON
	assert.JSONEq(t, expectedJSON, string(jsonResult))
}

func TestHandleTerraformTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    attr.Value
		expected interface{}
	}{
		{
			name:     "String",
			input:    types.StringValue("test"),
			expected: "test",
		},
		{
			name:     "Int64",
			input:    types.Int64Value(42),
			expected: int64(42),
		},
		{
			name:     "Float64",
			input:    types.Float64Value(3.14),
			expected: 3.14,
		},
		{
			name:     "Bool",
			input:    types.BoolValue(true),
			expected: true,
		},
		{
			name:     "Null",
			input:    types.StringNull(),
			expected: nil,
		},
		{
			name:     "Unknown",
			input:    types.StringUnknown(),
			expected: nil,
		},
		{
			name: "List",
			input: types.ListValueMust(
				types.StringType,
				[]attr.Value{types.StringValue("a"), types.StringValue("b")},
			),
			expected: []interface{}{"a", "b"},
		},
		{
			name: "Set",
			input: types.SetValueMust(
				types.Int64Type,
				[]attr.Value{types.Int64Value(1), types.Int64Value(2)},
			),
			expected: []interface{}{int64(1), int64(2)},
		},
		{
			name: "Map",
			input: types.MapValueMust(
				types.StringType,
				map[string]attr.Value{
					"key1": types.StringValue("value1"),
					"key2": types.StringValue("value2"),
				},
			),
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handleTerraformValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDebugPrintStruct(t *testing.T) {
	// Create a simple struct for testing
	testStruct := struct {
		Field1 types.String `tfsdk:"field1"`
		Field2 types.Int64  `tfsdk:"field2"`
	}{
		Field1: types.StringValue("test"),
		Field2: types.Int64Value(42),
	}

	// Call DebugPrintStruct
	ctx := context.Background()
	DebugPrintStruct(ctx, "TestStruct", testStruct)

	// Since DebugPrintStruct doesn't return anything, we can't directly test its output
	// However, we can ensure it doesn't panic
	// In a real scenario, you might want to capture the log output and assert on it
	// This would require modifying the DebugPrintStruct function to accept a custom logger
}
