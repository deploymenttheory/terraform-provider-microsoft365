package graphBetaNetworkTLSInspectionPolicyRule

import (
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Custom parsing preserves portal fields omitted from Graph metadata, notably action.
type tlsInspectionPolicyRuleResponse struct {
	id, name, description, action, status *string
	priority                              *int64
	conditions                            *tlsInspectionPolicyRuleConditionsResponse
}

func createTLSInspectionPolicyRuleResponseFromDiscriminatorValue(
	n s.ParseNode,
) (s.Parsable, error) {
	return &tlsInspectionPolicyRuleResponse{}, nil
}
func (r *tlsInspectionPolicyRuleResponse) Serialize(w s.SerializationWriter) error { return nil }

func (r *tlsInspectionPolicyRuleResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"id":          func(n s.ParseNode) error { v, e := n.GetStringValue(); r.id = v; return wrapResponseError(e) },
		"name":        func(n s.ParseNode) error { v, e := n.GetStringValue(); r.name = v; return wrapResponseError(e) },
		"description": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.description = v; return wrapResponseError(e) },
		"action":      func(n s.ParseNode) error { v, e := n.GetStringValue(); r.action = v; return wrapResponseError(e) },
		"priority":    func(n s.ParseNode) error { v, e := n.GetInt64Value(); r.priority = v; return wrapResponseError(e) },
		"settings": func(n s.ParseNode) error {
			v, e := n.GetObjectValue(
				func(s.ParseNode) (s.Parsable, error) { return &tlsInspectionPolicyRuleSettingsResponse{}, nil },
			)
			if e == nil && v != nil {
				r.status = v.(*tlsInspectionPolicyRuleSettingsResponse).status
			}
			return wrapResponseError(e)
		},
		"matchingConditions": func(n s.ParseNode) error {
			v, e := n.GetObjectValue(
				func(s.ParseNode) (s.Parsable, error) { return &tlsInspectionPolicyRuleConditionsResponse{}, nil },
			)
			if e == nil && v != nil {
				r.conditions = v.(*tlsInspectionPolicyRuleConditionsResponse)
			}
			return wrapResponseError(e)
		},
	}
}

type tlsInspectionPolicyRuleSettingsResponse struct{ status *string }

func (r *tlsInspectionPolicyRuleSettingsResponse) Serialize(w s.SerializationWriter) error {
	return nil
}

func (r *tlsInspectionPolicyRuleSettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"status": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.status = v; return wrapResponseError(e) },
	}
}

type tlsInspectionPolicyRuleConditionsResponse struct {
	destinations []*tlsInspectionPolicyRuleDestinationResponse
}

func (r *tlsInspectionPolicyRuleConditionsResponse) Serialize(w s.SerializationWriter) error {
	return nil
}

func (r *tlsInspectionPolicyRuleConditionsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{"destinations": func(n s.ParseNode) error {
		values, e := n.GetCollectionOfObjectValues(
			func(s.ParseNode) (s.Parsable, error) { return &tlsInspectionPolicyRuleDestinationResponse{}, nil },
		)
		if e != nil {
			return wrapResponseError(e)
		}
		for _, v := range values {
			if v != nil {
				r.destinations = append(
					r.destinations,
					v.(*tlsInspectionPolicyRuleDestinationResponse),
				)
			}
		}
		return nil
	}}
}

type tlsInspectionPolicyRuleDestinationResponse struct {
	odataType *string
	values    []string
}

func (r *tlsInspectionPolicyRuleDestinationResponse) Serialize(w s.SerializationWriter) error {
	return nil
}

func (r *tlsInspectionPolicyRuleDestinationResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.odataType = v; return wrapResponseError(e) },
		"values": func(n s.ParseNode) error {
			values, e := n.GetCollectionOfPrimitiveValues("string")
			if e != nil {
				return wrapResponseError(e)
			}
			for _, v := range values {
				if value, ok := v.(*string); ok && value != nil {
					r.values = append(r.values, *value)
				}
			}
			return nil
		},
	}
}

func wrapResponseError(err error) error {
	if err != nil {
		return fmt.Errorf("parse TLS inspection policy rule response: %w", err)
	}
	return nil
}
