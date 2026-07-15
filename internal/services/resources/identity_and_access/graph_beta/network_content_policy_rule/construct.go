package graphBetaNetworkContentPolicyRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	commonattr "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/attr"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

const (
	fileRuleODataType                         = "#microsoft.graph.networkaccess.fileRule"
	filePolicyWebCategoryDestinationODataType = "#microsoft.graph.networkaccess.filePolicyWebCategoryDestination"
	filePolicyFQDNDestinationODataType        = "#microsoft.graph.networkaccess.filePolicyFqdnDestination"
	filePolicyURLDestinationODataType         = "#microsoft.graph.networkaccess.filePolicyUrlDestination"
)

func constructResource(ctx context.Context, data *NetworkContentPolicyRuleResourceModel) (s.Parsable, error) {
	activities := commonattr.StringSetElements(data.Activities)
	contentTypes := commonattr.StringSetElements(data.ContentTypes)
	textContentTypes := commonattr.StringSetElements(data.TextContentTypes)
	sessionTypes := commonattr.StringSetElements(data.SessionTypes)
	destinations, err := contentPolicyRuleDestinationValues(ctx, data.Destinations)
	if err != nil {
		return nil, err
	}

	body := &contentPolicyRuleRequestBody{
		odataType:   fileRuleODataType,
		name:        data.Name.ValueStringPointer(),
		description: data.Description.ValueStringPointer(),
		action:      data.Action.ValueStringPointer(),
		priority:    data.Priority.ValueInt64Pointer(),
		settings:    &contentPolicyRuleSettingsRequestBody{status: data.Status.ValueStringPointer()},
		matchingConditions: &contentPolicyRuleMatchingConditionsRequestBody{
			fileAttributes: &contentPolicyRuleFileAttributesRequestBody{
				activities:       joinedContentPolicyRuleValues(activities),
				contentTypes:     contentTypes,
				textContentTypes: textContentTypes,
			},
			destinations: destinations,
			sources: &contentPolicyRuleSourcesRequestBody{
				sessionType: joinedContentPolicyRuleValues(sessionTypes),
			},
		},
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), body); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}
	return body, nil
}

func contentPolicyRuleDestinationValues(ctx context.Context, value types.List) ([]s.Parsable, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}
	var models []ContentPolicyRuleDestinationModel
	diags := value.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to read destinations: %s", diags.Errors()[0].Detail())
	}
	result := make([]s.Parsable, 0, len(models))
	for _, model := range models {
		values := commonattr.StringSetElements(model.Values)
		odataType, err := graphContentPolicyRuleDestinationType(model.Type.ValueString())
		if err != nil {
			return nil, err
		}
		result = append(result, &contentPolicyRuleDestinationRequestBody{odataType: odataType, values: values})
	}
	return result, nil
}

func graphContentPolicyRuleDestinationType(value string) (string, error) {
	switch value {
	case destinationTypeWebCategory:
		return filePolicyWebCategoryDestinationODataType, nil
	case destinationTypeFQDN:
		return filePolicyFQDNDestinationODataType, nil
	case destinationTypeURL:
		return filePolicyURLDestinationODataType, nil
	default:
		return "", fmt.Errorf("unsupported destination type %q", value)
	}
}

func terraformContentPolicyRuleDestinationType(value string) string {
	switch value {
	case filePolicyWebCategoryDestinationODataType:
		return destinationTypeWebCategory
	case filePolicyURLDestinationODataType:
		return destinationTypeURL
	default:
		return destinationTypeFQDN
	}
}

func joinedContentPolicyRuleValues(values []string) *string {
	if len(values) == 0 {
		return nil
	}
	value := helpers.JoinWithSeparator(values, ",")
	return &value
}

type contentPolicyRuleRequestBody struct {
	odataType          string
	name               *string
	description        *string
	action             *string
	priority           *int64
	settings           *contentPolicyRuleSettingsRequestBody
	matchingConditions *contentPolicyRuleMatchingConditionsRequestBody
}

func (b *contentPolicyRuleRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("@odata.type", &b.odataType); err != nil {
		return err
	}
	if err := writer.WriteStringValue("name", b.name); err != nil {
		return err
	}
	if err := writer.WriteStringValue("description", b.description); err != nil {
		return err
	}
	if err := writer.WriteStringValue("action", b.action); err != nil {
		return err
	}
	if err := writer.WriteInt64Value("priority", b.priority); err != nil {
		return err
	}
	if err := writer.WriteObjectValue("settings", b.settings); err != nil {
		return err
	}
	return writer.WriteObjectValue("matchingConditions", b.matchingConditions)
}

func (b *contentPolicyRuleRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type contentPolicyRuleSettingsRequestBody struct{ status *string }

func (b *contentPolicyRuleSettingsRequestBody) Serialize(writer s.SerializationWriter) error {
	return writer.WriteStringValue("status", b.status)
}
func (b *contentPolicyRuleSettingsRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type contentPolicyRuleMatchingConditionsRequestBody struct {
	fileAttributes *contentPolicyRuleFileAttributesRequestBody
	destinations   []s.Parsable
	sources        *contentPolicyRuleSourcesRequestBody
}

func (b *contentPolicyRuleMatchingConditionsRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteObjectValue("fileAttributes", b.fileAttributes); err != nil {
		return err
	}
	if err := writer.WriteCollectionOfObjectValues("destinations", b.destinations); err != nil {
		return err
	}
	return writer.WriteObjectValue("sources", b.sources)
}
func (b *contentPolicyRuleMatchingConditionsRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type contentPolicyRuleFileAttributesRequestBody struct {
	activities       *string
	contentTypes     []string
	textContentTypes []string
}

func (b *contentPolicyRuleFileAttributesRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("activities", b.activities); err != nil {
		return err
	}
	if len(b.contentTypes) > 0 {
		if err := writer.WriteCollectionOfStringValues("contentTypes", b.contentTypes); err != nil {
			return err
		}
	}
	if len(b.textContentTypes) > 0 {
		if err := writer.WriteCollectionOfStringValues("textContentTypes", b.textContentTypes); err != nil {
			return err
		}
	}
	return nil
}
func (b *contentPolicyRuleFileAttributesRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type contentPolicyRuleDestinationRequestBody struct {
	odataType string
	values    []string
}

func (b *contentPolicyRuleDestinationRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("@odata.type", &b.odataType); err != nil {
		return err
	}
	return writer.WriteCollectionOfStringValues("values", b.values)
}
func (b *contentPolicyRuleDestinationRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type contentPolicyRuleSourcesRequestBody struct{ sessionType *string }

func (b *contentPolicyRuleSourcesRequestBody) Serialize(writer s.SerializationWriter) error {
	return writer.WriteStringValue("sessionType", b.sessionType)
}
func (b *contentPolicyRuleSourcesRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}
