package attribute

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var atLeastOneOfTestTypes = map[string]attr.Type{
	"operator":                types.StringType,
	"built_in_controls":       types.SetType{ElemType: types.StringType},
	"terms_of_use":            types.ListType{ElemType: types.StringType},
	"authentication_strength": types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType}},
}

func atLeastOneOfTestObject(builtIn, terms, strength attr.Value) types.Object {
	return types.ObjectValueMust(atLeastOneOfTestTypes, map[string]attr.Value{
		"operator":                types.StringValue("OR"),
		"built_in_controls":       builtIn,
		"terms_of_use":            terms,
		"authentication_strength": strength,
	})
}

func TestAtLeastOneOfNonEmpty(t *testing.T) {
	t.Parallel()

	strengthType := map[string]attr.Type{"id": types.StringType}
	emptySet := types.SetValueMust(types.StringType, []attr.Value{})
	mfaSet := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("mfa")})
	nullSet := types.SetNull(types.StringType)
	emptyList := types.ListValueMust(types.StringType, []attr.Value{})
	nullList := types.ListNull(types.StringType)
	termsList := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("terms-id")})
	nullStrength := types.ObjectNull(strengthType)
	strength := types.ObjectValueMust(strengthType, map[string]attr.Value{"id": types.StringValue("00000000-0000-0000-0000-000000000002")})

	tests := []struct {
		name      string
		value     types.Object
		wantError bool
	}{
		{name: "null block", value: types.ObjectNull(atLeastOneOfTestTypes)},
		{name: "unknown block", value: types.ObjectUnknown(atLeastOneOfTestTypes)},
		{name: "only a non-listed attribute set", value: atLeastOneOfTestObject(nullSet, nullList, nullStrength), wantError: true},
		{name: "empty set and empty list", value: atLeastOneOfTestObject(emptySet, emptyList, nullStrength), wantError: true},
		{name: "non-empty set", value: atLeastOneOfTestObject(mfaSet, nullList, nullStrength)},
		{name: "non-empty list", value: atLeastOneOfTestObject(emptySet, termsList, nullStrength)},
		{name: "nested object", value: atLeastOneOfTestObject(emptySet, emptyList, strength)},
		{name: "unknown set", value: atLeastOneOfTestObject(types.SetUnknown(types.StringType), emptyList, nullStrength)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := validator.ObjectRequest{Path: path.Root("grant_controls"), ConfigValue: tt.value}
			resp := &validator.ObjectResponse{}

			AtLeastOneOfNonEmpty("Remove the block.", "built_in_controls", "terms_of_use", "authentication_strength", "not_in_schema").ValidateObject(context.Background(), req, resp)

			if resp.Diagnostics.HasError() != tt.wantError {
				t.Fatalf("expected error %v, got diagnostics: %v", tt.wantError, resp.Diagnostics)
			}
			if tt.wantError && resp.Diagnostics.Errors()[0].Detail() != "Remove the block." {
				t.Errorf("expected the configured message, got %q", resp.Diagnostics.Errors()[0].Detail())
			}
		})
	}
}

func TestAtLeastOneOfNonEmpty_DefaultMessageAndDescriptions(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := AtLeastOneOfNonEmpty("", "a", "b")

	if got, want := v.Description(ctx), "When this block is set, at least one of [a, b] must be set to a non-empty value"; got != want {
		t.Errorf("Description() = %q, want %q", got, want)
	}
	if got, want := v.MarkdownDescription(ctx), "When this block is set, at least one of `a`, `b` must be set to a non-empty value"; got != want {
		t.Errorf("MarkdownDescription() = %q, want %q", got, want)
	}

	value := types.ObjectValueMust(
		map[string]attr.Type{"a": types.MapType{ElemType: types.StringType}, "b": types.StringType},
		map[string]attr.Value{"a": types.MapValueMust(types.StringType, map[string]attr.Value{}), "b": types.StringNull()},
	)
	resp := &validator.ObjectResponse{}
	v.ValidateObject(ctx, validator.ObjectRequest{Path: path.Root("block"), ConfigValue: value}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error for an empty map and a null string")
	}
	if got, want := resp.Diagnostics.Errors()[0].Detail(), "At least one of these attributes must be set to a non-empty value: a, b."; got != want {
		t.Errorf("detail = %q, want %q", got, want)
	}

	value = types.ObjectValueMust(
		map[string]attr.Type{"a": types.MapType{ElemType: types.StringType}, "b": types.StringType},
		map[string]attr.Value{"a": types.MapNull(types.StringType), "b": types.StringValue("set")},
	)
	resp = &validator.ObjectResponse{}
	v.ValidateObject(ctx, validator.ObjectRequest{Path: path.Root("block"), ConfigValue: value}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no error when a string attribute is set, got: %v", resp.Diagnostics)
	}
}
