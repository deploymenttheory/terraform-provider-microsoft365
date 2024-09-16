package graphBetaWinGetApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WinGetAppResourceModel, remoteResource models.WinGetAppable) {
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
	data.Publisher = types.StringValue(state.StringPtrToString(remoteResource.GetPublisher()))
	data.IsFeatured = state.BoolPtrToTypeBool(remoteResource.GetIsFeatured())
	data.PrivacyInformationUrl = types.StringValue(state.StringPtrToString(remoteResource.GetPrivacyInformationUrl()))
	data.InformationUrl = types.StringValue(state.StringPtrToString(remoteResource.GetInformationUrl()))
	data.Owner = types.StringValue(state.StringPtrToString(remoteResource.GetOwner()))
	data.Developer = types.StringValue(state.StringPtrToString(remoteResource.GetDeveloper()))
	data.Notes = types.StringValue(state.StringPtrToString(remoteResource.GetNotes()))
	data.ManifestHash = types.StringValue(state.StringPtrToString(remoteResource.GetManifestHash()))
	data.PackageIdentifier = types.StringValue(state.StringPtrToString(remoteResource.GetPackageIdentifier()))
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.UploadState = state.Int32PtrToTypeInt64(remoteResource.GetUploadState())
	data.PublishingState = state.EnumPtrToTypeString(remoteResource.GetPublishingState())
	data.IsAssigned = state.BoolPtrToTypeBool(remoteResource.GetIsAssigned())
	data.DependentAppCount = state.Int32PtrToTypeInt64(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersededAppCount())

	// Handle LargeIcon
	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.LargeIcon = &MimeContentModel{
			Type:  types.StringValue(state.StringPtrToString(largeIcon.GetTypeEscaped())),
			Value: types.StringValue(string(largeIcon.GetValue())),
		}
	}

	// Handle InstallExperience
	if installExperience := remoteResource.GetInstallExperience(); installExperience != nil {
		data.InstallExperience = &WinGetAppInstallExperienceModel{
			RunAsAccount: state.EnumPtrToTypeString(installExperience.GetRunAsAccount()),
		}
	}

	// Handle RoleScopeTagIds
	roleScopeTags := remoteResource.GetRoleScopeTagIds()
	if len(roleScopeTags) == 0 {
		data.RoleScopeTagIds = []types.String{}
	} else {
		data.RoleScopeTagIds = make([]types.String, len(roleScopeTags))
		for i, tag := range roleScopeTags {
			data.RoleScopeTagIds[i] = types.StringValue(tag)
		}
	}

	// Handle Assignments
	if assignments := remoteResource.GetAssignments(); assignments != nil {
		data.MobileAppAssignments = make([]MobileAppAssignmentModel, len(assignments))
		for i, assignment := range assignments {
			data.MobileAppAssignments[i] = mapAssignmentToModel(assignment)
		}
	} else {
		data.MobileAppAssignments = []MobileAppAssignmentModel{}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}

func mapAssignmentToModel(assignment models.MobileAppAssignmentable) MobileAppAssignmentModel {
	model := MobileAppAssignmentModel{
		Intent: state.EnumPtrToTypeString(assignment.GetIntent()),
	}

	// Map Target
	if target := assignment.GetTarget(); target != nil {
		model.Target = mapTargetToModel(target)
	}

	// Map Settings
	if settings := assignment.GetSettings(); settings != nil {
		if winGetSettings, ok := settings.(models.WinGetAppAssignmentSettingsable); ok {
			model.Settings = WinGetAppAssignmentSettingsModel{
				Notifications: state.EnumPtrToTypeString(winGetSettings.GetNotifications()),
			}
		}
	}

	return model
}

func mapTargetToModel(target models.DeviceAndAppManagementAssignmentTargetable) AssignmentTargetModel {
	model := AssignmentTargetModel{}

	switch target.(type) {
	case models.AllLicensedUsersAssignmentTargetable:
		model.Type = types.StringValue("allLicensedUsers")
	case models.AllDevicesAssignmentTargetable:
		model.Type = types.StringValue("allDevices")
	case models.GroupAssignmentTargetable:
		model.Type = types.StringValue("group")
		if groupTarget, ok := target.(models.GroupAssignmentTargetable); ok {
			model.GroupID = types.StringValue(state.StringPtrToString(groupTarget.GetGroupId()))
		}
	}

	return model
}
