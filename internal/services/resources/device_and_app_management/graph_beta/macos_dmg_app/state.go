package graphBetaMacOSDmgApp

import (
	"context"
	"fmt"

	attribute "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/attr"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the properties of a MacOSDmgApp to Terraform state.
//
// This function handles the conversion of Graph API model properties to Terraform state.
// It follows a direct mapping approach with proper logging.
//
// Parameters:
//   - ctx: Context for logging or cancellation
//   - data: Terraform model to populate
//   - remoteResource: Graph API model containing source data
func MapRemoteResourceStateToTerraform(ctx context.Context, data *MacOSDmgAppResourceModel, remoteResource graphmodels.MacOSDmgAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]interface{}{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
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

	// Handle AppIcon - preserve original configuration values
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

	// Handle collection fields
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.Categories = sharedstater.MapMobileAppCategoriesStateToTerraform(ctx, remoteResource.GetCategories())

	// Set fields that are not currently mapped to null
	data.Relationships = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                            types.StringType,
			"source_display_name":           types.StringType,
			"source_display_version":        types.StringType,
			"source_id":                     types.StringType,
			"source_publisher_display_name": types.StringType,
			"target_display_name":           types.StringType,
			"target_display_version":        types.StringType,
			"target_id":                     types.StringType,
			"target_publisher":              types.StringType,
			"target_publisher_display_name": types.StringType,
			"target_type":                   types.StringType,
		},
	})
	data.AppInstaller = types.ObjectNull(map[string]attr.Type{})
	data.ContentVersion = types.ListNull(types.StringType)

	// Handle nested MacOSDmgApp object
	if data.MacOSDmgApp == nil {
		data.MacOSDmgApp = &MacOSDmgAppDetailsResourceModel{}
	}
	mapMacOSDMGAppStateToTerraform(ctx, data.MacOSDmgApp, remoteResource)

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapMacOSDMGAppStateToTerraform handles fields specific to macOS DMG apps
func mapMacOSDMGAppStateToTerraform(ctx context.Context, data *MacOSDmgAppDetailsResourceModel, remoteResource graphmodels.MacOSDmgAppable) {
	if data == nil {
		data = &MacOSDmgAppDetailsResourceModel{}
	}

	data.PrimaryBundleId = convert.GraphToFrameworkString(remoteResource.GetPrimaryBundleId())
	data.PrimaryBundleVersion = convert.GraphToFrameworkString(remoteResource.GetPrimaryBundleVersion())
	data.IgnoreVersionDetection = convert.GraphToFrameworkBool(remoteResource.GetIgnoreVersionDetection())

	apps := remoteResource.GetIncludedApps()
	data.IncludedApps = attribute.ObjectSetFromSlice(
		ctx,
		map[string]attr.Type{
			"bundle_id":      types.StringType,
			"bundle_version": types.StringType,
		},
		func(i int) map[string]attr.Value {
			app := apps[i]
			return map[string]attr.Value{
				"bundle_id":      convert.GraphToFrameworkString(app.GetBundleId()),
				"bundle_version": convert.GraphToFrameworkString(app.GetBundleVersion()),
			}
		},
		len(apps),
	)

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
