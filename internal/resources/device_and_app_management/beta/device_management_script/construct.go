package graphbetadevicemanagementscript

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *DeviceManagementScriptResourceModel) (models.DeviceManagementScriptable, error) {
	tflog.Debug(ctx, "Constructing DeviceManagementScript resource")

	script := models.NewDeviceManagementScript()

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		script.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		script.SetDescription(&description)
	}

	if !data.ScriptContent.IsNull() && !data.ScriptContent.IsUnknown() {
		encodedContent := base64.StdEncoding.EncodeToString([]byte(data.ScriptContent.ValueString()))
		scriptContent := []byte(encodedContent)
		script.SetScriptContent(scriptContent)
	}

	if !data.RunAsAccount.IsNull() && !data.RunAsAccount.IsUnknown() {
		runAsAccountStr := data.RunAsAccount.ValueString()
		runAsAccountAny, err := models.ParseRunAsAccountType(runAsAccountStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing RunAsAccount: %v", err)
		}
		if runAsAccountAny != nil {
			runAsAccount, ok := runAsAccountAny.(*models.RunAsAccountType)
			if !ok {
				return nil, fmt.Errorf("unexpected type for RunAsAccount: %T", runAsAccountAny)
			}
			script.SetRunAsAccount(runAsAccount)
		}
	}

	if !data.EnforceSignatureCheck.IsNull() && !data.EnforceSignatureCheck.IsUnknown() {
		enforceSignatureCheck := data.EnforceSignatureCheck.ValueBool()
		script.SetEnforceSignatureCheck(&enforceSignatureCheck)
	}

	if !data.FileName.IsNull() && !data.FileName.IsUnknown() {
		fileName := data.FileName.ValueString()
		script.SetFileName(&fileName)
	}

	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, 0, len(data.RoleScopeTagIds))
		for _, v := range data.RoleScopeTagIds {
			if !v.IsNull() && !v.IsUnknown() {
				roleScopeTagIds = append(roleScopeTagIds, v.ValueString())
			}
		}
		if len(roleScopeTagIds) > 0 {
			script.SetRoleScopeTagIds(roleScopeTagIds)
		}
	}

	if !data.RunAs32Bit.IsNull() && !data.RunAs32Bit.IsUnknown() {
		runAs32Bit := data.RunAs32Bit.ValueBool()
		script.SetRunAs32Bit(&runAs32Bit)
	}

	requestBodyJSON, err := json.MarshalIndent(script, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed DeviceManagementScript resource:\n"+string(requestBodyJSON))

	return script, nil
}

func constructAssignments(ctx context.Context, assignments []DeviceManagementScriptAssignmentResourceModel) ([]models.DeviceManagementScriptAssignmentable, error) {
	var constructedAssignments []models.DeviceManagementScriptAssignmentable

	for _, assignment := range assignments {
		newAssignment := models.NewDeviceManagementScriptAssignment()
		target := models.NewDeviceAndAppManagementAssignmentTarget()

		if !assignment.Target.DeviceAndAppManagementAssignmentFilterType.IsNull() && !assignment.Target.DeviceAndAppManagementAssignmentFilterType.IsUnknown() {
			filterTypeStr := assignment.Target.DeviceAndAppManagementAssignmentFilterType.ValueString()
			filterTypeAny, err := models.ParseDeviceAndAppManagementAssignmentFilterType(filterTypeStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing filter type: %v", err)
			}
			if filterTypeAny != nil {
				filterType, ok := filterTypeAny.(*models.DeviceAndAppManagementAssignmentFilterType)
				if !ok {
					return nil, fmt.Errorf("unexpected type for filter type: %T", filterTypeAny)
				}
				target.SetDeviceAndAppManagementAssignmentFilterType(filterType)
			}
		}

		if !assignment.Target.DeviceAndAppManagementAssignmentFilterId.IsNull() && !assignment.Target.DeviceAndAppManagementAssignmentFilterId.IsUnknown() {
			filterIdStr := assignment.Target.DeviceAndAppManagementAssignmentFilterId.ValueString()
			target.SetDeviceAndAppManagementAssignmentFilterId(&filterIdStr)
		}

		// TODO - raise bug. this appears to be missing from the SDK, but is in the data model.
		// https://learn.microsoft.com/en-us/graph/api/intune-devices-devicemanagementscriptassignment-get?view=graph-rest-beta
		//if !assignment.Target.EntraObjectId.IsNull() && !assignment.Target.EntraObjectId.IsUnknown() {
		//	entraObjectIdStr := assignment.Target.EntraObjectId.ValueString()
		//	target.SetEntraObjectId(&entraObjectIdStr)
		//}

		newAssignment.SetTarget(target)
		constructedAssignments = append(constructedAssignments, newAssignment)
	}

	assignmentsJSON, err := json.MarshalIndent(constructedAssignments, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling assignments to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed DeviceManagementScript assignments:\n"+string(assignmentsJSON))

	return constructedAssignments, nil
}

func constructGroupAssignments(ctx context.Context, groupAssignments []DeviceManagementScriptGroupAssignmentResourceModel) ([]models.DeviceManagementScriptGroupAssignmentable, error) {
	var constructedGroupAssignments []models.DeviceManagementScriptGroupAssignmentable

	for _, groupAssignment := range groupAssignments {
		newGroupAssignment := models.NewDeviceManagementScriptGroupAssignment()

		if !groupAssignment.TargetGroupId.IsNull() && !groupAssignment.TargetGroupId.IsUnknown() {
			targetGroupIdStr := groupAssignment.TargetGroupId.ValueString()
			newGroupAssignment.SetTargetGroupId(&targetGroupIdStr)
		}

		constructedGroupAssignments = append(constructedGroupAssignments, newGroupAssignment)
	}

	groupAssignmentsJSON, err := json.MarshalIndent(constructedGroupAssignments, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling group assignments to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed DeviceManagementScript group assignments:\n"+string(groupAssignmentsJSON))

	return constructedGroupAssignments, nil
}
