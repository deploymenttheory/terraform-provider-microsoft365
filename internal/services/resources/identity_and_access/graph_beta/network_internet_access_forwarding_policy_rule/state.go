package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkInternetAccessForwardingPolicyRuleResourceModel, remote *internetAccessForwardingRuleResponse) {
	if remote == nil {
		tflog.Debug(ctx, "Remote internet access forwarding policy rule is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remote.id)
	data.Name = convert.GraphToFrameworkString(remote.name)
	data.Action = convert.GraphToFrameworkString(remote.action)
	if remote.ruleType != nil {
		data.RuleType = types.StringValue(terraformRuleType(*remote.ruleType))
	} else {
		data.RuleType = types.StringNull()
	}
	data.ClientFallbackAction = convert.GraphToFrameworkString(remote.clientFallbackAction)
	data.Protocol = convert.GraphToFrameworkString(remote.protocol)
	data.Ports = convert.GraphToFrameworkStringSet(ctx, remote.ports)

	destinations := make([]attr.Value, 0, len(remote.destinations))
	for _, destination := range remote.destinations {
		destinations = append(destinations, destinationObjectValue(destination))
	}
	data.Destinations = types.ListValueMust(destinationObjectType(), destinations)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with id %s", ResourceName, data.ID.ValueString()))
}

func destinationObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":          types.StringType,
			"value":         types.StringType,
			"begin_address": types.StringType,
			"end_address":   types.StringType,
		},
	}
}

func destinationObjectValue(destination ruleDestinationResponse) attr.Value {
	destinationType := types.StringNull()
	if destination.odataType != nil {
		destinationType = types.StringValue(terraformDestinationType(*destination.odataType))
	}
	return types.ObjectValueMust(destinationObjectType().AttrTypes, map[string]attr.Value{
		"type":          destinationType,
		"value":         convert.GraphToFrameworkString(destination.value),
		"begin_address": convert.GraphToFrameworkString(destination.beginAddress),
		"end_address":   convert.GraphToFrameworkString(destination.endAddress),
	})
}

func terraformDestinationType(odataType string) string {
	switch odataType {
	case ipAddressODataType:
		return ruleTypeIPAddress
	case ipRangeODataType:
		return ruleTypeIPRange
	case ipSubnetODataType:
		return ruleTypeIPSubnet
	default:
		return ruleTypeFQDN
	}
}
