package graphBetaMacOSSoftwareUpdateConfiguration

import (
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_management"
)

// ValidateAssignmentConfiguration validates the assignment configuration
// - ExcludeGroupIds must be in alphanumeric order by group_id
// - IncludeGroups must be in alphanumeric order by group_id
// - No group is used more than once across include and exclude assignments
// - AllDevices and IncludeGroups cannot be used at the same time
// - AllUsers and IncludeGroups cannot be used at the same time
func validateAssignmentConfig(config *sharedmodels.DeviceManagementScriptAssignmentResourceModel) error {
	// Validate include_group_ids are in alphanumeric order
	if len(config.IncludeGroupIds) > 1 {
		for i := 0; i < len(config.IncludeGroupIds)-1; i++ {
			current := config.IncludeGroupIds[i].ValueString()
			next := config.IncludeGroupIds[i+1].ValueString()
			if current > next {
				return fmt.Errorf("include_group_ids must be in alphanumeric order. Found %s before %s",
					current, next)
			}
		}
	}

	// Validate exclude_group_ids are in alphanumeric order
	if len(config.ExcludeGroupIds) > 1 {
		for i := 0; i < len(config.ExcludeGroupIds)-1; i++ {
			current := config.ExcludeGroupIds[i].ValueString()
			next := config.ExcludeGroupIds[i+1].ValueString()
			if current > next {
				return fmt.Errorf("exclude_group_ids must be in alphanumeric order. Found %s before %s",
					current, next)
			}
		}
	}

	// Validate no group is used more than once across include and exclude assignments
	for _, includeGroupId := range config.IncludeGroupIds {
		for _, excludeGroupId := range config.ExcludeGroupIds {
			if !includeGroupId.IsNull() && !excludeGroupId.IsNull() &&
				includeGroupId.ValueString() == excludeGroupId.ValueString() {
				return fmt.Errorf("group %s is used in both include and exclude assignments. Each group assignment can only be used once across all assignment rules",
					includeGroupId.ValueString())
			}
		}
	}

	// Validate AllDevices cannot be used with Include Groups
	if !config.AllDevices.IsNull() && config.AllDevices.ValueBool() && len(config.IncludeGroupIds) > 0 {
		return fmt.Errorf("cannot assign to All Devices and Include Groups at the same time")
	}

	// Validate AllUsers cannot be used with Include Groups
	if !config.AllUsers.IsNull() && config.AllUsers.ValueBool() && len(config.IncludeGroupIds) > 0 {
		return fmt.Errorf("cannot assign to All Users and Include Groups at the same time")
	}

	return nil
}
