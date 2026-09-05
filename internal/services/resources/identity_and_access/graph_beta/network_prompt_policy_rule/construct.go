package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Explicit discriminators are required on PATCH: Graph can otherwise silently ignore fields.
func constructResource(
	ctx context.Context,
	data *NetworkPromptPolicyRuleResourceModel,
) (s.Parsable, error) {
	b := &promptPolicyRuleRequestBody{
		name:           data.Name.ValueStringPointer(),
		description:    data.Description.ValueStringPointer(),
		descriptionSet: !data.Description.IsNull(),
		action:         data.Action.ValueStringPointer(),
		promptLogging:  data.PromptLogging.ValueStringPointer(),
		priority:       data.Priority.ValueInt32Pointer(),
		settings:       ruleSettings(data.Enabled.ValueBool()),
	}
	conditions, err := ruleConditions(ctx, data)
	b.conditions = conditions
	return b, err
}

func constructUpdateResource(
	ctx context.Context,
	plan, state *NetworkPromptPolicyRuleResourceModel,
) (*promptPolicyRuleRequestBody, error) {
	b := &promptPolicyRuleRequestBody{}
	if !plan.Name.Equal(state.Name) {
		b.name = plan.Name.ValueStringPointer()
	}
	if !plan.Description.Equal(state.Description) {
		b.description = plan.Description.ValueStringPointer()
		b.descriptionSet = true
	}
	if !plan.Action.Equal(state.Action) {
		b.action = plan.Action.ValueStringPointer()
	}
	if !plan.Priority.Equal(state.Priority) {
		b.priority = plan.Priority.ValueInt32Pointer()
	}
	if !plan.Enabled.Equal(state.Enabled) {
		b.settings = ruleSettings(plan.Enabled.ValueBool())
	}
	if !plan.PromptLogging.Equal(state.PromptLogging) {
		b.promptLogging = plan.PromptLogging.ValueStringPointer()
	}
	if !plan.ConversationSchemes.Equal(state.ConversationSchemes) ||
		!plan.ScanResult.Equal(state.ScanResult) {
		conditions, err := ruleConditions(ctx, plan)
		if err != nil {
			return nil, err
		}
		b.conditions = conditions
	}
	return b, nil
}

func ruleSettings(enabled bool) *promptPolicyRuleSettings {
	status := "disabled"
	if enabled {
		status = "enabled"
	}
	return &promptPolicyRuleSettings{status: &status}
}

func ruleConditions(
	ctx context.Context,
	data *NetworkPromptPolicyRuleResourceModel,
) (*promptPolicyRuleConditions, error) {
	var schemes []ConversationSchemeModel
	if diags := data.ConversationSchemes.ElementsAs(ctx, &schemes, false); diags.HasError() {
		return nil, fmt.Errorf("%w: %v", errInvalidSchemes, diags)
	}
	b := &promptPolicyRuleConditions{
		scanResult: data.ScanResult.ValueStringPointer(),
		schemes:    make([]s.Parsable, 0, len(schemes)),
	}
	for _, scheme := range schemes {
		item := &conversationSchemeRequest{}
		switch scheme.Type.ValueString() {
		case schemeTypeCustom:
			item.odataType = "#microsoft.graph.networkaccess.customConversationScheme"
			item.url = scheme.URL.ValueStringPointer()
			jsonPath := scheme.JSONPath.ValueString()
			item.jsonPath = &jsonPath
		case schemeTypePredefined:
			item.odataType = "#microsoft.graph.networkaccess.predefinedConversationScheme"
			item.schemeName = scheme.SchemeName.ValueStringPointer()
		default:
			return nil, fmt.Errorf(
				"%w: unsupported type %q",
				errInvalidSchemes,
				scheme.Type.ValueString(),
			)
		}
		b.schemes = append(b.schemes, item)
	}
	return b, nil
}

type promptPolicyRuleRequestBody struct {
	name, description, action, promptLogging *string
	descriptionSet                           bool
	priority                                 *int32
	settings                                 *promptPolicyRuleSettings
	conditions                               *promptPolicyRuleConditions
}

func (b *promptPolicyRuleRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

func (b *promptPolicyRuleRequestBody) Serialize(w s.SerializationWriter) error {
	discriminator := "#microsoft.graph.networkaccess.promptRule"
	if err := w.WriteStringValue("@odata.type", &discriminator); err != nil {
		return wrapSerializationError(err)
	}
	if err := w.WriteStringValue("name", b.name); err != nil {
		return wrapSerializationError(err)
	}
	if b.descriptionSet {
		if b.description == nil {
			if err := w.WriteNullValue("description"); err != nil {
				return wrapSerializationError(err)
			}
		} else if err := w.WriteStringValue("description", b.description); err != nil {
			return wrapSerializationError(err)
		}
	}
	if err := w.WriteStringValue("action", b.action); err != nil {
		return wrapSerializationError(err)
	}
	if err := w.WriteStringValue("promptLogging", b.promptLogging); err != nil {
		return wrapSerializationError(err)
	}
	if err := w.WriteInt32Value("priority", b.priority); err != nil {
		return wrapSerializationError(err)
	}
	if b.settings != nil {
		if err := w.WriteObjectValue("settings", b.settings); err != nil {
			return wrapSerializationError(err)
		}
	}
	if b.conditions != nil {
		if err := w.WriteObjectValue("matchingConditions", b.conditions); err != nil {
			return wrapSerializationError(err)
		}
	}
	return nil
}

type promptPolicyRuleSettings struct{ status *string }

func (b *promptPolicyRuleSettings) Serialize(w s.SerializationWriter) error {
	return wrapSerializationError(w.WriteStringValue("status", b.status))
}

func (b *promptPolicyRuleSettings) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type promptPolicyRuleConditions struct {
	scanResult *string
	schemes    []s.Parsable
}

func (b *promptPolicyRuleConditions) Serialize(w s.SerializationWriter) error {
	if err := w.WriteStringValue("scanResult", b.scanResult); err != nil {
		return wrapSerializationError(err)
	}
	return wrapSerializationError(w.WriteCollectionOfObjectValues("conversationSchemes", b.schemes))
}

func (b *promptPolicyRuleConditions) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type conversationSchemeRequest struct {
	odataType                 string
	url, jsonPath, schemeName *string
}

func (b *conversationSchemeRequest) Serialize(w s.SerializationWriter) error {
	for key, value := range map[string]*string{"@odata.type": &b.odataType, "url": b.url, "jsonPath": b.jsonPath, "schemeName": b.schemeName} {
		if err := w.WriteStringValue(key, value); err != nil {
			return wrapSerializationError(err)
		}
	}
	return nil
}

func (b *conversationSchemeRequest) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

// hasChanges excludes Terraform-only changes such as timeouts from empty Graph PATCH requests.
func (b *promptPolicyRuleRequestBody) hasChanges() bool {
	return b.name != nil || b.descriptionSet || b.settings != nil || b.action != nil ||
		b.priority != nil || b.promptLogging != nil ||
		b.conditions != nil
}

func wrapSerializationError(err error) error {
	if err != nil {
		return fmt.Errorf("serialize prompt policy rule: %w", err)
	}
	return nil
}
