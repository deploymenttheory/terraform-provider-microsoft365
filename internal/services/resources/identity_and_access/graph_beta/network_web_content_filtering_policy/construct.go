package graphBetaNetworkWebContentFilteringPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// constructResource builds the portal-observed web content filtering policy
// payload for /networkaccess/webFilteringPolicies.
//
// Microsoft Graph beta currently documents the older/generic
// networkaccess.filteringPolicy shape, not this webFilteringPolicy surface:
// https://learn.microsoft.com/graph/api/resources/networkaccess-filteringpolicy
//
// The Entra Global Secure Access Web content filtering blade sends
// settings.defaultAction.@odata.type plus an empty policyRules array on create.
// policyRules is intentionally create-only here because subsequent rule
// management happens through the child /policyRules resource.
func constructResource(ctx context.Context, data *NetworkWebContentFilteringPolicyResourceModel, includePolicyRules bool) (s.Parsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := &webContentFilteringPolicyRequestBody{
		settings: &webContentFilteringPolicySettingsRequestBody{
			defaultAction: &webContentFilteringPolicyDefaultActionRequestBody{
				odataType: graphDefaultActionODataType(data.DefaultAction.ValueString()),
			},
		},
		includePolicyRules: includePolicyRules,
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		requestBody.name = data.Name.ValueStringPointer()
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		requestBody.description = data.Description.ValueStringPointer()
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

type webContentFilteringPolicyRequestBody struct {
	name               *string
	description        *string
	settings           *webContentFilteringPolicySettingsRequestBody
	includePolicyRules bool
}

func (b *webContentFilteringPolicyRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("name", b.name); err != nil {
		return err
	}
	if err := writer.WriteStringValue("description", b.description); err != nil {
		return err
	}
	if err := writer.WriteObjectValue("settings", b.settings); err != nil {
		return err
	}
	if b.includePolicyRules {
		if err := writer.WriteCollectionOfObjectValues("policyRules", []s.Parsable{}); err != nil {
			return err
		}
	}

	return nil
}

func (b *webContentFilteringPolicyRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type webContentFilteringPolicySettingsRequestBody struct {
	defaultAction *webContentFilteringPolicyDefaultActionRequestBody
}

func (b *webContentFilteringPolicySettingsRequestBody) Serialize(writer s.SerializationWriter) error {
	return writer.WriteObjectValue("defaultAction", b.defaultAction)
}

func (b *webContentFilteringPolicySettingsRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type webContentFilteringPolicyDefaultActionRequestBody struct {
	odataType string
}

func (b *webContentFilteringPolicyDefaultActionRequestBody) Serialize(writer s.SerializationWriter) error {
	return writer.WriteStringValue("@odata.type", &b.odataType)
}

func (b *webContentFilteringPolicyDefaultActionRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}
