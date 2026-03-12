package graphBetaCrossTenantAccessPartnerUserSyncSettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CrossTenantAccessPartnerUserSyncSettingsResourceModel represents the schema for the
// Cross-Tenant Access Partner User Sync Settings resource.
// This resource manages identity synchronization configuration for a specific partner tenant.
type CrossTenantAccessPartnerUserSyncSettingsResourceModel struct {
	ID              types.String                    `tfsdk:"id"`
	TenantID        types.String                    `tfsdk:"tenant_id"`
	DisplayName     types.String                    `tfsdk:"display_name"`
	UserSyncInbound *CrossTenantUserSyncInbound     `tfsdk:"user_sync_inbound"`
	Timeouts        timeouts.Value                  `tfsdk:"timeouts"`
}

// CrossTenantUserSyncInbound represents the inbound user synchronization settings
type CrossTenantUserSyncInbound struct {
	IsSyncAllowed types.Bool `tfsdk:"is_sync_allowed"`
}
