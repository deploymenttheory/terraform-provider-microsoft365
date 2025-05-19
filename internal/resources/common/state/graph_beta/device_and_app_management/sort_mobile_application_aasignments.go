package sharedStater

import (
	"sort"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
)

// sortMobileAppAssignments sorts a slice of mobile app assignments
// The sort order is as follows:
// 1. First tier: Sort by intent alphabetically
// 2. Second tier: Within same intent, sort by target_type alphabetically
// 3. Third tier: Within same target_type, sort by group_id alphabetically
// In mobile_application_assignments.go, update SortMobileAppAssignments:
func SortMobileAppAssignments(assignments []sharedmodels.MobileAppAssignmentResourceModel) {
	// Define intent order priorities based on observed API behavior
	intentOrder := map[string]int{
		"required":                   1,
		"available":                  2,
		"uninstall":                  3,
		"availableWithoutEnrollment": 4,
	}

	// Define target type order priorities based on observed API behavior
	targetTypeOrder := map[string]int{
		"groupAssignment":          1,
		"exclusionGroupAssignment": 2,
		"allLicensedUsers":         3,
		"allDevices":               4,
	}

	// Define filter type order priorities based on observed API behavior
	filterTypeOrder := map[string]int{
		"exclude": 1,
		"include": 2,
		"none":    3,
	}

	sort.SliceStable(assignments, func(i, j int) bool {
		// First sort by intent according to API's priority
		iIntent := intentOrder[assignments[i].Intent.ValueString()]
		jIntent := intentOrder[assignments[j].Intent.ValueString()]
		if iIntent != jIntent {
			return iIntent < jIntent
		}

		// Then sort by target type according to API's priority
		iTargetType := targetTypeOrder[assignments[i].Target.TargetType.ValueString()]
		jTargetType := targetTypeOrder[assignments[j].Target.TargetType.ValueString()]
		if iTargetType != jTargetType {
			return iTargetType < jTargetType
		}

		// For same target type, sort by filter type
		iFilterType := filterTypeOrder[assignments[i].Target.DeviceAndAppManagementAssignmentFilterType.ValueString()]
		jFilterType := filterTypeOrder[assignments[j].Target.DeviceAndAppManagementAssignmentFilterType.ValueString()]
		if iFilterType != jFilterType {
			return iFilterType < jFilterType
		}

		// For group assignments with same intent and filter type, sort by group ID
		if assignments[i].Target.TargetType.ValueString() == "groupAssignment" &&
			assignments[j].Target.TargetType.ValueString() == "groupAssignment" &&
			!assignments[i].Target.GroupId.IsNull() &&
			!assignments[j].Target.GroupId.IsNull() {
			return assignments[i].Target.GroupId.ValueString() <
				assignments[j].Target.GroupId.ValueString()
		}

		// Default to original order for equivalent assignments
		return i < j
	})
}
