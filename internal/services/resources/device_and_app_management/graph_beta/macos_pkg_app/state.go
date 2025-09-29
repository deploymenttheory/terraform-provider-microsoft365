package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the properties of a MacOSPkgApp to Terraform state.
//
// This function handles the conversion of Graph API model properties to Terraform state.
// It follows a direct mapping approach with proper logging.
//
// Parameters:
//   - ctx: Context for logging or cancellation
//   - data: Terraform model to populate
//   - remoteResource: Graph API model containing source data
func MapRemoteResourceStateToTerraform(ctx context.Context, data *MacOSPKGAppResourceModel, remoteResource graphmodels.MacOSPkgAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
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
	//data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
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

	// Initialize Relationships as null list since we don't currently map this field
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

	if data.MacOSPkgApp == nil {
		data.MacOSPkgApp = &MacOSPKGAppDetailsResourceModel{}
	}
	mapMacOSPKGAppStateToTerraform(data.MacOSPkgApp, remoteResource)

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}

// mapMacOSPKGAppStateToTerraform handle fields specific to macOs pkgs
func mapMacOSPKGAppStateToTerraform(data *MacOSPKGAppDetailsResourceModel, remoteResource graphmodels.MacOSPkgAppable) {
	if data == nil {
		data = &MacOSPKGAppDetailsResourceModel{}
	}

	data.PrimaryBundleId = convert.GraphToFrameworkString(remoteResource.GetPrimaryBundleId())
	data.PrimaryBundleVersion = convert.GraphToFrameworkString(remoteResource.GetPrimaryBundleVersion())
	data.IgnoreVersionDetection = convert.GraphToFrameworkBool(remoteResource.GetIgnoreVersionDetection())

	apps := remoteResource.GetIncludedApps()
	data.IncludedApps = buildObjectSetFromSlice(
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

	if preScript := remoteResource.GetPreInstallScript(); preScript != nil {
		if data.PreInstallScript == nil {
			data.PreInstallScript = &MacOSAppScriptResourceModel{}
		}
		data.PreInstallScript.ScriptContent = convert.GraphToFrameworkString(preScript.GetScriptContent())
	}

	if postScript := remoteResource.GetPostInstallScript(); postScript != nil {
		if data.PostInstallScript == nil {
			data.PostInstallScript = &MacOSAppScriptResourceModel{}
		}
		data.PostInstallScript.ScriptContent = convert.GraphToFrameworkString(postScript.GetScriptContent())
	}
}

// buildObjectSetFromSlice is a helper function to build a set of objects from a slice
func buildObjectSetFromSlice(attrTypes map[string]attr.Type, valueFunc func(int) map[string]attr.Value, length int) types.Set {
	if length == 0 {
		return types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}

	values := make([]attr.Value, length)
	for i := 0; i < length; i++ {
		obj, _ := types.ObjectValue(attrTypes, valueFunc(i))
		values[i] = obj
	}

	result, _ := types.SetValue(types.ObjectType{AttrTypes: attrTypes}, values)
	return result
}
