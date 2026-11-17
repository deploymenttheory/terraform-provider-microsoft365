// REF: https://learn.microsoft.com/en-us/graph/api/resources/sharepointsettings?view=graph-rest-beta
package graphBetaM365AdminSharepointSettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SharepointSettingsResourceModel represents the schema for the SharepointSettings resource.
// IdleSessionSignOut is a types.Object because it is Optional+Computed: when it is not configured it is unknown
// in the plan and is filled from the current tenant settings.
type SharepointSettingsResourceModel struct {
	ID                                              types.String   `tfsdk:"id"`
	AllowedDomainGuidsForSyncApp                    types.Set      `tfsdk:"allowed_domain_guids_for_sync_app"`
	AvailableManagedPathsForSiteCreation            types.Set      `tfsdk:"available_managed_paths_for_site_creation"`
	DeletedUserPersonalSiteRetentionPeriodInDays    types.Int32    `tfsdk:"deleted_user_personal_site_retention_period_in_days"`
	ExcludedFileExtensionsForSyncApp                types.Set      `tfsdk:"excluded_file_extensions_for_sync_app"`
	IdleSessionSignOut                              types.Object   `tfsdk:"idle_session_sign_out"`
	ImageTaggingOption                              types.String   `tfsdk:"image_tagging_option"`
	IsCommentingOnSitePagesEnabled                  types.Bool     `tfsdk:"is_commenting_on_site_pages_enabled"`
	IsFileActivityNotificationEnabled               types.Bool     `tfsdk:"is_file_activity_notification_enabled"`
	IsLegacyAuthProtocolsEnabled                    types.Bool     `tfsdk:"is_legacy_auth_protocols_enabled"`
	IsLoopEnabled                                   types.Bool     `tfsdk:"is_loop_enabled"`
	IsMacSyncAppEnabled                             types.Bool     `tfsdk:"is_mac_sync_app_enabled"`
	IsRequireAcceptingUserToMatchInvitedUserEnabled types.Bool     `tfsdk:"is_require_accepting_user_to_match_invited_user_enabled"`
	IsResharingByExternalUsersEnabled               types.Bool     `tfsdk:"is_resharing_by_external_users_enabled"`
	IsSharePointMobileNotificationEnabled           types.Bool     `tfsdk:"is_sharepoint_mobile_notification_enabled"`
	IsSharePointNewsfeedEnabled                     types.Bool     `tfsdk:"is_sharepoint_newsfeed_enabled"`
	IsSiteCreationEnabled                           types.Bool     `tfsdk:"is_site_creation_enabled"`
	IsSiteCreationUIEnabled                         types.Bool     `tfsdk:"is_site_creation_ui_enabled"`
	IsSitePagesCreationEnabled                      types.Bool     `tfsdk:"is_site_pages_creation_enabled"`
	IsSitesStorageLimitAutomatic                    types.Bool     `tfsdk:"is_sites_storage_limit_automatic"`
	IsSyncButtonHiddenOnPersonalSite                types.Bool     `tfsdk:"is_sync_button_hidden_on_personal_site"`
	IsUnmanagedSyncAppForTenantRestricted           types.Bool     `tfsdk:"is_unmanaged_sync_app_for_tenant_restricted"`
	PersonalSiteDefaultStorageLimitInMB             types.Int64    `tfsdk:"personal_site_default_storage_limit_in_mb"`
	SharingAllowedDomainList                        types.Set      `tfsdk:"sharing_allowed_domain_list"`
	SharingBlockedDomainList                        types.Set      `tfsdk:"sharing_blocked_domain_list"`
	SharingCapability                               types.String   `tfsdk:"sharing_capability"`
	SharingDomainRestrictionMode                    types.String   `tfsdk:"sharing_domain_restriction_mode"`
	SiteCreationDefaultManagedPath                  types.String   `tfsdk:"site_creation_default_managed_path"`
	SiteCreationDefaultStorageLimitInMB             types.Int32    `tfsdk:"site_creation_default_storage_limit_in_mb"`
	TenantDefaultTimezone                           types.String   `tfsdk:"tenant_default_timezone"`
	Timeouts                                        timeouts.Value `tfsdk:"timeouts"`
}

// IdleSessionSignOutModel represents the idle session sign-out policy for SharePoint and OneDrive.
type IdleSessionSignOutModel struct {
	IsEnabled             types.Bool  `tfsdk:"is_enabled"`
	SignOutAfterInSeconds types.Int64 `tfsdk:"sign_out_after_in_seconds"`
	WarnAfterInSeconds    types.Int64 `tfsdk:"warn_after_in_seconds"`
}

var idleSessionSignOutAttrTypes = map[string]attr.Type{
	"is_enabled":                types.BoolType,
	"sign_out_after_in_seconds": types.Int64Type,
	"warn_after_in_seconds":     types.Int64Type,
}
