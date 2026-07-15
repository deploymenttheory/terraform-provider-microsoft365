package graphBetaNetworkContentPolicy

import s "github.com/microsoft/kiota-abstractions-go/serialization"

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
