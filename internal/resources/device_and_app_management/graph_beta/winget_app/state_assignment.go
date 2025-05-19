package graphBetaWinGetApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// StateWinGetAppAssignments maps remote assignments to our new model structure with sets
// StateWinGetAppAssignments maps remote assignments to our new model structure with sets
func StateWinGetAppAssignments(ctx context.Context, remoteAssignmentsResponse graphmodels.MobileAppAssignmentCollectionResponseable) (*WinGetAppAssignmentsResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Initialize the assignments model
	assignments := &WinGetAppAssignmentsResourceModel{}

	// Define the object types for nested objects
	targetAttrTypes := map[string]attr.Type{
		"target_type":            types.StringType,
		"group_id":               types.StringType,
		"collection_id":          types.StringType,
		"assignment_filter_id":   types.StringType,
		"assignment_filter_type": types.StringType,
	}

	installTimeSettingsAttrTypes := map[string]attr.Type{
		"use_local_time":     types.BoolType,
		"deadline_date_time": types.StringType,
	}

	restartSettingsAttrTypes := map[string]attr.Type{
		"grace_period_in_minutes":                         types.Int32Type,
		"countdown_display_before_restart_in_minutes":     types.Int32Type,
		"restart_notification_snooze_duration_in_minutes": types.Int32Type,
	}

	settingsAttrTypes := map[string]attr.Type{
		"notifications": types.StringType,
		"install_time_settings": types.ObjectType{
			AttrTypes: installTimeSettingsAttrTypes,
		},
		"restart_settings": types.ObjectType{
			AttrTypes: restartSettingsAttrTypes,
		},
	}

	// Define the assignment object type - this is the correct element type for our sets
	assignmentObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":        types.StringType,
			"source":    types.StringType,
			"source_id": types.StringType,
			"target": types.ObjectType{
				AttrTypes: targetAttrTypes,
			},
			"settings": types.ObjectType{
				AttrTypes: settingsAttrTypes,
			},
		},
	}

	// Initialize empty sets with proper types
	// Create empty sets with the correct element type
	var setDiags diag.Diagnostics
	assignments.Required, setDiags = types.SetValue(assignmentObjectType, []attr.Value{})
	diags.Append(setDiags...)

	assignments.Available, setDiags = types.SetValue(assignmentObjectType, []attr.Value{})
	diags.Append(setDiags...)

	assignments.Uninstall, setDiags = types.SetValue(assignmentObjectType, []attr.Value{})
	diags.Append(setDiags...)

	// Return early if no assignments
	if remoteAssignmentsResponse == nil || remoteAssignmentsResponse.GetValue() == nil || len(remoteAssignmentsResponse.GetValue()) == 0 {
		tflog.Debug(ctx, "Remote assignments response is empty")
		return assignments, diags
	}

	remoteAssignments := remoteAssignmentsResponse.GetValue()

	// Prepare sets for each assignment type
	requiredAssignments := []attr.Value{}
	availableAssignments := []attr.Value{}
	uninstallAssignments := []attr.Value{}

	// Process each remote assignment
	for _, remoteAssignment := range remoteAssignments {
		// Skip if intent is missing
		if remoteAssignment.GetIntent() == nil {
			continue
		}

		// Create the assignment object
		assignmentObj, objDiags := createAssignmentObject(ctx, remoteAssignment, assignmentObjectType.AttrTypes, targetAttrTypes, settingsAttrTypes, installTimeSettingsAttrTypes, restartSettingsAttrTypes)
		diags.Append(objDiags...)
		if objDiags.HasError() {
			continue
		}

		// Add to the appropriate set based on intent
		if intent := remoteAssignment.GetIntent(); intent != nil {
			intentEnum := *intent
			intentValue := intentEnum.String()

			switch intentValue {
			case "required":
				requiredAssignments = append(requiredAssignments, assignmentObj)
			case "available":
				availableAssignments = append(availableAssignments, assignmentObj)
			case "uninstall":
				uninstallAssignments = append(uninstallAssignments, assignmentObj)
			default:
				tflog.Debug(ctx, fmt.Sprintf("Ignoring assignment with unknown intent: %s", intentValue))
			}
		}
	}

	// Create the sets only if there are assignments
	if len(requiredAssignments) > 0 {
		requiredSet, setDiags := types.SetValue(assignmentObjectType, requiredAssignments)
		diags.Append(setDiags...)
		if !setDiags.HasError() {
			assignments.Required = requiredSet
		}
	}

	if len(availableAssignments) > 0 {
		availableSet, setDiags := types.SetValue(assignmentObjectType, availableAssignments)
		diags.Append(setDiags...)
		if !setDiags.HasError() {
			assignments.Available = availableSet
		}
	}

	if len(uninstallAssignments) > 0 {
		uninstallSet, setDiags := types.SetValue(assignmentObjectType, uninstallAssignments)
		diags.Append(setDiags...)
		if !setDiags.HasError() {
			assignments.Uninstall = uninstallSet
		}
	}

	return assignments, diags
}

// createAssignmentObject creates a Terraform object from a remote assignment
func createAssignmentObject(
	ctx context.Context,
	remoteAssignment graphmodels.MobileAppAssignmentable,
	assignmentAttrTypes map[string]attr.Type,
	targetAttrTypes map[string]attr.Type,
	settingsAttrTypes map[string]attr.Type,
	installTimeSettingsAttrTypes map[string]attr.Type,
	restartSettingsAttrTypes map[string]attr.Type,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create attribute map
	attrs := map[string]attr.Value{
		"id":        state.StringPointerValue(remoteAssignment.GetId()),
		"source":    state.EnumPtrToTypeString(remoteAssignment.GetSource()),
		"source_id": state.StringPointerValue(remoteAssignment.GetSourceId()),
	}

	// Process target
	if remoteTarget := remoteAssignment.GetTarget(); remoteTarget != nil {
		targetObj, targetDiags := createTargetObject(ctx, remoteTarget, targetAttrTypes)
		diags.Append(targetDiags...)
		if !targetDiags.HasError() {
			attrs["target"] = targetObj
		}
	} else {
		// Create empty target object if missing
		attrs["target"] = types.ObjectNull(targetAttrTypes)
	}

	// Process settings
	if remoteSettings := remoteAssignment.GetSettings(); remoteSettings != nil {
		// We only handle WinGet settings in this resource
		if winGetSettings, ok := remoteSettings.(*graphmodels.WinGetAppAssignmentSettings); ok {
			settingsObj, settingsDiags := createWinGetSettingsObject(ctx, winGetSettings, settingsAttrTypes, installTimeSettingsAttrTypes, restartSettingsAttrTypes)
			diags.Append(settingsDiags...)
			if !settingsDiags.HasError() {
				attrs["settings"] = settingsObj
			}
		} else {
			// If it's not WinGet settings, create empty settings object
			attrs["settings"] = types.ObjectNull(settingsAttrTypes)
		}
	} else {
		// Create empty settings object if missing
		attrs["settings"] = types.ObjectNull(settingsAttrTypes)
	}

	// Create the object
	obj, objDiags := types.ObjectValue(assignmentAttrTypes, attrs)
	diags.Append(objDiags...)

	return obj, diags
}

// createTargetObject creates a Terraform object from a remote target
func createTargetObject(ctx context.Context, remoteTarget graphmodels.DeviceAndAppManagementAssignmentTargetable, attrTypes map[string]attr.Type) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Prepare attributes
	attrs := map[string]attr.Value{
		"target_type":            types.StringNull(),
		"group_id":               types.StringNull(),
		"collection_id":          types.StringNull(),
		"assignment_filter_id":   types.StringPointerValue(remoteTarget.GetDeviceAndAppManagementAssignmentFilterId()),
		"assignment_filter_type": state.EnumPtrToTypeString(remoteTarget.GetDeviceAndAppManagementAssignmentFilterType()),
	}

	// Determine target type and set specific attributes
	switch v := remoteTarget.(type) {
	case *graphmodels.GroupAssignmentTarget:
		attrs["target_type"] = types.StringValue("groupAssignment")
		attrs["group_id"] = types.StringPointerValue(v.GetGroupId())
	case *graphmodels.ExclusionGroupAssignmentTarget:
		attrs["target_type"] = types.StringValue("exclusionGroupAssignment")
		attrs["group_id"] = types.StringPointerValue(v.GetGroupId())
	case *graphmodels.ConfigurationManagerCollectionAssignmentTarget:
		attrs["target_type"] = types.StringValue("configurationManagerCollection")
		attrs["collection_id"] = types.StringPointerValue(v.GetCollectionId())
	case *graphmodels.AllDevicesAssignmentTarget:
		attrs["target_type"] = types.StringValue("allDevices")
	case *graphmodels.AllLicensedUsersAssignmentTarget:
		attrs["target_type"] = types.StringValue("allLicensedUsers")
	default:
		attrs["target_type"] = types.StringValue("unknown")
	}

	// Create the object
	obj, objDiags := types.ObjectValue(attrTypes, attrs)
	diags.Append(objDiags...)

	return obj, diags
}

// createWinGetSettingsObject creates a Terraform object from WinGet settings
func createWinGetSettingsObject(
	ctx context.Context,
	remoteSettings *graphmodels.WinGetAppAssignmentSettings,
	settingsAttrTypes map[string]attr.Type,
	installTimeSettingsAttrTypes map[string]attr.Type,
	restartSettingsAttrTypes map[string]attr.Type,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrs := map[string]attr.Value{
		"notifications": state.EnumPtrToTypeString(remoteSettings.GetNotifications()),
	}

	// Process install time settings
	if installTimeSettings := remoteSettings.GetInstallTimeSettings(); installTimeSettings != nil {
		installTimeAttrs := map[string]attr.Value{
			"use_local_time":     types.BoolPointerValue(installTimeSettings.GetUseLocalTime()),
			"deadline_date_time": state.TimeToString(installTimeSettings.GetDeadlineDateTime()),
		}

		installTimeObj, installTimeDiags := types.ObjectValue(installTimeSettingsAttrTypes, installTimeAttrs)
		diags.Append(installTimeDiags...)
		if !installTimeDiags.HasError() {
			attrs["install_time_settings"] = installTimeObj
		}
	} else {
		attrs["install_time_settings"] = types.ObjectNull(installTimeSettingsAttrTypes)
	}

	// Process restart settings
	if restartSettings := remoteSettings.GetRestartSettings(); restartSettings != nil {
		restartAttrs := map[string]attr.Value{
			"grace_period_in_minutes":                         state.Int32PtrToTypeInt32(restartSettings.GetGracePeriodInMinutes()),
			"countdown_display_before_restart_in_minutes":     state.Int32PtrToTypeInt32(restartSettings.GetCountdownDisplayBeforeRestartInMinutes()),
			"restart_notification_snooze_duration_in_minutes": state.Int32PtrToTypeInt32(restartSettings.GetRestartNotificationSnoozeDurationInMinutes()),
		}

		restartObj, restartDiags := types.ObjectValue(restartSettingsAttrTypes, restartAttrs)
		diags.Append(restartDiags...)
		if !restartDiags.HasError() {
			attrs["restart_settings"] = restartObj
		}
	} else {
		attrs["restart_settings"] = types.ObjectNull(restartSettingsAttrTypes)
	}

	// Create the object
	obj, objDiags := types.ObjectValue(settingsAttrTypes, attrs)
	diags.Append(objDiags...)

	return obj, diags
}
