package graphBetaApplications

import (
	"bytes"
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	archiver "github.com/mholt/archives"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"howett.net/plist"
)

func constructMacOSPkgAppResource(ctx context.Context, data *MacOSPkgAppResourceModel, baseApp graphmodels.MacOSPkgAppable) (graphmodels.MacOSPkgAppable, error) {
	// Extract metadata if file source provided
	if !data.PackageInstallerFileSource.IsNull() {
		bundleId, bundleVersion, includedApps, err := extractPkgMetadata(data.PackageInstallerFileSource.ValueString())
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("Failed to extract pkg metadata: %v", err))
		} else {
			// Set values from extracted metadata if not already set
			if data.PrimaryBundleId.IsNull() {
				data.PrimaryBundleId = types.StringValue(bundleId)
			}
			if data.PrimaryBundleVersion.IsNull() {
				data.PrimaryBundleVersion = types.StringValue(bundleVersion)
			}
			if len(data.IncludedApps) == 0 {
				data.IncludedApps = includedApps
			}
		}
	}

	// Set base properties
	constructors.SetBoolProperty(data.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)
	constructors.SetStringProperty(data.PrimaryBundleId, baseApp.SetPrimaryBundleId)
	constructors.SetStringProperty(data.PrimaryBundleVersion, baseApp.SetPrimaryBundleVersion)

	// Validate and set included apps
	if len(data.IncludedApps) > 500 {
		return nil, fmt.Errorf("included_apps exceeds maximum limit of 500")
	}

	if len(data.IncludedApps) > 0 {
		includedApps := make([]graphmodels.MacOSIncludedAppable, 0, len(data.IncludedApps))
		for _, app := range data.IncludedApps {
			includedApp := graphmodels.NewMacOSIncludedApp()
			constructors.SetStringProperty(app.BundleId, includedApp.SetBundleId)
			constructors.SetStringProperty(app.BundleVersion, includedApp.SetBundleVersion)
			includedApps = append(includedApps, includedApp)
		}
		baseApp.SetIncludedApps(includedApps)
	}

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

	// Set pre/post install scripts if provided
	if data.PreInstallScript != nil {
		script := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PreInstallScript.ScriptContent, script.SetScriptContent)
		baseApp.SetPreInstallScript(script)
	}

	if data.PostInstallScript != nil {
		script := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PostInstallScript.ScriptContent, script.SetScriptContent)
		baseApp.SetPostInstallScript(script)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s MacOSPkgApp resource values", ResourceName))

	return baseApp, nil
}

func extractPkgMetadata(filePath string) (bundleId string, bundleVersion string, includedApps []MacOSIncludedAppResourceModel, err error) {
	// Open pkg file as archive
	archive := archiver.NewPkgArchive()
	err = archive.Open(filePath)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to open pkg archive: %v", err)
	}
	defer archive.Close()

	// Extract Contents/Info.plist
	var plistData []byte
	err = archive.Extract("Contents/Info.plist", &plistData)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to extract Info.plist: %v", err)
	}

	// Parse plist data
	var data struct {
		CFBundleIdentifier         string `plist:"CFBundleIdentifier"`
		CFBundleShortVersionString string `plist:"CFBundleShortVersionString"`
	}

	decoder := plist.NewDecoder(bytes.NewReader(plistData))
	if err := decoder.Decode(&data); err != nil {
		return "", "", nil, fmt.Errorf("failed to decode plist: %v", err)
	}

	return data.CFBundleIdentifier,
		data.CFBundleShortVersionString,
		[]MacOSIncludedAppResourceModel{{
			BundleId:      types.StringValue(data.CFBundleIdentifier),
			BundleVersion: types.StringValue(data.CFBundleShortVersionString),
		}},
		nil
}
