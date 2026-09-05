package graphBetaNetworkPromptPolicyRule

import (
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Custom parsing preserves portal fields omitted from Graph metadata, notably action.
type promptPolicyRuleResponse struct {
	id, name, description, action, status, promptLogging, odataType *string
	priority                                                        *int64
	conditions                                                      *promptPolicyRuleConditionsResponse
}

func createPromptPolicyRuleResponseFromDiscriminatorValue(
	n s.ParseNode,
) (s.Parsable, error) {
	return &promptPolicyRuleResponse{}, nil
}
func (r *promptPolicyRuleResponse) Serialize(w s.SerializationWriter) error { return nil }

func (r *promptPolicyRuleResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.odataType = v; return wrapResponseError(e) },
		"promptLogging": func(n s.ParseNode) error {
			v, e := n.GetStringValue()
			r.promptLogging = v
			return wrapResponseError(e)
		},
		"id":          func(n s.ParseNode) error { v, e := n.GetStringValue(); r.id = v; return wrapResponseError(e) },
		"name":        func(n s.ParseNode) error { v, e := n.GetStringValue(); r.name = v; return wrapResponseError(e) },
		"description": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.description = v; return wrapResponseError(e) },
		"action":      func(n s.ParseNode) error { v, e := n.GetStringValue(); r.action = v; return wrapResponseError(e) },
		"priority":    func(n s.ParseNode) error { v, e := n.GetInt64Value(); r.priority = v; return wrapResponseError(e) },
		"settings": func(n s.ParseNode) error {
			v, e := n.GetObjectValue(
				func(s.ParseNode) (s.Parsable, error) { return &promptPolicyRuleSettingsResponse{}, nil },
			)
			if e == nil && v != nil {
				r.status = v.(*promptPolicyRuleSettingsResponse).status
			}
			return wrapResponseError(e)
		},
		"matchingConditions": func(n s.ParseNode) error {
			v, e := n.GetObjectValue(
				func(s.ParseNode) (s.Parsable, error) { return &promptPolicyRuleConditionsResponse{}, nil },
			)
			if e == nil && v != nil {
				r.conditions = v.(*promptPolicyRuleConditionsResponse)
			}
			return wrapResponseError(e)
		},
	}
}

type promptPolicyRuleSettingsResponse struct{ status *string }

func (r *promptPolicyRuleSettingsResponse) Serialize(w s.SerializationWriter) error {
	return nil
}

func (r *promptPolicyRuleSettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"status": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.status = v; return wrapResponseError(e) },
	}
}

type promptPolicyRuleConditionsResponse struct {
	scanResult     *string
	schemes        []*conversationSchemeResponse
	schemesPresent bool
}

func (r *promptPolicyRuleConditionsResponse) Serialize(w s.SerializationWriter) error { return nil }
func (r *promptPolicyRuleConditionsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"scanResult": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.scanResult = v; return wrapResponseError(e) },
		"conversationSchemes": func(n s.ParseNode) error {
			values, e := n.GetCollectionOfObjectValues(func(s.ParseNode) (s.Parsable, error) { return &conversationSchemeResponse{}, nil })
			if e != nil {
				return wrapResponseError(e)
			}
			r.schemesPresent = true
			for _, v := range values {
				if v == nil {
					return fmt.Errorf("%w: null conversation scheme", errInvalidResponse)
				}
				r.schemes = append(r.schemes, v.(*conversationSchemeResponse))
			}
			return nil
		},
	}
}

type conversationSchemeResponse struct{ odataType, url, jsonPath, schemeName *string }

func (r *conversationSchemeResponse) Serialize(w s.SerializationWriter) error { return nil }
func (r *conversationSchemeResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.odataType = v; return wrapResponseError(e) },
		"url":         func(n s.ParseNode) error { v, e := n.GetStringValue(); r.url = v; return wrapResponseError(e) },
		"jsonPath":    func(n s.ParseNode) error { v, e := n.GetStringValue(); r.jsonPath = v; return wrapResponseError(e) },
		"schemeName":  func(n s.ParseNode) error { v, e := n.GetStringValue(); r.schemeName = v; return wrapResponseError(e) },
	}
}

func wrapResponseError(err error) error {
	if err != nil {
		return fmt.Errorf("parse prompt policy rule response: %w", err)
	}
	return nil
}
