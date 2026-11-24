package graphBetaNetworkCloudFirewallPolicyRule

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	commonattr "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/attr"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errInvalidRequest = errors.New("invalid cloud firewall request")

// constructResource builds the cloud firewall rule body. Create passes a nil state and sends every field; Update passes the
// prior state and sends only the fields that changed.
func constructResource(ctx context.Context, plan, state *NetworkCloudFirewallPolicyRuleResourceModel) (graph.CloudFirewallRuleable, error) {
	// sortedValues returns a set's strings sorted, and an empty slice rather than nil so the API receives [] for an empty set.
	sortedValues := func(value types.Set) []string {
		values := append([]string{}, commonattr.StringSetElements(value)...)
		sort.Strings(values)
		return values
	}

	// send reports whether a field goes in the body: every field on Create, only the fields that changed on Update.
	prior := state
	if prior == nil {
		prior = &NetworkCloudFirewallPolicyRuleResourceModel{}
	}
	send := func(planned, previous attr.Value) bool { return state == nil || !planned.Equal(previous) }

	body := graph.NewCloudFirewallRule()
	if send(plan.Name, prior.Name) {
		convert.FrameworkToGraphString(plan.Name, body.SetName)
	}
	if send(plan.Description, prior.Description) {
		convert.FrameworkToGraphString(plan.Description, body.SetDescription)
		if plan.Description.IsNull() {
			// The generated setter omits nil strings, so an explicit JSON null goes through additional data.
			body.SetAdditionalData(map[string]any{"description": nil})
		}
	}
	if send(plan.Priority, prior.Priority) {
		priority := int64(plan.Priority.ValueInt32())
		body.SetPriority(&priority)
	}
	if send(plan.Action, prior.Action) {
		if err := convert.FrameworkToGraphEnum(plan.Action, graph.ParseCloudFirewallAction, body.SetAction); err != nil {
			return nil, fmt.Errorf("%w: action: %w", errInvalidRequest, err)
		}
		if body.GetAction() == nil {
			return nil, fmt.Errorf("%w: action is required", errInvalidRequest)
		}
	}
	if send(plan.Enabled, prior.Enabled) {
		settings := graph.NewCloudFirewallRuleSettings()
		status := graph.DISABLED_SECURITYRULESTATUS
		if plan.Enabled.ValueBool() {
			status = graph.ENABLED_SECURITYRULESTATUS
		}
		settings.SetStatus(&status)
		body.SetSettings(settings)
	}
	conditions := graph.NewCloudFirewallMatchingConditions()
	// Accumulate explicit nulls before SetAdditionalData, which replaces the map.
	explicitNulls := map[string]any{}
	conditionsChanged := state == nil
	for _, side := range []struct {
		name    string
		changed bool
		object  types.Object
	}{
		{"sources", send(plan.Sources, prior.Sources), plan.Sources},
		{"destinations", send(plan.Destinations, prior.Destinations), plan.Destinations},
	} {
		if !side.changed {
			continue
		}
		conditionsChanged = true
		if side.object.IsNull() {
			explicitNulls[side.name] = nil
			continue
		}
		// Address groups are sorted by type so the request body is stable.
		var addresses []ruleAddressModel
		// A null or unknown set has no elements, so it decodes to no address groups.
		if value := side.object.Attributes()["addresses"].(types.Set); len(value.Elements()) > 0 {
			if diags := value.ElementsAs(ctx, &addresses, false); diags.HasError() {
				return nil, fmt.Errorf("%w: invalid addresses: %s", errInvalidRequest, diags.Errors()[0].Detail())
			}
			sort.Slice(addresses, func(i, j int) bool { return addresses[i].Type.ValueString() < addresses[j].Type.ValueString() })
		}
		ports := sortedValues(side.object.Attributes()["ports"].(types.Set))
		if side.name == "sources" {
			sources := graph.NewCloudFirewallSourceMatching()
			values := make([]graph.CloudFirewallSourceAddressable, 0, len(addresses))
			for _, address := range addresses {
				if address.Type.ValueString() != "ip" {
					return nil, fmt.Errorf("%w: unsupported source address type", errInvalidRequest)
				}
				source := graph.NewCloudFirewallSourceIpAddress()
				source.SetValues(sortedValues(address.Values))
				values = append(values, source)
			}
			sources.SetAddresses(values)
			sources.SetPorts(ports)
			// Branch conditions are not authored by this resource. Read rejects nonempty remote branchIds.
			sources.SetAdditionalData(map[string]any{"branchIds": []string{}})
			conditions.SetSources(sources)
			continue
		}
		destinations := graph.NewCloudFirewallDestinationMatching()
		values := make([]graph.CloudFirewallDestinationAddressable, 0, len(addresses))
		for _, address := range addresses {
			switch address.Type.ValueString() {
			case "ip":
				destination := graph.NewCloudFirewallDestinationIpAddress()
				destination.SetValues(sortedValues(address.Values))
				values = append(values, destination)
			case "fqdn":
				destination := graph.NewCloudFirewallDestinationFqdnAddress()
				destination.SetValues(sortedValues(address.Values))
				values = append(values, destination)
			default:
				return nil, fmt.Errorf("%w: unsupported destination address type", errInvalidRequest)
			}
		}
		destinations.SetAddresses(values)
		destinations.SetPorts(ports)
		protocols := side.object.Attributes()["protocols"].(types.Set)
		if err := convert.FrameworkToGraphBitmaskEnumFromSet(ctx, protocols, graph.ParseCloudFirewallProtocol, destinations.SetProtocols); err != nil {
			return nil, fmt.Errorf("%w: destination protocols: %w", errInvalidRequest, err)
		}
		// The shared converter skips empty/null sets; this API requires a protocol.
		if destinations.GetProtocols() == nil {
			return nil, fmt.Errorf("%w: destinations.protocols must contain tcp, udp, or both", errInvalidRequest)
		}
		conditions.SetDestinations(destinations)
	}
	if conditionsChanged {
		conditions.SetAdditionalData(explicitNulls)
		body.SetMatchingConditions(conditions)
	}
	return body, nil
}
