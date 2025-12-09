package graphBetaAgentIdentity

import (
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AgentIdentityResponse represents the response from creating an agent identity
type AgentIdentityResponse struct {
	id                        *string
	displayName               *string
	agentIdentityBlueprintId  *string
	accountEnabled            *bool
	createdByAppId            *string
	createdDateTime           *string
	servicePrincipalType      *string
	disabledByMicrosoftStatus *string
	tags                      []string
	additionalData            map[string]any
}

// CreateAgentIdentityResponseFactory creates a new instance of AgentIdentityResponse
func CreateAgentIdentityResponseFactory() s.ParsableFactory {
	return func(parseNode s.ParseNode) (s.Parsable, error) {
		return NewAgentIdentityResponse(), nil
	}
}

// NewAgentIdentityResponse creates a new AgentIdentityResponse
func NewAgentIdentityResponse() *AgentIdentityResponse {
	return &AgentIdentityResponse{
		additionalData: make(map[string]any),
	}
}

// GetFieldDeserializers returns the deserialization information for this object
func (r *AgentIdentityResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"id": func(n s.ParseNode) error {
			val, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.id = val
			return nil
		},
		"displayName": func(n s.ParseNode) error {
			val, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.displayName = val
			return nil
		},
		"agentIdentityBlueprintId": func(n s.ParseNode) error {
			val, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.agentIdentityBlueprintId = val
			return nil
		},
		"accountEnabled": func(n s.ParseNode) error {
			val, err := n.GetBoolValue()
			if err != nil {
				return err
			}
			r.accountEnabled = val
			return nil
		},
		"createdByAppId": func(n s.ParseNode) error {
			val, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.createdByAppId = val
			return nil
		},
		"createdDateTime": func(n s.ParseNode) error {
			val, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.createdDateTime = val
			return nil
		},
		"servicePrincipalType": func(n s.ParseNode) error {
			val, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.servicePrincipalType = val
			return nil
		},
		"disabledByMicrosoftStatus": func(n s.ParseNode) error {
			val, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.disabledByMicrosoftStatus = val
			return nil
		},
		"tags": func(n s.ParseNode) error {
			val, err := n.GetCollectionOfPrimitiveValues("string")
			if err != nil {
				return err
			}
			if val != nil {
				tags := make([]string, len(val))
				for i, v := range val {
					if str, ok := v.(*string); ok && str != nil {
						tags[i] = *str
					}
				}
				r.tags = tags
			}
			return nil
		},
	}
}

// Serialize writes the object to the given writer
func (r *AgentIdentityResponse) Serialize(writer s.SerializationWriter) error {
	return nil // Read-only response, no serialization needed
}

// GetAdditionalData gets the additional data for this object
func (r *AgentIdentityResponse) GetAdditionalData() map[string]any {
	return r.additionalData
}

// SetAdditionalData sets the additional data for this object
func (r *AgentIdentityResponse) SetAdditionalData(value map[string]any) {
	r.additionalData = value
}

// GetId returns the id property
func (r *AgentIdentityResponse) GetId() *string {
	return r.id
}

// GetDisplayName returns the displayName property
func (r *AgentIdentityResponse) GetDisplayName() *string {
	return r.displayName
}

// GetAgentIdentityBlueprintId returns the agentIdentityBlueprintId property
func (r *AgentIdentityResponse) GetAgentIdentityBlueprintId() *string {
	return r.agentIdentityBlueprintId
}

// GetAccountEnabled returns the accountEnabled property
func (r *AgentIdentityResponse) GetAccountEnabled() *bool {
	return r.accountEnabled
}

// GetCreatedByAppId returns the createdByAppId property
func (r *AgentIdentityResponse) GetCreatedByAppId() *string {
	return r.createdByAppId
}

// GetCreatedDateTime returns the createdDateTime property
func (r *AgentIdentityResponse) GetCreatedDateTime() *string {
	return r.createdDateTime
}

// GetServicePrincipalType returns the servicePrincipalType property
func (r *AgentIdentityResponse) GetServicePrincipalType() *string {
	return r.servicePrincipalType
}

// GetDisabledByMicrosoftStatus returns the disabledByMicrosoftStatus property
func (r *AgentIdentityResponse) GetDisabledByMicrosoftStatus() *string {
	return r.disabledByMicrosoftStatus
}

// GetTags returns the tags property
func (r *AgentIdentityResponse) GetTags() []string {
	return r.tags
}
