// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-wip-intunebrandingprofile?view=graph-rest-beta
package graphBetaDeviceManagementIntuneBrandingProfile

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IntuneBrandingProfileResourceModel represents the Terraform resource model for an Intune branding profile.
type IntuneBrandingProfileResourceModel struct {
	ID                                        types.String                               `tfsdk:"id"`
	ProfileName                               types.String                               `tfsdk:"profile_name"`
	ProfileDescription                        types.String                               `tfsdk:"profile_description"`
	DisplayName                               types.String                               `tfsdk:"display_name"`
	ThemeColor                                *RgbColorResourceModel                     `tfsdk:"theme_color"`
	ShowLogo                                  types.Bool                                 `tfsdk:"show_logo"`
	ShowDisplayNameNextToLogo                 types.Bool                                 `tfsdk:"show_display_name_next_to_logo"`
	ThemeColorLogo                            *sharedmodels.ImageResourceModel           `tfsdk:"theme_color_logo"`
	LightBackgroundLogo                       *sharedmodels.ImageResourceModel           `tfsdk:"light_background_logo"`
	LandingPageCustomizedImage                *sharedmodels.ImageResourceModel           `tfsdk:"landing_page_customized_image"`
	ContactITName                             types.String                               `tfsdk:"contact_it_name"`
	ContactITPhoneNumber                      types.String                               `tfsdk:"contact_it_phone_number"`
	ContactITEmailAddress                     types.String                               `tfsdk:"contact_it_email_address"`
	ContactITNotes                            types.String                               `tfsdk:"contact_it_notes"`
	OnlineSupportSiteUrl                      types.String                               `tfsdk:"online_support_site_url"`
	OnlineSupportSiteName                     types.String                               `tfsdk:"online_support_site_name"`
	PrivacyUrl                                types.String                               `tfsdk:"privacy_url"`
	CustomPrivacyMessage                      types.String                               `tfsdk:"custom_privacy_message"`
	CustomCanSeePrivacyMessage                types.String                               `tfsdk:"custom_can_see_privacy_message"`
	CustomCantSeePrivacyMessage               types.String                               `tfsdk:"custom_cant_see_privacy_message"`
	IsRemoveDeviceDisabled                    types.Bool                                 `tfsdk:"is_remove_device_disabled"`
	IsFactoryResetDisabled                    types.Bool                                 `tfsdk:"is_factory_reset_disabled"`
	CompanyPortalBlockedActions               []*CompanyPortalBlockedActionResourceModel `tfsdk:"company_portal_blocked_actions"`
	ShowAzureADEnterpriseApps                 types.Bool                                 `tfsdk:"show_azure_ad_enterprise_apps"`
	ShowOfficeWebApps                         types.Bool                                 `tfsdk:"show_office_web_apps"`
	ShowConfigurationManagerApps              types.Bool                                 `tfsdk:"show_configuration_manager_apps"`
	DisableDeviceCategorySelection            types.Bool                                 `tfsdk:"disable_device_category_selection"`
	SendDeviceOwnershipChangePushNotification types.Bool                                 `tfsdk:"send_device_ownership_change_push_notification"`
	EnrollmentAvailability                    types.String                               `tfsdk:"enrollment_availability"`
	DisableClientTelemetry                    types.Bool                                 `tfsdk:"disable_client_telemetry"`
	IsDefaultProfile                          types.Bool                                 `tfsdk:"is_default_profile"`
	CreatedDateTime                           types.String                               `tfsdk:"created_date_time"`
	LastModifiedDateTime                      types.String                               `tfsdk:"last_modified_date_time"`
	RoleScopeTagIds                           types.Set                                  `tfsdk:"role_scope_tag_ids"`
	Assignments                               types.Set                                  `tfsdk:"assignments"`
	Timeouts                                  timeouts.Value                             `tfsdk:"timeouts"`
}

// RgbColorResourceModel represents an RGB color.
type RgbColorResourceModel struct {
	R types.Int32 `tfsdk:"r"`
	G types.Int32 `tfsdk:"g"`
	B types.Int32 `tfsdk:"b"`
}

// CompanyPortalBlockedActionResourceModel represents a blocked action for the Company Portal.
type CompanyPortalBlockedActionResourceModel struct {
	Platform  types.String `tfsdk:"platform"`
	OwnerType types.String `tfsdk:"owner_type"`
	Action    types.String `tfsdk:"action"`
}

// IntuneBrandingProfileAssignmentModel defines the schema for a Intune Branding Profile assignment.
type IntuneBrandingProfileAssignmentModel struct {
	Type    types.String `tfsdk:"type"`     // "groupAssignmentTarget", "exclusionGroupAssignmentTarget"
	GroupId types.String `tfsdk:"group_id"` // For group targets (both include and exclude)
}
