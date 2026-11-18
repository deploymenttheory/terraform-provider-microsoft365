package attribute

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExactlyOneOfMixedTypes_Description(t *testing.T) {
	t.Parallel()

	v := ExactlyOneOfMixedTypes("string_value", "int_value", "bool_value")

	expected := "Exactly one of [string_value, int_value, bool_value] must be specified"
	if desc := v.Description(context.Background()); desc != expected {
		t.Errorf("Expected description %q, got %q", expected, desc)
	}
}

func TestExactlyOneOfMixedTypes_MarkdownDescription(t *testing.T) {
	t.Parallel()

	v := ExactlyOneOfMixedTypes("string_value", "int_value", "bool_value")

	expected := "Exactly one of `string_value`, `int_value`, `bool_value` must be specified"
	if desc := v.MarkdownDescription(context.Background()); desc != expected {
		t.Errorf("Expected markdown description %q, got %q", expected, desc)
	}
}

func TestExactlyOneOfMixedTypes_NullObject(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	objValue := types.ObjectNull(map[string]attr.Type{
		"string_value": types.StringType,
		"bool_value":   types.BoolType,
	})

	req := validator.ObjectRequest{
		Path:        path.Root("test"),
		ConfigValue: objValue,
	}

	resp := &validator.ObjectResponse{}

	v := ExactlyOneOfMixedTypes("string_value", "bool_value")
	v.ValidateObject(ctx, req, resp)

	// Null objects should not produce errors (they're skipped)
	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors for null object, got: %v", resp.Diagnostics.Errors())
	}
}

func TestExactlyOneOfMixedTypes_UnknownObject(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	objValue := types.ObjectUnknown(map[string]attr.Type{
		"string_value": types.StringType,
		"bool_value":   types.BoolType,
	})

	req := validator.ObjectRequest{
		Path:        path.Root("test"),
		ConfigValue: objValue,
	}

	resp := &validator.ObjectResponse{}

	v := ExactlyOneOfMixedTypes("string_value", "bool_value")
	v.ValidateObject(ctx, req, resp)

	// Unknown objects should not produce errors (they're skipped)
	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors for unknown object, got: %v", resp.Diagnostics.Errors())
	}
}

func TestExactlyOneOfMixedTypes_SingleAttribute(t *testing.T) {
	t.Parallel()

	// Test with only one attribute name - should just work without issues
	v := ExactlyOneOfMixedTypes("string_value")

	expected := "Exactly one of [string_value] must be specified"
	if desc := v.Description(context.Background()); desc != expected {
		t.Errorf("Expected description %q, got %q", expected, desc)
	}
}

func TestExactlyOneOfMixedTypes_EmptyAttributes(t *testing.T) {
	t.Parallel()

	// Test with no attribute names
	v := ExactlyOneOfMixedTypes()

	expected := "Exactly one of [] must be specified"
	if desc := v.Description(context.Background()); desc != expected {
		t.Errorf("Expected description %q, got %q", expected, desc)
	}
}

func TestAtLeastOneOfNonEmpty(t *testing.T) {
	t.Parallel()

	strengthType := types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType}}
	attrTypes := map[string]attr.Type{
		"operator":          types.StringType,
		"built_in_controls": types.SetType{ElemType: types.StringType},
		"names":             types.ListType{ElemType: types.StringType},
		"tags":              types.MapType{ElemType: types.StringType},
		"strength":          strengthType,
	}
	emptySet := types.SetValueMust(types.StringType, []attr.Value{})
	emptyList := types.ListValueMust(types.StringType, []attr.Value{})
	emptyMap := types.MapValueMust(types.StringType, map[string]attr.Value{})
	noStrength := types.ObjectNull(strengthType.AttrTypes)
	object := func(values map[string]attr.Value) types.Object {
		all := map[string]attr.Value{
			"operator":          types.StringValue("OR"),
			"built_in_controls": types.SetNull(types.StringType),
			"names":             types.ListNull(types.StringType),
			"tags":              types.MapNull(types.StringType),
			"strength":          noStrength,
		}
		for name, value := range values {
			all[name] = value
		}
		return types.ObjectValueMust(attrTypes, all)
	}

	for name, tc := range map[string]struct {
		value     types.Object
		wantError bool
	}{
		"null object":    {value: types.ObjectNull(attrTypes)},
		"unknown object": {value: types.ObjectUnknown(attrTypes)},
		"nothing set":    {value: object(nil), wantError: true},
		"empty collections": {
			value:     object(map[string]attr.Value{"built_in_controls": emptySet, "names": emptyList, "tags": emptyMap}),
			wantError: true,
		},
		"unknown after empty": {
			value: object(map[string]attr.Value{"built_in_controls": emptySet, "names": types.ListUnknown(types.StringType)}),
		},
		"non-empty set": {
			value: object(map[string]attr.Value{"built_in_controls": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("mfa")})}),
		},
		"non-empty list": {
			value: object(map[string]attr.Value{"names": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a")})}),
		},
		"nested object": {
			value: object(map[string]attr.Value{
				"built_in_controls": emptySet,
				"strength":          types.ObjectValueMust(strengthType.AttrTypes, map[string]attr.Value{"id": types.StringValue("x")}),
			}),
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &validator.ObjectResponse{}
			AtLeastOneOfNonEmpty("built_in_controls", "names", "tags", "strength").ValidateObject(
				context.Background(),
				validator.ObjectRequest{Path: path.Root("test"), ConfigValue: tc.value},
				resp,
			)
			if resp.Diagnostics.HasError() != tc.wantError {
				t.Fatalf("expected error %v, got %v", tc.wantError, resp.Diagnostics)
			}
		})
	}
}
