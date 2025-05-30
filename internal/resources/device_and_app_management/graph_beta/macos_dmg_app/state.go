package graphBetaMacOSDmgApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the GraphServiceClient result to the Terraform provider model
func MapRemoteStateToTerraform(ctx context.Context, data *MacOSDmgAppResourceModel, remoteResource graphmodels.MacOSDmgAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform model", map[string]interface{}{
		"resourceId": state.StringPointerValue(remoteResource.GetId()).ValueString(),
	})

	// Map base resource properties
	data.ID = state.StringPointerValue(remoteResource.GetId())
	data.DisplayName = state.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = state.StringPointerValue(remoteResource.GetDescription())
	data.Publisher = state.StringPointerValue(remoteResource.GetPublisher())
	data.Developer = state.StringPointerValue(remoteResource.GetDeveloper())
	data.Owner = state.StringPointerValue(remoteResource.GetOwner())
	data.Notes = state.StringPointerValue(remoteResource.GetNotes())
	data.InformationUrl = state.StringPointerValue(remoteResource.GetInformationUrl())
	data.PrivacyInformationUrl = state.StringPointerValue(remoteResource.GetPrivacyInformationUrl())
	data.IsFeatured = state.BoolPointerValue(remoteResource.GetIsFeatured())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.UploadState = state.Int32PointerValue(remoteResource.GetUploadState())
	data.PublishingState = state.EnumPtrToTypeString(remoteResource.GetPublishingState())
	data.IsAssigned = state.BoolPointerValue(remoteResource.GetIsAssigned())
	data.DependentAppCount = state.Int32PointerValue(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = state.Int32PointerValue(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = state.Int32PointerValue(remoteResource.GetSupersededAppCount())

	// Map role scope tag IDs
	if roleScopeTagIds := remoteResource.GetRoleScopeTagIds(); roleScopeTagIds != nil {
		roleScopeTagIdsList, diags := types.SetValueFrom(ctx, types.StringType, roleScopeTagIds)
		if diags.HasError() {
			tflog.Error(ctx, "Error mapping role scope tag IDs", map[string]interface{}{
				"error": diags.Errors(),
			})
		} else {
			data.RoleScopeTagIds = roleScopeTagIdsList
		}
	}

	// Map categories if available
	if categories := remoteResource.GetCategories(); categories != nil {
		var categoryIds []string
		for _, category := range categories {
			if category.GetId() != nil {
				categoryIds = append(categoryIds, *category.GetId())
			}
		}
		if len(categoryIds) > 0 {
			categoriesList, diags := types.SetValueFrom(ctx, types.StringType, categoryIds)
			if diags.HasError() {
				tflog.Error(ctx, "Error mapping categories", map[string]interface{}{
					"error": diags.Errors(),
				})
			} else {
				data.Categories = categoriesList
			}
		}
	}

	// Map macOS DMG app specific properties
	if data.MacOSDmgApp == nil {
		data.MacOSDmgApp = &MacOSDmgAppDetailsResourceModel{}
	}

	data.MacOSDmgApp.IgnoreVersionDetection = state.BoolPointerValue(remoteResource.GetIgnoreVersionDetection())
	data.MacOSDmgApp.PrimaryBundleId = state.StringPointerValue(remoteResource.GetPrimaryBundleId())
	data.MacOSDmgApp.PrimaryBundleVersion = state.StringPointerValue(remoteResource.GetPrimaryBundleVersion())

	// Map minimum supported operating system
	if minOS := remoteResource.GetMinimumSupportedOperatingSystem(); minOS != nil {
		if data.MacOSDmgApp.MinimumSupportedOperatingSystem == nil {
			data.MacOSDmgApp.MinimumSupportedOperatingSystem = &MacOSMinimumOperatingSystemResourceModel{}
		}
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V107 = state.BoolPointerValue(minOS.GetV107())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V108 = state.BoolPointerValue(minOS.GetV108())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V109 = state.BoolPointerValue(minOS.GetV109())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1010 = state.BoolPointerValue(minOS.GetV1010())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1011 = state.BoolPointerValue(minOS.GetV1011())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1012 = state.BoolPointerValue(minOS.GetV1012())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1013 = state.BoolPointerValue(minOS.GetV1013())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1014 = state.BoolPointerValue(minOS.GetV1014())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1015 = state.BoolPointerValue(minOS.GetV1015())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V110 = state.BoolPointerValue(minOS.GetV110())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V120 = state.BoolPointerValue(minOS.GetV120())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V130 = state.BoolPointerValue(minOS.GetV130())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V140 = state.BoolPointerValue(minOS.GetV140())
		data.MacOSDmgApp.MinimumSupportedOperatingSystem.V150 = state.BoolPointerValue(minOS.GetV150())
	}

	// Map included apps
	if includedApps := remoteResource.GetIncludedApps(); includedApps != nil {
		var includedAppElements []attr.Value
		for _, app := range includedApps {
			if app == nil {
				continue
			}

			appValues := map[string]attr.Value{
				"bundle_id":      state.StringPointerValue(app.GetBundleId()),
				"bundle_version": state.StringPointerValue(app.GetBundleVersion()),
			}

			appObj, diags := types.ObjectValue(map[string]attr.Type{
				"bundle_id":      types.StringType,
				"bundle_version": types.StringType,
			}, appValues)

			if diags.HasError() {
				tflog.Error(ctx, "Error creating included app object", map[string]interface{}{
					"error": diags.Errors(),
				})
				continue
			}

			includedAppElements = append(includedAppElements, appObj)
		}

		includedAppsSet, diags := types.SetValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"bundle_id":      types.StringType,
				"bundle_version": types.StringType,
			},
		}, includedAppElements)

		if diags.HasError() {
			tflog.Error(ctx, "Error creating included apps set", map[string]interface{}{
				"error": diags.Errors(),
			})
			data.MacOSDmgApp.IncludedApps = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"bundle_id":      types.StringType,
					"bundle_version": types.StringType,
				},
			})
		} else {
			data.MacOSDmgApp.IncludedApps = includedAppsSet
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform model")
}
