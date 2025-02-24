package graphBetaMacOSPKGApp

import (
	"context"
	"encoding/base64"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
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
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.PublishingState = state.EnumPtrToTypeString(remoteResource.GetPublishingState())

	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.LargeIcon = &LargeIconResourceModel{
			Type:  state.StringValue(state.StringPtrToString(largeIcon.GetTypeEscaped())),
			Value: state.StringValue(base64.StdEncoding.EncodeToString(largeIcon.GetValue())),
		}
	} else {
		data.LargeIcon = nil
	}

	var roleScopeTagIds []attr.Value
	for _, v := range state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds()) {
		roleScopeTagIds = append(roleScopeTagIds, v)
	}
	data.RoleScopeTagIds = types.ListValueMust(types.StringType, roleScopeTagIds)

	categories := remoteResource.GetCategories()
	if len(categories) > 0 {
		categoriesValues := make([]MobileAppCategoryResourceModel, len(categories))
		for i, category := range categories {
			categoriesValues[i] = MobileAppCategoryResourceModel{
				ID:          state.StringPointerValue(category.GetId()),
				DisplayName: state.StringPointerValue(category.GetDisplayName()),
			}
		}
		data.Categories = categoriesValues
	}

	// Initialize the MacOSPkgApp struct if it's nil
	if data.MacOSPkgApp == nil {
		data.MacOSPkgApp = &MacOSPkgAppResourceModel{}
	}

	data.MacOSPkgApp.PrimaryBundleId = state.StringPointerValue(remoteResource.GetPrimaryBundleId())
	data.MacOSPkgApp.PrimaryBundleVersion = state.StringPointerValue(remoteResource.GetPrimaryBundleVersion())
	data.MacOSPkgApp.IgnoreVersionDetection = state.BoolPointerValue(remoteResource.GetIgnoreVersionDetection())

	includedApps := remoteResource.GetIncludedApps()
	if len(includedApps) > 0 {
		includedAppsValues := make([]MacOSIncludedAppResourceModel, len(includedApps))
		for i, app := range includedApps {
			includedAppsValues[i] = MacOSIncludedAppResourceModel{
				BundleId:      state.StringPointerValue(app.GetBundleId()),
				BundleVersion: state.StringPointerValue(app.GetBundleVersion()),
			}
		}
		data.MacOSPkgApp.IncludedApps = includedAppsValues
	}

	if minOS := remoteResource.GetMinimumSupportedOperatingSystem(); minOS != nil {
		if data.MacOSPkgApp.MinimumSupportedOperatingSystem == nil {
			data.MacOSPkgApp.MinimumSupportedOperatingSystem = &MacOSMinimumOperatingSystemResourceModel{}
		}

		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V107 = state.BoolPointerValue(minOS.GetV107())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V108 = state.BoolPointerValue(minOS.GetV108())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V109 = state.BoolPointerValue(minOS.GetV109())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1010 = state.BoolPointerValue(minOS.GetV1010())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1011 = state.BoolPointerValue(minOS.GetV1011())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1012 = state.BoolPointerValue(minOS.GetV1012())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1013 = state.BoolPointerValue(minOS.GetV1013())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1014 = state.BoolPointerValue(minOS.GetV1014())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1015 = state.BoolPointerValue(minOS.GetV1015())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V110 = state.BoolPointerValue(minOS.GetV110())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V120 = state.BoolPointerValue(minOS.GetV120())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V130 = state.BoolPointerValue(minOS.GetV130())
		data.MacOSPkgApp.MinimumSupportedOperatingSystem.V140 = state.BoolPointerValue(minOS.GetV140())
	}

	if preScript := remoteResource.GetPreInstallScript(); preScript != nil {
		if data.MacOSPkgApp.PreInstallScript == nil {
			data.MacOSPkgApp.PreInstallScript = &MacOSAppScriptResourceModel{}
		}

		if scriptContent := preScript.GetScriptContent(); scriptContent != nil {
			data.MacOSPkgApp.PreInstallScript.ScriptContent = state.DecodeBase64ToString(ctx, *scriptContent)
		}
	}

	if postScript := remoteResource.GetPostInstallScript(); postScript != nil {
		if data.MacOSPkgApp.PostInstallScript == nil {
			data.MacOSPkgApp.PostInstallScript = &MacOSAppScriptResourceModel{}
		}

		if scriptContent := postScript.GetScriptContent(); scriptContent != nil {
			data.MacOSPkgApp.PostInstallScript.ScriptContent = state.DecodeBase64ToString(ctx, *scriptContent)
		}
	}

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId":        data.ID.ValueString(),
		"displayName":       data.DisplayName.ValueString(),
		"includedAppsCount": len(data.MacOSPkgApp.IncludedApps),
	})
}
