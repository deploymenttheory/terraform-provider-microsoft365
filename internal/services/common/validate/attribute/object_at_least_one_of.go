package attribute

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// atLeastOneOfNonEmptyValidator validates that a configured object sets at least one of the named attributes to a
// non-empty value.
type atLeastOneOfNonEmptyValidator struct {
	attributeNames    []string
	validationMessage string
}

// Description describes the validation in plain text formatting.
func (v atLeastOneOfNonEmptyValidator) Description(_ context.Context) string {
	return fmt.Sprintf("When this block is set, at least one of [%s] must be set to a non-empty value", strings.Join(v.attributeNames, ", "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v atLeastOneOfNonEmptyValidator) MarkdownDescription(_ context.Context) string {
	names := make([]string, len(v.attributeNames))
	for i, name := range v.attributeNames {
		names[i] = fmt.Sprintf("`%s`", name)
	}
	return fmt.Sprintf("When this block is set, at least one of %s must be set to a non-empty value", strings.Join(names, ", "))
}

// ValidateObject performs the validation. A null or unknown block, or a block where one of the named attributes is
// still unknown, is not rejected: it cannot be shown to be empty yet.
func (v atLeastOneOfNonEmptyValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attributes := req.ConfigValue.Attributes()
	for _, name := range v.attributeNames {
		value, ok := attributes[name]
		if !ok || value == nil || value.IsNull() {
			continue
		}
		if value.IsUnknown() || !isEmptyCollection(ctx, value) {
			return
		}
	}

	message := v.validationMessage
	if message == "" {
		message = fmt.Sprintf("At least one of these attributes must be set to a non-empty value: %s.", strings.Join(v.attributeNames, ", "))
	}

	resp.Diagnostics.AddAttributeError(req.Path, "Missing Required Attribute", message)
}

// isEmptyCollection reports whether value is a known set, list or map with no elements.
func isEmptyCollection(ctx context.Context, value attr.Value) bool {
	switch collection := value.(type) {
	case basetypes.SetValuable:
		set, diags := collection.ToSetValue(ctx)
		return !diags.HasError() && len(set.Elements()) == 0
	case basetypes.ListValuable:
		list, diags := collection.ToListValue(ctx)
		return !diags.HasError() && len(list.Elements()) == 0
	case basetypes.MapValuable:
		mapValue, diags := collection.ToMapValue(ctx)
		return !diags.HasError() && len(mapValue.Elements()) == 0
	default:
		return false
	}
}

// AtLeastOneOfNonEmpty returns an object validator which ensures that a configured block sets at least one of the
// named attributes. A null attribute, or an empty set, list or map, does not count; any other known value, including
// a nested object, does. validationMessage replaces the default error detail when it is not empty.
//
// Use it for a block that Graph stores as null when none of its settings carries a value, so an empty block could
// not be read back.
//
// Example usage:
//
//	attribute.AtLeastOneOfNonEmpty(
//	    "grant_controls sets no control. Remove the block.",
//	    "built_in_controls", "custom_authentication_factors", "terms_of_use", "authentication_strength",
//	)
func AtLeastOneOfNonEmpty(validationMessage string, attributeNames ...string) validator.Object {
	return &atLeastOneOfNonEmptyValidator{
		attributeNames:    attributeNames,
		validationMessage: validationMessage,
	}
}
