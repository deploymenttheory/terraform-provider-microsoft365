package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

const (
	// Microsoft Learn documents policyLink as an abstract type with these Graph beta
	// derived types:
	// https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-policylink?view=graph-rest-beta
	//
	// webFilteringPolicyLink is intentionally included even though it is not in the
	// current Microsoft Graph beta Go SDK policyLink discriminator. Entra admin center
	// XHR was observed linking V2 web filtering policies by POSTing this explicit
	// @odata.type to /networkAccess/filteringProfiles/{filteringProfileId}/policies.
	//
	// prompt_policy, content_policy, and netskope_dlp_policy are portal-first cases:
	// the Entra admin center JavaScript bundle defines promptPolicyLink,
	// filePolicyLink, contentPolicyLink, and securityProviderPolicyLink, but these
	// types were not present in Microsoft Graph beta $metadata or the generated Go
	// SDK when this resource was added. Observed Entra admin center XHR for the UI's
	// content policy option sends filePolicyLink/filePolicy.
	//
	// Netskope DLP policy selection was observed querying
	// /networkAccess/securityProviderPolicies with a Netskope schema filter, so
	// Terraform exposes that UI concept as netskope_dlp_policy while keeping the
	// generic securityProviderPolicy Graph type internal.
	//
	// Keep these values as strings instead of SDK constructors so portal-first policy
	// link types can be sent before the generated SDK catches up.
	filteringPolicyLinkODataType          = "#microsoft.graph.networkaccess.filteringPolicyLink"
	filteringPolicyODataType              = "#microsoft.graph.networkaccess.filteringPolicy"
	webFilteringPolicyLinkODataType       = "#microsoft.graph.networkaccess.webFilteringPolicyLink"
	webFilteringPolicyODataType           = "#microsoft.graph.networkaccess.webFilteringPolicy"
	promptPolicyLinkODataType             = "#microsoft.graph.networkaccess.promptPolicyLink"
	promptPolicyODataType                 = "#microsoft.graph.networkaccess.promptPolicy"
	filePolicyLinkODataType               = "#microsoft.graph.networkaccess.filePolicyLink"
	filePolicyODataType                   = "#microsoft.graph.networkaccess.filePolicy"
	securityProviderPolicyLinkODataType   = "#microsoft.graph.networkaccess.securityProviderPolicyLink"
	securityProviderPolicyODataType       = "#microsoft.graph.networkaccess.securityProviderPolicy"
	cloudFirewallPolicyLinkODataType      = "#microsoft.graph.networkaccess.cloudFirewallPolicyLink"
	cloudFirewallPolicyODataType          = "#microsoft.graph.networkaccess.cloudFirewallPolicy"
	threatIntelligencePolicyLinkODataType = "#microsoft.graph.networkaccess.threatIntelligencePolicyLink"
	threatIntelligencePolicyODataType     = "#microsoft.graph.networkaccess.threatIntelligencePolicy"
	tlsInspectionPolicyLinkODataType      = "#microsoft.graph.networkaccess.tlsInspectionPolicyLink"
	tlsInspectionPolicyODataType          = "#microsoft.graph.networkaccess.tlsInspectionPolicy"
)

type policyODataTypes struct {
	link   string
	policy string
}

func resolvePolicyODataTypes(data *NetworkFilteringProfilePolicyLinkResourceModel) (policyODataTypes, error) {
	if odataTypes, ok := policyTypeToODataTypes[data.PolicyType.ValueString()]; ok {
		return odataTypes, nil
	}

	return policyODataTypes{}, fmt.Errorf("unsupported policy_type %q", data.PolicyType.ValueString())
}

var policyTypeToODataTypes = map[string]policyODataTypes{
	policyTypeFiltering: {
		link:   filteringPolicyLinkODataType,
		policy: filteringPolicyODataType,
	},
	policyTypeWebFiltering: {
		link:   webFilteringPolicyLinkODataType,
		policy: webFilteringPolicyODataType,
	},
	policyTypePrompt: {
		link:   promptPolicyLinkODataType,
		policy: promptPolicyODataType,
	},
	policyTypeContent: {
		link:   filePolicyLinkODataType,
		policy: filePolicyODataType,
	},
	policyTypeNetskopeDlp: {
		link:   securityProviderPolicyLinkODataType,
		policy: securityProviderPolicyODataType,
	},
	policyTypeCloudFirewall: {
		link:   cloudFirewallPolicyLinkODataType,
		policy: cloudFirewallPolicyODataType,
	},
	policyTypeThreatIntelligence: {
		link:   threatIntelligencePolicyLinkODataType,
		policy: threatIntelligencePolicyODataType,
	},
	policyTypeTlsInspection: {
		link:   tlsInspectionPolicyLinkODataType,
		policy: tlsInspectionPolicyODataType,
	},
}

func policyTypeFromODataTypes(linkODataType, policyODataType string) (types.String, bool) {
	for policyType, odataTypes := range policyTypeToODataTypes {
		if linkODataType != "" && odataTypes.link == linkODataType {
			return types.StringValue(policyType), true
		}
		if policyODataType != "" && odataTypes.policy == policyODataType {
			return types.StringValue(policyType), true
		}
	}

	return types.StringNull(), false
}

func populateComputedRequestFields(data *NetworkFilteringProfilePolicyLinkResourceModel) error {
	_, err := resolvePolicyODataTypes(data)
	if err != nil {
		return err
	}

	if data.PolicyType.ValueString() == policyTypeFiltering {
		if data.Priority.IsNull() || data.Priority.IsUnknown() {
			data.Priority = types.Int64Value(100)
		}
		if data.LoggingState.IsNull() || data.LoggingState.IsUnknown() {
			data.LoggingState = types.StringValue("enabled")
		}
		return nil
	}

	data.Priority = types.Int64Null()
	data.LoggingState = types.StringNull()
	return nil
}

func constructCreateResource(ctx context.Context, data *NetworkFilteringProfilePolicyLinkResourceModel) (s.Parsable, error) {
	_ = ctx

	odataTypes, err := resolvePolicyODataTypes(data)
	if err != nil {
		return nil, err
	}

	body := &policyLinkRequestBody{
		State:           data.State.ValueString(),
		PolicyID:        data.PolicyID.ValueString(),
		ODataType:       odataTypes.link,
		PolicyODataType: odataTypes.policy,
	}

	switch data.PolicyType.ValueString() {
	case policyTypeFiltering:
		body.Priority = int64Ptr(100)
		if !data.Priority.IsNull() && !data.Priority.IsUnknown() {
			body.Priority = int64Ptr(data.Priority.ValueInt64())
		}
		body.LoggingState = "enabled"
		if !data.LoggingState.IsNull() && !data.LoggingState.IsUnknown() {
			body.LoggingState = data.LoggingState.ValueString()
		}
	}

	return body, nil
}

func constructUpdateResource(ctx context.Context, data *NetworkFilteringProfilePolicyLinkResourceModel) (s.Parsable, error) {
	_ = ctx

	odataTypes, err := resolvePolicyODataTypes(data)
	if err != nil {
		return nil, err
	}

	body := &policyLinkRequestBody{
		State:     data.State.ValueString(),
		ODataType: odataTypes.link,
	}

	return body, nil
}

type policyLinkRequestBody struct {
	ODataType       string
	State           string
	Priority        *int64
	LoggingState    string
	PolicyID        string
	PolicyODataType string
}

func (b *policyLinkRequestBody) Serialize(writer s.SerializationWriter) error {
	if b.ODataType != "" {
		if err := writer.WriteStringValue("@odata.type", &b.ODataType); err != nil {
			return err
		}
	}
	if b.Priority != nil {
		if err := writer.WriteInt64Value("priority", b.Priority); err != nil {
			return err
		}
	}
	if b.State != "" {
		if err := writer.WriteStringValue("state", &b.State); err != nil {
			return err
		}
	}
	if b.LoggingState != "" {
		if err := writer.WriteStringValue("loggingState", &b.LoggingState); err != nil {
			return err
		}
	}
	if b.PolicyID != "" {
		policy := &policyReferenceRequestBody{
			ID:        b.PolicyID,
			ODataType: b.PolicyODataType,
		}
		if err := writer.WriteObjectValue("policy", policy); err != nil {
			return err
		}
	}
	return nil
}

func (b *policyLinkRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

func (b *policyLinkRequestBody) GetAdditionalData() map[string]any {
	return nil
}

type policyReferenceRequestBody struct {
	ID        string
	ODataType string
}

func (b *policyReferenceRequestBody) Serialize(writer s.SerializationWriter) error {
	if b.ID != "" {
		if err := writer.WriteStringValue("id", &b.ID); err != nil {
			return err
		}
	}
	if b.ODataType != "" {
		if err := writer.WriteStringValue("@odata.type", &b.ODataType); err != nil {
			return err
		}
	}
	return nil
}

func (b *policyReferenceRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

func (b *policyReferenceRequestBody) GetAdditionalData() map[string]any {
	return nil
}

func int64Ptr(v int64) *int64 {
	return &v
}
