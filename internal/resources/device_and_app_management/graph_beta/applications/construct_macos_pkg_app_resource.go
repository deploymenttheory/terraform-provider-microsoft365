package graphBetaApplications

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	pkg "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructMacOSPkgAppResource(ctx context.Context, data *MacOSPkgAppResourceModel, baseApp graphmodels.MacOSPkgAppable) (graphmodels.MacOSPkgAppable, error) {

	if data.PackageInstallerFileSource.IsNull() || data.PackageInstallerFileSource.ValueString() == "" {
		return nil, fmt.Errorf("package_installer_file_source is required but not provided")
	}

	var includedApps []MacOSIncludedAppResourceModel
	if !data.PackageInstallerFileSource.IsNull() {
		bundleId, bundleVersion, bundledApps, err := pkg.ExtractmacOSPkgMetadata(ctx, data.PackageInstallerFileSource.ValueString())
		if err != nil {
			return nil, fmt.Errorf("failed to extract metadata from the provided .pkg file at '%s': %w. Ensure the file path is correct and accessible", data.PackageInstallerFileSource.ValueString(), err)
		}

		// Populate metadata if successfully extracted
		if len(bundledApps) > 0 {
			// Convert types during append
			for _, app := range bundledApps {
				includedApps = append(includedApps, MacOSIncludedAppResourceModel{
					BundleId:      app.BundleId,
					BundleVersion: app.BundleVersion,
				})
			}

			if data.PrimaryBundleId.IsNull() {
				data.PrimaryBundleId = types.StringValue(bundleId)
			}
			if data.PrimaryBundleVersion.IsNull() {
				data.PrimaryBundleVersion = types.StringValue(bundleVersion)
			}
		}
	}

	// Validate the included apps
	if len(includedApps) == 0 && len(data.IncludedApps) == 0 {
		return nil, fmt.Errorf("no valid IncludedApps found in metadata or provided manually; at least one app is required")
	}

	// Merge user-provided IncludedApps with extracted metadata
	if len(data.IncludedApps) > 0 {
		includedApps = append(includedApps, data.IncludedApps...)
	}

	// Validate final IncludedApps list size
	if len(includedApps) > 500 {
		return nil, fmt.Errorf("too many IncludedApps provided; maximum limit is 500")
	}

	// Set IncludedApps in the base app model
	if len(includedApps) > 0 {
		graphIncludedApps := make([]graphmodels.MacOSIncludedAppable, 0, len(includedApps))
		for _, app := range includedApps {
			includedApp := graphmodels.NewMacOSIncludedApp()
			constructors.SetStringProperty(app.BundleId, includedApp.SetBundleId)
			constructors.SetStringProperty(app.BundleVersion, includedApp.SetBundleVersion)
			graphIncludedApps = append(graphIncludedApps, includedApp)
		}
		baseApp.SetIncludedApps(graphIncludedApps)
	}

	// Set additional base properties
	constructors.SetBoolProperty(data.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)
	constructors.SetStringProperty(data.PrimaryBundleId, baseApp.SetPrimaryBundleId)
	constructors.SetStringProperty(data.PrimaryBundleVersion, baseApp.SetPrimaryBundleVersion)

	// Set minimum OS requirements if provided
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

	// Set pre-install script if provided
	if data.PreInstallScript != nil {
		preScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PreInstallScript.ScriptContent, preScript.SetScriptContent)
		baseApp.SetPreInstallScript(preScript)
	}

	// Set post-install script if provided
	if data.PostInstallScript != nil {
		postScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PostInstallScript.ScriptContent, postScript.SetScriptContent)
		baseApp.SetPostInstallScript(postScript)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully constructed MacOSPkgApp resource values with %d IncludedApps", len(includedApps)))
	return baseApp, nil
}
