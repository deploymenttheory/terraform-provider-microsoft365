package graphBetaMacOSLobApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the properties of a MacOSLobApp to Terraform state.
//
// This function handles the conversion of Graph API model properties to Terraform state.
// It follows a direct mapping approach with proper logging.
//
// Parameters:
//   - ctx: Context for logging or cancellation
//   - data: Terraform model to populate
//   - remoteResource: Graph API model containing source data
func MapRemoteResourceStateToTerraform(ctx context.Context, data *MacOSLobAppResourceModel, remoteResource graphmodels.MacOSLobAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(remoteResource.GetPublisher())
	data.InformationUrl = convert.GraphToFrameworkString(remoteResource.GetInformationUrl())
	data.PrivacyInformationUrl = convert.GraphToFrameworkString(remoteResource.GetPrivacyInformationUrl())
	data.Owner = convert.GraphToFrameworkString(remoteResource.GetOwner())
	data.Developer = convert.GraphToFrameworkString(remoteResource.GetDeveloper())
	data.Notes = convert.GraphToFrameworkString(remoteResource.GetNotes())
	data.IsFeatured = convert.GraphToFrameworkBool(remoteResource.GetIsFeatured())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.PublishingState = convert.GraphToFrameworkEnum(remoteResource.GetPublishingState())
	data.DependentAppCount = convert.GraphToFrameworkInt32(remoteResource.GetDependentAppCount())
	data.IsAssigned = convert.GraphToFrameworkBool(remoteResource.GetIsAssigned())
	data.SupersededAppCount = convert.GraphToFrameworkInt32(remoteResource.GetSupersededAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(remoteResource.GetSupersedingAppCount())
	data.UploadState = convert.GraphToFrameworkInt32(remoteResource.GetUploadState())

	if data.AppIcon != nil {
		tflog.Debug(ctx, "Preserving original app_icon values from configuration")
	} else if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.AppIcon = &sharedmodels.MobileAppIconResourceModel{
			IconFilePathSource: types.StringNull(),
			IconURLSource:      types.StringNull(),
		}
	} else {
		data.AppIcon = nil
	}

	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	data.Categories = sharedstater.MapMobileAppCategoriesStateToTerraform(ctx, remoteResource.GetCategories())

	if data.MacOSLobApp == nil {
		data.MacOSLobApp = &MacOSLobAppDetailsResourceModel{}
	}
	mapMacOSLobAppStateToTerraform(data.MacOSLobApp, remoteResource)

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapMacOSLobAppStateToTerraform handles fields specific to macOS LOB apps
func mapMacOSLobAppStateToTerraform(data *MacOSLobAppDetailsResourceModel, remoteResource graphmodels.MacOSLobAppable) {
	if data == nil {
		data = &MacOSLobAppDetailsResourceModel{}
	}

	data.BundleId = convert.GraphToFrameworkString(remoteResource.GetBundleId())
	data.BuildNumber = convert.GraphToFrameworkString(remoteResource.GetBuildNumber())
	data.VersionNumber = convert.GraphToFrameworkString(remoteResource.GetVersionNumber())
	data.IgnoreVersionDetection = convert.GraphToFrameworkBool(remoteResource.GetIgnoreVersionDetection())
	data.InstallAsManaged = convert.GraphToFrameworkBool(remoteResource.GetInstallAsManaged())
	data.MD5HashChunkSize = convert.GraphToFrameworkInt32(remoteResource.GetMd5HashChunkSize())

	if md5Hashes := remoteResource.GetMd5Hash(); md5Hashes != nil {
		data.MD5Hash = convert.GraphToFrameworkStringList(md5Hashes)
	}

	if childApps := remoteResource.GetChildApps(); len(childApps) > 0 {
		var mappedChildApps []MacOSLobChildAppResourceModel
		for _, childApp := range childApps {
			mappedChildApp := MacOSLobChildAppResourceModel{
				BundleId:      convert.GraphToFrameworkString(childApp.GetBundleId()),
				BuildNumber:   convert.GraphToFrameworkString(childApp.GetBuildNumber()),
				VersionNumber: convert.GraphToFrameworkString(childApp.GetVersionNumber()),
			}
			mappedChildApps = append(mappedChildApps, mappedChildApp)
		}
		data.ChildApps = mappedChildApps
	}

	// Map minimum supported operating system
	if minOS := remoteResource.GetMinimumSupportedOperatingSystem(); minOS != nil {
		if data.MinimumSupportedOperatingSystem == nil {
			data.MinimumSupportedOperatingSystem = &MacOSMinimumOperatingSystemResourceModel{}
		}

		data.MinimumSupportedOperatingSystem.V107 = convert.GraphToFrameworkBool(minOS.GetV107())
		data.MinimumSupportedOperatingSystem.V108 = convert.GraphToFrameworkBool(minOS.GetV108())
		data.MinimumSupportedOperatingSystem.V109 = convert.GraphToFrameworkBool(minOS.GetV109())
		data.MinimumSupportedOperatingSystem.V1010 = convert.GraphToFrameworkBool(minOS.GetV1010())
		data.MinimumSupportedOperatingSystem.V1011 = convert.GraphToFrameworkBool(minOS.GetV1011())
		data.MinimumSupportedOperatingSystem.V1012 = convert.GraphToFrameworkBool(minOS.GetV1012())
		data.MinimumSupportedOperatingSystem.V1013 = convert.GraphToFrameworkBool(minOS.GetV1013())
		data.MinimumSupportedOperatingSystem.V1014 = convert.GraphToFrameworkBool(minOS.GetV1014())
		data.MinimumSupportedOperatingSystem.V1015 = convert.GraphToFrameworkBool(minOS.GetV1015())
		data.MinimumSupportedOperatingSystem.V110 = convert.GraphToFrameworkBool(minOS.GetV110())
		data.MinimumSupportedOperatingSystem.V120 = convert.GraphToFrameworkBool(minOS.GetV120())
		data.MinimumSupportedOperatingSystem.V130 = convert.GraphToFrameworkBool(minOS.GetV130())
		data.MinimumSupportedOperatingSystem.V140 = convert.GraphToFrameworkBool(minOS.GetV140())
		data.MinimumSupportedOperatingSystem.V150 = convert.GraphToFrameworkBool(minOS.GetV150())
	}
}
