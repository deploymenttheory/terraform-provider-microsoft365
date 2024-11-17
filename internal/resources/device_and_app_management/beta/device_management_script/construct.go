package graphBetaDeviceManagementScript

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *DeviceManagementScriptResourceModel) (models.DeviceManagementScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewDeviceManagementScript()

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.ScriptContent.IsNull() && !data.ScriptContent.IsUnknown() {
		encodedContent := base64.StdEncoding.EncodeToString([]byte(data.ScriptContent.ValueString()))
		scriptContent := []byte(encodedContent)
		requestBody.SetScriptContent(scriptContent)
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
			requestBody.SetRunAsAccount(runAsAccount)
		}
	}

	if !data.EnforceSignatureCheck.IsNull() && !data.EnforceSignatureCheck.IsUnknown() {
		enforceSignatureCheck := data.EnforceSignatureCheck.ValueBool()
		requestBody.SetEnforceSignatureCheck(&enforceSignatureCheck)
	}

	if !data.FileName.IsNull() && !data.FileName.IsUnknown() {
		fileName := data.FileName.ValueString()
		requestBody.SetFileName(&fileName)
	}

	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, 0, len(data.RoleScopeTagIds))
		for _, v := range data.RoleScopeTagIds {
			if !v.IsNull() && !v.IsUnknown() {
				roleScopeTagIds = append(roleScopeTagIds, v.ValueString())
			}
		}
		if len(roleScopeTagIds) > 0 {
			requestBody.SetRoleScopeTagIds(roleScopeTagIds)
		}
	}

	if !data.RunAs32Bit.IsNull() && !data.RunAs32Bit.IsUnknown() {
		runAs32Bit := data.RunAs32Bit.ValueBool()
		requestBody.SetRunAs32Bit(&runAs32Bit)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

func constructAssignments(ctx context.Context, assignments []DeviceManagementScriptAssignmentResourceModel) ([]models.DeviceManagementScriptAssignmentable, error) {
	tflog.Debug(ctx, "Constructing DeviceManagementScript assignments", map[string]interface{}{"count": len(assignments)})
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

	return constructedAssignments, nil
}

func constructGroupAssignments(ctx context.Context, groupAssignments []DeviceManagementScriptGroupAssignmentResourceModel) ([]models.DeviceManagementScriptGroupAssignmentable, error) {
	tflog.Debug(ctx, "Constructing DeviceManagementScript group assignments", map[string]interface{}{"count": len(groupAssignments)})
	var constructedGroupAssignments []models.DeviceManagementScriptGroupAssignmentable

	for _, groupAssignment := range groupAssignments {
		newGroupAssignment := models.NewDeviceManagementScriptGroupAssignment()

		if !groupAssignment.TargetGroupId.IsNull() && !groupAssignment.TargetGroupId.IsUnknown() {
			targetGroupIdStr := groupAssignment.TargetGroupId.ValueString()
			newGroupAssignment.SetTargetGroupId(&targetGroupIdStr)
		}

		constructedGroupAssignments = append(constructedGroupAssignments, newGroupAssignment)
	}

	return constructedGroupAssignments, nil
}
