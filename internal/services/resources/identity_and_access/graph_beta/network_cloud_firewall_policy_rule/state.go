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
	sources, err := sourceState(ctx, remote.GetMatchingConditions().GetSources())
	if err != nil {
		return err
	}
	destinations, err := destinationState(ctx, remote.GetMatchingConditions().GetDestinations())
	if err != nil {
		return err
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

func addressState(ctx context.Context, typ string, values []string) types.Object {
	return types.ObjectValueMust(addressObjectType().AttrTypes, map[string]attr.Value{"type": types.StringValue(typ), "values": stringSet(ctx, values)})
}
func stringSet(ctx context.Context, values []string) types.Set {
	values = append([]string{}, values...)
	slices.Sort(values)
	values = slices.Compact(values)
	return convert.GraphToFrameworkStringSetPreserveEmpty(ctx, values)
}

func sourceState(ctx context.Context, source graph.CloudFirewallSourceMatchingable) (types.Object, error) {
	empty := types.ObjectNull(matchingObjectType(false).AttrTypes)
	if source == nil {
		return empty, nil
	}
	if branches, ok := source.GetAdditionalData()["branchIds"]; ok && branches != nil {
		nonempty := true
		switch v := branches.(type) {
		case []any:
			nonempty = len(v) > 0
		case []string:
			nonempty = len(v) > 0
		}
		if nonempty {
			return empty, fmt.Errorf("%w: this rule contains unsupported branchIds; refusing to discard existing branch conditions", errInvalidResponse)
		}
	}
	addresses := make([]attr.Value, 0, len(source.GetAddresses()))
	for _, address := range source.GetAddresses() {
		v, ok := address.(graph.CloudFirewallSourceIpAddressable)
		if !ok || address.GetOdataType() == nil || *address.GetOdataType() != "#microsoft.graph.networkaccess.cloudFirewallSourceIpAddress" {
			return empty, fmt.Errorf("%w: unsupported cloud firewall source address type", errInvalidResponse)
		}
		addresses = append(addresses, addressState(ctx, "ip", v.GetValues()))
	}
	if len(addresses) > 1 {
		return empty, fmt.Errorf("%w: multiple source address groups with the same type are not supported", errInvalidResponse)
	}
	return types.ObjectValueMust(matchingObjectType(false).AttrTypes, map[string]attr.Value{"addresses": types.SetValueMust(addressObjectType(), addresses), "ports": stringSet(ctx, source.GetPorts())}), nil
}

func destinationState(ctx context.Context, destination graph.CloudFirewallDestinationMatchingable) (types.Object, error) {
	empty := types.ObjectNull(matchingObjectType(true).AttrTypes)
	if destination == nil {
		return empty, nil
	}
	if destination.GetProtocols() == nil {
		return empty, fmt.Errorf("%w: destination protocols are missing; the API can return invalid data after a null protocol update", errInvalidResponse)
	}
	protocol := destination.GetProtocols().String()
	if protocol != "tcp" && protocol != "udp" && protocol != "tcp,udp" {
		return empty, fmt.Errorf("%w: unsupported destination protocols %q", errInvalidResponse, protocol)
	}
	addresses := make([]attr.Value, 0, len(destination.GetAddresses()))
	seen := map[string]bool{}
	for _, address := range destination.GetAddresses() {
		var typ string
		var values []string
		// Inspect the discriminator because generated address interfaces share the same methods.
		if address == nil || address.GetOdataType() == nil {
			return empty, fmt.Errorf("%w: destination address discriminator is missing", errInvalidResponse)
		}
		switch *address.GetOdataType() {
		case "#microsoft.graph.networkaccess.cloudFirewallDestinationIpAddress":
			v, ok := address.(graph.CloudFirewallDestinationIpAddressable)
			if !ok {
				return empty, fmt.Errorf("%w: invalid destination IP address", errInvalidResponse)
			}
			typ = "ip"
			values = v.GetValues()
		case "#microsoft.graph.networkaccess.cloudFirewallDestinationFqdnAddress":
			v, ok := address.(graph.CloudFirewallDestinationFqdnAddressable)
			if !ok {
				return empty, fmt.Errorf("%w: invalid destination FQDN address", errInvalidResponse)
			}
			typ = "fqdn"
			values = v.GetValues()
		default:
			return empty, fmt.Errorf("%w: unsupported destination address type %q", errInvalidResponse, *address.GetOdataType())
		}
		if seen[typ] {
			return empty, fmt.Errorf("%w: duplicate destination address type %q", errInvalidResponse, typ)
		}
		seen[typ] = true
		addresses = append(addresses, addressState(ctx, typ, values))
	}
	return types.ObjectValueMust(matchingObjectType(true).AttrTypes, map[string]attr.Value{"addresses": types.SetValueMust(addressObjectType(), addresses), "ports": stringSet(ctx, destination.GetPorts()), "protocols": convert.GraphToFrameworkBitmaskEnumAsSet(ctx, destination.GetProtocols())}), nil
}
