// Base resource REF: https://learn.microsoft.com/en-us/graph/api/intune-enrollment-appleenrollmentprofileassignment-get?view=graph-rest-beta
package graphBetaAppleUserInitiatedEnrollmentProfileAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppleUserInitiatedEnrollmentProfileAssignmentResourceModel represents the Terraform resource model for Apple User Initiated Enrollment Profile Assignment
type AppleUserInitiatedEnrollmentProfileAssignmentResourceModel struct {
	AppleUserInitiatedEnrollmentProfileId types.String                  `tfsdk:"apple_user_initiated_enrollment_profile_id"`
	ID                                    types.String                  `tfsdk:"id"`
	Target                                AssignmentTargetResourceModel `tfsdk:"target"`
	Timeouts                              timeouts.Value                `tfsdk:"timeouts"`
}

// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-alldevicesassignmenttarget?view=graph-rest-beta
type AssignmentTargetResourceModel struct {
	TargetType                                 types.String `tfsdk:"target_type"` // allDevices, allLicensedUsers, exclusionGroupAssignment, groupAssignment
	DeviceAndAppManagementAssignmentFilterId   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"`
	GroupId                                    types.String `tfsdk:"group_id"`
	EntraObjectId                              types.String `tfsdk:"entra_object_id"`
}
