package graphBetaAppControlForBusinessBuiltInControls

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

// constructAssignment constructs and returns a DeviceHealthScriptsItemAssignPostRequestBody
func constructAssignment(ctx context.Context, data *AppControlForBusinessResourceBuiltInControlsModel) (devicemanagement.ConfigurationPoliciesItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Starting Device Health Script assignment construction")

	requestBody := devicemanagement.NewConfigurationPoliciesItemAssignPostRequestBody()
	scriptAssignments := make([]graphmodels.DeviceManagementConfigurationPolicyAssignmentable, 0)

	if data.Assignments.IsNull() || data.Assignments.IsUnknown() {
		tflog.Debug(ctx, "Assignments is null or unknown, creating empty assignments array")
		requestBody.SetAssignments(scriptAssignments)
		return requestBody, nil
	}

	var terraformAssignments []sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel
	diags := data.Assignments.ElementsAs(ctx, &terraformAssignments, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract assignments: %v", diags.Errors())
	}

	for idx, assignment := range terraformAssignments {
		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"index": idx,
		})

		graphAssignment := graphmodels.NewDeviceManagementConfigurationPolicyAssignment()

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

		scriptAssignments = append(scriptAssignments, graphAssignment)
	}

	tflog.Debug(ctx, "Completed assignment construction", map[string]interface{}{
		"totalAssignments": len(scriptAssignments),
	})

	requestBody.SetAssignments(scriptAssignments)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignment request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructTarget creates the appropriate target based on the target type
func constructTarget(ctx context.Context, targetType string, assignment sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel) graphmodels.DeviceAndAppManagementAssignmentTargetable {
	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch targetType {
	case "allDevicesAssignmentTarget":
		target = graphmodels.NewAllDevicesAssignmentTarget()
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
	case "exclusionGroupAssignmentTarget":
		exclusionTarget := graphmodels.NewExclusionGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.GroupId, exclusionTarget.SetGroupId)
		} else {
			tflog.Error(ctx, "Exclusion group assignment target missing required group_id", map[string]interface{}{
				"targetType": targetType,
			})
			return nil
		}
		target = exclusionTarget
	default:
		tflog.Error(ctx, "Unsupported target type", map[string]interface{}{
			"targetType": targetType,
		})
		return nil
	}

	// Set filter if provided and meaningful (not default values)
	if !assignment.FilterId.IsNull() && !assignment.FilterId.IsUnknown() &&
		assignment.FilterId.ValueString() != "" &&
		assignment.FilterId.ValueString() != "00000000-0000-0000-0000-000000000000" {

		convert.FrameworkToGraphString(assignment.FilterId, target.SetDeviceAndAppManagementAssignmentFilterId)

		if !assignment.FilterType.IsNull() && !assignment.FilterType.IsUnknown() &&
			assignment.FilterType.ValueString() != "" && assignment.FilterType.ValueString() != "none" {

			filterType := assignment.FilterType.ValueString()
			var filterTypeEnum graphmodels.DeviceAndAppManagementAssignmentFilterType
			switch filterType {
			case "include":
				filterTypeEnum = graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			case "exclude":
				filterTypeEnum = graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			default:
				tflog.Warn(ctx, "Unknown filter type, not setting filter", map[string]interface{}{
					"filterType": filterType,
				})
				return target
			}
			target.SetDeviceAndAppManagementAssignmentFilterType(&filterTypeEnum)
		}
	}

	return target
}
