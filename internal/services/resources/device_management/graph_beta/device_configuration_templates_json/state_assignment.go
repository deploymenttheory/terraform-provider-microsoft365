package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var assignmentType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type": types.StringType, "group_id": types.StringType,
	"filter_id": types.StringType, "filter_type": types.StringType,
}}

func mapAssignment(
	target graphmodels.DeviceAndAppManagementAssignmentTargetable,
) (attr.Value, error) {
	values := map[string]attr.Value{
		"group_id":    types.StringNull(),
		"filter_id":   types.StringValue(noFilterID),
		"filter_type": types.StringValue("none"),
	}
	switch target := target.(type) {
	case *graphmodels.AllDevicesAssignmentTarget:
		values["type"] = types.StringValue("allDevicesAssignmentTarget")
	case *graphmodels.AllLicensedUsersAssignmentTarget:
		values["type"] = types.StringValue("allLicensedUsersAssignmentTarget")
	case *graphmodels.ExclusionGroupAssignmentTarget:
		values["type"] = types.StringValue("exclusionGroupAssignmentTarget")
		values["group_id"] = convert.GraphToFrameworkString(target.GetGroupId())
	case *graphmodels.GroupAssignmentTarget:
		values["type"] = types.StringValue("groupAssignmentTarget")
		values["group_id"] = convert.GraphToFrameworkString(target.GetGroupId())
	default:
		return nil, fmt.Errorf("%w: unsupported assignment target %T", errInvalidAssignment, target)
	}
	if id := target.GetDeviceAndAppManagementAssignmentFilterId(); id != nil && *id != "" {
		values["filter_id"] = types.StringValue(*id)
	}
	if filter := target.GetDeviceAndAppManagementAssignmentFilterType(); filter != nil {
		values["filter_type"] = types.StringValue(filter.String())
	}
	value, diags := types.ObjectValue(assignmentType.AttrTypes, values)
	if diags.HasError() {
		return nil, fmt.Errorf("%w: map assignment: %v", errInvalidAssignment, diags.Errors())
	}
	return value, nil
}

// MapAssignmentsToTerraform maps assignment targets using the shared assignment schema.
func MapAssignmentsToTerraform(
	ctx context.Context,
	data *DeviceConfigurationTemplatesJsonResourceModel,
	assignments []graphmodels.DeviceConfigurationAssignmentable,
) diag.Diagnostics {
	var diags diag.Diagnostics
	values := make([]attr.Value, 0, len(assignments))
	for _, assignment := range assignments {
		if assignment == nil || assignment.GetTarget() == nil {
			diags.AddError(
				"Invalid assignment response",
				"Graph returned an assignment without a target.",
			)
			return diags
		}
		value, err := mapAssignment(assignment.GetTarget())
		if err != nil {
			diags.AddError("Cannot read assignment", err.Error())
			return diags
		}
		values = append(values, value)
	}
	if len(values) == 0 && data.Assignments.IsNull() {
		data.Assignments = types.SetNull(assignmentType)
	} else {
		var setDiags diag.Diagnostics
		data.Assignments, setDiags = types.SetValue(assignmentType, values)
		diags.Append(setDiags...)
	}
	return diags
}
