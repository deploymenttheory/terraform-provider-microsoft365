package graphBetaApplications

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	utility "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/macos_pkg"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructMacOSPkgAppResource(ctx context.Context, data *MacOSPkgAppResourceModel, baseApp graphmodels.MacOSPkgAppable) (graphmodels.MacOSPkgAppable, error) {
	if data.PackageInstallerFileSource.IsNull() || data.PackageInstallerFileSource.ValueString() == "" {
		return nil, fmt.Errorf("package_installer_file_source is required but not provided")
	}

	// Define fields to extract from Info.plist files
	fields := []utility.Field{
		{Key: "CFBundleIdentifier", Required: true},
		{Key: "CFBundleShortVersionString", Required: true},
	}

	// Extract fields from all Info.plist files in the package
	extractedFields, err := utility.ExtractFieldsFromFiles(
		ctx,
		data.PackageInstallerFileSource.ValueString(),
		"Info.plist",
		fields,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata from pkg file: %w", err)
	}

	// First Info.plist becomes primary bundle
	primaryBundleId := extractedFields[0].Values["CFBundleIdentifier"]
	primaryBundleVersion := extractedFields[0].Values["CFBundleShortVersionString"]

	constructors.SetStringProperty(types.StringValue(primaryBundleId), baseApp.SetPrimaryBundleId)
	constructors.SetStringProperty(types.StringValue(primaryBundleVersion), baseApp.SetPrimaryBundleVersion)

	// All entries are set as included apps (including primary)
	var includedApps []graphmodels.MacOSIncludedAppable
	for _, fields := range extractedFields {
		includedApp := graphmodels.NewMacOSIncludedApp()
		constructors.SetStringProperty(
			types.StringValue(fields.Values["CFBundleIdentifier"]),
			includedApp.SetBundleId,
		)
		constructors.SetStringProperty(
			types.StringValue(fields.Values["CFBundleShortVersionString"]),
			includedApp.SetBundleVersion,
		)
		includedApps = append(includedApps, includedApp)
	}

	baseApp.SetIncludedApps(includedApps)

	constructors.SetBoolProperty(data.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)

	if data.MinimumSupportedOperatingSystem != nil {
		minOS := graphmodels.NewMacOSMinimumOperatingSystem()
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
		baseApp.SetMinimumSupportedOperatingSystem(minOS)
	}

	if data.PreInstallScript != nil {
		preScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PreInstallScript.ScriptContent, preScript.SetScriptContent)
		baseApp.SetPreInstallScript(preScript)
	}

	if data.PostInstallScript != nil {
		postScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PostInstallScript.ScriptContent, postScript.SetScriptContent)
		baseApp.SetPostInstallScript(postScript)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully constructed MacOSPkgApp resource with %d included apps", len(includedApps)))
	return baseApp, nil
}
