// Base resource REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicyassignment?view=graph-rest-beta
package graphBetaDeviceManagementConfigurationPolicyAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceManagementConfigurationPolicyAssignmentResourceModel struct {
	ConfigurationPolicyId types.String                  `tfsdk:"settings_catalog_id"`
	ID                    types.String                  `tfsdk:"id"`
	Target                AssignmentTargetResourceModel `tfsdk:"target"`
	Source                types.String                  `tfsdk:"source"`
	SourceId              types.String                  `tfsdk:"source_id"`
	Timeouts              timeouts.Value                `tfsdk:"timeouts"`
}

// Target models
type AssignmentTargetResourceModel struct {
	TargetType                                 types.String `tfsdk:"target_type"` // allDevices, allLicensedUsers, configurationManagerCollection, exclusionGroupAssignment, groupAssignment
	DeviceAndAppManagementAssignmentFilterId   types.String `tfsdk:"device_and_app_management_assignment_filter_id"`
	DeviceAndAppManagementAssignmentFilterType types.String `tfsdk:"device_and_app_management_assignment_filter_type"`
	GroupId                                    types.String `tfsdk:"group_id"`
	CollectionId                               types.String `tfsdk:"collection_id"`
}
