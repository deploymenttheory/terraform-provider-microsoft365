package graphbetamacospkgapp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *MacOSPkgAppResourceModel) (models.MacOSPkgAppable, error) {
	tflog.Debug(ctx, "Constructing MacOSPkgApp resource")

	app := models.NewMacOSPkgApp()

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		app.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		app.SetDescription(&description)
	}

	if !data.Publisher.IsNull() && !data.Publisher.IsUnknown() {
		publisher := data.Publisher.ValueString()
		app.SetPublisher(&publisher)
	}

	if !data.PrivacyInformationUrl.IsNull() && !data.PrivacyInformationUrl.IsUnknown() {
		privacyUrl := data.PrivacyInformationUrl.ValueString()
		app.SetPrivacyInformationUrl(&privacyUrl)
	}

	if !data.InformationUrl.IsNull() && !data.InformationUrl.IsUnknown() {
		infoUrl := data.InformationUrl.ValueString()
		app.SetInformationUrl(&infoUrl)
	}

	if !data.Owner.IsNull() && !data.Owner.IsUnknown() {
		owner := data.Owner.ValueString()
		app.SetOwner(&owner)
	}

	if !data.Developer.IsNull() && !data.Developer.IsUnknown() {
		developer := data.Developer.ValueString()
		app.SetDeveloper(&developer)
	}

	if !data.Notes.IsNull() && !data.Notes.IsUnknown() {
		notes := data.Notes.ValueString()
		app.SetNotes(&notes)
	}

	if !data.FileName.IsNull() && !data.FileName.IsUnknown() {
		fileName := data.FileName.ValueString()
		app.SetFileName(&fileName)
	}

	if !data.PrimaryBundleId.IsNull() && !data.PrimaryBundleId.IsUnknown() {
		primaryBundleId := data.PrimaryBundleId.ValueString()
		app.SetPrimaryBundleId(&primaryBundleId)
	}

	if !data.PrimaryBundleVersion.IsNull() && !data.PrimaryBundleVersion.IsUnknown() {
		primaryBundleVersion := data.PrimaryBundleVersion.ValueString()
		app.SetPrimaryBundleVersion(&primaryBundleVersion)
	}

	if !data.IgnoreVersionDetection.IsNull() && !data.IgnoreVersionDetection.IsUnknown() {
		ignoreVersionDetection := data.IgnoreVersionDetection.ValueBool()
		app.SetIgnoreVersionDetection(&ignoreVersionDetection)
	}

	if !data.IsFeatured.IsNull() && !data.IsFeatured.IsUnknown() {
		isFeatured := data.IsFeatured.ValueBool()
		app.SetIsFeatured(&isFeatured)
	}

	if len(data.RoleScopeTagIds) > 0 {
		roleScopeTagIds := make([]string, 0, len(data.RoleScopeTagIds))
		for _, v := range data.RoleScopeTagIds {
			if !v.IsNull() && !v.IsUnknown() {
				roleScopeTagIds = append(roleScopeTagIds, v.ValueString())
			}
		}
		if len(roleScopeTagIds) > 0 {
			app.SetRoleScopeTagIds(roleScopeTagIds)
		}
	}

	if data.LargeIcon.Type != types.StringNull() && data.LargeIcon.Value != types.StringNull() {
		largeIcon := models.NewMimeContent()
		//largeIcon.SetType(data.LargeIcon.Type.ValueStringPointer()) // TODO: field not in sdk yet, but is in data model
		largeIcon.SetValue([]byte(data.LargeIcon.Value.ValueString()))
		app.SetLargeIcon(largeIcon)
	}

	if len(data.IncludedApps) > 0 {
		includedApps := make([]models.MacOSIncludedAppable, 0, len(data.IncludedApps))
		for _, v := range data.IncludedApps {
			includedApp := models.NewMacOSIncludedApp()
			includedApp.SetBundleId(v.BundleId.ValueStringPointer())
			includedApp.SetBundleVersion(v.BundleVersion.ValueStringPointer())
			includedApps = append(includedApps, includedApp)
		}
		app.SetIncludedApps(includedApps)
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
	app.SetMinimumSupportedOperatingSystem(minOS)

	if data.PreInstallScript.ScriptContent != types.StringNull() {
		preInstallScript := models.NewMacOSAppScript()
		preInstallScript.SetScriptContent(data.PreInstallScript.ScriptContent.ValueStringPointer())
		app.SetPreInstallScript(preInstallScript)
	}

	if data.PostInstallScript.ScriptContent != types.StringNull() {
		postInstallScript := models.NewMacOSAppScript()
		postInstallScript.SetScriptContent(data.PostInstallScript.ScriptContent.ValueStringPointer())
		app.SetPostInstallScript(postInstallScript)
	}

	requestBodyJSON, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed MacOSPkgApp resource:\n"+string(requestBodyJSON))

	return app, nil
}
