package graphBetaDeviceAndAppManagementAppAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Helper function to determine if an assignment matches our resource
func matchesAssignment(ctx context.Context, object MobileAppAssignmentResourceModel, assign any) bool {
	tflog.Debug(ctx, "Matching assignment against configuration", map[string]any{
		"mobile_app_id": object.MobileAppId.ValueString(),
		"intent":        object.Intent.ValueString(),
	})

	assignment, ok := assign.(graphmodels.MobileAppAssignmentable)
	if !ok {
		return false
	}

	if assignment.GetIntent() != nil {
		intentStr := assignment.GetIntent().String()
		if !object.Intent.IsNull() && object.Intent.ValueString() != intentStr {
			return false
		}
	}

	target := assignment.GetTarget()
	if target != nil {
		odataType := target.GetOdataType()
		if odataType != nil {
			targetType := getTargetTypeFromOdataType(*odataType)
			if !object.Target.TargetType.IsNull() && object.Target.TargetType.ValueString() != targetType {
				return false
			}
		}

		// Check group ID if applicable
		if !object.Target.GroupId.IsNull() {
			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				if groupTarget.GetGroupId() == nil || object.Target.GroupId.ValueString() != *groupTarget.GetGroupId() {
					return false
				}
			}
		}

		// Check filter ID if applicable
		if !object.Target.DeviceAndAppManagementAssignmentFilterId.IsNull() {
			filterId := target.GetDeviceAndAppManagementAssignmentFilterId()
			if filterId == nil || object.Target.DeviceAndAppManagementAssignmentFilterId.ValueString() != *filterId {
				return false
			}
		}
	}

	return true
}

// Helper to get target type from odata.type
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
	case "#microsoft.graph.androidFotaDeploymentAssignmentTarget":
		return "androidFotaDeployment"
	case "#microsoft.graph.configurationManagerCollectionAssignmentTarget":
		return "configurationManagerCollection"
	default:
		return ""
	}
}
