package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SetNestedFieldNullWhenValidator validates that within each element of a set-nested attribute,
// a specific field must be null whenever another field equals a given trigger value.
type SetNestedFieldNullWhenValidator struct {
	// ParentPath is a guard path — validation is skipped when the parent block is absent.
	ParentPath path.Path
	// SetPath is the full path to the set-nested attribute.
	SetPath path.Path
	// TriggerField is the field whose value conditionally forbids ForbiddenField.
	TriggerField string
	// TriggerValue is the value of TriggerField that activates the constraint.
	TriggerValue string
	// ForbiddenField is the field that must be null when TriggerField equals TriggerValue.
	ForbiddenField string
}

// SetNestedFieldNullWhen returns a resource.ConfigValidator that ensures a field within each
// element of a set-nested attribute is null whenever another field equals a trigger value.
//
// Example usage — threshold must not be set when action is "offerFallback":
//
//	resourcevalidator.SetNestedFieldNullWhen(
//	    path.Root("settings"),
//	    path.Root("settings").AtName("monitoring").AtName("monitoring_rules"),
//	    "action",
//	    "offerFallback",
//	    "threshold",
//	)
func SetNestedFieldNullWhen(parentPath, setPath path.Path, triggerField, triggerValue, forbiddenField string) resource.ConfigValidator {
	return &SetNestedFieldNullWhenValidator{
		ParentPath:     parentPath,
		SetPath:        setPath,
		TriggerField:   triggerField,
		TriggerValue:   triggerValue,
		ForbiddenField: forbiddenField,
	}
}

// Description describes the validation in plain text formatting.
func (v SetNestedFieldNullWhenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("%s: %q must not be set when %s is %q",
		v.SetPath, v.ForbiddenField, v.TriggerField, v.TriggerValue)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v SetNestedFieldNullWhenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateResource performs the validation.
func (v SetNestedFieldNullWhenValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if req.Config.Raw.IsNull() {
		return
	}

	var parentObj types.Object
	if diags := req.Config.GetAttribute(ctx, v.ParentPath, &parentObj); diags.HasError() || parentObj.IsNull() || parentObj.IsUnknown() {
		return
	}

	var rulesSet types.Set
	if diags := req.Config.GetAttribute(ctx, v.SetPath, &rulesSet); diags.HasError() || rulesSet.IsNull() || rulesSet.IsUnknown() {
		return
	}

	for _, elem := range rulesSet.Elements() {
		obj, ok := elem.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			continue
		}

		attrs := obj.Attributes()

		triggerVal, ok := attrs[v.TriggerField].(types.String)
		if !ok || triggerVal.IsNull() || triggerVal.IsUnknown() {
			continue
		}

		if triggerVal.ValueString() != v.TriggerValue {
			continue
		}

		// TriggerField matches — ForbiddenField must be null or absent.
		forbidden, exists := attrs[v.ForbiddenField]
		if !exists {
			continue
		}

		if isNullOrUnknown(forbidden) {
			continue
		}

		resp.Diagnostics.AddAttributeError(
			v.SetPath,
			"Invalid field combination in set element",
			fmt.Sprintf(
				"%q must not be set when %s is %q.",
				v.ForbiddenField, v.TriggerField, v.TriggerValue,
			),
		)
	}
}

// isNullOrUnknown returns true if the attr.Value is null or unknown.
func isNullOrUnknown(val interface{ IsNull() bool }) bool {
	type unknowner interface {
		IsUnknown() bool
	}
	if val.IsNull() {
		return true
	}
	if u, ok := val.(unknowner); ok && u.IsUnknown() {
		return true
	}
	return false
}
