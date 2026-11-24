package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"fmt"
	"net/netip"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// A port entry is a single port or an inclusive range such as "1000-2000".
// Graph accepts "1000 - 2000" with spaces, but the portal rejects it. Leading
// zeros are rejected so a value Graph might normalise cannot drift.
var portRegex = regexp.MustCompile(`^(0|[1-9]\d{0,4})(?:-(0|[1-9]\d{0,4}))?$`)

// portValidator checks one ports element, matching the Entra portal's rule.
type portValidator struct{}

func (v portValidator) Description(ctx context.Context) string {
	return `must be a port such as "443" or a range such as "1000-2000", between 0 and 65535, with no spaces or leading zeros`
}

func (v portValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v portValidator) ValidateString(
	ctx context.Context,
	req validator.StringRequest,
	resp *validator.StringResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	valid := false
	if match := portRegex.FindStringSubmatch(value); match != nil {
		begin, beginErr := strconv.Atoi(match[1])
		end, endErr := begin, error(nil)
		if match[2] != "" {
			end, endErr = strconv.Atoi(match[2])
		}
		valid = beginErr == nil && endErr == nil && end <= 65535 && begin <= end
	}
	if !valid {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid port",
			fmt.Sprintf("Port %q %s.", value, v.Description(ctx)),
		)
	}
}

var _ validator.String = portValidator{}

// destinationsValidator reports at plan time what Graph rejects with long or
// opaque errors: destinations whose type differs from rule_type, missing or
// misplaced fields, and anything other than IPv4 (Graph rejects IPv6).
type destinationsValidator struct{}

func (v destinationsValidator) Description(ctx context.Context) string {
	return "destinations must match rule_type and use IPv4 addresses"
}

func (v destinationsValidator) MarkdownDescription(ctx context.Context) string {
	return "`destinations` must match `rule_type` and use IPv4 addresses"
}

func (v destinationsValidator) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var ruleType types.String
	var destinations types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("rule_type"), &ruleType)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("destinations"), &destinations)...)
	if resp.Diagnostics.HasError() || destinations.IsNull() || destinations.IsUnknown() {
		return
	}

	known := func(value types.String) bool {
		return !value.IsNull() && !value.IsUnknown()
	}
	stringAttr := func(attrs map[string]attr.Value, name string) types.String {
		if value, ok := attrs[name].(types.String); ok {
			return value
		}
		return types.StringNull()
	}
	requireSet := func(value types.String, attrPath path.Path, destinationType string) {
		if value.IsUnknown() || (!value.IsNull() && value.ValueString() != "") {
			return
		}
		resp.Diagnostics.AddAttributeError(attrPath, "Missing destination field", fmt.Sprintf("This field is required when the destination type is %q.", destinationType))
	}
	requireUnset := func(value types.String, attrPath path.Path, destinationType string) {
		if !known(value) {
			return
		}
		resp.Diagnostics.AddAttributeError(attrPath, "Unexpected destination field", fmt.Sprintf("This field cannot be set when the destination type is %q.", destinationType))
	}
	parseIPv4 := func(value types.String, attrPath path.Path) (netip.Addr, bool) {
		if !known(value) || value.ValueString() == "" {
			return netip.Addr{}, false
		}
		addr, err := netip.ParseAddr(value.ValueString())
		if err != nil || !addr.Is4() {
			resp.Diagnostics.AddAttributeError(attrPath, "Invalid IP address", fmt.Sprintf("%q must be an IPv4 address; Graph rejects IPv6.", value.ValueString()))
			return netip.Addr{}, false
		}
		return addr, true
	}
	checkDestinationFields := func(destinationType string, attrs map[string]attr.Value, itemPath path.Path) {
		value := stringAttr(attrs, "value")
		begin := stringAttr(attrs, "begin_address")
		end := stringAttr(attrs, "end_address")

		if destinationType == ruleTypeIPRange {
			requireUnset(value, itemPath.AtName("value"), destinationType)
			requireSet(begin, itemPath.AtName("begin_address"), destinationType)
			requireSet(end, itemPath.AtName("end_address"), destinationType)
			beginAddr, beginOK := parseIPv4(begin, itemPath.AtName("begin_address"))
			endAddr, endOK := parseIPv4(end, itemPath.AtName("end_address"))
			if beginOK && endOK && endAddr.Less(beginAddr) {
				resp.Diagnostics.AddAttributeError(itemPath.AtName("end_address"), "Invalid IP range", fmt.Sprintf("end_address %s must not be lower than begin_address %s.", endAddr, beginAddr))
			}
			return
		}

		requireSet(value, itemPath.AtName("value"), destinationType)
		requireUnset(begin, itemPath.AtName("begin_address"), destinationType)
		requireUnset(end, itemPath.AtName("end_address"), destinationType)
		if !known(value) || value.ValueString() == "" {
			return
		}
		switch destinationType {
		case ruleTypeIPAddress:
			parseIPv4(value, itemPath.AtName("value"))
		case ruleTypeIPSubnet:
			prefix, err := netip.ParsePrefix(value.ValueString())
			if err != nil || !prefix.Addr().Is4() {
				resp.Diagnostics.AddAttributeError(itemPath.AtName("value"), "Invalid IP subnet", fmt.Sprintf("%q must be an IPv4 CIDR such as 198.51.100.0/24; Graph rejects IPv6.", value.ValueString()))
			}
		}
	}

	for i, element := range destinations.Elements() {
		destination, ok := element.(types.Object)
		if !ok || destination.IsNull() || destination.IsUnknown() {
			continue
		}
		attrs := destination.Attributes()
		itemPath := path.Root("destinations").AtListIndex(i)
		destinationType, _ := attrs["type"].(types.String)
		if destinationType.IsUnknown() || destinationType.IsNull() {
			continue
		}
		if !ruleType.IsUnknown() && !ruleType.IsNull() &&
			destinationType.ValueString() != ruleType.ValueString() {
			resp.Diagnostics.AddAttributeError(
				itemPath.AtName("type"),
				"Destination type does not match rule_type",
				fmt.Sprintf(
					"Every destination must have the rule's type. Got %q, but rule_type is %q.",
					destinationType.ValueString(),
					ruleType.ValueString(),
				),
			)
		}
		checkDestinationFields(destinationType.ValueString(), attrs, itemPath)
	}
}

var _ resource.ConfigValidator = destinationsValidator{}
