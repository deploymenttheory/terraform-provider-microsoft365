package graphBetaSettingsCatalogInventoryPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func inventoryPolicyAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"group_id":    types.StringType,
			"filter_id":   types.StringType,
			"filter_type": types.StringType,
		},
	}
}

func MapAssignmentsToTerraform(ctx context.Context, data *InventoryPolicyResourceModel, assignments []graphmodels.DeviceManagementConfigurationPolicyAssignmentable) {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process")
		data.Assignments = types.SetNull(inventoryPolicyAssignmentType())
		return
	}

	assignmentValues := []attr.Value{}

	for i, assignment := range assignments {
		target := assignment.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment target is nil, skipping", map[string]any{"index": i})
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			tflog.Warn(ctx, "Assignment target OData type is nil, skipping", map[string]any{"index": i})
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
			assignmentObj["group_id"] = types.StringNull()

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

		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			assignmentObj["type"] = types.StringValue("exclusionGroupAssignmentTarget")
			if groupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
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
			assignmentObj["group_id"] = types.StringNull()
		}

		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		if filterID != nil && *filterID != "" && *filterID != "00000000-0000-0000-0000-000000000000" {
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

		objValue, diags := types.ObjectValue(inventoryPolicyAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
		if !diags.HasError() {
			assignmentValues = append(assignmentValues, objValue)
		}
	}

	if len(assignmentValues) > 0 {
		setVal, diags := types.SetValue(inventoryPolicyAssignmentType(), assignmentValues)
		if diags.HasError() {
			data.Assignments = types.SetNull(inventoryPolicyAssignmentType())
		} else {
			data.Assignments = setVal
		}
	} else {
		data.Assignments = types.SetNull(inventoryPolicyAssignmentType())
	}
}
