package graphBetaNetworkCloudFirewallPolicyRule

import (
	"context"
	"fmt"
	"net/netip"
	"regexp"
	"strconv"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var portRangePattern = regexp.MustCompile(constants.PortRangeRegex)

type portValidator struct{}

func (portValidator) Description(context.Context) string {
	return "A decimal port in 0–65535 or ascending inclusive range, without leading zeros."
}
func (v portValidator) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }
func (v portValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	// Reuse the shared bounds/format regex; represent a single port as an equal range.
	candidate := value
	if !strings.Contains(value, "-") {
		candidate = value + "-" + value
	}
	valid := portRangePattern.MatchString(candidate)
	if valid {
		parts := strings.Split(candidate, "-")
		first, _ := strconv.Atoi(parts[0]) // The shared regex guarantees decimal values in range.
		last, _ := strconv.Atoi(parts[1])
		valid = first <= last // Ordering is a cloud firewall API constraint beyond the regex.
	}
	if !valid {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid port", v.Description(ctx))
	}
}

func validIPv4(value string) bool {
	if strings.Contains(value, "/") {
		p, err := netip.ParsePrefix(value)
		return err == nil && p.Addr().Is4()
	}
	if strings.Contains(value, "-") {
		parts := strings.Split(value, "-")
		if len(parts) != 2 {
			return false
		}
		a, e1 := netip.ParseAddr(parts[0])
		b, e2 := netip.ParseAddr(parts[1])
		return e1 == nil && e2 == nil && a.Is4() && b.Is4() && a.Compare(b) <= 0
	}
	a, err := netip.ParseAddr(value)
	return err == nil && a.Is4()
}

func (r *NetworkCloudFirewallPolicyRuleResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data NetworkCloudFirewallPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	for key, object := range map[string]types.Object{"sources": data.Sources, "destinations": data.Destinations} {
		if object.IsNull() || object.IsUnknown() {
			continue
		}
		validateCollectionElements(object, key, resp)
		addressSet := object.Attributes()["addresses"].(types.Set)
		if addressSet.IsUnknown() || addressSet.IsNull() {
			continue
		}
		unknown := false
		for _, element := range addressSet.Elements() {
			if element.IsUnknown() || element.IsNull() {
				unknown = true
			}
		}
		if unknown {
			continue
		}
		addresses, err := addressModels(ctx, object)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root(key), "Invalid address groups", err.Error())
			continue
		}
		seen := map[string]bool{}
		for _, address := range addresses {
			if address.Type.IsUnknown() {
				continue
			}
			typ := address.Type.ValueString()
			if seen[typ] {
				resp.Diagnostics.AddAttributeError(path.Root(key).AtName("addresses"), "Duplicate address type", fmt.Sprintf("Address type %q can appear only once.", typ))
			}
			seen[typ] = true
			if address.Values.IsNull() || address.Values.IsUnknown() {
				continue
			}
			for _, element := range address.Values.Elements() {
				v := element.(types.String)
				if v.IsNull() {
					resp.Diagnostics.AddAttributeError(path.Root(key).AtName("addresses"), "Null condition element", "Address values must not contain null.")
					continue
				}
				if v.IsUnknown() {
					continue
				}
				value := v.ValueString()
				if typ == "ip" && !validIPv4(value) {
					resp.Diagnostics.AddAttributeError(path.Root(key).AtName("addresses"), "Invalid IPv4 condition", fmt.Sprintf("%q must be an IPv4 address, CIDR, or ascending IPv4 range.", value))
				}
				if typ == "fqdn" && value != "*" && (strings.ContainsAny(value, " \t\r\n/*") || !strings.Contains(value, ".") || strings.HasSuffix(value, ".")) {
					resp.Diagnostics.AddAttributeError(path.Root(key).AtName("addresses"), "Invalid FQDN condition", "Use a fully qualified domain name without a scheme, trailing dot, whitespace, or wildcard prefix, or a standalone *.")
				}
			}
		}
	}
}

func validateCollectionElements(object types.Object, key string, resp *resource.ValidateConfigResponse) {
	for _, field := range []string{"ports", "protocols"} {
		value, exists := object.Attributes()[field]
		if !exists || value.IsNull() || value.IsUnknown() {
			continue
		}
		for _, element := range value.(types.Set).Elements() {
			if element.IsNull() {
				resp.Diagnostics.AddAttributeError(path.Root(key).AtName(field), "Null condition element", "Collection elements must not be null.")
			}
		}
	}
}
