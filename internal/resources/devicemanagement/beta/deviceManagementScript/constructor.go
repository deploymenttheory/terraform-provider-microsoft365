// objectConstruction creates a new DeviceManagementScript object from the provided data.
func objectConstruction(data deviceManagementScriptData) (*models.DeviceManagementScript, error) {
	script := models.NewDeviceManagementScript()
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
	script.SetScriptContent([]byte(detectionScriptContent))

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

// assignmentObjectConstruction creates a new DeviceManagementScriptAssignmentable object from the provided data.
func assignmentObjectConstruction(data deviceManagementScriptData) ([]models.DeviceManagementScriptAssignmentable, error) {
	if data.Assignments.IsNull() || data.Assignments.IsUnknown() {
		return nil, nil
	}

	var assignments []models.DeviceManagementScriptAssignmentable
	for _, assignment := range data.Assignments.Elements() {
		var a deviceManagementScriptAssignmentData
		diags := assignment.As(&a, nil)
		if diags.HasError() {
			return nil, fmt.Errorf("error parsing assignment data: %v", diags)
		}

		assignmentObj := models.NewDeviceManagementScriptAssignment()

		// Setting the target
		target := models.NewDeviceAndAppManagementAssignmentTarget()
		target.SetId(a.TargetGroupID.ValueString())
		assignmentObj.SetTarget(target)

		// Setting the run schedule if provided
		if !a.RunSchedule.IsNull() {
			var rs runSchedule
			diags := a.RunSchedule.As(&rs, nil)
			if diags.HasError() {
				return nil, fmt.Errorf("error parsing run schedule: %v", diags)
			}

			schedule := models.NewDeviceManagementScriptRunSchedule()
			schedule.SetInterval(rs.Interval.ValueInt64())
			schedule.SetTime(rs.Time.ValueString())
			schedule.SetUseUtc(rs.UseUtc.ValueBool())

			assignmentObj.SetRunSchedule(schedule)
		}

		assignments = append(assignments, assignmentObj)
	}

	return assignments, nil
}
