package graphBetaCrossTenantAccessPartnerSettings

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var inheritablePartnerSettings = []struct {
	terraformName string
	graphName     string
	value         func(*CrossTenantAccessPartnerSettingsResourceModel) types.Object
}{
	{
		"b2b_collaboration_inbound",
		"b2bCollaborationInbound",
		func(m *CrossTenantAccessPartnerSettingsResourceModel) types.Object { return m.B2bCollaborationInbound },
	},
	{
		"b2b_collaboration_outbound",
		"b2bCollaborationOutbound",
		func(m *CrossTenantAccessPartnerSettingsResourceModel) types.Object { return m.B2bCollaborationOutbound },
	},
	{
		"b2b_direct_connect_inbound",
		"b2bDirectConnectInbound",
		func(m *CrossTenantAccessPartnerSettingsResourceModel) types.Object { return m.B2bDirectConnectInbound },
	},
	{
		"b2b_direct_connect_outbound",
		"b2bDirectConnectOutbound",
		func(m *CrossTenantAccessPartnerSettingsResourceModel) types.Object { return m.B2bDirectConnectOutbound },
	},
	{
		"inbound_trust",
		"inboundTrust",
		func(m *CrossTenantAccessPartnerSettingsResourceModel) types.Object { return m.InboundTrust },
	},
	{
		"automatic_user_consent_settings",
		"automaticUserConsentSettings",
		func(m *CrossTenantAccessPartnerSettingsResourceModel) types.Object {
			return m.AutomaticUserConsentSettings
		},
	},
	{
		"tenant_restrictions",
		"tenantRestrictions",
		func(m *CrossTenantAccessPartnerSettingsResourceModel) types.Object { return m.TenantRestrictions },
	},
}
