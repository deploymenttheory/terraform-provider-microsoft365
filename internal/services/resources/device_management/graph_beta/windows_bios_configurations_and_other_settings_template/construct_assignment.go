package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
)

var (
	errAssignmentMissingType    = errors.New("assignment is missing a target type")
	errAssignmentMissingGroupId = errors.New("assignment target is missing the required group_id")
	errUnsupportedTargetType    = errors.New("unsupported assignment target type")
	errExtractAssignments       = errors.New("failed to extract assignments")
	errConstructTarget          = errors.New("failed to construct assignment target")
)

// constructAssignment constructs and returns a HardwareConfigurationsItemAssignPostRequestBody
func constructAssignment(
	ctx context.Context,
	data *WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel,
) (devicemanagement.HardwareConfigurationsItemAssignPostRequestBodyable, error) {
	tflog.Debug(ctx, "Starting hardware configuration assignment construction")

	requestBody := devicemanagement.NewHardwareConfigurationsItemAssignPostRequestBody()
	hardwareConfigurationAssignments := make([]graphmodels.HardwareConfigurationAssignmentable, 0)

	if data.Assignments.IsNull() || data.Assignments.IsUnknown() {
		tflog.Debug(ctx, "Assignments is null or unknown, creating empty assignments array")
		requestBody.SetHardwareConfigurationAssignments(hardwareConfigurationAssignments)
		return requestBody, nil
	}

	var terraformAssignments []sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel
	diags := data.Assignments.ElementsAs(ctx, &terraformAssignments, false)
	if diags.HasError() {
		return nil, fmt.Errorf("%w: %v", errExtractAssignments, diags.Errors())
	}

	for idx, assignment := range terraformAssignments {
		tflog.Debug(ctx, "Processing assignment", map[string]any{
			"index": idx,
		})

		if assignment.Type.IsNull() || assignment.Type.IsUnknown() {
			return nil, fmt.Errorf("%w: index %d", errAssignmentMissingType, idx)
		}

		targetType := assignment.Type.ValueString()

		target, err := constructTarget(ctx, targetType, assignment)
		if err != nil {
			return nil, fmt.Errorf("%w at index %d: %w", errConstructTarget, idx, err)
		}

		graphAssignment := graphmodels.NewHardwareConfigurationAssignment()
		graphAssignment.SetTarget(target)

		hardwareConfigurationAssignments = append(hardwareConfigurationAssignments, graphAssignment)
	}

	tflog.Debug(ctx, "Completed assignment construction", map[string]any{
		"totalAssignments": len(hardwareConfigurationAssignments),
	})

	requestBody.SetHardwareConfigurationAssignments(hardwareConfigurationAssignments)

	if err := constructors.DebugLogGraphObject(
		ctx,
		"Constructed assignment request body",
		requestBody,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment request body", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructTarget creates the appropriate target based on the target type
func constructTarget(
	ctx context.Context,
	targetType string,
	assignment sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel,
) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	hasGroupId := !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() &&
		assignment.GroupId.ValueString() != ""

	switch targetType {
	case "groupAssignmentTarget":
		if !hasGroupId {
			return nil, fmt.Errorf("%w: %s", errAssignmentMissingGroupId, targetType)
		}
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		convert.FrameworkToGraphString(assignment.GroupId, groupTarget.SetGroupId)
		target = groupTarget
	case "exclusionGroupAssignmentTarget":
		if !hasGroupId {
			return nil, fmt.Errorf("%w: %s", errAssignmentMissingGroupId, targetType)
		}
		exclusionTarget := graphmodels.NewExclusionGroupAssignmentTarget()
		convert.FrameworkToGraphString(assignment.GroupId, exclusionTarget.SetGroupId)
		target = exclusionTarget
	default:
		return nil, fmt.Errorf("%w: %q", errUnsupportedTargetType, targetType)
	}

	// Set filter if provided and meaningful (not default values)
	if !assignment.FilterId.IsNull() && !assignment.FilterId.IsUnknown() &&
		assignment.FilterId.ValueString() != "" &&
		assignment.FilterId.ValueString() != "00000000-0000-0000-0000-000000000000" {

		convert.FrameworkToGraphString(
			assignment.FilterId,
			target.SetDeviceAndAppManagementAssignmentFilterId,
		)

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
				tflog.Warn(ctx, "Unknown filter type, not setting filter", map[string]any{
					"filterType": filterType,
				})
				return target, nil
			}
			target.SetDeviceAndAppManagementAssignmentFilterType(&filterTypeEnum)
		}
	}

	return target, nil
}
