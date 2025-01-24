package graphBetaApplications

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote WinGetApp resource to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ApplicationsResourceModel, remoteResource graphmodels.WinGetAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	// Handle PackageIdentifier value to support case-insensitive comparison
	if data != nil && !data.WinGetApp.PackageIdentifier.IsNull() &&
		strings.EqualFold(data.WinGetApp.PackageIdentifier.ValueString(), state.StringPtrToString(remoteResource.GetPackageIdentifier())) {
	} else {
		data.WinGetApp.PackageIdentifier = types.StringPointerValue(remoteResource.GetPackageIdentifier())
	}

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.Publisher = types.StringPointerValue(remoteResource.GetPublisher())
	data.IsFeatured = types.BoolPointerValue(remoteResource.GetIsFeatured())
	data.PrivacyInformationUrl = types.StringPointerValue(remoteResource.GetPrivacyInformationUrl())
	data.InformationUrl = types.StringPointerValue(remoteResource.GetInformationUrl())
	data.Owner = types.StringPointerValue(remoteResource.GetOwner())
	data.Developer = types.StringPointerValue(remoteResource.GetDeveloper())
	data.Notes = types.StringPointerValue(remoteResource.GetNotes())
	data.WinGetApp.ManifestHash = types.StringPointerValue(remoteResource.GetManifestHash())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.UploadState = state.Int32PtrToTypeInt64(remoteResource.GetUploadState())
	data.PublishingState = state.EnumPtrToTypeString(remoteResource.GetPublishingState())
	data.IsAssigned = types.BoolPointerValue(remoteResource.GetIsAssigned())
	data.DependentAppCount = state.Int32PtrToTypeInt64(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersededAppCount())

	// Handle InstallExperience
	if installExperience := remoteResource.GetInstallExperience(); installExperience != nil {
		data.WinGetApp.InstallExperience = &WinGetAppInstallExperienceResourceModel{
			RunAsAccount: state.EnumPtrToTypeString(installExperience.GetRunAsAccount()),
		}
	}

	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.LargeIcon = types.ObjectValueMust(
			map[string]attr.Type{
				"type":  types.StringType,
				"value": types.StringType,
			},
			map[string]attr.Value{
				"type":  types.StringValue(state.StringPtrToString(largeIcon.GetTypeEscaped())),
				"value": types.StringValue(base64.StdEncoding.EncodeToString(largeIcon.GetValue())),
			},
		)
	} else {
		data.LargeIcon = types.ObjectNull(
			map[string]attr.Type{
				"type":  types.StringType,
				"value": types.StringType,
			},
		)
	}

	var roleScopeTagIds []attr.Value
	for _, v := range state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds()) {
		roleScopeTagIds = append(roleScopeTagIds, v)
	}

	data.RoleScopeTagIds = types.ListValueMust(
		types.StringType,
		roleScopeTagIds,
	)

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
