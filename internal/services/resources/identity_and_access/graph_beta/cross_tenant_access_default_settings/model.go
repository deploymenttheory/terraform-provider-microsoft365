// REF: https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationdefault?view=graph-rest-beta
package graphBetaCrossTenantAccessDefaultSettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// tenantRestrictionsB2BSetting implements CrossTenantAccessPolicyTenantRestrictionsable using the base
// B2BSetting constructor so that no @odata.type is serialised. The derived constructor
// NewCrossTenantAccessPolicyTenantRestrictions() auto-sets
// @odata.type = "#microsoft.graph.crossTenantAccessPolicyTenantRestrictions" which the default
// settings endpoint rejects.
type tenantRestrictionsB2BSetting struct {
	graphmodels.CrossTenantAccessPolicyB2BSetting
}

func (t *tenantRestrictionsB2BSetting) GetDevices() graphmodels.DevicesFilterable { return nil }
func (t *tenantRestrictionsB2BSetting) SetDevices(value graphmodels.DevicesFilterable)() {}

// CrossTenantAccessDefaultSettingsResourceModel represents the schema for the Cross Tenant Access Default Settings resource.
// This is a singleton resource — one default configuration exists per tenant and cannot be created or deleted via the API.
type CrossTenantAccessDefaultSettingsResourceModel struct {
	ID                                               types.String   `tfsdk:"id"`
	IsServiceDefault                                 types.Bool     `tfsdk:"is_service_default"`
	B2bCollaborationInbound                          types.Object   `tfsdk:"b2b_collaboration_inbound"`
	B2bCollaborationOutbound                         types.Object   `tfsdk:"b2b_collaboration_outbound"`
	B2bDirectConnectInbound                          types.Object   `tfsdk:"b2b_direct_connect_inbound"`
	B2bDirectConnectOutbound                         types.Object   `tfsdk:"b2b_direct_connect_outbound"`
	InboundTrust                                     types.Object   `tfsdk:"inbound_trust"`
	InvitationRedemptionIdentityProviderConfiguration types.Object   `tfsdk:"invitation_redemption_identity_provider_configuration"`
	TenantRestrictions                               types.Object   `tfsdk:"tenant_restrictions"`
	AutomaticUserConsentSettings                     types.Object   `tfsdk:"automatic_user_consent_settings"`
	RestoreDefaultsOnDestroy                         types.Bool     `tfsdk:"restore_defaults_on_destroy"`
	Timeouts                                         timeouts.Value `tfsdk:"timeouts"`
}

// CrossTenantAccessPolicyB2BSetting represents B2B collaboration or direct connect settings
type CrossTenantAccessPolicyB2BSetting struct {
	UsersAndGroups types.Object `tfsdk:"users_and_groups"`
	Applications   types.Object `tfsdk:"applications"`
}

// CrossTenantAccessPolicyTargetConfiguration represents target configuration for users/groups or applications
type CrossTenantAccessPolicyTargetConfiguration struct {
	AccessType types.String `tfsdk:"access_type"`
	Targets    types.Set    `tfsdk:"targets"`
}

// CrossTenantAccessPolicyTarget represents a target (user, group, or application)
type CrossTenantAccessPolicyTarget struct {
	Target     types.String `tfsdk:"target"`
	TargetType types.String `tfsdk:"target_type"`
}

// CrossTenantAccessPolicyInboundTrust represents inbound trust settings
type CrossTenantAccessPolicyInboundTrust struct {
	IsMfaAccepted                      types.Bool `tfsdk:"is_mfa_accepted"`
	IsCompliantDeviceAccepted          types.Bool `tfsdk:"is_compliant_device_accepted"`
	IsHybridAzureADJoinedDeviceAccepted types.Bool `tfsdk:"is_hybrid_azure_ad_joined_device_accepted"`
}

// DefaultInvitationRedemptionIdentityProviderConfiguration represents invitation redemption configuration
type DefaultInvitationRedemptionIdentityProviderConfiguration struct {
	PrimaryIdentityProviderPrecedenceOrder types.List   `tfsdk:"primary_identity_provider_precedence_order"`
	FallbackIdentityProvider               types.String `tfsdk:"fallback_identity_provider"`
}

// InboundOutboundPolicyConfiguration represents automatic user consent settings
type InboundOutboundPolicyConfiguration struct {
	InboundAllowed  types.Bool `tfsdk:"inbound_allowed"`
	OutboundAllowed types.Bool `tfsdk:"outbound_allowed"`
}
