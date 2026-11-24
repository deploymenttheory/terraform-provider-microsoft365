package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// destinationObjectType is the Terraform object type of one destinations element.
var destinationObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type":          types.StringType,
	"value":         types.StringType,
	"begin_address": types.StringType,
	"end_address":   types.StringType,
}}

// MapRemoteStateToTerraform maps the remote Microsoft Graph response to Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *NetworkInternetAccessForwardingPolicyRuleResourceModel, remote networkaccess.InternetAccessForwardingRuleable) {
	if remote == nil {
		tflog.Debug(ctx, "Remote internet access forwarding policy rule is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.Name = convert.GraphToFrameworkString(remote.GetName())
	data.Action = convert.GraphToFrameworkEnum(remote.GetAction())
	data.RuleType = convert.GraphToFrameworkEnum(remote.GetRuleType())
	switch data.RuleType.ValueString() {
	case "ipAddress":
		data.RuleType = types.StringValue(ruleTypeIPAddress)
	case "ipRange":
		data.RuleType = types.StringValue(ruleTypeIPRange)
	case "ipSubnet":
		data.RuleType = types.StringValue(ruleTypeIPSubnet)
	}
	data.ClientFallbackAction = convert.GraphToFrameworkEnum(remote.GetClientFallbackAction())
	data.Protocol = convert.GraphToFrameworkEnum(remote.GetProtocol())
	data.Ports = convert.GraphToFrameworkStringSet(ctx, remote.GetPorts())

	// A destination without an OData type keeps a null type; any type other than the IP ones maps to fqdn.
	destinations := make([]attr.Value, 0, len(remote.GetDestinations()))
	for _, destination := range remote.GetDestinations() {
		destinationType := types.StringNull()
		if odataType := destination.GetOdataType(); odataType != nil {
			switch *odataType {
			case ipAddressODataType:
				destinationType = types.StringValue(ruleTypeIPAddress)
			case ipRangeODataType:
				destinationType = types.StringValue(ruleTypeIPRange)
			case ipSubnetODataType:
				destinationType = types.StringValue(ruleTypeIPSubnet)
			default:
				destinationType = types.StringValue(ruleTypeFQDN)
			}
		}
		attributes := map[string]attr.Value{"type": destinationType, "value": types.StringNull(), "begin_address": types.StringNull(), "end_address": types.StringNull()}
		if withValue, ok := destination.(interface{ GetValue() *string }); ok {
			attributes["value"] = convert.GraphToFrameworkString(withValue.GetValue())
		}
		if ipRange, ok := destination.(networkaccess.IpRangeable); ok {
			attributes["begin_address"] = convert.GraphToFrameworkString(ipRange.GetBeginAddress())
			attributes["end_address"] = convert.GraphToFrameworkString(ipRange.GetEndAddress())
		}
		destinations = append(destinations, types.ObjectValueMust(destinationObjectType.AttrTypes, attributes))
	}
	data.Destinations = types.ListValueMust(destinationObjectType, destinations)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with id %s", ResourceName, data.ID.ValueString()))
}
