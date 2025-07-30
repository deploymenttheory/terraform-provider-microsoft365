package graphBetaWindowsQualityUpdatePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the Graph API model into the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsQualityUpdatePolicyResourceModel, remoteResource graphmodels.WindowsQualityUpdatePolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{"resourceId": remoteResource.GetId()})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.HotpatchEnabled = convert.GraphToFrameworkBool(remoteResource.GetHotpatchEnabled())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, setting assignments to null", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(WindowsQualityUpdatePolicyAssignmentType())
	} else {
		tflog.Debug(ctx, "Starting assignment mapping process", map[string]interface{}{
			"resourceId":      data.ID.ValueString(),
			"assignmentCount": len(assignments),
		})
		MapAssignmentsToTerraform(ctx, data, assignments)
		tflog.Debug(ctx, "Completed assignment mapping process", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// WindowsUpdateRingAssignmentType returns the object type for WindowsRemediationScriptAssignmentModel
func WindowsQualityUpdatePolicyAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":     types.StringType,
			"group_id": types.StringType,
		},
	}
}

// MapAssignmentsToTerraform maps the remote Windows Quality Update Policy assignments to Terraform state
func MapAssignmentsToTerraform(ctx context.Context, data *WindowsQualityUpdatePolicyResourceModel, assignments []graphmodels.WindowsQualityUpdatePolicyAssignmentable) {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process")
		data.Assignments = types.SetNull(WindowsQualityUpdatePolicyAssignmentType())
		return
	}

	tflog.Debug(ctx, "Starting assignment mapping process", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	assignmentValues := []attr.Value{}

	for i, assignment := range assignments {
		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		target := assignment.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment target is nil, skipping assignment", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			tflog.Warn(ctx, "Assignment target OData type is nil, skipping assignment", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		tflog.Debug(ctx, "Processing assignment target", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"targetType":      *odataType,
			"resourceId":      data.ID.ValueString(),
		})

		assignmentObj := map[string]attr.Value{
			"type":     types.StringNull(),
			"group_id": types.StringNull(),
		}

		switch *odataType {
		case "#microsoft.graph.groupAssignmentTarget":
			tflog.Debug(ctx, "Mapping groupAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("groupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					tflog.Debug(ctx, "Setting group ID for group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"groupId":         *groupId,
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(ctx, "Group ID is nil/empty for group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = types.StringNull()
				}
			} else {
				tflog.Error(ctx, "Failed to cast target to GroupAssignmentTargetable", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["group_id"] = types.StringNull()
			}

		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			tflog.Debug(ctx, "Mapping exclusionGroupAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("exclusionGroupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					tflog.Debug(ctx, "Setting group ID for exclusion group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"groupId":         *groupId,
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(ctx, "Group ID is nil/empty for exclusion group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = types.StringNull()
				}
			} else {
				tflog.Error(ctx, "Failed to cast target to ExclusionGroupAssignmentTargetable", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["group_id"] = types.StringNull()
			}

		default:
			tflog.Warn(ctx, "Unknown target type encountered", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"targetType":      *odataType,
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["group_id"] = types.StringNull()
		}

		tflog.Debug(ctx, "Creating assignment object value", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		objValue, diags := types.ObjectValue(WindowsQualityUpdatePolicyAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
		if !diags.HasError() {
			tflog.Debug(ctx, "Successfully created assignment object", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentValues = append(assignmentValues, objValue)
		} else {
			tflog.Error(ctx, "Failed to create assignment object value", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"errors":          diags.Errors(),
				"resourceId":      data.ID.ValueString(),
			})
		}
	}

	tflog.Debug(ctx, "Creating assignments set", map[string]interface{}{
		"processedAssignments": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})

	if len(assignmentValues) > 0 {
		setVal, diags := types.SetValue(WindowsQualityUpdatePolicyAssignmentType(), assignmentValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]interface{}{
				"errors":     diags.Errors(),
				"resourceId": data.ID.ValueString(),
			})
			data.Assignments = types.SetNull(WindowsQualityUpdatePolicyAssignmentType())
		} else {
			tflog.Debug(ctx, "Successfully created assignments set", map[string]interface{}{
				"assignmentCount": len(assignmentValues),
				"resourceId":      data.ID.ValueString(),
			})
			data.Assignments = setVal
		}
	} else {
		tflog.Debug(ctx, "No valid assignments processed, setting assignments to null", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(WindowsQualityUpdatePolicyAssignmentType())
	}

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]interface{}{
		"finalAssignmentCount": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})
}
