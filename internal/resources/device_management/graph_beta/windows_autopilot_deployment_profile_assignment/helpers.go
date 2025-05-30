package graphBetaWindowsAutopilotDeploymentProfileAssignment

import (
	"context"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Helper function to determine if an assignment matches our resource
func matchesAssignment(ctx context.Context, object WindowsAutopilotDeploymentProfileAssignmentResourceModel, assign graphmodels.WindowsAutopilotDeploymentProfileAssignmentable) bool {
	if assign == nil {
		return false
	}

	// Compare source
	if assign.GetSource() != nil {
		sourceValue := assign.GetSource().String()
		if sourceValue != object.Source.ValueString() {
			return false
		}
	}

	// Compare target type
	if target := assign.GetTarget(); target != nil {
		targetType := getTargetTypeFromOdataType(*target.GetOdataType())
		if targetType != object.Target.TargetType.ValueString() {
			return false
		}

		// For group assignments, also compare group ID
		if targetType == "groupAssignment" || targetType == "exclusionGroupAssignment" {
			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				if groupTarget.GetGroupId() != nil && *groupTarget.GetGroupId() != object.Target.GroupId.ValueString() {
					return false
				}
			}
		}
	}

	return true
}

// Helper function to extract target type from OData type
func getTargetTypeFromOdataType(odataType string) string {
	switch odataType {
	case "#microsoft.graph.allDevicesAssignmentTarget":
		return "allDevices"
	case "#microsoft.graph.allLicensedUsersAssignmentTarget":
		return "allLicensedUsers"
	case "#microsoft.graph.groupAssignmentTarget":
		return "groupAssignment"
	case "#microsoft.graph.exclusionGroupAssignmentTarget":
		return "exclusionGroupAssignment"
	default:
		return ""
	}
}
