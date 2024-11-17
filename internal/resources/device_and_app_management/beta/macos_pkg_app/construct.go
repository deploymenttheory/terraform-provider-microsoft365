package graphbetamacospkgapp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a MacOSPkgApp resource using data from the Terraform model.
func constructResource(ctx context.Context, typeName string, data *MacOSPkgAppResourceModel) (models.MacOSPkgAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", typeName))

	requestBody := models.NewMacOSPkgApp()

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.Publisher.IsNull() && !data.Publisher.IsUnknown() {
		publisher := data.Publisher.ValueString()
		requestBody.SetPublisher(&publisher)
	}

	if !data.PrivacyInformationUrl.IsNull() && !data.PrivacyInformationUrl.IsUnknown() {
		privacyUrl := data.PrivacyInformationUrl.ValueString()
		requestBody.SetPrivacyInformationUrl(&privacyUrl)
	}

	if !data.InformationUrl.IsNull() && !data.InformationUrl.IsUnknown() {
		infoUrl := data.InformationUrl.ValueString()
		requestBody.SetInformationUrl(&infoUrl)
	}

	if !data.Owner.IsNull() && !data.Owner.IsUnknown() {
		owner := data.Owner.ValueString()
		requestBody.SetOwner(&owner)
	}

	if !data.Developer.IsNull() && !data.Developer.IsUnknown() {
		developer := data.Developer.ValueString()
		requestBody.SetDeveloper(&developer)
	}

	if !data.Notes.IsNull() && !data.Notes.IsUnknown() {
		notes := data.Notes.ValueString()
		requestBody.SetNotes(&notes)
	}

	if !data.FileName.IsNull() && !data.FileName.IsUnknown() {
		fileName := data.FileName.ValueString()
		requestBody.SetFileName(&fileName)
	}

	if !data.PrimaryBundleId.IsNull() && !data.PrimaryBundleId.IsUnknown() {
		primaryBundleId := data.PrimaryBundleId.ValueString()
		requestBody.SetPrimaryBundleId(&primaryBundleId)
	}

	if !data.PrimaryBundleVersion.IsNull() && !data.PrimaryBundleVersion.IsUnknown() {
		primaryBundleVersion := data.PrimaryBundleVersion.ValueString()
		requestBody.SetPrimaryBundleVersion(&primaryBundleVersion)
	}

	if !data.IgnoreVersionDetection.IsNull() && !data.IgnoreVersionDetection.IsUnknown() {
		ignoreVersionDetection := data.IgnoreVersionDetection.ValueBool()
		requestBody.SetIgnoreVersionDetection(&ignoreVersionDetection)
	}

	if !data.IsFeatured.IsNull() && !data.IsFeatured.IsUnknown() {
		isFeatured := data.IsFeatured.ValueBool()
		requestBody.SetIsFeatured(&isFeatured)
	}

	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, 0, len(data.RoleScopeTagIds))
		for _, v := range data.RoleScopeTagIds {
			if !v.IsNull() && !v.IsUnknown() {
				roleScopeTagIds = append(roleScopeTagIds, v.ValueString())
			}
		}
		if len(roleScopeTagIds) > 0 {
			requestBody.SetRoleScopeTagIds(roleScopeTagIds)
		}
	}

	if data.LargeIcon.Type != types.StringNull() && data.LargeIcon.Value != types.StringNull() {
		largeIcon := models.NewMimeContent()
		//largeIcon.SetType(data.LargeIcon.Type.ValueStringPointer()) // TODO: field not in sdk yet, but is in data model
		largeIcon.SetValue([]byte(data.LargeIcon.Value.ValueString()))
		requestBody.SetLargeIcon(largeIcon)
	}

	if len(data.IncludedApps) > 0 {
		includedApps := make([]models.MacOSIncludedAppable, 0, len(data.IncludedApps))
		for _, v := range data.IncludedApps {
			includedApp := models.NewMacOSIncludedApp()
			includedApp.SetBundleId(v.BundleId.ValueStringPointer())
			includedApp.SetBundleVersion(v.BundleVersion.ValueStringPointer())
			includedApps = append(includedApps, includedApp)
		}
		requestBody.SetIncludedApps(includedApps)
	}

	minOS := models.NewMacOSMinimumOperatingSystem()
	minOS.SetV107(data.MinimumSupportedOperatingSystem.V10_7.ValueBoolPointer())
	minOS.SetV108(data.MinimumSupportedOperatingSystem.V10_8.ValueBoolPointer())
	minOS.SetV109(data.MinimumSupportedOperatingSystem.V10_9.ValueBoolPointer())
	minOS.SetV1010(data.MinimumSupportedOperatingSystem.V10_10.ValueBoolPointer())
	minOS.SetV1011(data.MinimumSupportedOperatingSystem.V10_11.ValueBoolPointer())
	minOS.SetV1012(data.MinimumSupportedOperatingSystem.V10_12.ValueBoolPointer())
	minOS.SetV1013(data.MinimumSupportedOperatingSystem.V10_13.ValueBoolPointer())
	minOS.SetV1014(data.MinimumSupportedOperatingSystem.V10_14.ValueBoolPointer())
	minOS.SetV1015(data.MinimumSupportedOperatingSystem.V10_15.ValueBoolPointer())
	minOS.SetV110(data.MinimumSupportedOperatingSystem.V11_0.ValueBoolPointer())
	minOS.SetV120(data.MinimumSupportedOperatingSystem.V12_0.ValueBoolPointer())
	minOS.SetV130(data.MinimumSupportedOperatingSystem.V13_0.ValueBoolPointer())
	minOS.SetV140(data.MinimumSupportedOperatingSystem.V14_0.ValueBoolPointer())
	requestBody.SetMinimumSupportedOperatingSystem(minOS)

	if data.PreInstallScript.ScriptContent != types.StringNull() {
		preInstallScript := models.NewMacOSAppScript()
		preInstallScript.SetScriptContent(data.PreInstallScript.ScriptContent.ValueStringPointer())
		requestBody.SetPreInstallScript(preInstallScript)
	}

	if data.PostInstallScript.ScriptContent != types.StringNull() {
		postInstallScript := models.NewMacOSAppScript()
		postInstallScript.SetScriptContent(data.PostInstallScript.ScriptContent.ValueStringPointer())
		requestBody.SetPostInstallScript(postInstallScript)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", typeName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", typeName))

	return requestBody, nil
}
