package graphBetaNetworkPromptPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	s "github.com/microsoft/kiota-abstractions-go/serialization"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
)

// constructResource builds the policy container payload observed from the
// Entra Global Secure Access prompt policy blade.
func constructResource(
	ctx context.Context,
	data *NetworkPromptPolicyResourceModel,
) (s.Parsable, error) {
	requestBody := &promptPolicyRequestBody{
		settings: &promptPolicySettingsRequestBody{
			defaultAction: data.DefaultAction.ValueStringPointer(),
		},
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		requestBody.name = data.Name.ValueStringPointer()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		requestBody.description = data.Description.ValueStringPointer()
		requestBody.descriptionSet = true
	}

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		requestBody,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody, nil
}

// constructUpdateResource sends only changed fields, matching the observed
// portal PATCH shape while still allowing default_action to be updated.
func constructUpdateResource(
	ctx context.Context,
	plan, state *NetworkPromptPolicyResourceModel,
) (*promptPolicyRequestBody, error) {
	requestBody := &promptPolicyRequestBody{}
	if !plan.Name.Equal(state.Name) {
		requestBody.name = plan.Name.ValueStringPointer()
	}
	if !plan.Description.Equal(state.Description) {
		requestBody.description = plan.Description.ValueStringPointer()
		requestBody.descriptionSet = true
	}
	if !plan.DefaultAction.Equal(state.DefaultAction) {
		requestBody.settings = &promptPolicySettingsRequestBody{
			defaultAction: plan.DefaultAction.ValueStringPointer(),
		}
	}

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		requestBody,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}
	return requestBody, nil
}

type promptPolicyRequestBody struct {
	name           *string
	description    *string
	descriptionSet bool
	settings       *promptPolicySettingsRequestBody
}

func (b *promptPolicyRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("name", b.name); err != nil {
		return wrapSerializationError(err)
	}
	if b.descriptionSet {
		if b.description == nil {
			if err := writer.WriteNullValue("description"); err != nil {
				return wrapSerializationError(err)
			}
		} else if err := writer.WriteStringValue("description", b.description); err != nil {
			return wrapSerializationError(err)
		}
	}
	if b.settings != nil {
		if err := writer.WriteObjectValue("settings", b.settings); err != nil {
			return wrapSerializationError(err)
		}
	}
	return nil
}

func (b *promptPolicyRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type promptPolicySettingsRequestBody struct {
	defaultAction *string
}

func (b *promptPolicySettingsRequestBody) Serialize(writer s.SerializationWriter) error {
	return wrapSerializationError(writer.WriteStringValue("defaultAction", b.defaultAction))
}

func (b *promptPolicySettingsRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

// hasChanges excludes Terraform-only changes such as timeouts from empty Graph PATCH requests.
func (b *promptPolicyRequestBody) hasChanges() bool {
	return b.name != nil || b.descriptionSet || b.settings != nil
}

func wrapSerializationError(err error) error {
	if err != nil {
		return fmt.Errorf("serialize prompt policy: %w", err)
	}
	return nil
}
