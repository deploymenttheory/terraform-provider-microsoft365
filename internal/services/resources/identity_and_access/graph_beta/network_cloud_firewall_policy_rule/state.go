package graphBetaNetworkCloudFirewallPolicyRule

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

var errInvalidResponse = errors.New("invalid cloud firewall response")

// MapRemoteStateToTerraform refuses unsupported or incomplete responses before changing state.
func MapRemoteStateToTerraform(ctx context.Context, data *NetworkCloudFirewallPolicyRuleResourceModel, response graph.PolicyRuleable) error {
	// stringSet maps API strings to a set, sorted and without duplicates, keeping an empty set distinct from null.
	stringSet := func(values []string) types.Set {
		values = append([]string{}, values...)
		slices.Sort(values)
		values = slices.Compact(values)
		return convert.GraphToFrameworkStringSetPreserveEmpty(ctx, values)
	}
	// addressState maps one address group.
	addressState := func(typ string, values []string) types.Object {
		return types.ObjectValueMust(addressObjectType.AttrTypes, map[string]attr.Value{"type": types.StringValue(typ), "values": stringSet(values)})
	}
	remote, ok := response.(graph.CloudFirewallRuleable)
	if !ok || remote == nil || remote.GetOdataType() == nil || *remote.GetOdataType() != "#microsoft.graph.networkaccess.cloudFirewallRule" || remote.GetId() == nil || *remote.GetId() == "" || remote.GetName() == nil || remote.GetPriority() == nil || remote.GetAction() == nil || remote.GetSettings() == nil || remote.GetSettings().GetStatus() == nil || remote.GetMatchingConditions() == nil {
		return fmt.Errorf("%w: cloud firewall rule response is missing required properties or has an unsupported type", errInvalidResponse)
	}
	if !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != *remote.GetId() {
		return fmt.Errorf("%w: rule response ID does not match requested ID", errInvalidResponse)
	}
	priority := *remote.GetPriority()
	if priority < math.MinInt32 || priority > math.MaxInt32 {
		return fmt.Errorf("%w: rule priority %d is outside the supported 32-bit range", errInvalidResponse, priority)
	}
	status := remote.GetSettings().GetStatus().String()
	if status != "enabled" && status != "disabled" {
		return fmt.Errorf("%w: unknown cloud firewall rule status %q; preserving the previous enabled value", errInvalidResponse, status)
	}
	action := remote.GetAction().String()
	if action != "allow" && action != "block" {
		return fmt.Errorf("%w: unsupported cloud firewall action %q", errInvalidResponse, action)
	}
	sources := types.ObjectNull(sourceMatchingObjectType.AttrTypes)
	if source := remote.GetMatchingConditions().GetSources(); source != nil {
		if branches, ok := source.GetAdditionalData()["branchIds"]; ok && branches != nil {
			nonempty := true
			switch v := branches.(type) {
			case []any:
				nonempty = len(v) > 0
			case []string:
				nonempty = len(v) > 0
			}
			if nonempty {
				return fmt.Errorf("%w: this rule contains unsupported branchIds; refusing to discard existing branch conditions", errInvalidResponse)
			}
		}
		addresses := make([]attr.Value, 0, len(source.GetAddresses()))
		for _, address := range source.GetAddresses() {
			v, ok := address.(graph.CloudFirewallSourceIpAddressable)
			if !ok || address.GetOdataType() == nil || *address.GetOdataType() != "#microsoft.graph.networkaccess.cloudFirewallSourceIpAddress" {
				return fmt.Errorf("%w: unsupported cloud firewall source address type", errInvalidResponse)
			}
			addresses = append(addresses, addressState("ip", v.GetValues()))
		}
		if len(addresses) > 1 {
			return fmt.Errorf("%w: multiple source address groups with the same type are not supported", errInvalidResponse)
		}
		sourceValues := map[string]attr.Value{"addresses": types.SetValueMust(addressObjectType, addresses), "ports": stringSet(source.GetPorts())}
		sources = types.ObjectValueMust(sourceMatchingObjectType.AttrTypes, sourceValues)
	}
	destinations := types.ObjectNull(destinationMatchingObjectType.AttrTypes)
	if destination := remote.GetMatchingConditions().GetDestinations(); destination != nil {
		if destination.GetProtocols() == nil {
			return fmt.Errorf("%w: destination protocols are missing; the API can return invalid data after a null protocol update", errInvalidResponse)
		}
		protocol := destination.GetProtocols().String()
		if protocol != "tcp" && protocol != "udp" && protocol != "tcp,udp" {
			return fmt.Errorf("%w: unsupported destination protocols %q", errInvalidResponse, protocol)
		}
		addresses := make([]attr.Value, 0, len(destination.GetAddresses()))
		seen := map[string]bool{}
		for _, address := range destination.GetAddresses() {
			var typ string
			var values []string
			// Inspect the discriminator because generated address interfaces share the same methods.
			if address == nil || address.GetOdataType() == nil {
				return fmt.Errorf("%w: destination address discriminator is missing", errInvalidResponse)
			}
			switch *address.GetOdataType() {
			case "#microsoft.graph.networkaccess.cloudFirewallDestinationIpAddress":
				v, ok := address.(graph.CloudFirewallDestinationIpAddressable)
				if !ok {
					return fmt.Errorf("%w: invalid destination IP address", errInvalidResponse)
				}
				typ = "ip"
				values = v.GetValues()
			case "#microsoft.graph.networkaccess.cloudFirewallDestinationFqdnAddress":
				v, ok := address.(graph.CloudFirewallDestinationFqdnAddressable)
				if !ok {
					return fmt.Errorf("%w: invalid destination FQDN address", errInvalidResponse)
				}
				typ = "fqdn"
				values = v.GetValues()
			default:
				return fmt.Errorf("%w: unsupported destination address type %q", errInvalidResponse, *address.GetOdataType())
			}
			if seen[typ] {
				return fmt.Errorf("%w: duplicate destination address type %q", errInvalidResponse, typ)
			}
			seen[typ] = true
			addresses = append(addresses, addressState(typ, values))
		}
		destinationValues := map[string]attr.Value{
			"addresses": types.SetValueMust(addressObjectType, addresses),
			"ports":     stringSet(destination.GetPorts()),
			"protocols": convert.GraphToFrameworkBitmaskEnumAsSet(ctx, destination.GetProtocols()),
		}
		destinations = types.ObjectValueMust(destinationMatchingObjectType.AttrTypes, destinationValues)
	}
	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.Name = convert.GraphToFrameworkString(remote.GetName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.Priority = types.Int32Value(int32(priority))
	data.Action = convert.GraphToFrameworkEnum(remote.GetAction())
	data.Status = convert.GraphToFrameworkEnum(remote.GetSettings().GetStatus())
	data.Enabled = types.BoolValue(status == "enabled")
	data.Sources = sources
	data.Destinations = destinations
	return nil
}
