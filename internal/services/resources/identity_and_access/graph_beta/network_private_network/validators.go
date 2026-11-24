package graphBetaNetworkPrivateNetwork

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// expectedIPResolutionObjectValidator checks that an expected_ip_resolutions element sets exactly the fields its
// type uses: value for ip_address and ip_subnet, begin_address and end_address for ip_range.
type expectedIPResolutionObjectValidator struct{}

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

	var required, forbidden []string
	switch resolutionType.ValueString() {
	case expectedIPResolutionTypeIPAddress, expectedIPResolutionTypeIPSubnet:
		required, forbidden = []string{"value"}, []string{"begin_address", "end_address"}
	case expectedIPResolutionTypeIPRange:
		required, forbidden = []string{"begin_address", "end_address"}, []string{"value"}
	}

	for _, fieldName := range required {
		value, ok := attrs[fieldName].(types.String)
		if !ok || value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
			detail := fmt.Sprintf("expected_ip_resolutions type %q requires %q to be set.", resolutionType.ValueString(), fieldName)
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid expected IP resolution", detail)
		}
	}
	for _, fieldName := range forbidden {
		value, ok := attrs[fieldName].(types.String)
		if ok && !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
			detail := fmt.Sprintf("expected_ip_resolutions type %q must not set %q.", resolutionType.ValueString(), fieldName)
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid expected IP resolution", detail)
		}
	}
}
