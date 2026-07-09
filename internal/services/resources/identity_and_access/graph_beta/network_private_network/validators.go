package graphBetaNetworkPrivateNetwork

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type expectedIPResolutionObjectValidator struct{}

func expectedIPResolutionObject() validator.Object {
	return expectedIPResolutionObjectValidator{}
}

func (v expectedIPResolutionObjectValidator) Description(_ context.Context) string {
	return "validates expected_ip_resolutions field combinations for the selected type"
}

func (v expectedIPResolutionObjectValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v expectedIPResolutionObjectValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()
	resolutionType, ok := attrs["type"].(types.String)
	if !ok || resolutionType.IsNull() || resolutionType.IsUnknown() {
		return
	}

	switch resolutionType.ValueString() {
	case expectedIPResolutionTypeIPAddress, expectedIPResolutionTypeIPSubnet:
		requireConfiguredString(resp, req.Path, attrs, "value", resolutionType.ValueString())
		forbidConfiguredString(resp, req.Path, attrs, "begin_address", resolutionType.ValueString())
		forbidConfiguredString(resp, req.Path, attrs, "end_address", resolutionType.ValueString())
	case expectedIPResolutionTypeIPRange:
		requireConfiguredString(resp, req.Path, attrs, "begin_address", resolutionType.ValueString())
		requireConfiguredString(resp, req.Path, attrs, "end_address", resolutionType.ValueString())
		forbidConfiguredString(resp, req.Path, attrs, "value", resolutionType.ValueString())
	}
}

func requireConfiguredString(resp *validator.ObjectResponse, attrPath path.Path, attrs map[string]attr.Value, fieldName string, resolutionType string) {
	value, ok := attrs[fieldName].(types.String)
	if ok && !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
		return
	}

	resp.Diagnostics.AddAttributeError(
		attrPath,
		"Invalid expected IP resolution",
		fmt.Sprintf("expected_ip_resolutions type %q requires %q to be set.", resolutionType, fieldName),
	)
}

func forbidConfiguredString(resp *validator.ObjectResponse, attrPath path.Path, attrs map[string]attr.Value, fieldName string, resolutionType string) {
	value, ok := attrs[fieldName].(types.String)
	if !ok || value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return
	}

	resp.Diagnostics.AddAttributeError(
		attrPath,
		"Invalid expected IP resolution",
		fmt.Sprintf("expected_ip_resolutions type %q must not set %q.", resolutionType, fieldName),
	)
}
