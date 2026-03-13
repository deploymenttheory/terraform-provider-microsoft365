package graphBetaCrossTenantAccessPartnerSettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CrossTenantAccessPartnerSettingsResourceModel struct {
	ID                           types.String   `tfsdk:"id"`
	TenantID                     types.String   `tfsdk:"tenant_id"`
	IsServiceProvider            types.Bool     `tfsdk:"is_service_provider"`
	IsInMultiTenantOrganization  types.Bool     `tfsdk:"is_in_multi_tenant_organization"`
	B2bCollaborationInbound      types.Object   `tfsdk:"b2b_collaboration_inbound"`
	B2bCollaborationOutbound     types.Object   `tfsdk:"b2b_collaboration_outbound"`
	B2bDirectConnectInbound      types.Object   `tfsdk:"b2b_direct_connect_inbound"`
	B2bDirectConnectOutbound     types.Object   `tfsdk:"b2b_direct_connect_outbound"`
	InboundTrust                 types.Object   `tfsdk:"inbound_trust"`
	AutomaticUserConsentSettings types.Object   `tfsdk:"automatic_user_consent_settings"`
	TenantRestrictions           types.Object   `tfsdk:"tenant_restrictions"`
	HardDelete                   types.Bool     `tfsdk:"hard_delete"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
}

type CrossTenantAccessPolicyB2BSetting struct {
	UsersAndGroups types.Object `tfsdk:"users_and_groups"`
	Applications   types.Object `tfsdk:"applications"`
}

type CrossTenantAccessPolicyTargetConfiguration struct {
	AccessType types.String `tfsdk:"access_type"`
	Targets    types.Set    `tfsdk:"targets"`
}

type CrossTenantAccessPolicyTarget struct {
	Target     types.String `tfsdk:"target"`
	TargetType types.String `tfsdk:"target_type"`
}

type CrossTenantAccessPolicyInboundTrust struct {
	IsMfaAccepted                       types.Bool `tfsdk:"is_mfa_accepted"`
	IsCompliantDeviceAccepted           types.Bool `tfsdk:"is_compliant_device_accepted"`
	IsHybridAzureADJoinedDeviceAccepted types.Bool `tfsdk:"is_hybrid_azure_ad_joined_device_accepted"`
}

type InboundOutboundPolicyConfiguration struct {
	InboundAllowed  types.Bool `tfsdk:"inbound_allowed"`
	OutboundAllowed types.Bool `tfsdk:"outbound_allowed"`
}

type CrossTenantAccessPolicyTenantRestrictions struct {
	UsersAndGroups types.Object `tfsdk:"users_and_groups"`
	Applications   types.Object `tfsdk:"applications"`
}
