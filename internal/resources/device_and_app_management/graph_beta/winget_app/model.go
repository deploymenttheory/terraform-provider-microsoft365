// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-wingetapp?view=graph-rest-beta

package graphBetaWinGetApp

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WinGetAppResourceModel represents the Terraform resource model for a WinGetApp
type WinGetAppResourceModel struct {
	ID                            types.String                             `tfsdk:"id"`
	DisplayName                   types.String                             `tfsdk:"display_name"`
	Description                   types.String                             `tfsdk:"description"`
	Publisher                     types.String                             `tfsdk:"publisher"`
	Categories                    types.Set                                `tfsdk:"categories"`
	LargeIcon                     types.Object                             `tfsdk:"large_icon"`
	CreatedDateTime               types.String                             `tfsdk:"created_date_time"`
	LastModifiedDateTime          types.String                             `tfsdk:"last_modified_date_time"`
	IsFeatured                    types.Bool                               `tfsdk:"is_featured"`
	PrivacyInformationUrl         types.String                             `tfsdk:"privacy_information_url"`
	InformationUrl                types.String                             `tfsdk:"information_url"`
	Owner                         types.String                             `tfsdk:"owner"`
	Developer                     types.String                             `tfsdk:"developer"`
	Notes                         types.String                             `tfsdk:"notes"`
	UploadState                   types.Int64                              `tfsdk:"upload_state"`
	PublishingState               types.String                             `tfsdk:"publishing_state"`
	IsAssigned                    types.Bool                               `tfsdk:"is_assigned"`
	RoleScopeTagIds               types.Set                                `tfsdk:"role_scope_tag_ids"`
	DependentAppCount             types.Int64                              `tfsdk:"dependent_app_count"`
	SupersedingAppCount           types.Int64                              `tfsdk:"superseding_app_count"`
	SupersededAppCount            types.Int64                              `tfsdk:"superseded_app_count"`
	ManifestHash                  types.String                             `tfsdk:"manifest_hash"`
	PackageIdentifier             types.String                             `tfsdk:"package_identifier"`
	AutomaticallyGenerateMetadata types.Bool                               `tfsdk:"automatically_generate_metadata"`
	InstallExperience             *WinGetAppInstallExperienceResourceModel `tfsdk:"install_experience"`
	Assignments                   *WinGetAppAssignmentsResourceModel       `tfsdk:"assignments"`
	Timeouts                      timeouts.Value                           `tfsdk:"timeouts"`
}

// WinGetAppInstallExperienceModel represents the install experience structure
type WinGetAppInstallExperienceResourceModel struct {
	RunAsAccount types.String `tfsdk:"run_as_account"`
}

// Base resource REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappassignment?view=graph-rest-beta
// WinGetAppAssignmentsResourceModel represents the assignments structure
// This struct replaces the intent field and uses a field to make the intent implicit.
type WinGetAppAssignmentsResourceModel struct {
	Required  types.Set `tfsdk:"required"`
	Available types.Set `tfsdk:"available"`
	Uninstall types.Set `tfsdk:"uninstall"`
}

type MobileAppAssignmentResourceModel struct {
	Id       types.String                              `tfsdk:"id"`
	Target   AssignmentTargetResourceModel             `tfsdk:"target"`
	Settings *WinGetAppAssignmentSettingsResourceModel `tfsdk:"settings"`
	Source   types.String                              `tfsdk:"source"`
	SourceId types.String                              `tfsdk:"source_id"`
}

// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-alldevicesassignmenttarget?view=graph-rest-beta
type AssignmentTargetResourceModel struct {
	TargetType           types.String `tfsdk:"target_type"` // allDevices, allLicensedUsers, androidFotaDeploymentAssignment, configurationManagerCollectionAssignment, exclusionGroupAssignment, groupAssignment
	AssignmentFilterId   types.String `tfsdk:"assignment_filter_id"`
	AssignmentFilterType types.String `tfsdk:"assignment_filter_type"`
	GroupId              types.String `tfsdk:"group_id"`
	CollectionId         types.String `tfsdk:"collection_id"`
}

// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-wingetappassignmentsettings?view=graph-rest-beta
type WinGetAppAssignmentSettingsResourceModel struct {
	InstallTimeSettings *WinGetAppInstallTimeSettingsResourceModel `tfsdk:"install_time_settings"`
	Notifications       types.String                               `tfsdk:"notifications"` // Values: showAll, showReboot, hideAll, unknownFutureValue
	RestartSettings     *WinGetAppRestartSettingsResourceModel     `tfsdk:"restart_settings"`
}

type WinGetAppInstallTimeSettingsResourceModel struct {
	DeadlineDateTime types.String `tfsdk:"deadline_date_time"`
	UseLocalTime     types.Bool   `tfsdk:"use_local_time"`
}

type WinGetAppRestartSettingsResourceModel struct {
	CountdownDisplayBeforeRestartInMinutes     types.Int32 `tfsdk:"countdown_display_before_restart_in_minutes"`
	GracePeriodInMinutes                       types.Int32 `tfsdk:"grace_period_in_minutes"`
	RestartNotificationSnoozeDurationInMinutes types.Int32 `tfsdk:"restart_notification_snooze_duration_in_minutes"`
}
