package graphBetaDeviceEnrollmentNotificationConfiguration

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func mapRemoteStateToTerraform(ctx context.Context, data *DeviceEnrollmentNotificationConfigurationResourceModel, remoteResource models.DeviceEnrollmentConfigurationable) {
	if remoteResource == nil {
		return
	}

	if notificationConfig, ok := remoteResource.(models.DeviceEnrollmentNotificationConfigurationable); ok {
		data.ID = state.StringPointerValue(notificationConfig.GetId())
		data.DisplayName = state.StringPointerValue(notificationConfig.GetDisplayName())
		data.Description = state.StringPointerValue(notificationConfig.GetDescription())
		data.Priority = state.Int32PointerValue(notificationConfig.GetPriority())
		data.CreatedDateTime = state.TimeToString(notificationConfig.GetCreatedDateTime())
		data.LastModifiedDateTime = state.TimeToString(notificationConfig.GetLastModifiedDateTime())
		data.Version = state.Int32PointerValue(notificationConfig.GetVersion())

		if configType := notificationConfig.GetDeviceEnrollmentConfigurationType(); configType != nil {
			data.DeviceEnrollmentConfigurationType = state.EnumPtrToTypeString(configType)
		}

		if platformType := notificationConfig.GetPlatformType(); platformType != nil {
			data.PlatformType = state.EnumPtrToTypeString(platformType)
		}

		if data.TemplateTypes.IsNull() || data.TemplateTypes.IsUnknown() {
			if templateType := notificationConfig.GetTemplateType(); templateType != nil {
				data.TemplateTypes = state.StringSliceToSet(ctx, []string{templateType.String()})
			}
		}

		if brandingOptions := notificationConfig.GetBrandingOptions(); brandingOptions != nil {
			brandingString := brandingOptions.String()
			var brandingStrings []string

			if brandingString != "" && brandingString != "none" {
				parts := strings.Split(brandingString, ",")
				for _, part := range parts {
					trimmed := strings.TrimSpace(part)
					if trimmed != "" && trimmed != "none" {
						brandingStrings = append(brandingStrings, trimmed)
					}
				}
				// Sort for stable diffs
				sort.Strings(brandingStrings)
				data.BrandingOptions = state.StringSliceToSet(ctx, brandingStrings)
			}
		}

		if data.NotificationMessageTemplateId.IsNull() || data.NotificationMessageTemplateId.IsUnknown() {
			if templateId := notificationConfig.GetNotificationMessageTemplateId(); templateId != nil {
				data.NotificationMessageTemplateId = types.StringValue(templateId.String())
			}
		}

		if data.NotificationTemplates.IsNull() || data.NotificationTemplates.IsUnknown() {
			if notificationTemplates := notificationConfig.GetNotificationTemplates(); notificationTemplates != nil {
				data.NotificationTemplates = state.StringSliceToSet(ctx, notificationTemplates)
			}
		}

		if roleScopeTagIds := notificationConfig.GetRoleScopeTagIds(); roleScopeTagIds != nil {
			data.RoleScopeTagIds = state.StringSliceToSet(ctx, roleScopeTagIds)
		}
	}
}

func mapAssignmentsToState(ctx context.Context, assignments []models.EnrollmentConfigurationAssignmentable, model *DeviceEnrollmentNotificationConfigurationResourceModel) {
	// Direct print to stdout for debugging
	fmt.Printf("DIRECT DEBUG: Entering mapAssignmentsToState with %d assignments\n", len(assignments))

	// Add multiple debug statements at the beginning
	tflog.Debug(ctx, "==== ENTERING mapAssignmentsToState function ====")
	tflog.Debug(ctx, fmt.Sprintf("Assignments length: %d", len(assignments)))
	tflog.Debug(ctx, fmt.Sprintf("Model pointer: %p", model))
	tflog.Debug(ctx, "Starting mapAssignmentsToState with assignments")

	if len(assignments) == 0 {
		fmt.Println("DIRECT DEBUG: No assignments found, setting to nil")
		tflog.Debug(ctx, "No assignments found, setting Assignments to nil")
		model.Assignments = nil
		return
	}

	// Sort assignments by type to ensure consistent state
	var groupAssignments []models.EnrollmentConfigurationAssignmentable
	var allDevicesAssignments []models.EnrollmentConfigurationAssignmentable
	var allLicensedUsersAssignments []models.EnrollmentConfigurationAssignmentable

	// Categorize assignments by target type using OData type
	for i, assignment := range assignments {
		if assignment == nil {
			tflog.Debug(ctx, fmt.Sprintf("Assignment %d is nil, skipping", i))
			continue
		}

		target := assignment.GetTarget()
		if target == nil {
			tflog.Debug(ctx, fmt.Sprintf("Assignment %d target is nil, skipping", i))
			continue
		}

		odataType := ""
		if target.GetOdataType() != nil {
			odataType = *target.GetOdataType()
		}
		tflog.Debug(ctx, fmt.Sprintf("Assignment %d target OData type: %s", i, odataType))

		switch {
		case strings.Contains(odataType, "groupAssignmentTarget"):
			groupAssignments = append(groupAssignments, assignment)
		case strings.Contains(odataType, "allDevicesAssignmentTarget"):
			allDevicesAssignments = append(allDevicesAssignments, assignment)
		case strings.Contains(odataType, "allLicensedUsersAssignmentTarget"):
			allLicensedUsersAssignments = append(allLicensedUsersAssignments, assignment)
		default:
			tflog.Warn(ctx, fmt.Sprintf("Unsupported assignment target type: %s for assignment %d", odataType, i))
		}
	}

	// Create state assignments
	stateAssignments := make([]AssignmentModel, 0)

	// Add group assignments
	for i, assignment := range groupAssignments {
		target := assignment.GetTarget()
		groupTarget, ok := target.(models.GroupAssignmentTargetable)
		if !ok {
			tflog.Warn(ctx, fmt.Sprintf("Group assignment %d could not be cast to GroupAssignmentTargetable", i))
			continue
		}

		groupID := groupTarget.GetGroupId()
		if groupID == nil {
			tflog.Warn(ctx, fmt.Sprintf("Group assignment %d has nil group ID", i))
			continue
		}

		tflog.Debug(ctx, fmt.Sprintf("Processing group assignment %d with group ID: %s", i, *groupID))

		stateTarget := AssignmentTargetModel{
			TargetType:                               types.StringValue("group"),
			GroupId:                                  types.StringValue(*groupID),
			DeviceAndAppManagementAssignmentFilterId: types.StringValue(""),
			DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
		}

		// Set filter values if they exist
		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()

		if filterID != nil && *filterID != "" {
			stateTarget.DeviceAndAppManagementAssignmentFilterId = types.StringValue(*filterID)
			tflog.Debug(ctx, fmt.Sprintf("Group assignment %d has filter ID: %s", i, *filterID))
		}

		if filterType != nil {
			stateTarget.DeviceAndAppManagementAssignmentFilterType = types.StringValue(filterType.String())
			tflog.Debug(ctx, fmt.Sprintf("Group assignment %d has filter type: %s", i, filterType.String()))
		}

		stateAssignments = append(stateAssignments, AssignmentModel{
			Target: &stateTarget,
		})

		tflog.Debug(ctx, fmt.Sprintf("Added group assignment %d to state", i))
	}

	// Add allDevices assignments
	for i, assignment := range allDevicesAssignments {
		target := assignment.GetTarget()

		tflog.Debug(ctx, fmt.Sprintf("Processing allDevices assignment %d", i))

		stateTarget := AssignmentTargetModel{
			TargetType:                                 types.StringValue("allDevices"),
			DeviceAndAppManagementAssignmentFilterId:   types.StringValue(""),
			DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
		}

		// Set filter values if they exist
		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()

		if filterID != nil && *filterID != "" {
			stateTarget.DeviceAndAppManagementAssignmentFilterId = types.StringValue(*filterID)
			tflog.Debug(ctx, fmt.Sprintf("AllDevices assignment %d has filter ID: %s", i, *filterID))
		}

		if filterType != nil {
			stateTarget.DeviceAndAppManagementAssignmentFilterType = types.StringValue(filterType.String())
			tflog.Debug(ctx, fmt.Sprintf("AllDevices assignment %d has filter type: %s", i, filterType.String()))
		}

		stateAssignments = append(stateAssignments, AssignmentModel{
			Target: &stateTarget,
		})

		tflog.Debug(ctx, fmt.Sprintf("Added allDevices assignment %d to state", i))
	}

	// Add allLicensedUsers assignments
	for i, assignment := range allLicensedUsersAssignments {
		target := assignment.GetTarget()

		tflog.Debug(ctx, fmt.Sprintf("Processing allLicensedUsers assignment %d", i))

		stateTarget := AssignmentTargetModel{
			TargetType:                                 types.StringValue("allLicensedUsers"),
			DeviceAndAppManagementAssignmentFilterId:   types.StringValue(""),
			DeviceAndAppManagementAssignmentFilterType: types.StringValue("none"),
		}

		// Set filter values if they exist
		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()

		if filterID != nil && *filterID != "" {
			stateTarget.DeviceAndAppManagementAssignmentFilterId = types.StringValue(*filterID)
			tflog.Debug(ctx, fmt.Sprintf("AllLicensedUsers assignment %d has filter ID: %s", i, *filterID))
		}

		if filterType != nil {
			stateTarget.DeviceAndAppManagementAssignmentFilterType = types.StringValue(filterType.String())
			tflog.Debug(ctx, fmt.Sprintf("AllLicensedUsers assignment %d has filter type: %s", i, filterType.String()))
		}

		stateAssignments = append(stateAssignments, AssignmentModel{
			Target: &stateTarget,
		})

		tflog.Debug(ctx, fmt.Sprintf("Added allLicensedUsers assignment %d to state", i))
	}

	// Ensure consistent order of assignments for stable diffs
	if len(stateAssignments) > 0 {
		// Sort by target type first, then by group ID if applicable
		sort.Slice(stateAssignments, func(i, j int) bool {
			// First compare by target type
			if stateAssignments[i].Target.TargetType.ValueString() != stateAssignments[j].Target.TargetType.ValueString() {
				return stateAssignments[i].Target.TargetType.ValueString() < stateAssignments[j].Target.TargetType.ValueString()
			}

			// If target types are the same and they're group targets, compare by group ID
			if stateAssignments[i].Target.TargetType.ValueString() == "group" {
				return stateAssignments[i].Target.GroupId.ValueString() < stateAssignments[j].Target.GroupId.ValueString()
			}

			return false
		})
	}

	model.Assignments = stateAssignments
	fmt.Printf("DIRECT DEBUG: Completed mapAssignmentsToState with %d assignments in state\n", len(stateAssignments))
	tflog.Debug(ctx, fmt.Sprintf("Completed mapAssignmentsToState with %d assignments in state", len(stateAssignments)))
}
