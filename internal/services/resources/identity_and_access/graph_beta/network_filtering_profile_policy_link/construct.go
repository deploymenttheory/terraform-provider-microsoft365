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
	// Keep these values as strings instead of SDK constructors so portal-first policy
	// link types can be sent before the generated SDK catches up.
	filteringPolicyLinkODataType          = "#microsoft.graph.networkaccess.filteringPolicyLink"
	filteringPolicyODataType              = "#microsoft.graph.networkaccess.filteringPolicy"
	webFilteringPolicyLinkODataType       = "#microsoft.graph.networkaccess.webFilteringPolicyLink"
	webFilteringPolicyODataType           = "#microsoft.graph.networkaccess.webFilteringPolicy"
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
	switch data.PolicyType.ValueString() {
	case policyTypeFiltering:
		return policyODataTypes{link: filteringPolicyLinkODataType, policy: filteringPolicyODataType}, nil
	case policyTypeWebFiltering:
		return policyODataTypes{link: webFilteringPolicyLinkODataType, policy: webFilteringPolicyODataType}, nil
	case policyTypeCloudFirewall:
		return policyODataTypes{link: cloudFirewallPolicyLinkODataType, policy: cloudFirewallPolicyODataType}, nil
	case policyTypeThreatIntelligence:
		return policyODataTypes{link: threatIntelligencePolicyLinkODataType, policy: threatIntelligencePolicyODataType}, nil
	case policyTypeTlsInspection:
		return policyODataTypes{link: tlsInspectionPolicyLinkODataType, policy: tlsInspectionPolicyODataType}, nil
	case policyTypeCustom:
		// Custom exists for portal-first Global Secure Access policy link types that
		// are visible in Entra admin center XHR before they are published in Learn or
		// generated into the Microsoft Graph beta SDK.
		if data.PolicyLinkODataType.IsNull() || data.PolicyLinkODataType.IsUnknown() || data.PolicyLinkODataType.ValueString() == "" {
			return policyODataTypes{}, fmt.Errorf("policy_link_odata_type is required when policy_type is custom")
		}
		if data.PolicyODataType.IsNull() || data.PolicyODataType.IsUnknown() || data.PolicyODataType.ValueString() == "" {
			return policyODataTypes{}, fmt.Errorf("policy_odata_type is required when policy_type is custom")
		}
		return policyODataTypes{link: data.PolicyLinkODataType.ValueString(), policy: data.PolicyODataType.ValueString()}, nil
	default:
		return policyODataTypes{}, fmt.Errorf("unsupported policy_type %q", data.PolicyType.ValueString())
	}
}

func populateComputedRequestFields(data *NetworkFilteringProfilePolicyLinkResourceModel) error {
	odataTypes, err := resolvePolicyODataTypes(data)
	if err != nil {
		return err
	}

	data.PolicyLinkODataType = types.StringValue(odataTypes.link)
	data.PolicyODataType = types.StringValue(odataTypes.policy)

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
