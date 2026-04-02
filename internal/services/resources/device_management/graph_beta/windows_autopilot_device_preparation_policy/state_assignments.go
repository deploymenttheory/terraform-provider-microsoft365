package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// WindowsAutopilotDevicePreparationPolicyAssignmentType returns the object type for assignment model
func WindowsAutopilotDevicePreparationPolicyAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"group_id":    types.StringType,
			"filter_id":   types.StringType,
			"filter_type": types.StringType,
		},
	}
}

// mapAssignmentsToState maps the assignments response to the state model.
func mapAssignmentsToState(
	ctx context.Context,
	stateModel *WindowsAutopilotDevicePreparationPolicyResourceModel,
	assignmentsResponse graphmodels.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable,
) {
	if assignmentsResponse == nil {
		stateModel.Assignments = types.SetNull(
			WindowsAutopilotDevicePreparationPolicyAssignmentType(),
		)
		return
	}

	assignments := assignmentsResponse.GetValue()
	if len(assignments) == 0 {
		stateModel.Assignments = types.SetNull(
			WindowsAutopilotDevicePreparationPolicyAssignmentType(),
		)
		return
	}

	tflog.Debug(ctx, "Starting assignment mapping process", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      stateModel.ID.ValueString(),
	})

	assignmentValues := []attr.Value{}

	for i, assignment := range assignments {
		target := assignment.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment target is nil, skipping assignment", map[string]any{
				"assignmentIndex": i,
			})
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			tflog.Warn(
				ctx,
				"Assignment target OData type is nil, skipping assignment",
				map[string]any{
					"assignmentIndex": i,
				},
			)
			continue
		}

		assignmentObj := map[string]attr.Value{
			"type":        types.StringNull(),
			"group_id":    types.StringNull(),
			"filter_id":   types.StringNull(),
			"filter_type": types.StringNull(),
		}

		switch *odataType {
		case "#microsoft.graph.allLicensedUsersAssignmentTarget":
			assignmentObj["type"] = types.StringValue("allLicensedUsersAssignmentTarget")
			assignmentObj["group_id"] = types.StringNull()

		case "#microsoft.graph.groupAssignmentTarget":
			assignmentObj["type"] = types.StringValue("groupAssignmentTarget")
			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					assignmentObj["group_id"] = types.StringNull()
				}
			} else {
				assignmentObj["group_id"] = types.StringNull()
			}

		default:
			tflog.Warn(ctx, "Unknown target type encountered", map[string]any{
				"assignmentIndex": i,
				"targetType":      *odataType,
			})
			continue
		}

		// Process filter ID
		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		if filterID != nil && *filterID != "" &&
			*filterID != "00000000-0000-0000-0000-000000000000" {
			assignmentObj["filter_id"] = convert.GraphToFrameworkString(filterID)
		} else {
			assignmentObj["filter_id"] = types.StringNull()
		}

		// Process filter type
		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()
		if filterType != nil {
			switch *filterType {
			case graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				assignmentObj["filter_type"] = types.StringValue("include")
			case graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				assignmentObj["filter_type"] = types.StringValue("exclude")
			case graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				assignmentObj["filter_type"] = types.StringValue("none")
			default:
				assignmentObj["filter_type"] = types.StringValue("none")
			}
		} else {
			assignmentObj["filter_type"] = types.StringValue("none")
		}

		objValue, diags := types.ObjectValue(
			WindowsAutopilotDevicePreparationPolicyAssignmentType().(types.ObjectType).AttrTypes,
			assignmentObj,
		)
		if !diags.HasError() {
			assignmentValues = append(assignmentValues, objValue)
		} else {
			tflog.Error(ctx, "Failed to create assignment object value", map[string]any{
				"assignmentIndex": i,
				"errors":          diags.Errors(),
			})
		}
	}

	if len(assignmentValues) > 0 {
		setVal, diags := types.SetValue(
			WindowsAutopilotDevicePreparationPolicyAssignmentType(),
			assignmentValues,
		)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]any{
				"errors": diags.Errors(),
			})
			stateModel.Assignments = types.SetNull(
				WindowsAutopilotDevicePreparationPolicyAssignmentType(),
			)
		} else {
			stateModel.Assignments = setVal
		}
	} else {
		stateModel.Assignments = types.SetNull(
			WindowsAutopilotDevicePreparationPolicyAssignmentType(),
		)
	}

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]any{
		"finalAssignmentCount": len(assignmentValues),
		"originalAssignments":  len(assignments),
	})
}
