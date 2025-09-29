package graphBetaWinGetApp

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote WinGetApp resource to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WinGetAppResourceModel, remoteResource graphmodels.WinGetAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// Handle PackageIdentifier value to support case-insensitive comparison
	if data != nil && !data.PackageIdentifier.IsNull() &&
		strings.EqualFold(data.PackageIdentifier.ValueString(), convert.GraphToFrameworkString(remoteResource.GetPackageIdentifier()).ValueString()) {
	} else {
		data.PackageIdentifier = convert.GraphToFrameworkString(remoteResource.GetPackageIdentifier())
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(remoteResource.GetPublisher())
	data.IsFeatured = convert.GraphToFrameworkBool(remoteResource.GetIsFeatured())
	data.PrivacyInformationUrl = convert.GraphToFrameworkString(remoteResource.GetPrivacyInformationUrl())
	data.InformationUrl = convert.GraphToFrameworkString(remoteResource.GetInformationUrl())
	data.Owner = convert.GraphToFrameworkString(remoteResource.GetOwner())
	data.Developer = convert.GraphToFrameworkString(remoteResource.GetDeveloper())
	data.Notes = convert.GraphToFrameworkString(remoteResource.GetNotes())
	data.ManifestHash = convert.GraphToFrameworkString(remoteResource.GetManifestHash())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.UploadState = convert.GraphToFrameworkInt32(remoteResource.GetUploadState())
	data.PublishingState = convert.GraphToFrameworkEnum(remoteResource.GetPublishingState())
	data.IsAssigned = convert.GraphToFrameworkBool(remoteResource.GetIsAssigned())
	data.DependentAppCount = convert.GraphToFrameworkInt32(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = convert.GraphToFrameworkInt32(remoteResource.GetSupersededAppCount())

	if installExperience := remoteResource.GetInstallExperience(); installExperience != nil {
		data.InstallExperience = &WinGetAppInstallExperienceResourceModel{
			RunAsAccount: convert.GraphToFrameworkEnum(installExperience.GetRunAsAccount()),
		}
	}

	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.LargeIcon = types.ObjectValueMust(
			map[string]attr.Type{
				"type":  types.StringType,
				"value": types.StringType,
			},
			map[string]attr.Value{
				"type":  convert.GraphToFrameworkString(largeIcon.GetTypeEscaped()),
				"value": convert.GraphToFrameworkBytes(largeIcon.GetValue()),
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

	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	data.Categories = sharedstater.MapMobileAppCategoriesStateToTerraform(ctx, remoteResource.GetCategories())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
