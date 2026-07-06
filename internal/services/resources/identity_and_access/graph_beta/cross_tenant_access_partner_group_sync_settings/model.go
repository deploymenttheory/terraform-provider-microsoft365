package graphBetaCrossTenantAccessPartnerGroupSyncSettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CrossTenantAccessPartnerGroupSyncSettingsResourceModel represents the schema for the
// Cross-Tenant Access Partner Group Sync Settings resource.
// This resource manages group synchronization configuration for a specific partner tenant.
type CrossTenantAccessPartnerGroupSyncSettingsResourceModel struct {
	ID               types.String                 `tfsdk:"id"`
	TenantID         types.String                 `tfsdk:"tenant_id"`
	DisplayName      types.String                 `tfsdk:"display_name"`
	GroupSyncInbound *CrossTenantGroupSyncInbound `tfsdk:"group_sync_inbound"`
	Timeouts         timeouts.Value               `tfsdk:"timeouts"`
}

// CrossTenantGroupSyncInbound represents the inbound group synchronization settings
type CrossTenantGroupSyncInbound struct {
	IsSyncAllowed types.Bool `tfsdk:"is_sync_allowed"`
}
