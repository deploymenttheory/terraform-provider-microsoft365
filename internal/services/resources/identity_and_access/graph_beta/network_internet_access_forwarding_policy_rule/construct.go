package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

const (
	internetAccessForwardingRuleODataType = "#microsoft.graph.networkaccess.internetAccessForwardingRule"
	fqdnODataType                         = "#microsoft.graph.networkaccess.fqdn"
	ipAddressODataType                    = "#microsoft.graph.networkaccess.ipAddress"
	ipRangeODataType                      = "#microsoft.graph.networkaccess.ipRange"
	ipSubnetODataType                     = "#microsoft.graph.networkaccess.ipSubnet"
)

func constructResource(ctx context.Context, data *NetworkInternetAccessForwardingPolicyRuleResourceModel, includeID bool) (s.Parsable, error) {
	ports, err := stringSetValues(ctx, data.Ports)
	if err != nil {
		return nil, fmt.Errorf("failed to read ports: %w", err)
	}
	if err := validateDestinationsMatchRuleType(ctx, data.RuleType.ValueString(), data.Destinations); err != nil {
		return nil, err
	}
	destinations, err := destinationValues(ctx, data.Destinations)
	if err != nil {
		return nil, err
	}
	if len(destinations) == 0 {
		return nil, fmt.Errorf("at least one destination must be specified")
	}

	body := &internetAccessForwardingRuleRequestBody{
		ODataType:    internetAccessForwardingRuleODataType,
		RuleType:     graphRuleType(data.RuleType.ValueString()),
		Ports:        ports,
		Protocol:     data.Protocol.ValueString(),
		Destinations: destinations,
	}
	if includeID && !data.ID.IsNull() && !data.ID.IsUnknown() {
		body.ID = data.ID.ValueString()
	}
	if !includeID {
		body.Name = data.Name.ValueString()
		body.Action = data.Action.ValueString()
	}
	return body, nil
}

func stringSetValues(ctx context.Context, value types.Set) ([]string, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}
	var result []string
	diags := value.ElementsAs(ctx, &result, false)
	if diags.HasError() {
		return nil, fmt.Errorf("%s", diags.Errors()[0].Detail())
	}
	return result, nil
}

func destinationValues(ctx context.Context, value types.List) ([]s.Parsable, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}
	var terraformValues []RuleDestinationModel
	diags := value.ElementsAs(ctx, &terraformValues, false)
	if diags.HasError() {
		return nil, fmt.Errorf("%s", diags.Errors()[0].Detail())
	}

	result := make([]s.Parsable, 0, len(terraformValues))
	for _, item := range terraformValues {
		destination, err := destinationRequestBody(item)
		if err != nil {
			return nil, err
		}
		result = append(result, destination)
	}
	return result, nil
}

func validateDestinationsMatchRuleType(ctx context.Context, ruleType string, destinations types.List) error {
	if destinations.IsNull() || destinations.IsUnknown() {
		return nil
	}
	var terraformValues []RuleDestinationModel
	diags := destinations.ElementsAs(ctx, &terraformValues, false)
	if diags.HasError() {
		return fmt.Errorf("%s", diags.Errors()[0].Detail())
	}
	for _, destination := range terraformValues {
		if destination.Type.ValueString() != ruleType {
			return fmt.Errorf("destination type %q must match rule_type %q", destination.Type.ValueString(), ruleType)
		}
	}
	return nil
}

func destinationRequestBody(item RuleDestinationModel) (s.Parsable, error) {
	switch item.Type.ValueString() {
	case ruleTypeFQDN:
		if item.Value.IsNull() || item.Value.ValueString() == "" {
			return nil, fmt.Errorf("destination value is required when type is %s", ruleTypeFQDN)
		}
		return &ruleDestinationRequestBody{ODataType: fqdnODataType, Value: item.Value.ValueString()}, nil
	case ruleTypeIPAddress:
		if item.Value.IsNull() || item.Value.ValueString() == "" {
			return nil, fmt.Errorf("destination value is required when type is %s", ruleTypeIPAddress)
		}
		return &ruleDestinationRequestBody{ODataType: ipAddressODataType, Value: item.Value.ValueString()}, nil
	case ruleTypeIPSubnet:
		if item.Value.IsNull() || item.Value.ValueString() == "" {
			return nil, fmt.Errorf("destination value is required when type is %s", ruleTypeIPSubnet)
		}
		return &ruleDestinationRequestBody{ODataType: ipSubnetODataType, Value: item.Value.ValueString()}, nil
	case ruleTypeIPRange:
		if item.BeginAddress.IsNull() || item.BeginAddress.ValueString() == "" || item.EndAddress.IsNull() || item.EndAddress.ValueString() == "" {
			return nil, fmt.Errorf("begin_address and end_address are required when destination type is %s", ruleTypeIPRange)
		}
		return &ruleDestinationRequestBody{
			ODataType:    ipRangeODataType,
			BeginAddress: item.BeginAddress.ValueString(),
			EndAddress:   item.EndAddress.ValueString(),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported destination type %q", item.Type.ValueString())
	}
}

func graphRuleType(value string) string {
	switch value {
	case ruleTypeIPAddress:
		return "ipAddress"
	case ruleTypeIPRange:
		return "ipRange"
	case ruleTypeIPSubnet:
		return "ipSubnet"
	default:
		return value
	}
}

func terraformRuleType(value string) string {
	switch value {
	case "ipAddress":
		return ruleTypeIPAddress
	case "ipRange":
		return ruleTypeIPRange
	case "ipSubnet":
		return ruleTypeIPSubnet
	default:
		return value
	}
}

type internetAccessForwardingRuleRequestBody struct {
	ODataType    string
	ID           string
	Name         string
	RuleType     string
	Action       string
	Ports        []string
	Protocol     string
	Destinations []s.Parsable
}

func (b *internetAccessForwardingRuleRequestBody) Serialize(writer s.SerializationWriter) error {
	if b.ODataType != "" {
		if err := writer.WriteStringValue("@odata.type", &b.ODataType); err != nil {
			return err
		}
	}
	if b.ID != "" {
		if err := writer.WriteStringValue("id", &b.ID); err != nil {
			return err
		}
	}
	if b.Name != "" {
		if err := writer.WriteStringValue("name", &b.Name); err != nil {
			return err
		}
	}
	if b.Action != "" {
		if err := writer.WriteStringValue("action", &b.Action); err != nil {
			return err
		}
	}
	if b.RuleType != "" {
		if err := writer.WriteStringValue("ruleType", &b.RuleType); err != nil {
			return err
		}
	}
	if len(b.Ports) > 0 {
		if err := writer.WriteCollectionOfStringValues("ports", b.Ports); err != nil {
			return err
		}
	}
	if b.Protocol != "" {
		if err := writer.WriteStringValue("protocol", &b.Protocol); err != nil {
			return err
		}
	}
	if len(b.Destinations) > 0 {
		if err := writer.WriteCollectionOfObjectValues("destinations", b.Destinations); err != nil {
			return err
		}
	}
	return nil
}

func (b *internetAccessForwardingRuleRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type ruleDestinationRequestBody struct {
	ODataType    string
	Value        string
	BeginAddress string
	EndAddress   string
}

func (b *ruleDestinationRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("@odata.type", &b.ODataType); err != nil {
		return err
	}
	if b.Value != "" {
		if err := writer.WriteStringValue("value", &b.Value); err != nil {
			return err
		}
	}
	if b.BeginAddress != "" {
		if err := writer.WriteStringValue("beginAddress", &b.BeginAddress); err != nil {
			return err
		}
	}
	if b.EndAddress != "" {
		if err := writer.WriteStringValue("endAddress", &b.EndAddress); err != nil {
			return err
		}
	}
	return nil
}

func (b *ruleDestinationRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}
