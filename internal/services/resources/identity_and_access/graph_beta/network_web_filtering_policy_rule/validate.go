package graphBetaNetworkWebFilteringPolicyRule

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// uniqueCustomHeaderNamesValidator catches at plan time what Graph rejects with
// "Duplicate header name". Graph compares names case-insensitively.
type uniqueCustomHeaderNamesValidator struct{}

func (v uniqueCustomHeaderNamesValidator) Description(ctx context.Context) string {
	return "custom header names must be unique (case-insensitive)"
}

func (v uniqueCustomHeaderNamesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v uniqueCustomHeaderNamesValidator) ValidateList(
	ctx context.Context,
	req validator.ListRequest,
	resp *validator.ListResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var headers []CustomHeaderModel
	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &headers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	seen := make(map[string]int, len(headers))
	for i, header := range headers {
		if header.Name.IsNull() || header.Name.IsUnknown() {
			continue
		}
		key := strings.ToLower(header.Name.ValueString())
		if first, ok := seen[key]; ok {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtListIndex(i).AtName("name"),
				"Duplicate custom header name",
				fmt.Sprintf(
					"Header name %q duplicates custom_headers[%d]. Each header name must be unique within a rule (case-insensitive).",
					header.Name.ValueString(),
					first,
				),
			)
			continue
		}
		seen[key] = i
	}
}

var _ validator.List = uniqueCustomHeaderNamesValidator{}

// customHeadersRequireAllowValidator rejects custom_headers on block rules,
// which Graph refuses with "Invalid request parameters".
type customHeadersRequireAllowValidator struct{}

func (v customHeadersRequireAllowValidator) Description(ctx context.Context) string {
	return "custom_headers can only be set when action is allow"
}

func (v customHeadersRequireAllowValidator) MarkdownDescription(ctx context.Context) string {
	return "`custom_headers` can only be set when `action` is `allow`"
}

func (v customHeadersRequireAllowValidator) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var action types.String
	var headers types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("action"), &action)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("custom_headers"), &headers)...)
	if resp.Diagnostics.HasError() || action.IsNull() || action.IsUnknown() || headers.IsNull() ||
		headers.IsUnknown() {
		return
	}

	if action.ValueString() != "allow" {
		resp.Diagnostics.AddAttributeError(
			path.Root("custom_headers"),
			"Custom headers require an allow action",
			fmt.Sprintf(
				"custom_headers can only be set when action is \"allow\", got %q.",
				action.ValueString(),
			),
		)
	}
}

var _ resource.ConfigValidator = customHeadersRequireAllowValidator{}
