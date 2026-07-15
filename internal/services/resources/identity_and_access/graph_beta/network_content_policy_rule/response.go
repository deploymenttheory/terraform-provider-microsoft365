package graphBetaNetworkContentPolicyRule

import (
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

type contentPolicyRuleResponse struct {
	id               *string
	name             *string
	description      *string
	action           *string
	priority         *int64
	status           *string
	activities       []string
	contentTypes     []string
	textContentTypes []string
	destinations     []contentPolicyRuleDestinationResponse
	sessionTypes     []string
}

func createContentPolicyRuleResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &contentPolicyRuleResponse{}, nil
}

func (r *contentPolicyRuleResponse) Serialize(writer s.SerializationWriter) error { return nil }

func (r *contentPolicyRuleResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"id":          func(n s.ParseNode) error { value, err := n.GetStringValue(); r.id = value; return err },
		"name":        func(n s.ParseNode) error { value, err := n.GetStringValue(); r.name = value; return err },
		"description": func(n s.ParseNode) error { value, err := n.GetStringValue(); r.description = value; return err },
		"action":      func(n s.ParseNode) error { value, err := n.GetStringValue(); r.action = value; return err },
		"priority":    func(n s.ParseNode) error { value, err := n.GetInt64Value(); r.priority = value; return err },
		"settings": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createContentPolicyRuleSettingsResponseFromDiscriminatorValue)
			if err != nil || value == nil {
				return err
			}
			if settings, ok := value.(*contentPolicyRuleSettingsResponse); ok {
				r.status = settings.status
			}
			return nil
		},
		"matchingConditions": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createContentPolicyRuleMatchingConditionsResponseFromDiscriminatorValue)
			if err != nil || value == nil {
				return err
			}
			conditions, ok := value.(*contentPolicyRuleMatchingConditionsResponse)
			if !ok {
				return nil
			}
			if conditions.fileAttributes != nil {
				r.activities = conditions.fileAttributes.activities
				r.contentTypes = conditions.fileAttributes.contentTypes
				r.textContentTypes = conditions.fileAttributes.textContentTypes
			}
			r.destinations = conditions.destinations
			if conditions.sources != nil {
				r.sessionTypes = conditions.sources.sessionTypes
			}
			return nil
		},
	}
}

type contentPolicyRuleSettingsResponse struct{ status *string }

func createContentPolicyRuleSettingsResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &contentPolicyRuleSettingsResponse{}, nil
}
func (r *contentPolicyRuleSettingsResponse) Serialize(writer s.SerializationWriter) error { return nil }
func (r *contentPolicyRuleSettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{"status": func(n s.ParseNode) error { value, err := n.GetStringValue(); r.status = value; return err }}
}

type contentPolicyRuleMatchingConditionsResponse struct {
	fileAttributes *contentPolicyRuleFileAttributesResponse
	destinations   []contentPolicyRuleDestinationResponse
	sources        *contentPolicyRuleSourcesResponse
}

func createContentPolicyRuleMatchingConditionsResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &contentPolicyRuleMatchingConditionsResponse{}, nil
}
func (r *contentPolicyRuleMatchingConditionsResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}
func (r *contentPolicyRuleMatchingConditionsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"fileAttributes": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createContentPolicyRuleFileAttributesResponseFromDiscriminatorValue)
			if err != nil || value == nil {
				return err
			}
			if result, ok := value.(*contentPolicyRuleFileAttributesResponse); ok {
				r.fileAttributes = result
			}
			return nil
		},
		"destinations": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfObjectValues(createContentPolicyRuleDestinationResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			for _, value := range values {
				if destination, ok := value.(*contentPolicyRuleDestinationResponse); ok {
					r.destinations = append(r.destinations, *destination)
				}
			}
			return nil
		},
		"sources": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createContentPolicyRuleSourcesResponseFromDiscriminatorValue)
			if err != nil || value == nil {
				return err
			}
			if result, ok := value.(*contentPolicyRuleSourcesResponse); ok {
				r.sources = result
			}
			return nil
		},
	}
}

type contentPolicyRuleFileAttributesResponse struct {
	activities       []string
	contentTypes     []string
	textContentTypes []string
}

func createContentPolicyRuleFileAttributesResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &contentPolicyRuleFileAttributesResponse{}, nil
}
func (r *contentPolicyRuleFileAttributesResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}
func (r *contentPolicyRuleFileAttributesResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"activities": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if value != nil {
				r.activities = splitContentPolicyRuleValues(*value)
			}
			return err
		},
		"contentTypes": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfPrimitiveValues("string")
			r.contentTypes = contentPolicyRulePrimitiveStrings(values)
			return err
		},
		"textContentTypes": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfPrimitiveValues("string")
			r.textContentTypes = contentPolicyRulePrimitiveStrings(values)
			return err
		},
	}
}

type contentPolicyRuleDestinationResponse struct {
	odataType *string
	values    []string
}

func createContentPolicyRuleDestinationResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &contentPolicyRuleDestinationResponse{}, nil
}
func (r *contentPolicyRuleDestinationResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}
func (r *contentPolicyRuleDestinationResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error { value, err := n.GetStringValue(); r.odataType = value; return err },
		"values": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfPrimitiveValues("string")
			r.values = contentPolicyRulePrimitiveStrings(values)
			return err
		},
	}
}

type contentPolicyRuleSourcesResponse struct{ sessionTypes []string }

func createContentPolicyRuleSourcesResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &contentPolicyRuleSourcesResponse{}, nil
}
func (r *contentPolicyRuleSourcesResponse) Serialize(writer s.SerializationWriter) error { return nil }
func (r *contentPolicyRuleSourcesResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{"sessionType": func(n s.ParseNode) error {
		value, err := n.GetStringValue()
		if value != nil {
			r.sessionTypes = splitContentPolicyRuleValues(*value)
		}
		return err
	}}
}

func splitContentPolicyRuleValues(value string) []string {
	parts := helpers.SplitCommaSeparatedString(value)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func contentPolicyRulePrimitiveStrings(values []any) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if stringValue, ok := value.(*string); ok && stringValue != nil {
			result = append(result, *stringValue)
		}
	}
	return result
}
