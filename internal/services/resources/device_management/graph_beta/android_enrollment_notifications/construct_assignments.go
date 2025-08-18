package graphBetaAndroidEnrollmentNotifications

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAssignments maps the Terraform assignment data to the SDK assignment model.
func constructAssignments(ctx context.Context, data *AndroidEnrollmentNotificationsResourceModel) (devicemanagement.DeviceEnrollmentConfigurationsItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Starting enrollment configuration assignment construction")

	requestBody := devicemanagement.NewDeviceEnrollmentConfigurationsItemAssignPostRequestBody()
	enrollmentAssignments := make([]graphmodels.EnrollmentConfigurationAssignmentable, 0)

	if data.Assignments.IsNull() || data.Assignments.IsUnknown() {
		tflog.Debug(ctx, "Assignments is null or unknown, creating empty assignments array")
		requestBody.SetEnrollmentConfigurationAssignments(enrollmentAssignments)
		return requestBody, nil
	}

	var terraformAssignments []sharedmodels.AndroidNotificationAssignmentModel
	diags := data.Assignments.ElementsAs(ctx, &terraformAssignments, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract assignments: %v", diags.Errors())
	}

	for idx, assignment := range terraformAssignments {
		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"index": idx,
		})

		graphAssignment := graphmodels.NewEnrollmentConfigurationAssignment()

		if assignment.Type.IsNull() || assignment.Type.IsUnknown() {
			tflog.Error(ctx, "Assignment target type is missing or invalid", map[string]interface{}{
				"index": idx,
			})
			continue
		}

		targetType := assignment.Type.ValueString()

		target := constructTarget(ctx, targetType, assignment)
		if target == nil {
			tflog.Error(ctx, "Failed to create target", map[string]interface{}{
				"index":      idx,
				"targetType": targetType,
			})
			continue
		}

		graphAssignment.SetTarget(target)

		enrollmentAssignments = append(enrollmentAssignments, graphAssignment)
	}

	requestBody.SetEnrollmentConfigurationAssignments(enrollmentAssignments)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructTarget creates the appropriate target based on the target type
func constructTarget(ctx context.Context, targetType string, assignment sharedmodels.AndroidNotificationAssignmentModel) graphmodels.DeviceAndAppManagementAssignmentTargetable {
	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch targetType {
	case "allLicensedUsersAssignmentTarget":
		target = graphmodels.NewAllLicensedUsersAssignmentTarget()
	case "groupAssignmentTarget":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.GroupId, groupTarget.SetGroupId)
		} else {
			tflog.Error(ctx, "Group assignment target missing required group_id", map[string]interface{}{
				"targetType": targetType,
			})
			return nil
		}
		target = groupTarget
	default:
		tflog.Error(ctx, "Unsupported target type", map[string]interface{}{
			"targetType": targetType,
		})
		return nil
	}

	return target
}
