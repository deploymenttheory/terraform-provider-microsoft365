package graphBetaWinGetApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	largeIconObj := types.ObjectNull(
		map[string]attr.Type{
			"type":  types.StringType,
			"value": types.StringType,
		},
	)

	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		largeIconObj, _ = types.ObjectValue(
			map[string]attr.Type{
				"type":  types.StringType,
				"value": types.StringType,
			},
			map[string]attr.Value{
				"type":  types.StringValue(state.StringPtrToString(largeIcon.GetTypeEscaped())),
				"value": types.StringValue(state.ByteToString(largeIcon.GetValue())),
			},
		)
	}

	data.LargeIcon = largeIconObj

	// Handle InstallExperience
	if installExperience := remoteResource.GetInstallExperience(); installExperience != nil {
		data.InstallExperience = &WinGetAppInstallExperienceModel{
			RunAsAccount: state.EnumPtrToTypeString(installExperience.GetRunAsAccount()),
		}
	}

	// Handle RoleScopeTagIds
	roleScopeTags := remoteResource.GetRoleScopeTagIds()

	if len(roleScopeTags) > 0 {
		data.RoleScopeTagIds = make([]types.String, len(roleScopeTags))
		for i, tag := range roleScopeTags {
			data.RoleScopeTagIds[i] = types.StringValue(tag)
		}
	} else {
		data.RoleScopeTagIds = nil
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
