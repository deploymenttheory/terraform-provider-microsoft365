package graphBetaNetworkCloudFirewallPolicyRule

import (
	"context"
	"errors"
	"fmt"
	"sort"

	commonattr "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/attr"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

var errInvalidRequest = errors.New("invalid cloud firewall request")

type ruleRequest struct {
	*graph.CloudFirewallRule
	changed bool
}

func constructRule(ctx context.Context, plan, state *NetworkCloudFirewallPolicyRuleResourceModel) (*ruleRequest, error) {
	b := &ruleRequest{CloudFirewallRule: graph.NewCloudFirewallRule()}
	if state == nil || !plan.Name.Equal(state.Name) {
		convert.FrameworkToGraphString(plan.Name, b.SetName)
		b.changed = true
	}
	if state == nil || !plan.Description.Equal(state.Description) {
		convert.FrameworkToGraphString(plan.Description, b.SetDescription)
		if plan.Description.IsNull() {
			b.SetAdditionalData(map[string]any{"description": nil})
		}
		b.changed = true
	}
	if state == nil || !plan.Priority.Equal(state.Priority) {
		v := int64(plan.Priority.ValueInt32())
		b.SetPriority(&v)
		b.changed = true
	}
	if state == nil || !plan.Action.Equal(state.Action) {
		if err := convert.FrameworkToGraphEnum(plan.Action, graph.ParseCloudFirewallAction, b.SetAction); err != nil {
			return nil, fmt.Errorf("%w: action: %w", errInvalidRequest, err)
		}
		if b.GetAction() == nil {
			return nil, fmt.Errorf("%w: action is required", errInvalidRequest)
		}
		b.changed = true
	}
	if state == nil || !plan.Enabled.Equal(state.Enabled) {
		settings := graph.NewCloudFirewallRuleSettings()
		status := graph.DISABLED_SECURITYRULESTATUS
		if plan.Enabled.ValueBool() {
			status = graph.ENABLED_SECURITYRULESTATUS
		}
		settings.SetStatus(&status)
		b.SetSettings(settings)
		b.changed = true
	}
	mc := graph.NewCloudFirewallMatchingConditions()
	// Accumulate explicit nulls before SetAdditionalData, which replaces the map.
	explicitNulls := map[string]any{}
	conditionsChanged := state == nil
	if state == nil || !plan.Sources.Equal(state.Sources) {
		conditionsChanged = true
		if plan.Sources.IsNull() {
			explicitNulls["sources"] = nil
		} else {
			sources := graph.NewCloudFirewallSourceMatching()
			addresses, err := addressModels(ctx, plan.Sources)
			if err != nil {
				return nil, err
			}
			values := make([]graph.CloudFirewallSourceAddressable, 0, len(addresses))
			for _, address := range addresses {
				if address.Type.ValueString() != "ip" {
					return nil, fmt.Errorf("%w: unsupported source address type", errInvalidRequest)
				}
				v := graph.NewCloudFirewallSourceIpAddress()
				v.SetValues(sortedValues(address.Values))
				values = append(values, v)
			}
			sources.SetAddresses(values)
			sources.SetPorts(objectStrings(plan.Sources, "ports"))
			// Branch conditions are not authored by this resource. Read rejects nonempty remote branchIds.
			sources.SetAdditionalData(map[string]any{"branchIds": []string{}})
			mc.SetSources(sources)
		}
	}
	if state == nil || !plan.Destinations.Equal(state.Destinations) {
		conditionsChanged = true
		if plan.Destinations.IsNull() {
			explicitNulls["destinations"] = nil
		} else {
			destinations := graph.NewCloudFirewallDestinationMatching()
			addresses, err := addressModels(ctx, plan.Destinations)
			if err != nil {
				return nil, err
			}
			values := make([]graph.CloudFirewallDestinationAddressable, 0, len(addresses))
			for _, address := range addresses {
				switch address.Type.ValueString() {
				case "ip":
					v := graph.NewCloudFirewallDestinationIpAddress()
					v.SetValues(sortedValues(address.Values))
					values = append(values, v)
				case "fqdn":
					v := graph.NewCloudFirewallDestinationFqdnAddress()
					v.SetValues(sortedValues(address.Values))
					values = append(values, v)
				default:
					return nil, fmt.Errorf("%w: unsupported destination address type", errInvalidRequest)
				}
			}
			destinations.SetAddresses(values)
			destinations.SetPorts(objectStrings(plan.Destinations, "ports"))
			if err := convert.FrameworkToGraphBitmaskEnumFromSet(ctx, plan.Destinations.Attributes()["protocols"].(types.Set), graph.ParseCloudFirewallProtocol, destinations.SetProtocols); err != nil {
				return nil, fmt.Errorf("%w: destination protocols: %w", errInvalidRequest, err)
			}
			// The shared converter skips empty/null sets; this API requires a protocol.
			if destinations.GetProtocols() == nil {
				return nil, fmt.Errorf("%w: destinations.protocols must contain tcp, udp, or both", errInvalidRequest)
			}
			mc.SetDestinations(destinations)
		}
	}
	if conditionsChanged {
		mc.SetAdditionalData(explicitNulls)
		b.SetMatchingConditions(mc)
		b.changed = true
	}
	return b, nil
}

func sortedValues(value types.Set) []string {
	values := commonattr.StringSetElements(value)
	if values == nil {
		values = []string{}
	}
	sort.Strings(values)
	return values
}
func objectStrings(value types.Object, key string) []string {
	return sortedValues(value.Attributes()[key].(types.Set))
}
func addressModels(ctx context.Context, object types.Object) ([]ruleAddressModel, error) {
	value := object.Attributes()["addresses"].(types.Set)
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}
	var models []ruleAddressModel
	if diags := value.ElementsAs(ctx, &models, false); diags.HasError() {
		return nil, fmt.Errorf("%w: invalid addresses: %s", errInvalidRequest, diags.Errors()[0].Detail())
	}
	sort.Slice(models, func(i, j int) bool { return models[i].Type.ValueString() < models[j].Type.ValueString() })
	return models, nil
}
