package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
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

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = state.StringPointerValue(remoteResource.GetId())
	data.DisplayName = state.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = state.StringPointerValue(remoteResource.GetDescription())
	data.Publisher = state.StringPointerValue(remoteResource.GetPublisher())
	data.InformationUrl = state.StringPointerValue(remoteResource.GetInformationUrl())
	data.PrivacyInformationUrl = state.StringPointerValue(remoteResource.GetPrivacyInformationUrl())
	data.Owner = state.StringPointerValue(remoteResource.GetOwner())
	data.Developer = state.StringPointerValue(remoteResource.GetDeveloper())
	data.Notes = state.StringPointerValue(remoteResource.GetNotes())
	data.IsFeatured = state.BoolPointerValue(remoteResource.GetIsFeatured())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	//data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.PublishingState = state.EnumPtrToTypeString(remoteResource.GetPublishingState())
	data.DependentAppCount = state.Int32PointerValue(remoteResource.GetDependentAppCount())
	data.IsAssigned = state.BoolPointerValue(remoteResource.GetIsAssigned())
	data.SupersededAppCount = state.Int32PointerValue(remoteResource.GetSupersededAppCount())
	data.SupersedingAppCount = state.Int32PointerValue(remoteResource.GetSupersedingAppCount())
	data.UploadState = state.Int32PointerValue(remoteResource.GetUploadState())

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

	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())

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
	mapMacOSPKGAppStateToTerraform(ctx, data.MacOSPkgApp, remoteResource)

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}

// mapMacOSPKGAppStateToTerraform handle fields specific to macOs pkgs
func mapMacOSPKGAppStateToTerraform(ctx context.Context, data *MacOSPKGAppDetailsResourceModel, remoteResource graphmodels.MacOSPkgAppable) {
	if data == nil {
		data = &MacOSPKGAppDetailsResourceModel{}
	}

	data.PrimaryBundleId = state.StringPointerValue(remoteResource.GetPrimaryBundleId())
	data.PrimaryBundleVersion = state.StringPointerValue(remoteResource.GetPrimaryBundleVersion())
	data.IgnoreVersionDetection = state.BoolPointerValue(remoteResource.GetIgnoreVersionDetection())

	apps := remoteResource.GetIncludedApps()
	data.IncludedApps = state.BuildObjectSetFromSlice(
		ctx,
		map[string]attr.Type{
			"bundle_id":      types.StringType,
			"bundle_version": types.StringType,
		},
		func(i int) map[string]attr.Value {
			app := apps[i]
			return map[string]attr.Value{
				"bundle_id":      types.StringPointerValue(app.GetBundleId()),
				"bundle_version": types.StringPointerValue(app.GetBundleVersion()),
			}
		},
		len(apps),
	)

	if minOS := remoteResource.GetMinimumSupportedOperatingSystem(); minOS != nil {
		if data.MinimumSupportedOperatingSystem == nil {
			data.MinimumSupportedOperatingSystem = &MacOSMinimumOperatingSystemResourceModel{}
		}

		data.MinimumSupportedOperatingSystem.V107 = state.BoolPointerValue(minOS.GetV107())
		data.MinimumSupportedOperatingSystem.V108 = state.BoolPointerValue(minOS.GetV108())
		data.MinimumSupportedOperatingSystem.V109 = state.BoolPointerValue(minOS.GetV109())
		data.MinimumSupportedOperatingSystem.V1010 = state.BoolPointerValue(minOS.GetV1010())
		data.MinimumSupportedOperatingSystem.V1011 = state.BoolPointerValue(minOS.GetV1011())
		data.MinimumSupportedOperatingSystem.V1012 = state.BoolPointerValue(minOS.GetV1012())
		data.MinimumSupportedOperatingSystem.V1013 = state.BoolPointerValue(minOS.GetV1013())
		data.MinimumSupportedOperatingSystem.V1014 = state.BoolPointerValue(minOS.GetV1014())
		data.MinimumSupportedOperatingSystem.V1015 = state.BoolPointerValue(minOS.GetV1015())
		data.MinimumSupportedOperatingSystem.V110 = state.BoolPointerValue(minOS.GetV110())
		data.MinimumSupportedOperatingSystem.V120 = state.BoolPointerValue(minOS.GetV120())
		data.MinimumSupportedOperatingSystem.V130 = state.BoolPointerValue(minOS.GetV130())
		data.MinimumSupportedOperatingSystem.V140 = state.BoolPointerValue(minOS.GetV140())
		data.MinimumSupportedOperatingSystem.V150 = state.BoolPointerValue(minOS.GetV150())
	}

	if preScript := remoteResource.GetPreInstallScript(); preScript != nil {
		if data.PreInstallScript == nil {
			data.PreInstallScript = &MacOSAppScriptResourceModel{}
		}
		data.PreInstallScript.ScriptContent = state.StringPointerValue(preScript.GetScriptContent())
	}

	if postScript := remoteResource.GetPostInstallScript(); postScript != nil {
		if data.PostInstallScript == nil {
			data.PostInstallScript = &MacOSAppScriptResourceModel{}
		}
		data.PostInstallScript.ScriptContent = state.StringPointerValue(postScript.GetScriptContent())
	}

}
