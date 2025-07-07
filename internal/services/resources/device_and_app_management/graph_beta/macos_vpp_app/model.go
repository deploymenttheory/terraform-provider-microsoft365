// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosvppapp?view=graph-rest-beta

package graphBetaMacOSVppApp

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MacOSVppAppResourceModel represents the root Terraform resource model for macOS VPP applications
type MacOSVppAppResourceModel struct {
	ID                         types.String                             `tfsdk:"id"`
	DisplayName                types.String                             `tfsdk:"display_name"`
	Description                types.String                             `tfsdk:"description"`
	Publisher                  types.String                             `tfsdk:"publisher"`
	AppIcon                    *sharedmodels.MobileAppIconResourceModel `tfsdk:"app_icon"`
	CreatedDateTime            types.String                             `tfsdk:"created_date_time"`
	LastModifiedDateTime       types.String                             `tfsdk:"last_modified_date_time"`
	IsFeatured                 types.Bool                               `tfsdk:"is_featured"`
	PrivacyInformationUrl      types.String                             `tfsdk:"privacy_information_url"`
	InformationUrl             types.String                             `tfsdk:"information_url"`
	Owner                      types.String                             `tfsdk:"owner"`
	Developer                  types.String                             `tfsdk:"developer"`
	Notes                      types.String                             `tfsdk:"notes"`
	UploadState                types.Int32                              `tfsdk:"upload_state"`
	PublishingState            types.String                             `tfsdk:"publishing_state"`
	IsAssigned                 types.Bool                               `tfsdk:"is_assigned"`
	RoleScopeTagIds            types.Set                                `tfsdk:"role_scope_tag_ids"`
	DependentAppCount          types.Int32                              `tfsdk:"dependent_app_count"`
	SupersedingAppCount        types.Int32                              `tfsdk:"superseding_app_count"`
	SupersededAppCount         types.Int32                              `tfsdk:"superseded_app_count"`
	UsedLicenseCount           types.Int32                              `tfsdk:"used_license_count"`
	TotalLicenseCount          types.Int32                              `tfsdk:"total_license_count"`
	ReleaseDateTime            types.String                             `tfsdk:"release_date_time"`
	AppStoreUrl                types.String                             `tfsdk:"app_store_url"`
	LicensingType              *VppLicensingTypeResourceModel           `tfsdk:"licensing_type"`
	VppTokenOrganizationName   types.String                             `tfsdk:"vpp_token_organization_name"`
	VppTokenAccountType        types.String                             `tfsdk:"vpp_token_account_type"`
	VppTokenAppleId            types.String                             `tfsdk:"vpp_token_apple_id"`
	BundleId                   types.String                             `tfsdk:"bundle_id"`
	VppTokenId                 types.String                             `tfsdk:"vpp_token_id"`
	VppTokenDisplayName        types.String                             `tfsdk:"vpp_token_display_name"`
	Categories                 types.Set                                `tfsdk:"categories"`
	Relationships              types.List                               `tfsdk:"relationships"`
	AssignedLicenses           types.List                               `tfsdk:"assigned_licenses"`
	RevokeLicenseActionResults types.List                               `tfsdk:"revoke_license_action_results"`
	Timeouts                   timeouts.Value                           `tfsdk:"timeouts"`
}

// VppLicensingTypeResourceModel represents the Terraform resource model for VPP licensing type
type VppLicensingTypeResourceModel struct {
	SupportUserLicensing    types.Bool `tfsdk:"support_user_licensing"`
	SupportDeviceLicensing  types.Bool `tfsdk:"support_device_licensing"`
	SupportsUserLicensing   types.Bool `tfsdk:"supports_user_licensing"`
	SupportsDeviceLicensing types.Bool `tfsdk:"supports_device_licensing"`
}

// MobileAppRelationshipResourceModel represents the Terraform resource model for a Mobile App Relationship
type MobileAppRelationshipResourceModel struct {
	ID                         types.String `tfsdk:"id"`
	SourceDisplayName          types.String `tfsdk:"source_display_name"`
	SourceDisplayVersion       types.String `tfsdk:"source_display_version"`
	SourceId                   types.String `tfsdk:"source_id"`
	SourcePublisherDisplayName types.String `tfsdk:"source_publisher_display_name"`
	TargetDisplayName          types.String `tfsdk:"target_display_name"`
	TargetDisplayVersion       types.String `tfsdk:"target_display_version"`
	TargetId                   types.String `tfsdk:"target_id"`
	TargetPublisher            types.String `tfsdk:"target_publisher"`
	TargetPublisherDisplayName types.String `tfsdk:"target_publisher_display_name"`
	TargetType                 types.String `tfsdk:"target_type"`
}

// MacOSVppAppAssignedLicenseResourceModel represents the Terraform resource model for assigned licenses
type MacOSVppAppAssignedLicenseResourceModel struct {
	UserId            types.String `tfsdk:"user_id"`
	UserEmailAddress  types.String `tfsdk:"user_email_address"`
	UserName          types.String `tfsdk:"user_name"`
	UserPrincipalName types.String `tfsdk:"user_principal_name"`
}

// MacOSVppAppRevokeLicensesActionResultResourceModel represents the result of a revoke license action
type MacOSVppAppRevokeLicensesActionResultResourceModel struct {
	UserId              types.String `tfsdk:"user_id"`
	ManagedDeviceId     types.String `tfsdk:"managed_device_id"`
	TotalLicensesCount  types.Int32  `tfsdk:"total_licenses_count"`
	FailedLicensesCount types.Int32  `tfsdk:"failed_licenses_count"`
	ActionFailureReason types.String `tfsdk:"action_failure_reason"`
	ActionName          types.String `tfsdk:"action_name"`
	ActionState         types.String `tfsdk:"action_state"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	LastUpdatedDateTime types.String `tfsdk:"last_updated_date_time"`
}
