package graphBetaMacOSDeviceEnrollmentPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MacOSDeviceEnrollmentPolicyAssignmentType returns the object type for the assignments set.
func MacOSDeviceEnrollmentPolicyAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"group_id":    types.StringType,
			"filter_id":   types.StringType,
			"filter_type": types.StringType,
		},
	}
}

// mapAssignmentsToState maps the remote configuration policy assignments to Terraform state.
func mapAssignmentsToState(ctx context.Context, stateModel *MacOSDeviceEnrollmentPolicyResourceModel, assignmentsResponse graphmodels.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable) {
	if assignmentsResponse == nil {
		stateModel.Assignments = types.SetNull(MacOSDeviceEnrollmentPolicyAssignmentType())
		return
	}

	assignments := assignmentsResponse.GetValue()
	if len(assignments) == 0 {
		stateModel.Assignments = types.SetNull(MacOSDeviceEnrollmentPolicyAssignmentType())
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
			tflog.Warn(ctx, "Assignment target is nil, skipping assignment", map[string]any{"assignmentIndex": i})
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			tflog.Warn(ctx, "Assignment target OData type is nil, skipping assignment", map[string]any{"assignmentIndex": i})
			continue
		}

		assignmentObj := map[string]attr.Value{
			"type":        types.StringNull(),
			"group_id":    types.StringNull(),
			"filter_id":   types.StringNull(),
			"filter_type": types.StringNull(),
		}

		switch *odataType {
		case "#microsoft.graph.allDevicesAssignmentTarget":
			assignmentObj["type"] = types.StringValue("allDevicesAssignmentTarget")

		case "#microsoft.graph.allLicensedUsersAssignmentTarget":
			assignmentObj["type"] = types.StringValue("allLicensedUsersAssignmentTarget")

		case "#microsoft.graph.groupAssignmentTarget":
			assignmentObj["type"] = types.StringValue("groupAssignmentTarget")
			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				if groupId := groupTarget.GetGroupId(); groupId != nil && *groupId != "" {
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				}
			}

		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			assignmentObj["type"] = types.StringValue("exclusionGroupAssignmentTarget")
			if groupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				if groupId := groupTarget.GetGroupId(); groupId != nil && *groupId != "" {
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				}
			}

		default:
			tflog.Warn(ctx, "Unknown target type encountered", map[string]any{"assignmentIndex": i, "targetType": *odataType})
			continue
		}

		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		if filterID != nil && *filterID != "" && *filterID != "00000000-0000-0000-0000-000000000000" {
			assignmentObj["filter_id"] = convert.GraphToFrameworkString(filterID)
		} else {
			assignmentObj["filter_id"] = types.StringValue("00000000-0000-0000-0000-000000000000")
		}

		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()
		switch {
		case filterType == nil:
			assignmentObj["filter_type"] = types.StringValue("none")
		case *filterType == graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
			assignmentObj["filter_type"] = types.StringValue("include")
		case *filterType == graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
			assignmentObj["filter_type"] = types.StringValue("exclude")
		default:
			assignmentObj["filter_type"] = types.StringValue("none")
		}

		objValue, diags := types.ObjectValue(MacOSDeviceEnrollmentPolicyAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignment object value", map[string]any{"assignmentIndex": i, "errors": diags.Errors()})
			continue
		}
		assignmentValues = append(assignmentValues, objValue)
	}

	if len(assignmentValues) == 0 {
		stateModel.Assignments = types.SetNull(MacOSDeviceEnrollmentPolicyAssignmentType())
		return
	}

	setVal, diags := types.SetValue(MacOSDeviceEnrollmentPolicyAssignmentType(), assignmentValues)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to create assignments set", map[string]any{"errors": diags.Errors()})
		stateModel.Assignments = types.SetNull(MacOSDeviceEnrollmentPolicyAssignmentType())
		return
	}
	stateModel.Assignments = setVal
}
