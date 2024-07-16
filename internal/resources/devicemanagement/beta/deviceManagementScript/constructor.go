package deviceManagementScript

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-m365/internal/resources/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// objectConstruction creates a new DeviceHealthScript object from the provided data.
func objectConstruction(data deviceManagementScriptData) (*models.DeviceHealthScript, error) {
	script := models.NewDeviceHealthScript()
	displayName := data.Name.ValueString()
	script.SetDisplayName(&displayName)

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		script.SetDescription(&description)
	}

	detectionScriptContent, err := helpers.Base64Encode(data.DetectionScriptContent.ValueString())
	if err != nil {
		return nil, fmt.Errorf("failed to encode detection script content: %v", err)
	}
	script.SetDetectionScriptContent([]byte(detectionScriptContent))

	remediationScriptContent, err := helpers.Base64Encode(data.RemediationScriptContent.ValueString())
	if err != nil {
		return nil, fmt.Errorf("failed to encode remediation script content: %v", err)
	}
	script.SetRemediationScriptContent([]byte(remediationScriptContent))
	if !data.Publisher.IsNull() {
		publisher := data.Publisher.ValueString()
		script.SetPublisher(&publisher)
	}

	if !data.RunAsAccount.IsNull() {
		runAsAccount, err := models.ParseRunAsAccountType(data.RunAsAccount.ValueString())
		if err != nil || runAsAccount == nil {
			return nil, fmt.Errorf("invalid RunAsAccount value: got %q, should be one of %q or %q", data.RunAsAccount.ValueString(), models.SYSTEM_RUNASACCOUNTTYPE.String(), models.USER_RUNASACCOUNTTYPE.String())
		}
		script.SetRunAsAccount(runAsAccount.(*models.RunAsAccountType))
	}

	if !data.RoleScopeTagIds.IsNull() {
		script.SetRoleScopeTagIds(expandStringList(data.RoleScopeTagIds))
	}

	if !data.RunAs32Bit.IsNull() {
		runAs32Bit := data.RunAs32Bit.ValueBool()
		script.SetRunAs32Bit(&runAs32Bit)
	}

	if !data.EnforceSignatureCheck.IsNull() {
		enforceSignatureCheck := data.EnforceSignatureCheck.ValueBool()
		script.SetEnforceSignatureCheck(&enforceSignatureCheck)
	}

	return script, nil
}

// assignmentObjectConstruction creates a new DeviceHealthScriptAssignmentable object from the provided data.
func assignmentObjectConstruction(data deviceManagementScriptData) ([]models.DeviceHealthScriptAssignmentable, error) {
	if data.Assignments.IsNull() || len(data.Assignments.Elements()) == 0 {
		return nil, nil
	}

	var assignments []models.DeviceHealthScriptAssignmentable
	for _, assignmentElem := range data.Assignments.Elements() {
		assignmentMap := assignmentElem.(types.Object).Attrs

		// Create a new DeviceHealthScriptAssignment instance
		deviceHealthScriptAssignment := models.NewDeviceHealthScriptAssignment()

		// Set target
		target := models.NewGroupAssignmentTarget()
		targetGroupID := assignmentMap["target_group_id"].(types.String).ValueString()
		target.SetGroupId(targetGroupID)
		deviceHealthScriptAssignment.SetTarget(target)

		// Set runRemediationScript
		runRemediationScript := assignmentMap["run_remediation_script"].(types.Bool).ValueBool()
		deviceHealthScriptAssignment.SetRunRemediationScript(&runRemediationScript)

		// Set runSchedule
		runScheduleMap := assignmentMap["run_schedule"].(types.Object).Attrs
		runSchedule := models.NewDeviceHealthScriptDailySchedule()

		interval := runScheduleMap["interval"].(types.Int64).ValueInt64()
		runSchedule.SetInterval(&interval)

		time := runScheduleMap["time"].(types.String).ValueString()
		runSchedule.SetTime(&time)

		useUtc := runScheduleMap["use_utc"].(types.Bool).ValueBool()
		runSchedule.SetUseUtc(&useUtc)

		deviceHealthScriptAssignment.SetRunSchedule(runSchedule)

		assignments = append(assignments, deviceHealthScriptAssignment)
	}

	return assignments, nil
}
