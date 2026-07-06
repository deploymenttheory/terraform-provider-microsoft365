package graphBetaNetworkWebFilteringPolicy

import (
	"strings"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

const (
	defaultActionAllowODataType = "#microsoft.graph.networkaccess.webFilteringActionAllow"
	defaultActionBlockODataType = "#microsoft.graph.networkaccess.webFilteringActionBlock"
)

type webFilteringPolicyResponse struct {
	id            *string
	name          *string
	description   *string
	defaultAction *string
}

func newWebFilteringPolicyResponse() *webFilteringPolicyResponse {
	return &webFilteringPolicyResponse{}
}

func createWebFilteringPolicyResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return newWebFilteringPolicyResponse(), nil
}

func (r *webFilteringPolicyResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *webFilteringPolicyResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
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
			value, err := n.GetObjectValue(createWebFilteringPolicySettingsResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			settings, ok := value.(*webFilteringPolicySettingsResponse)
			if !ok {
				return nil
			}
			r.defaultAction = settings.defaultAction
			return nil
		},
	}
}

type webFilteringPolicySettingsResponse struct {
	defaultAction *string
}

func createWebFilteringPolicySettingsResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &webFilteringPolicySettingsResponse{}, nil
}

func (r *webFilteringPolicySettingsResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *webFilteringPolicySettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"defaultAction": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createWebFilteringPolicyDefaultActionResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			action, ok := value.(*webFilteringPolicyDefaultActionResponse)
			if !ok {
				return nil
			}
			r.defaultAction = action.odataType
			return nil
		},
	}
}

type webFilteringPolicyDefaultActionResponse struct {
	odataType *string
}

func createWebFilteringPolicyDefaultActionResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &webFilteringPolicyDefaultActionResponse{}, nil
}

func (r *webFilteringPolicyDefaultActionResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *webFilteringPolicyDefaultActionResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.odataType = value
			return nil
		},
	}
}

func terraformDefaultAction(odataType *string) *string {
	if odataType == nil {
		return nil
	}

	value := strings.ToLower(*odataType)
	switch {
	case strings.HasSuffix(value, "webfilteringactionallow"):
		action := "allow"
		return &action
	case strings.HasSuffix(value, "webfilteringactionblock"):
		action := "block"
		return &action
	default:
		return odataType
	}
}

func graphDefaultActionODataType(action string) string {
	if action == "block" {
		return defaultActionBlockODataType
	}

	return defaultActionAllowODataType
}
