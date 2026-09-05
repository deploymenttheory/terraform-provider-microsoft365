package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Explicit discriminators are required on PATCH: Graph can otherwise silently ignore fields.
func constructResource(
	ctx context.Context,
	data *NetworkMCPPolicyRuleResourceModel,
) (s.Parsable, error) {
	b := &mcpPolicyRuleRequestBody{
		name:           data.Name.ValueStringPointer(),
		description:    data.Description.ValueStringPointer(),
		descriptionSet: !data.Description.IsNull(),
		action:         data.Action.ValueStringPointer(),
		priority:       data.Priority.ValueInt32Pointer(),
		settings:       ruleSettings(data.Enabled.ValueBool()),
	}
	conditions, err := ruleConditions(ctx, data)
	b.conditions = conditions
	b.conditionsSet = conditions != nil
	return b, err
}

func constructUpdateResource(
	ctx context.Context,
	plan, state *NetworkMCPPolicyRuleResourceModel,
) (*mcpPolicyRuleRequestBody, error) {
	b := &mcpPolicyRuleRequestBody{}
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
	if !plan.MatchingConditions.Equal(state.MatchingConditions) {
		conditions, err := ruleConditions(ctx, plan)
		if err != nil {
			return nil, err
		}
		b.conditions = conditions
		b.conditionsSet = true
	}
	return b, nil
}

func ruleSettings(enabled bool) *mcpPolicyRuleSettings {
	status := "disabled"
	if enabled {
		status = "enabled"
	}
	return &mcpPolicyRuleSettings{status: &status}
}

type mcpPolicyRuleRequestBody struct {
	name, description, action *string
	descriptionSet            bool
	priority                  *int32
	settings                  *mcpPolicyRuleSettings
	conditions                *conditionPayload
	conditionsSet             bool
}

func (b *mcpPolicyRuleRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

func (b *mcpPolicyRuleRequestBody) Serialize(w s.SerializationWriter) error {
	discriminator := "#microsoft.graph.networkaccess.mcpPolicyRule"
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
	if err := w.WriteInt32Value("priority", b.priority); err != nil {
		return wrapSerializationError(err)
	}
	if b.settings != nil {
		if err := w.WriteObjectValue("settings", b.settings); err != nil {
			return wrapSerializationError(err)
		}
	}
	if b.conditionsSet {
		if b.conditions == nil {
			return wrapSerializationError(w.WriteNullValue("matchingConditions"))
		}
		if err := w.WriteObjectValue("matchingConditions", b.conditions); err != nil {
			return wrapSerializationError(err)
		}
	}
	return nil
}

type mcpPolicyRuleSettings struct{ status *string }

func (b *mcpPolicyRuleSettings) Serialize(w s.SerializationWriter) error {
	return wrapSerializationError(w.WriteStringValue("status", b.status))
}

func (b *mcpPolicyRuleSettings) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

// hasChanges excludes Terraform-only changes such as timeouts from empty Graph PATCH requests.
func (b *mcpPolicyRuleRequestBody) hasChanges() bool {
	return b.name != nil || b.descriptionSet || b.settings != nil || b.action != nil ||
		b.priority != nil ||
		b.conditionsSet
}

func wrapSerializationError(err error) error {
	if err != nil {
		return fmt.Errorf("serialize MCP policy rule: %w", err)
	}
	return nil
}
