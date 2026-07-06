package graphBetaNetworkWebFilteringPolicyRule

import (
	"strings"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

const (
	// These OData discriminator values are observed from the Entra portal Web
	// content filtering blade. They are not exposed as webFiltering* models in
	// the Microsoft Graph beta SDK, so request and response bodies are parsed with
	// local Kiota Parsable implementations.
	webFilteringRuleODataType                   = "#microsoft.graph.networkaccess.webFilteringRule"
	webFilteringActionAllowODataType            = "#microsoft.graph.networkaccess.webFilteringActionAllow"
	webFilteringActionBlockODataType            = "#microsoft.graph.networkaccess.webFilteringActionBlock"
	webFilteringURLDestinationODataType         = "#microsoft.graph.networkaccess.webFilteringUrlDestination"
	webFilteringWebCategoryDestinationODataType = "#microsoft.graph.networkaccess.webFilteringWebCategoryDestination"
	headerModificationAddODataType              = "#microsoft.graph.networkaccess.headerModificationAdd"
)

type webFilteringPolicyRuleResponse struct {
	id            *string
	name          *string
	description   *string
	priority      *int64
	action        *string
	status        *string
	urlsOrFqdns   []string
	webCategories []string
	httpMethods   []string
	sessionTypes  []string
	customHeaders []customHeaderResponse
}

type customHeaderResponse struct {
	headerName  *string
	headerValue *string
}

func newWebFilteringPolicyRuleResponse() *webFilteringPolicyRuleResponse {
	return &webFilteringPolicyRuleResponse{}
}

func createWebFilteringPolicyRuleResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return newWebFilteringPolicyRuleResponse(), nil
}

func (r *webFilteringPolicyRuleResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *webFilteringPolicyRuleResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
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
		"priority": func(n s.ParseNode) error {
			value, err := n.GetInt64Value()
			if err != nil {
				return err
			}
			r.priority = value
			return nil
		},
		"action": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createActionResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			action, ok := value.(*actionResponse)
			if ok {
				r.action = terraformAction(action.odataType)
				r.customHeaders = action.customHeaders
			}
			return nil
		},
		"settings": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createRuleSettingsResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			settings, ok := value.(*ruleSettingsResponse)
			if ok {
				r.status = settings.status
			}
			return nil
		},
		"matchingConditions": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createMatchingConditionsResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			conditions, ok := value.(*matchingConditionsResponse)
			if !ok {
				return nil
			}
			if conditions.destinations != nil {
				r.urlsOrFqdns = conditions.destinations.urlsOrFqdns
				r.webCategories = conditions.destinations.webCategories
				r.httpMethods = conditions.destinations.httpMethods
			}
			if conditions.sources != nil {
				r.sessionTypes = conditions.sources.sessionTypes
			}
			return nil
		},
		"customHeaders": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfObjectValues(createCustomHeaderResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if values == nil {
				return nil
			}
			r.customHeaders = make([]customHeaderResponse, 0, len(values))
			for _, value := range values {
				if value == nil {
					continue
				}
				header, ok := value.(*customHeaderResponse)
				if !ok {
					continue
				}
				r.customHeaders = append(r.customHeaders, *header)
			}
			return nil
		},
	}
}

type actionResponse struct {
	odataType     *string
	customHeaders []customHeaderResponse
}

func createActionResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &actionResponse{}, nil
}

func (r *actionResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *actionResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.odataType = value
			return nil
		},
		"headerSettings": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createHeaderSettingsResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			headerSettings, ok := value.(*headerSettingsResponse)
			if ok {
				r.customHeaders = headerSettings.modifications
			}
			return nil
		},
	}
}

type headerSettingsResponse struct {
	modifications []customHeaderResponse
}

func createHeaderSettingsResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &headerSettingsResponse{}, nil
}

func (r *headerSettingsResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *headerSettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"modifications": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfObjectValues(createCustomHeaderResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if values == nil {
				return nil
			}
			r.modifications = make([]customHeaderResponse, 0, len(values))
			for _, value := range values {
				if value == nil {
					continue
				}
				header, ok := value.(*customHeaderResponse)
				if !ok {
					continue
				}
				r.modifications = append(r.modifications, *header)
			}
			return nil
		},
	}
}

type ruleSettingsResponse struct {
	status *string
}

func createRuleSettingsResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &ruleSettingsResponse{}, nil
}

func (r *ruleSettingsResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *ruleSettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"status": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.status = value
			return nil
		},
	}
}

type matchingConditionsResponse struct {
	destinations *destinationsResponse
	sources      *sourcesResponse
}

func createMatchingConditionsResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &matchingConditionsResponse{}, nil
}

func (r *matchingConditionsResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *matchingConditionsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"destinations": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createDestinationsResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			destinations, ok := value.(*destinationsResponse)
			if ok {
				r.destinations = destinations
			}
			return nil
		},
		"sources": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createSourcesResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if value == nil {
				return nil
			}
			sources, ok := value.(*sourcesResponse)
			if ok {
				r.sources = sources
			}
			return nil
		},
	}
}

type destinationsResponse struct {
	urlsOrFqdns   []string
	webCategories []string
	httpMethods   []string
}

func createDestinationsResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &destinationsResponse{}, nil
}

func (r *destinationsResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *destinationsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"targets": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfObjectValues(createDestinationTargetResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if values == nil {
				return nil
			}
			for _, value := range values {
				if value == nil {
					continue
				}
				target, ok := value.(*destinationTargetResponse)
				if !ok {
					continue
				}
				switch target.odataType {
				case webFilteringURLDestinationODataType:
					r.urlsOrFqdns = append(r.urlsOrFqdns, target.values...)
				case webFilteringWebCategoryDestinationODataType:
					r.webCategories = append(r.webCategories, target.values...)
				}
			}
			return nil
		},
		"httpRequestMethod": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.httpMethods = splitCommaValues(value)
			return nil
		},
	}
}

type destinationTargetResponse struct {
	odataType string
	values    []string
}

func createDestinationTargetResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &destinationTargetResponse{}, nil
}

func (r *destinationTargetResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *destinationTargetResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			if value != nil {
				r.odataType = *value
			}
			return nil
		},
		"values": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfPrimitiveValues("string")
			if err != nil {
				return err
			}
			if values == nil {
				return nil
			}
			r.values = make([]string, 0, len(values))
			for _, value := range values {
				if value == nil {
					continue
				}
				r.values = append(r.values, *(value.(*string)))
			}
			return nil
		},
	}
}

type sourcesResponse struct {
	sessionTypes []string
}

func createSourcesResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &sourcesResponse{}, nil
}

func (r *sourcesResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *sourcesResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"sessionType": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.sessionTypes = splitCommaValues(value)
			return nil
		},
	}
}

func createCustomHeaderResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &customHeaderResponse{}, nil
}

func (r *customHeaderResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *customHeaderResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"headerName": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.headerName = value
			return nil
		},
		"headerValue": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.headerValue = value
			return nil
		},
	}
}

func splitCommaValues(value *string) []string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}

	parts := strings.Split(*value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, strings.ToLower(trimmed))
		}
	}

	return result
}

func graphActionODataType(action string) string {
	if action == "block" {
		return webFilteringActionBlockODataType
	}

	return webFilteringActionAllowODataType
}

func terraformAction(odataType *string) *string {
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
