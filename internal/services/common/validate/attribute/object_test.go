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
