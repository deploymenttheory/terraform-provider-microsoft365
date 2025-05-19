package graphBetaWinGetApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructWinGetAppAssignments constructs and returns a MobileAppsItemAssignPostRequestBody
// for the new assignments structure with three sets
func ConstructWinGetAppAssignments(ctx context.Context, data *WinGetAppAssignmentsResourceModel) (deviceappmanagement.MobileAppsItemAssignPostRequestBodyable, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create the request body
	requestBody := deviceappmanagement.NewMobileAppsItemAssignPostRequestBody()
	var allAssignments []graphmodels.MobileAppAssignmentable

	// Check if assignments block exists
	if data == nil {
		tflog.Debug(ctx, "No assignments provided, will clear all existing assignments")
		requestBody.SetMobileAppAssignments([]graphmodels.MobileAppAssignmentable{})
		return requestBody, diags
	}

	tflog.Debug(ctx, "Starting WinGet app assignment construction")

	// Build required assignments
	if !data.Required.IsNull() && !data.Required.IsUnknown() {
		requiredAssignments, reqDiags := constructAssignmentSet(ctx, data.Required, "required")
		diags.Append(reqDiags...)
		allAssignments = append(allAssignments, requiredAssignments...)
	}

	// Build available assignments
	if !data.Available.IsNull() && !data.Available.IsUnknown() {
		availableAssignments, availDiags := constructAssignmentSet(ctx, data.Available, "available")
		diags.Append(availDiags...)
		allAssignments = append(allAssignments, availableAssignments...)
	}

	// Build uninstall assignments
	if !data.Uninstall.IsNull() && !data.Uninstall.IsUnknown() {
		uninstallAssignments, uninstDiags := constructAssignmentSet(ctx, data.Uninstall, "uninstall")
		diags.Append(uninstDiags...)
		allAssignments = append(allAssignments, uninstallAssignments...)
	}

	// Set all assignments in the request body
	requestBody.SetMobileAppAssignments(allAssignments)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, diags
}

// constructAssignmentSet processes a set of assignments with the same intent
func constructAssignmentSet(ctx context.Context, assignmentSet types.Set, intentType string) ([]graphmodels.MobileAppAssignmentable, diag.Diagnostics) {
	var diags diag.Diagnostics
	var assignments []graphmodels.MobileAppAssignmentable

	// Exit early if the set is empty
	if assignmentSet.IsNull() || assignmentSet.IsUnknown() || len(assignmentSet.Elements()) == 0 {
		return assignments, diags
	}

	// Process each element in the set
	for _, element := range assignmentSet.Elements() {
		// Convert element to object
		obj, ok := element.(types.Object)
		if !ok {
			diags.AddError(
				"Error Processing Assignment Set",
				fmt.Sprintf("Expected types.Object but got %T", element),
			)
			continue
		}

		// Build a single assignment
		assignment, asmtDiags := constructSingleAssignment(ctx, obj, intentType)
		diags.Append(asmtDiags...)
		if asmtDiags.HasError() {
			continue
		}

		assignments = append(assignments, assignment)
	}

	return assignments, diags
}

// constructSingleAssignment creates a single MobileAppAssignmentable from a types.Object
func constructSingleAssignment(ctx context.Context, assignmentObj types.Object, intentType string) (graphmodels.MobileAppAssignmentable, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create a new assignment
	assignment := graphmodels.NewMobileAppAssignment()

	// Set the intent based on the set type (required, available, uninstall)
	intentValue, err := graphmodels.ParseInstallIntent(intentType)
	if err != nil {
		diags.AddError(
			"Invalid Intent Type",
			fmt.Sprintf("Error parsing intent type '%s': %v", intentType, err),
		)
		return nil, diags
	}

	intent, ok := intentValue.(*graphmodels.InstallIntent)
	if !ok || intent == nil {
		diags.AddError(
			"Invalid Intent Type",
			fmt.Sprintf("Failed to convert intent value for type '%s'", intentType),
		)
		return nil, diags
	}

	assignment.SetIntent(intent)

	// Get all attributes from the object
	attrs := assignmentObj.Attributes()

	// Set ID if it exists (for updates)
	if idAttr, ok := attrs["id"].(types.String); ok && !idAttr.IsNull() {
		idStr := idAttr.ValueString()
		assignment.SetId(&idStr)
	}

	// Set source (required field)
	if sourceAttr, ok := attrs["source"].(types.String); ok && !sourceAttr.IsNull() {
		sourceStr := sourceAttr.ValueString()
		sourceEnum, err := graphmodels.ParseDeviceAndAppManagementAssignmentSource(sourceStr)
		if err != nil {
			diags.AddError(
				"Invalid Source Value",
				fmt.Sprintf("Error parsing source value '%s': %v", sourceStr, err),
			)
			return nil, diags
		}
		assignment.SetSource(sourceEnum.(*graphmodels.DeviceAndAppManagementAssignmentSource))
	} else {
		diags.AddError(
			"Missing Required Attribute",
			"source is a required attribute for assignments",
		)
		return nil, diags
	}

	// Set source_id if present
	if sourceIdAttr, ok := attrs["source_id"].(types.String); ok && !sourceIdAttr.IsNull() {
		sourceIdStr := sourceIdAttr.ValueString()
		assignment.SetSourceId(&sourceIdStr)
	}

	// Process target (required field)
	if targetObj, ok := attrs["target"].(types.Object); ok && !targetObj.IsNull() {
		target, targetDiags := constructAssignmentTarget(ctx, targetObj)
		diags.Append(targetDiags...)
		if targetDiags.HasError() {
			return nil, diags
		}
		assignment.SetTarget(target)
	} else {
		diags.AddError(
			"Missing Required Attribute",
			"target is a required attribute for assignments",
		)
		return nil, diags
	}

	// Process settings if present
	if settingsObj, ok := attrs["settings"].(types.Object); ok && !settingsObj.IsNull() {
		settings, settingsDiags := constructWinGetAssignmentSettings(ctx, settingsObj)
		diags.Append(settingsDiags...)
		if !settingsDiags.HasError() && settings != nil {
			assignment.SetSettings(settings)
		}
	}

	return assignment, diags
}

// constructAssignmentTarget creates a DeviceAndAppManagementAssignmentTargetable from a types.Object
func constructAssignmentTarget(ctx context.Context, targetObj types.Object) (graphmodels.DeviceAndAppManagementAssignmentTargetable, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Get all attributes
	attrs := targetObj.Attributes()

	// Get target_type (required)
	targetTypeAttr, ok := attrs["target_type"].(types.String)
	if !ok || targetTypeAttr.IsNull() {
		diags.AddError(
			"Missing Required Attribute",
			"target_type is a required attribute for assignment targets",
		)
		return nil, diags
	}

	targetType := targetTypeAttr.ValueString()
	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	// Create target based on type
	switch targetType {
	case "allDevices":
		target = graphmodels.NewAllDevicesAssignmentTarget()

	case "allLicensedUsers":
		target = graphmodels.NewAllLicensedUsersAssignmentTarget()

	case "groupAssignment":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		if groupIdAttr, ok := attrs["group_id"].(types.String); ok && !groupIdAttr.IsNull() {
			groupId := groupIdAttr.ValueString()
			groupTarget.SetGroupId(&groupId)
		} else {
			diags.AddError(
				"Missing Required Attribute",
				"group_id is required for groupAssignment target type",
			)
			return nil, diags
		}
		target = groupTarget

	case "exclusionGroupAssignment":
		exclusionTarget := graphmodels.NewExclusionGroupAssignmentTarget()
		if groupIdAttr, ok := attrs["group_id"].(types.String); ok && !groupIdAttr.IsNull() {
			groupId := groupIdAttr.ValueString()
			exclusionTarget.SetGroupId(&groupId)
		} else {
			diags.AddError(
				"Missing Required Attribute",
				"group_id is required for exclusionGroupAssignment target type",
			)
			return nil, diags
		}
		target = exclusionTarget

	case "configurationManagerCollection":
		configTarget := graphmodels.NewConfigurationManagerCollectionAssignmentTarget()
		if collectionIdAttr, ok := attrs["collection_id"].(types.String); ok && !collectionIdAttr.IsNull() {
			collectionId := collectionIdAttr.ValueString()
			configTarget.SetCollectionId(&collectionId)
		} else {
			diags.AddError(
				"Missing Required Attribute",
				"collection_id is required for configurationManagerCollection target type",
			)
			return nil, diags
		}
		target = configTarget

	default:
		diags.AddError(
			"Invalid Target Type",
			fmt.Sprintf("Invalid target_type value: %s", targetType),
		)
		return nil, diags
	}

	// Set filter attributes if present
	if filterIdAttr, ok := attrs["assignment_filter_id"].(types.String); ok && !filterIdAttr.IsNull() {
		filterId := filterIdAttr.ValueString()
		target.SetDeviceAndAppManagementAssignmentFilterId(&filterId)
	}

	if filterTypeAttr, ok := attrs["assignment_filter_type"].(types.String); ok && !filterTypeAttr.IsNull() && filterTypeAttr.ValueString() != "none" {
		filterTypeStr := filterTypeAttr.ValueString()
		filterType, err := graphmodels.ParseDeviceAndAppManagementAssignmentFilterType(filterTypeStr)
		if err != nil {
			diags.AddError(
				"Invalid Filter Type",
				fmt.Sprintf("Error parsing filter type '%s': %v", filterTypeStr, err),
			)
			return nil, diags
		}
		target.SetDeviceAndAppManagementAssignmentFilterType(filterType.(*graphmodels.DeviceAndAppManagementAssignmentFilterType))
	}

	return target, diags
}

// constructWinGetAssignmentSettings creates a WinGetAppAssignmentSettings from a types.Object
// constructWinGetAssignmentSettings creates a WinGetAppAssignmentSettings from a types.Object
func constructWinGetAssignmentSettings(ctx context.Context, settingsObj types.Object) (graphmodels.MobileAppAssignmentSettingsable, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Exit early if settings is null
	if settingsObj.IsNull() {
		return nil, diags
	}

	// Get all attributes
	attrs := settingsObj.Attributes()

	// Create WinGet settings
	settings := graphmodels.NewWinGetAppAssignmentSettings()

	// Set notifications if present
	if notificationsAttr, ok := attrs["notifications"].(types.String); ok && !notificationsAttr.IsNull() {
		err := constructors.SetEnumProperty(notificationsAttr, graphmodels.ParseWinGetAppNotification, settings.SetNotifications)
		if err != nil {
			diags.AddError(
				"Invalid Notification Setting",
				fmt.Sprintf("Error setting notification setting: %v", err),
			)
		}
	}

	// Process install_time_settings if present
	if installTimeObj, ok := attrs["install_time_settings"].(types.Object); ok && !installTimeObj.IsNull() {
		installTimeAttrs := installTimeObj.Attributes()
		installSettings := graphmodels.NewWinGetAppInstallTimeSettings()

		// Set use_local_time
		if useLocalTimeAttr, ok := installTimeAttrs["use_local_time"].(types.Bool); ok {
			constructors.SetBoolProperty(useLocalTimeAttr, installSettings.SetUseLocalTime)
		}

		// Set deadline_date_time
		if deadlineAttr, ok := installTimeAttrs["deadline_date_time"].(types.String); ok {
			err := constructors.StringToTime(deadlineAttr, installSettings.SetDeadlineDateTime)
			if err != nil {
				diags.AddError(
					"Invalid Deadline DateTime",
					fmt.Sprintf("Error parsing deadline date time: %v", err),
				)
			}
		}

		settings.SetInstallTimeSettings(installSettings)
	}

	// Process restart_settings if present
	if restartObj, ok := attrs["restart_settings"].(types.Object); ok && !restartObj.IsNull() {
		restartAttrs := restartObj.Attributes()
		restartSettings := graphmodels.NewWinGetAppRestartSettings()

		// Set grace_period_in_minutes
		if gracePeriodAttr, ok := restartAttrs["grace_period_in_minutes"].(types.Int32); ok {
			constructors.SetInt32Property(gracePeriodAttr, restartSettings.SetGracePeriodInMinutes)
		}

		// Set countdown_display_before_restart_in_minutes
		if countdownAttr, ok := restartAttrs["countdown_display_before_restart_in_minutes"].(types.Int32); ok {
			constructors.SetInt32Property(countdownAttr, restartSettings.SetCountdownDisplayBeforeRestartInMinutes)
		}

		// Set restart_notification_snooze_duration_in_minutes
		if snoozeDurationAttr, ok := restartAttrs["restart_notification_snooze_duration_in_minutes"].(types.Int32); ok {
			constructors.SetInt32Property(snoozeDurationAttr, restartSettings.SetRestartNotificationSnoozeDurationInMinutes)
		}

		settings.SetRestartSettings(restartSettings)
	}

	return settings, diags
}
