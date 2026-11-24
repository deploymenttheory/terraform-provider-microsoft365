package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"errors"
	"fmt"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

const (
	internetAccessForwardingRuleODataType = "#microsoft.graph.networkaccess.internetAccessForwardingRule"
	fqdnODataType                         = "#microsoft.graph.networkaccess.fqdn"
	ipAddressODataType                    = "#microsoft.graph.networkaccess.ipAddress"
	ipRangeODataType                      = "#microsoft.graph.networkaccess.ipRange"
	ipSubnetODataType                     = "#microsoft.graph.networkaccess.ipSubnet"
)

var (
	errReadAttribute              = errors.New("failed to read")
	errDestinationTypeMismatch    = errors.New("destination type must match rule_type")
	errDestinationValueRequired   = errors.New("destination value is required")
	errDestinationRangeRequired   = errors.New("begin_address and end_address are required")
	errUnsupportedDestinationType = errors.New("unsupported destination type")
	errNoDestinations             = errors.New("at least one destination must be specified")
)

// ruleTypeToGraph maps the Terraform rule and destination types to the Graph ruleType values.
var ruleTypeToGraph = map[string]string{
	ruleTypeFQDN:      "fqdn",
	ruleTypeIPAddress: "ipAddress",
	ruleTypeIPRange:   "ipRange",
	ruleTypeIPSubnet:  "ipSubnet",
}

// constructResource converts the Terraform resource model to a Microsoft Graph request body. The update body carries
// the rule id and leaves out name, action and clientFallbackAction, which Graph accepts only on create.
func constructResource(ctx context.Context, data *NetworkInternetAccessForwardingPolicyRuleResourceModel, includeID bool) (networkaccess.InternetAccessForwardingRuleable, error) {
	var ports []string
	if !data.Ports.IsNull() && !data.Ports.IsUnknown() {
		if diags := data.Ports.ElementsAs(ctx, &ports, false); diags.HasError() {
			return nil, fmt.Errorf("%w ports: %s", errReadAttribute, diags.Errors()[0].Detail())
		}
	}

	var terraformDestinations []RuleDestinationModel
	if !data.Destinations.IsNull() && !data.Destinations.IsUnknown() {
		if diags := data.Destinations.ElementsAs(ctx, &terraformDestinations, false); diags.HasError() {
			return nil, fmt.Errorf("%w destinations: %s", errReadAttribute, diags.Errors()[0].Detail())
		}
	}
	for _, destination := range terraformDestinations {
		if destination.Type.ValueString() != data.RuleType.ValueString() {
			return nil, fmt.Errorf("%w: destination type %q, rule_type %q", errDestinationTypeMismatch, destination.Type.ValueString(), data.RuleType.ValueString())
		}
	}

	destinations := make([]networkaccess.RuleDestinationable, 0, len(terraformDestinations))
	for _, item := range terraformDestinations {
		destinationType := item.Type.ValueString()
		if destinationType == ruleTypeIPRange {
			if item.BeginAddress.ValueString() == "" || item.EndAddress.ValueString() == "" {
				return nil, fmt.Errorf("%w when destination type is %s", errDestinationRangeRequired, ruleTypeIPRange)
			}
			destination := networkaccess.NewIpRange()
			destination.SetBeginAddress(item.BeginAddress.ValueStringPointer())
			destination.SetEndAddress(item.EndAddress.ValueStringPointer())
			destinations = append(destinations, destination)
			continue
		}

		// fqdn, ipAddress and ipSubnet destinations all carry a single value.
		var destination interface {
			networkaccess.RuleDestinationable
			SetValue(value *string)
		}
		switch destinationType {
		case ruleTypeFQDN:
			destination = networkaccess.NewFqdn()
		case ruleTypeIPAddress:
			destination = networkaccess.NewIpAddress()
		case ruleTypeIPSubnet:
			destination = networkaccess.NewIpSubnet()
		default:
			return nil, fmt.Errorf("%w %q", errUnsupportedDestinationType, destinationType)
		}
		if item.Value.ValueString() == "" {
			return nil, fmt.Errorf("%w when type is %s", errDestinationValueRequired, destinationType)
		}
		destination.SetValue(item.Value.ValueStringPointer())
		destinations = append(destinations, destination)
	}
	if len(destinations) == 0 {
		return nil, errNoDestinations
	}

	requestBody := networkaccess.NewInternetAccessForwardingRule()
	requestBody.SetDestinations(destinations)
	if len(ports) > 0 {
		requestBody.SetPorts(ports)
	}
	if ruleType, ok := ruleTypeToGraph[data.RuleType.ValueString()]; ok {
		parsed, err := networkaccess.ParseNetworkDestinationType(ruleType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rule_type: %w", err)
		}
		requestBody.SetRuleType(parsed.(*networkaccess.NetworkDestinationType))
	}
	if err := convert.FrameworkToGraphEnum(data.Protocol, networkaccess.ParseNetworkingProtocol, requestBody.SetProtocol); err != nil {
		return nil, fmt.Errorf("failed to set protocol: %w", err)
	}

	if includeID {
		if !data.ID.IsNull() && !data.ID.IsUnknown() {
			requestBody.SetId(data.ID.ValueStringPointer())
		}
		return requestBody, nil
	}

	requestBody.SetName(data.Name.ValueStringPointer())
	if err := convert.FrameworkToGraphEnum(data.Action, networkaccess.ParseForwardingRuleAction, requestBody.SetAction); err != nil {
		return nil, fmt.Errorf("failed to set action: %w", err)
	}
	// Graph accepts clientFallbackAction only on create.
	if err := convert.FrameworkToGraphEnum(data.ClientFallbackAction, networkaccess.ParseClientFallbackAction, requestBody.SetClientFallbackAction); err != nil {
		return nil, fmt.Errorf("failed to set client_fallback_action: %w", err)
	}

	return requestBody, nil
}
