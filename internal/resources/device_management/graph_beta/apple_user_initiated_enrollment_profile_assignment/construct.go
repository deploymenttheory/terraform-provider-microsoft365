package graphBetaAppleUserInitiatedEnrollmentProfileAssignment

import (
	"context"
	"fmt"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructAppleUserInitiatedEnrollmentProfileAssignment constructs an Apple User Initiated Enrollment Profile Assignment object for API requests
func ConstructAppleUserInitiatedEnrollmentProfileAssignment(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	data AppleUserInitiatedEnrollmentProfileAssignmentResourceModel,
	isUpdate bool,
) (graphmodels.AppleEnrollmentProfileAssignmentable, error) {

	assignment := graphmodels.NewAppleEnrollmentProfileAssignment()

	target, err := constructTarget(data.Target)
	if err != nil {
		return nil, fmt.Errorf("error constructing target: %v", err)
	}
	assignment.SetTarget(target)

	return assignment, nil
}

// constructTarget creates an assignment target based on the target type
func constructTarget(targetData AssignmentTargetResourceModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	targetType := targetData.TargetType.ValueString()

	switch targetType {
	case "allUsers":
		// Use AllLicensedUsersAssignmentTarget for Apple assignments
		target := graphmodels.NewAllLicensedUsersAssignmentTarget()
		setCommonTargetProperties(target, targetData)
		return target, nil

	case "group":
		target := graphmodels.NewGroupAssignmentTarget()
		if !targetData.GroupId.IsNull() && !targetData.GroupId.IsUnknown() {
			groupId := targetData.GroupId.ValueString()
			target.SetGroupId(&groupId)
		}
		setCommonTargetProperties(target, targetData)
		return target, nil

	case "exclusionGroup":
		target := graphmodels.NewExclusionGroupAssignmentTarget()
		if !targetData.GroupId.IsNull() && !targetData.GroupId.IsUnknown() {
			groupId := targetData.GroupId.ValueString()
			target.SetGroupId(&groupId)
		}
		setCommonTargetProperties(target, targetData)
		return target, nil

	case "user":
		// For user targeting in Apple assignments, we'll use GroupAssignmentTarget
		// with the user ID as group ID (this is how the API expects it)
		target := graphmodels.NewGroupAssignmentTarget()
		if !targetData.EntraObjectId.IsNull() && !targetData.EntraObjectId.IsUnknown() {
			userId := targetData.EntraObjectId.ValueString()
			target.SetGroupId(&userId)
		}
		setCommonTargetProperties(target, targetData)
		return target, nil

	default:
		return nil, fmt.Errorf("unsupported target type: %s", targetType)
	}
}

// setCommonTargetProperties sets properties common to all target types
func setCommonTargetProperties(target graphmodels.DeviceAndAppManagementAssignmentTargetable, targetData AssignmentTargetResourceModel) {
	if !targetData.DeviceAndAppManagementAssignmentFilterId.IsNull() && !targetData.DeviceAndAppManagementAssignmentFilterId.IsUnknown() {
		filterId := targetData.DeviceAndAppManagementAssignmentFilterId.ValueString()
		target.SetDeviceAndAppManagementAssignmentFilterId(&filterId)
	}

	if !targetData.DeviceAndAppManagementAssignmentFilterType.IsNull() && !targetData.DeviceAndAppManagementAssignmentFilterType.IsUnknown() {
		filterType := targetData.DeviceAndAppManagementAssignmentFilterType.ValueString()
		if parsedFilterType, err := graphmodels.ParseDeviceAndAppManagementAssignmentFilterType(filterType); err == nil {
			target.SetDeviceAndAppManagementAssignmentFilterType(parsedFilterType.(*graphmodels.DeviceAndAppManagementAssignmentFilterType))
		}
	}
}
