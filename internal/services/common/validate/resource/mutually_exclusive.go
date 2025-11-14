package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// MutuallyExclusiveAttributesValidator validates that only one of the specified attributes is configured
type MutuallyExclusiveAttributesValidator struct {
	AttributePaths []path.Path
	AttributeNames []string
}

// MutuallyExclusiveAttributes returns a resource validator which ensures that
// only one of the specified attributes can be configured at a time.
//
// Example usage:
//
//	func (r *MyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
//	    return []resource.ConfigValidator{
//	        resource_level.MutuallyExclusiveAttributes(
//	            []path.Path{
//	                path.Root("encoded_setting_xml"),
//	                path.Root("settings"),
//	            },
//	            []string{
//	                "encoded_setting_xml",
//	                "settings",
//	            },
//	        ),
//	    }
//	}
func MutuallyExclusiveAttributes(attributePaths []path.Path, attributeNames []string) resource.ConfigValidator {
	return &MutuallyExclusiveAttributesValidator{
		AttributePaths: attributePaths,
		AttributeNames: attributeNames,
	}
}

// Description describes the validation in plain text formatting.
func (v MutuallyExclusiveAttributesValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Only one of %s can be configured at a time", strings.Join(v.AttributeNames, ", "))
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v MutuallyExclusiveAttributesValidator) MarkdownDescription(ctx context.Context) string {
	attrNames := make([]string, len(v.AttributeNames))
	for i, name := range v.AttributeNames {
		attrNames[i] = fmt.Sprintf("`%s`", name)
	}
	return fmt.Sprintf("Only one of %s can be configured at a time", strings.Join(attrNames, ", "))
}

// ValidateResource performs the validation.
func (v MutuallyExclusiveAttributesValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if len(v.AttributePaths) == 0 {
		return
	}

	// Skip validation if config is null (e.g., during destroy)
	if req.Config.Raw.IsNull() {
		return
	}

	var configuredAttrs []string
	var configuredPaths []path.Path

	for i, attrPath := range v.AttributePaths {
		// Check if attribute is configured by examining the config
		if isConfigValueSet(ctx, req.Config, attrPath) {
			configuredAttrs = append(configuredAttrs, v.AttributeNames[i])
			configuredPaths = append(configuredPaths, attrPath)
		}
	}

	// If more than one attribute is configured, add error diagnostics
	if len(configuredAttrs) > 1 {
		errorMsg := fmt.Sprintf("Only one of %s can be configured at a time. Currently configured: %s",
			strings.Join(v.AttributeNames, ", "),
			strings.Join(configuredAttrs, ", "))

		for _, attrPath := range configuredPaths {
			resp.Diagnostics.AddAttributeError(
				attrPath,
				"Mutually Exclusive Attributes",
				errorMsg,
			)
		}
	}
}

// isConfigValueSet checks if an attribute has a configured value (not null, not unknown, and has content)
func isConfigValueSet(ctx context.Context, config tfsdk.Config, attrPath path.Path) bool {
	// Try as String
	var stringVal basetypes.StringValue
	if diags := config.GetAttribute(ctx, attrPath, &stringVal); !diags.HasError() {
		return !stringVal.IsNull() && !stringVal.IsUnknown() && stringVal.ValueString() != ""
	}

	// Try as Set with attr.Value elements
	var setVal basetypes.SetValue
	if diags := config.GetAttribute(ctx, attrPath, &setVal); !diags.HasError() {
		return !setVal.IsNull() && !setVal.IsUnknown() && len(setVal.Elements()) > 0
	}

	// Try as List
	var listVal basetypes.ListValue
	if diags := config.GetAttribute(ctx, attrPath, &listVal); !diags.HasError() {
		return !listVal.IsNull() && !listVal.IsUnknown() && len(listVal.Elements()) > 0
	}

	// Try as Map
	var mapVal basetypes.MapValue
	if diags := config.GetAttribute(ctx, attrPath, &mapVal); !diags.HasError() {
		return !mapVal.IsNull() && !mapVal.IsUnknown() && len(mapVal.Elements()) > 0
	}

	// Try as Object
	var objVal basetypes.ObjectValue
	if diags := config.GetAttribute(ctx, attrPath, &objVal); !diags.HasError() {
		return !objVal.IsNull() && !objVal.IsUnknown() && len(objVal.Attributes()) > 0
	}

	// Try as Bool
	var boolVal basetypes.BoolValue
	if diags := config.GetAttribute(ctx, attrPath, &boolVal); !diags.HasError() {
		return !boolVal.IsNull() && !boolVal.IsUnknown()
	}

	// Try as Int64
	var int64Val basetypes.Int64Value
	if diags := config.GetAttribute(ctx, attrPath, &int64Val); !diags.HasError() {
		return !int64Val.IsNull() && !int64Val.IsUnknown()
	}

	// Try as Float64
	var float64Val basetypes.Float64Value
	if diags := config.GetAttribute(ctx, attrPath, &float64Val); !diags.HasError() {
		return !float64Val.IsNull() && !float64Val.IsUnknown()
	}

	// Try as Number
	var numberVal basetypes.NumberValue
	if diags := config.GetAttribute(ctx, attrPath, &numberVal); !diags.HasError() {
		return !numberVal.IsNull() && !numberVal.IsUnknown()
	}

	return false
}
