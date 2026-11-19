package graphBetaNetworkPromptPolicy

import (
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

type promptPolicyResponse struct {
	id                   *string
	name                 *string
	description          *string
	defaultAction        *string
	version              *string
	lastModifiedDateTime *string
}

func createPromptPolicyResponseFromDiscriminatorValue(
	parseNode s.ParseNode,
) (s.Parsable, error) {
	return &promptPolicyResponse{}, nil
}

func (r *promptPolicyResponse) Serialize(writer s.SerializationWriter) error { return nil }

func (r *promptPolicyResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"id": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return wrapResponseError(err)
			}
			r.id = value
			return nil
		},
		"name": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return wrapResponseError(err)
			}
			r.name = value
			return nil
		},
		"description": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return wrapResponseError(err)
			}
			r.description = value
			return nil
		},
		"version": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.version = v; return wrapResponseError(e) },
		"lastModifiedDateTime": func(n s.ParseNode) error {
			v, e := n.GetStringValue()
			r.lastModifiedDateTime = v
			return wrapResponseError(e)
		},
		"settings": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(
				createPromptPolicySettingsResponseFromDiscriminatorValue,
			)
			if err != nil || value == nil {
				return wrapResponseError(err)
			}
			if settings, ok := value.(*promptPolicySettingsResponse); ok {
				r.defaultAction = settings.defaultAction
			}
			return nil
		},
	}
}

type promptPolicySettingsResponse struct{ defaultAction *string }

func createPromptPolicySettingsResponseFromDiscriminatorValue(
	parseNode s.ParseNode,
) (s.Parsable, error) {
	return &promptPolicySettingsResponse{}, nil
}

func (r *promptPolicySettingsResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *promptPolicySettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"defaultAction": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return wrapResponseError(err)
			}
			r.defaultAction = value
			return nil
		},
	}
}

func wrapResponseError(err error) error {
	if err != nil {
		return fmt.Errorf("parse prompt policy response: %w", err)
	}
	return nil
}

// decodeSendPromptPolicyResponse validates and decodes a Microsoft Graph response.
func decodeSendPromptPolicyResponse(result any, err error) (*promptPolicyResponse, error) {
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	policy, ok := result.(*promptPolicyResponse)
	if !ok {
		return nil, fmt.Errorf("%w: received %T", errInvalidResponse, result)
	}
	return policy, nil
}

// decodeCreatePromptPolicyResponse validates and decodes the Microsoft Graph response.
func decodeCreatePromptPolicyResponse(transportResult any, err error) (*promptPolicyResponse, error) {
	result, err := decodeSendPromptPolicyResponse(transportResult, err)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	return result, nil
}
