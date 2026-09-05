package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Explicit discriminators are required on PATCH: Graph can otherwise silently ignore fields.
func constructResource(
	ctx context.Context,
	data *NetworkTLSInspectionPolicyRuleResourceModel,
) (s.Parsable, error) {
	b := &tlsInspectionPolicyRuleRequestBody{
		name:           data.Name.ValueStringPointer(),
		description:    data.Description.ValueStringPointer(),
		descriptionSet: !data.Description.IsNull(),
		action:         data.Action.ValueStringPointer(),
		priority:       data.Priority.ValueInt32Pointer(),
		settings:       ruleSettings(data.Enabled.ValueBool()),
	}
	conditions, err := ruleConditions(ctx, data)
	b.conditions = conditions
	return b, err
}

func constructUpdateResource(
	ctx context.Context,
	plan, state *NetworkTLSInspectionPolicyRuleResourceModel,
) (*tlsInspectionPolicyRuleRequestBody, error) {
	b := &tlsInspectionPolicyRuleRequestBody{}
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
	if !plan.Destinations.Equal(state.Destinations) {
		conditions, err := ruleConditions(ctx, plan)
		if err != nil {
			return nil, err
		}
		b.conditions = conditions
	}
	return b, nil
}

func ruleSettings(enabled bool) *tlsInspectionPolicyRuleSettings {
	status := "disabled"
	if enabled {
		status = "enabled"
	}
	return &tlsInspectionPolicyRuleSettings{status: &status}
}

func ruleConditions(
	ctx context.Context,
	data *NetworkTLSInspectionPolicyRuleResourceModel,
) (*tlsInspectionPolicyRuleConditions, error) {
	var groups []TLSInspectionPolicyRuleDestinationModel
	if diags := data.Destinations.ElementsAs(ctx, &groups, false); diags.HasError() {
		return nil, fmt.Errorf("%w: decode groups: %v", errInvalidDestinations, diags)
	}
	b := &tlsInspectionPolicyRuleConditions{destinations: make([]s.Parsable, 0, len(groups))}
	for _, group := range groups {
		var discriminator string
		switch group.Type.ValueString() {
		case destinationTypeFQDN:
			discriminator = "#microsoft.graph.networkaccess.tlsInspectionFqdnDestination"
		case destinationTypeWebCategory:
			discriminator = "#microsoft.graph.networkaccess.tlsInspectionWebCategoryDestination"
		default:
			return nil, fmt.Errorf(
				"%w: unsupported type %q",
				errInvalidDestinations,
				group.Type.ValueString(),
			)
		}
		var values []string
		if diags := group.Values.ElementsAs(ctx, &values, false); diags.HasError() {
			return nil, fmt.Errorf("%w: decode values: %v", errInvalidDestinations, diags)
		}
		b.destinations = append(
			b.destinations,
			&tlsInspectionPolicyRuleDestination{odataType: &discriminator, values: values},
		)
	}
	return b, nil
}

type tlsInspectionPolicyRuleRequestBody struct {
	name, description, action *string
	descriptionSet            bool
	priority                  *int32
	settings                  *tlsInspectionPolicyRuleSettings
	conditions                *tlsInspectionPolicyRuleConditions
}

func (b *tlsInspectionPolicyRuleRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

func (b *tlsInspectionPolicyRuleRequestBody) Serialize(w s.SerializationWriter) error {
	discriminator := "#microsoft.graph.networkaccess.tlsInspectionRule"
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
	if b.conditions != nil {
		if err := w.WriteObjectValue("matchingConditions", b.conditions); err != nil {
			return wrapSerializationError(err)
		}
	}
	return nil
}

type tlsInspectionPolicyRuleSettings struct{ status *string }

func (b *tlsInspectionPolicyRuleSettings) Serialize(w s.SerializationWriter) error {
	return wrapSerializationError(w.WriteStringValue("status", b.status))
}

func (b *tlsInspectionPolicyRuleSettings) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type tlsInspectionPolicyRuleConditions struct{ destinations []s.Parsable }

func (b *tlsInspectionPolicyRuleConditions) Serialize(w s.SerializationWriter) error {
	return wrapSerializationError(w.WriteCollectionOfObjectValues("destinations", b.destinations))
}

func (b *tlsInspectionPolicyRuleConditions) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type tlsInspectionPolicyRuleDestination struct {
	odataType *string
	values    []string
}

func (b *tlsInspectionPolicyRuleDestination) Serialize(w s.SerializationWriter) error {
	if err := w.WriteStringValue("@odata.type", b.odataType); err != nil {
		return wrapSerializationError(err)
	}
	return wrapSerializationError(w.WriteCollectionOfStringValues("values", b.values))
}

func (b *tlsInspectionPolicyRuleDestination) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

// hasChanges excludes Terraform-only changes such as timeouts from empty Graph PATCH requests.
func (b *tlsInspectionPolicyRuleRequestBody) hasChanges() bool {
	return b.name != nil || b.descriptionSet || b.settings != nil || b.action != nil ||
		b.priority != nil ||
		b.conditions != nil
}

func wrapSerializationError(err error) error {
	if err != nil {
		return fmt.Errorf("serialize TLS inspection policy rule: %w", err)
	}
	return nil
}
