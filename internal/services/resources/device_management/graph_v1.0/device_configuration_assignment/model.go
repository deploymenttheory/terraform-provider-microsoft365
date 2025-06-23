// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceconfigurationassignment?view=graph-rest-1.0
package graphDeviceConfigurationAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceConfigurationAssignmentResourceModel represents the Terraform resource model for DeviceConfigurationAssignment
type DeviceConfigurationAssignmentResourceModel struct {
	// Unique identifier for the assignment (computed)
	ID types.String `tfsdk:"id"`

	// ID of the device configuration policy to assign (required, forces replacement)
	DeviceConfigurationId types.String `tfsdk:"device_configuration_id"`

	// Target type for the assignment (required)
	// Possible values: allDevices, allLicensedUsers, configurationManagerCollection, exclusionGroupAssignment, groupAssignment
	TargetType types.String `tfsdk:"target_type"`

	// Group ID for group-based assignments (optional)
	// Required when target_type is groupAssignment, exclusionGroupAssignment, or configurationManagerCollection
	GroupId types.String `tfsdk:"group_id"`

	// Timeouts for CRUD operations
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
