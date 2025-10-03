package graphBetaWindowsBackupAndRestore

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapAssignmentsToTerraform maps assignments from the Graph API to Terraform state
func MapAssignmentsToTerraform(ctx context.Context, data *WindowsBackupAndRestoreResourceModel, assignments []graphmodels.EnrollmentConfigurationAssignmentable) {
	tflog.Debug(ctx, "Starting assignment mapping", map[string]any{
		"assignmentCount": len(assignments),
	})

	if len(assignments) == 0 {
		data.Assignments = types.SetNull(WindowsBackupAndRestoreAssignmentType())
		return
	}

	assignmentValues := make([]attr.Value, 0, len(assignments))

	for _, assignment := range assignments {
		target := assignment.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment target is nil, skipping assignment")
			continue
		}

		targetAttrs := map[string]attr.Value{
			"device_and_app_management_assignment_filter": types.ObjectNull(map[string]attr.Type{
				"filter_id":   types.StringType,
				"filter_type": types.StringType,
			}),
			"device_and_app_management_assignment_filter_id": types.StringNull(),
			"group_id": types.StringNull(),
			"intent":   types.StringNull(),
			"target":   types.StringNull(),
		}

		// Map target based on type
		switch targetType := target.(type) {
		case *graphmodels.AllDevicesAssignmentTarget:
			targetAttrs["target"] = types.StringValue("allDevices")
		case *graphmodels.AllLicensedUsersAssignmentTarget:
			targetAttrs["target"] = types.StringValue("allLicensedUsers")
		case *graphmodels.GroupAssignmentTarget:
			targetAttrs["target"] = types.StringValue("group")
			if groupId := targetType.GetGroupId(); groupId != nil {
				targetAttrs["group_id"] = types.StringValue(*groupId)
			}
		case *graphmodels.ExclusionGroupAssignmentTarget:
			targetAttrs["target"] = types.StringValue("exclusionGroup")
			if groupId := targetType.GetGroupId(); groupId != nil {
				targetAttrs["group_id"] = types.StringValue(*groupId)
			}
		default:
			tflog.Warn(ctx, "Unknown assignment target type", map[string]any{
				"targetType": fmt.Sprintf("%T", targetType),
			})
			continue
		}

		assignmentAttrs := map[string]attr.Value{
			"target": types.ObjectValueMust(map[string]attr.Type{
				"device_and_app_management_assignment_filter": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"filter_id":   types.StringType,
						"filter_type": types.StringType,
					},
				},
				"device_and_app_management_assignment_filter_id": types.StringType,
				"group_id": types.StringType,
				"intent":   types.StringType,
				"target":   types.StringType,
			}, targetAttrs),
		}

		assignmentValue, _ := types.ObjectValue(WindowsBackupAndRestoreAssignmentType().AttrTypes, assignmentAttrs)
		assignmentValues = append(assignmentValues, assignmentValue)
	}

	if len(assignmentValues) > 0 {
		set, diags := types.SetValue(WindowsBackupAndRestoreAssignmentType(), assignmentValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]any{
				"error": diags.Errors(),
			})
			data.Assignments = types.SetNull(WindowsBackupAndRestoreAssignmentType())
		} else {
			data.Assignments = set
		}
	} else {
		data.Assignments = types.SetNull(WindowsBackupAndRestoreAssignmentType())
	}

	tflog.Debug(ctx, "Completed assignment mapping", map[string]any{
		"assignmentCount": len(assignmentValues),
	})
}
