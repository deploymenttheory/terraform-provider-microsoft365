package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	commonattr "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/attr"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// WindowsBiosConfigurationsAndOtherSettingsTemplateAssignmentType returns the object type used by the assignments set
func WindowsBiosConfigurationsAndOtherSettingsTemplateAssignmentType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"group_id":    types.StringType,
			"filter_id":   types.StringType,
			"filter_type": types.StringType,
		},
	}
}

// MapAssignmentsToTerraform maps the remote hardware configuration assignments to Terraform state
func MapAssignmentsToTerraform(
	ctx context.Context,
	data *WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel,
	assignments []graphmodels.HardwareConfigurationAssignmentable,
) {
	assignmentType := WindowsBiosConfigurationsAndOtherSettingsTemplateAssignmentType()

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process")
		data.Assignments = types.SetNull(assignmentType)
		return
	}

	tflog.Debug(ctx, "Starting assignment mapping process", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	assignmentObjects := make([]map[string]attr.Value, 0, len(assignments))

	for i, assignment := range assignments {
		target := assignment.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment target is nil, skipping assignment", map[string]any{
				"assignmentIndex": i,
				"resourceId":      data.ID.ValueString(),
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
					"resourceId":      data.ID.ValueString(),
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
		case "#microsoft.graph.groupAssignmentTarget":
			assignmentObj["type"] = types.StringValue("groupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(
						ctx,
						"Group ID is nil/empty for group assignment target",
						map[string]any{
							"assignmentIndex": i,
							"resourceId":      data.ID.ValueString(),
						},
					)
				}
			} else {
				tflog.Error(
					ctx,
					"Failed to cast target to GroupAssignmentTargetable",
					map[string]any{
						"assignmentIndex": i,
						"resourceId":      data.ID.ValueString(),
					},
				)
			}

		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			assignmentObj["type"] = types.StringValue("exclusionGroupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(
						ctx,
						"Group ID is nil/empty for exclusion group assignment target",
						map[string]any{
							"assignmentIndex": i,
							"resourceId":      data.ID.ValueString(),
						},
					)
				}
			} else {
				tflog.Error(
					ctx,
					"Failed to cast target to ExclusionGroupAssignmentTargetable",
					map[string]any{
						"assignmentIndex": i,
						"resourceId":      data.ID.ValueString(),
					},
				)
			}

		default:
			tflog.Warn(ctx, "Unknown target type encountered, skipping assignment", map[string]any{
				"assignmentIndex": i,
				"targetType":      *odataType,
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		// The sentinel values below must match the schema defaults in
		// commonschemagraphbeta.HardwareConfigurationAssignmentsSchema, otherwise every plan
		// reports a diff on assignments that carry no filter.
		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		if filterID != nil && *filterID != "" &&
			*filterID != "00000000-0000-0000-0000-000000000000" {
			assignmentObj["filter_id"] = convert.GraphToFrameworkString(filterID)
		} else {
			assignmentObj["filter_id"] = types.StringValue("00000000-0000-0000-0000-000000000000")
		}

		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()
		if filterType != nil {
			switch *filterType {
			case graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				assignmentObj["filter_type"] = types.StringValue("include")
			case graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				assignmentObj["filter_type"] = types.StringValue("exclude")
			default:
				assignmentObj["filter_type"] = types.StringValue("none")
			}
		} else {
			assignmentObj["filter_type"] = types.StringValue("none")
		}

		assignmentObjects = append(assignmentObjects, assignmentObj)
	}

	data.Assignments = commonattr.ObjectSetFromSlice(ctx, assignmentType.AttrTypes,
		func(i int) map[string]attr.Value { return assignmentObjects[i] }, len(assignmentObjects))

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]any{
		"finalAssignmentCount": len(assignmentObjects),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})
}
