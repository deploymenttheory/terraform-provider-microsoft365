package graphBetaNetworkWebContentFilteringPolicy

import (
	"strings"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

const (
	defaultActionAllowODataType = "#microsoft.graph.networkaccess.webFilteringActionAllow"
	defaultActionBlockODataType = "#microsoft.graph.networkaccess.webFilteringActionBlock"
)

type webContentFilteringPolicyResponse struct {
	id            *string
	name          *string
	description   *string
	defaultAction *string
}

func newWebContentFilteringPolicyResponse() *webContentFilteringPolicyResponse {
	return &webContentFilteringPolicyResponse{}
}

func createWebContentFilteringPolicyResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return newWebContentFilteringPolicyResponse(), nil
}

func (r *webContentFilteringPolicyResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *webContentFilteringPolicyResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
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
			value, err := n.GetObjectValue(createWebContentFilteringPolicySettingsResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			settings, ok := value.(*webContentFilteringPolicySettingsResponse)
			if !ok {
				return nil
			}
			r.defaultAction = settings.defaultAction
			return nil
		},
	}
}

type webContentFilteringPolicySettingsResponse struct {
	defaultAction *string
}

func createWebContentFilteringPolicySettingsResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &webContentFilteringPolicySettingsResponse{}, nil
}

func (r *webContentFilteringPolicySettingsResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *webContentFilteringPolicySettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"defaultAction": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createWebContentFilteringPolicyDefaultActionResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			action, ok := value.(*webContentFilteringPolicyDefaultActionResponse)
			if !ok {
				return nil
			}
			r.defaultAction = action.odataType
			return nil
		},
	}
}

type webContentFilteringPolicyDefaultActionResponse struct {
	odataType *string
}

func createWebContentFilteringPolicyDefaultActionResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &webContentFilteringPolicyDefaultActionResponse{}, nil
}

func (r *webContentFilteringPolicyDefaultActionResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *webContentFilteringPolicyDefaultActionResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
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
