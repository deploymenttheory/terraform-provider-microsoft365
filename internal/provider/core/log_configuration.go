package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// redactedProviderConfiguration retains useful configuration details for logging while
// masking attributes marked Sensitive in the provider schema. It never changes the model.
func redactedProviderConfiguration(ctx context.Context, config M365ProviderModel) map[string]any {
	var response provider.SchemaResponse
	(&M365Provider{}).Schema(ctx, provider.SchemaRequest{}, &response)
	value, diags := types.ObjectValueFrom(
		ctx,
		schemaToAttrTypes(response.Schema.Attributes),
		config,
	)
	if response.Diagnostics.HasError() || diags.HasError() {
		return map[string]any{"configuration": "[unavailable]"}
	}
	return redactConfigurationAttributes(ctx, value.Attributes(), response.Schema.Attributes)
}

// redactConfigurationAttributes masks sensitive values before formatting them, including
// nested authentication and client options, so secrets never reach the log formatter.
func redactConfigurationAttributes(
	ctx context.Context,
	values map[string]attr.Value,
	attributes map[string]schema.Attribute,
) map[string]any {
	result := make(map[string]any, len(values))
	for name, value := range values {
		attribute, ok := attributes[name]
		if !ok || attribute.IsSensitive() {
			result[name] = "[REDACTED]"
			continue
		}
		if value.IsNull() {
			result[name] = nil
			continue
		}
		if value.IsUnknown() {
			result[name] = "[unknown]"
			continue
		}
		switch value := value.(type) {
		case types.Object:
			if nested, ok := attribute.(schema.SingleNestedAttribute); ok {
				result[name] = redactConfigurationAttributes(
					ctx,
					value.Attributes(),
					nested.Attributes,
				)
			} else {
				result[name] = "[REDACTED]"
			}
		case types.String:
			result[name] = value.ValueString()
		case types.Bool:
			result[name] = value.ValueBool()
		case types.Int64:
			result[name] = value.ValueInt64()
		case types.List:
			// The provider's only list is additionally_allowed_tenants, a list of strings.
			// Mask other element types rather than risk logging nested sensitive fields.
			if value.ElementType(ctx).Equal(types.StringType) {
				result[name] = value.String()
			} else {
				result[name] = "[REDACTED]"
			}
		default:
			result[name] = "[REDACTED]"
		}
	}
	return result
}
