package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteAssignmentsToTerraform maps Graph API assignments to Terraform state
func mapRemoteAssignmentsToTerraform(ctx context.Context, assignments []graphmodels.WindowsAutopilotDeploymentProfileAssignmentable) (types.Set, error) {
	if len(assignments) == 0 {
		return types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
		}), nil
	}

	var assignmentObjects []attr.Value

	for _, assignment := range assignments {
		if assignment == nil {
			continue
		}

		target := assignment.GetTarget()
		if target == nil {
			continue
		}

		var assignmentType string
		var groupID types.String = types.StringNull()

		// Determine assignment type and extract group ID if applicable
		odataType := target.GetOdataType()
		if odataType != nil {
			switch *odataType {
			case "#microsoft.graph.groupAssignmentTarget":
				assignmentType = "groupAssignmentTarget"
				if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
					if groupTarget.GetGroupId() != nil {
						groupID = types.StringValue(*groupTarget.GetGroupId())
					}
				}

			case "#microsoft.graph.exclusionGroupAssignmentTarget":
				assignmentType = "exclusionGroupAssignmentTarget"
				if exclusionTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
					if exclusionTarget.GetGroupId() != nil {
						groupID = types.StringValue(*exclusionTarget.GetGroupId())
					}
				}

			case "#microsoft.graph.allDevicesAssignmentTarget":
				assignmentType = "allDevicesAssignmentTarget"
				// No group ID for all devices assignment

			default:
				tflog.Warn(ctx, "Unknown assignment target type", map[string]interface{}{
					"odata_type": *odataType,
				})
				continue
			}
		} else {
			tflog.Warn(ctx, "Assignment target has no OData type")
			continue
		}

		assignmentObj, diags := types.ObjectValue(
			map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
			map[string]attr.Value{
				"type":     types.StringValue(assignmentType),
				"group_id": groupID,
			},
		)

		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignment object", map[string]interface{}{
				"diagnostics": diags,
			})
			continue
		}

		assignmentObjects = append(assignmentObjects, assignmentObj)
	}

	if len(assignmentObjects) == 0 {
		return types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
		}), nil
	}

	assignmentSet, diags := types.SetValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
		},
		assignmentObjects,
	)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to create assignments set", map[string]interface{}{
			"diagnostics": diags,
		})
		return types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
		}), nil
	}

	return assignmentSet, nil
}