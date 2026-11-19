package graphBetaNetworkContentPolicy

import (
	stderrors "errors"
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

type contentPolicyResponse struct {
	id            *string
	name          *string
	description   *string
	defaultAction *string
}

func createContentPolicyResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &contentPolicyResponse{}, nil
}

func (r *contentPolicyResponse) Serialize(writer s.SerializationWriter) error { return nil }

func (r *contentPolicyResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
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
		"description": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.description = value
			return nil
		},
		"settings": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createContentPolicySettingsResponseFromDiscriminatorValue)
			if err != nil || value == nil {
				return err
			}
			if settings, ok := value.(*contentPolicySettingsResponse); ok {
				r.defaultAction = settings.defaultAction
			}
			return nil
		},
	}
}

type contentPolicySettingsResponse struct{ defaultAction *string }

func createContentPolicySettingsResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &contentPolicySettingsResponse{}, nil
}

func (r *contentPolicySettingsResponse) Serialize(writer s.SerializationWriter) error { return nil }

func (r *contentPolicySettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"defaultAction": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.defaultAction = value
			return nil
		},
	}
}

// decodeSendContentPolicyResponse validates and decodes a Microsoft Graph response.
func decodeSendContentPolicyResponse(result any, err error) (*contentPolicyResponse, error) {
	if err != nil || result == nil {
		return nil, err
	}
	policy, ok := result.(*contentPolicyResponse)
	if !ok {
		return nil, fmt.Errorf("content policy request returned %T, %w", result, errUnexpectedResponseType)
	}
	return policy, nil
}

// decodeCreateContentPolicyResponse validates and decodes the Microsoft Graph response.
func decodeCreateContentPolicyResponse(transportResult any, err error) (*contentPolicyResponse, error) {
	result, err := decodeSendContentPolicyResponse(transportResult, err)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errCreateResponseNil
	}
	return result, nil
}

// decodeUpdateContentPolicyResponse validates and decodes the Microsoft Graph response.
func decodeUpdateContentPolicyResponse(result any, err error) error {
	if err != nil {
		return err
	}
	if result == nil {
		return errUpdateResponseNil
	}
	return nil
}

var (
	errUnexpectedResponseType = stderrors.New("expected contentPolicyResponse")
	errCreateResponseNil      = stderrors.New("create content policy returned nil response")
	errUpdateResponseNil      = stderrors.New("update content policy returned nil response")
)
