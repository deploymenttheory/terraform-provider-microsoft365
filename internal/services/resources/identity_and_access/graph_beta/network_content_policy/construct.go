package graphBetaNetworkContentPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
)

// constructResource builds the policy container payload observed from the
// Entra Global Secure Access content policy blade.
func constructResource(ctx context.Context, data *NetworkContentPolicyResourceModel) (s.Parsable, error) {
	requestBody := &contentPolicyRequestBody{
		settings: &contentPolicySettingsRequestBody{defaultAction: data.DefaultAction.ValueStringPointer()},
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		requestBody.name = data.Name.ValueStringPointer()
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		requestBody.description = data.Description.ValueStringPointer()
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody, nil
}

// constructUpdateResource sends only changed fields, matching the observed
// portal PATCH shape while still allowing default_action to be updated.
func constructUpdateResource(ctx context.Context, plan, state *NetworkContentPolicyResourceModel) (s.Parsable, error) {
	requestBody := &contentPolicyRequestBody{}
	if !plan.Name.Equal(state.Name) {
		requestBody.name = plan.Name.ValueStringPointer()
	}
	if !plan.Description.Equal(state.Description) {
		requestBody.description = plan.Description.ValueStringPointer()
	}
	if !plan.DefaultAction.Equal(state.DefaultAction) {
		requestBody.settings = &contentPolicySettingsRequestBody{defaultAction: plan.DefaultAction.ValueStringPointer()}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}
	return requestBody, nil
}

type contentPolicyRequestBody struct {
	name        *string
	description *string
	settings    *contentPolicySettingsRequestBody
}

func (b *contentPolicyRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("name", b.name); err != nil {
		return err
	}
	if err := writer.WriteStringValue("description", b.description); err != nil {
		return err
	}
	if b.settings != nil {
		if err := writer.WriteObjectValue("settings", b.settings); err != nil {
			return err
		}
	}
	return nil
}

func (b *contentPolicyRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type contentPolicySettingsRequestBody struct {
	defaultAction *string
}

func (b *contentPolicySettingsRequestBody) Serialize(writer s.SerializationWriter) error {
	return writer.WriteStringValue("defaultAction", b.defaultAction)
}

func (b *contentPolicySettingsRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

const (
	contentPoliciesURLTemplate   = "{+baseurl}/networkaccess/filePolicies"
	contentPolicyItemURLTemplate = contentPoliciesURLTemplate + "/{contentPolicyId}"
)

func newContentPolicyRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	policyID string,
	requestBody s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{"baseurl": adapter.GetBaseUrl()}
	urlTemplate := contentPoliciesURLTemplate
	if policyID != "" {
		urlTemplate = contentPolicyItemURLTemplate
		pathParameters["contentPolicyId"] = policyID
	}
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")
	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set content policy request content: %w", err)
		}
	}
	return requestInfo, nil
}
