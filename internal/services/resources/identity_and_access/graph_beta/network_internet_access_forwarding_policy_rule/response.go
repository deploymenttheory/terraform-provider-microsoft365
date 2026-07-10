package graphBetaNetworkInternetAccessForwardingPolicyRule

import s "github.com/microsoft/kiota-abstractions-go/serialization"

type internetAccessForwardingRuleResponse struct {
	id                   *string
	name                 *string
	ruleType             *string
	action               *string
	clientFallbackAction *string
	ports                []string
	protocol             *string
	destinations         []ruleDestinationResponse
}

func createInternetAccessForwardingRuleResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &internetAccessForwardingRuleResponse{}, nil
}

func (r *internetAccessForwardingRuleResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *internetAccessForwardingRuleResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"id": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.id = value
			return nil
		},
		"name": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.name = value
			return nil
		},
		"ruleType": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.ruleType = value
			return nil
		},
		"action": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.action = value
			return nil
		},
		"clientFallbackAction": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.clientFallbackAction = value
			return nil
		},
		"ports": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfPrimitiveValues("string")
			if err != nil {
				return err
			}
			r.ports = primitiveStringSlice(values)
			return nil
		},
		"protocol": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.protocol = value
			return nil
		},
		"destinations": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfObjectValues(createRuleDestinationResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			r.destinations = nil
			for _, value := range values {
				if typed, ok := value.(*ruleDestinationResponse); ok {
					r.destinations = append(r.destinations, *typed)
				}
			}
			return nil
		},
	}
}

func primitiveStringSlice(values []any) []string {
	if values == nil {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		switch v := value.(type) {
		case string:
			result = append(result, v)
		case *string:
			if v != nil {
				result = append(result, *v)
			}
		}
	}
	return result
}

type ruleDestinationResponse struct {
	odataType    *string
	value        *string
	beginAddress *string
	endAddress   *string
}

func createRuleDestinationResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &ruleDestinationResponse{}, nil
}

func (r *ruleDestinationResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *ruleDestinationResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.odataType = value
			return nil
		},
		"value": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.value = value
			return nil
		},
		"beginAddress": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.beginAddress = value
			return nil
		},
		"endAddress": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.endAddress = value
			return nil
		},
	}
}
