package graphbetadevicemanagementscript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *DeviceManagementScriptResourceModel, remoteResource models.DeviceManagementScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.ScriptContent = types.StringValue(state.ByteToString(remoteResource.GetScriptContent()))
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.RunAsAccount = state.EnumPtrToTypeString(remoteResource.GetRunAsAccount())
	data.EnforceSignatureCheck = types.BoolValue(state.BoolPtrToBool(remoteResource.GetEnforceSignatureCheck()))
	data.FileName = types.StringValue(state.StringPtrToString(remoteResource.GetFileName()))
	data.RunAs32Bit = types.BoolValue(state.BoolPtrToBool(remoteResource.GetRunAs32Bit()))

	// Handle RoleScopeTagIds
	roleScopeTagIds := remoteResource.GetRoleScopeTagIds()
	if len(roleScopeTagIds) == 0 {
		data.RoleScopeTagIds = []types.String{}
	} else {
		data.RoleScopeTagIds = make([]types.String, len(roleScopeTagIds))
		for i, tag := range roleScopeTagIds {
			data.RoleScopeTagIds[i] = types.StringValue(tag)
		}
	}

	// Handle Assignments
	assignments := remoteResource.GetAssignments()
	if len(assignments) == 0 {
		data.Assignments = []DeviceManagementScriptAssignmentResourceModel{}
	} else {
		data.Assignments = make([]DeviceManagementScriptAssignmentResourceModel, len(assignments))
		for i, assignment := range assignments {
			data.Assignments[i] = MapAssignmentsRemoteStateToTerraform(assignment)
		}
	}

	// Handle GroupAssignments
	groupAssignments := remoteResource.GetGroupAssignments()
	if len(groupAssignments) == 0 {
		data.GroupAssignments = []DeviceManagementScriptGroupAssignmentResourceModel{}
	} else {
		data.GroupAssignments = make([]DeviceManagementScriptGroupAssignmentResourceModel, len(groupAssignments))
		for i, groupAssignment := range groupAssignments {
			data.GroupAssignments[i] = MapGroupAssignmentsRemoteStateToTerraform(groupAssignment)
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

func MapAssignmentsRemoteStateToTerraform(assignment models.DeviceManagementScriptAssignmentable) DeviceManagementScriptAssignmentResourceModel {
	return DeviceManagementScriptAssignmentResourceModel{
		ID: types.StringValue(state.StringPtrToString(assignment.GetId())),
		Target: Target{
			DeviceAndAppManagementAssignmentFilterId:   types.StringValue(state.StringPtrToString(assignment.GetTarget().GetDeviceAndAppManagementAssignmentFilterId())),
			DeviceAndAppManagementAssignmentFilterType: state.EnumPtrToTypeString(assignment.GetTarget().GetDeviceAndAppManagementAssignmentFilterType()),
			TargetType: types.StringValue(state.StringPtrToString(assignment.GetTarget().GetOdataType())),
			// TODO - field is currently missing from the msft SDK
			//EntraObjectId: types.StringValue(state.StringPtrToString(assignment.GetTarget().GetGroupId())),
		},
	}
}

func MapGroupAssignmentsRemoteStateToTerraform(groupAssignment models.DeviceManagementScriptGroupAssignmentable) DeviceManagementScriptGroupAssignmentResourceModel {
	return DeviceManagementScriptGroupAssignmentResourceModel{
		ID:            types.StringValue(state.StringPtrToString(groupAssignment.GetId())),
		TargetGroupId: types.StringValue(state.StringPtrToString(groupAssignment.GetTargetGroupId())),
	}
}
