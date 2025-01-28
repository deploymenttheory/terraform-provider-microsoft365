package graphBetaApplications

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"howett.net/plist"
)

func constructMacOSPkgAppResource(ctx context.Context, data *MacOSPkgAppResourceModel, baseApp graphmodels.MacOSPkgAppable) (graphmodels.MacOSPkgAppable, error) {
	// Ensure a valid file path is provided and extract metadata
	if data.PackageInstallerFileSource.IsNull() || data.PackageInstallerFileSource.ValueString() == "" {
		return nil, fmt.Errorf("package_installer_file_source is required but not provided")
	}

	var includedApps []MacOSIncludedAppResourceModel
	if !data.PackageInstallerFileSource.IsNull() {
		// Attempt to extract metadata from the provided .pkg file path
		bundleId, bundleVersion, extractedApps, err := extractmacOSPkgMetadata(ctx, data.PackageInstallerFileSource.ValueString())
		if err != nil {
			return nil, fmt.Errorf("failed to extract metadata from the provided .pkg file at '%s': %w. Ensure the file path is correct and accessible", data.PackageInstallerFileSource.ValueString(), err)
		}

		// Populate metadata if successfully extracted
		if len(extractedApps) > 0 {
			includedApps = append(includedApps, extractedApps...)
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

// extractmacOSPkgMetadata extracts metadata from the `PackageInfo` file within a macOS `.pkg` file.
// Specifically, it retrieves the `CFBundleShortVersionString`, `CFBundleVersion`, and application bundle details
// for apps installed in the `/Applications` folder.
// This function ensures compliance with Intune app requirements for detection and monitoring,
// as described in the official documentation:
// https://learn.microsoft.com/en-us/mem/intune/apps/lob-apps-macos-dmg#step-3--detection-rules
//
// Example PackageInfo file:
// <?xml version="1.0" encoding="utf-8"?>
// <pkg-info overwrite-permissions="true" relocatable="false" identifier="org.mozilla.firefox" postinstall-action="none" version="134.0.0" format-version="2" generator-version="InstallCmds-681 (18F132)" install-location="/Applications" auth="root">
//
//	<payload numberOfFiles="161" installKBytes="410855"/>
//	<bundle path="./Firefox.app" id="org.mozilla.firefox" CFBundleShortVersionString="134.0" CFBundleVersion="13424.12.30"/>
//	<bundle-version>
//	    <bundle id="org.mozilla.firefox"/>
//	</bundle-version>
//	<upgrade-bundle>
//	    <bundle id="org.mozilla.firefox"/>
//	</upgrade-bundle>
//	<update-bundle/>
//	<atomic-update-bundle/>
//	<strict-identifier>
//	    <bundle id="org.mozilla.firefox"/>
//	</strict-identifier>
//	<relocate>
//	    <bundle id="org.mozilla.firefox"/>
//	</relocate>
//
// </pkg-info>
//
// In this example:
// - The `CFBundleShortVersionString` is "134.0"
// - The `CFBundleVersion` is "13424.12.30"
// - The `id` for the included app is "org.mozilla.firefox"
// These details are extracted and used to construct a valid list of included apps, ensuring compliance with monitoring and detection rules.
// extractmacOSPkgMetadata extracts metadata from the `PackageInfo` file within a macOS `.pkg` file.
func extractmacOSPkgMetadata(ctx context.Context, filePath string) (bundleId string, bundleVersion string, includedApps []MacOSIncludedAppResourceModel, err error) {
	// Create a temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "pkg_extract_*")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	tflog.Debug(ctx, "Created temporary directory for pkg extraction", map[string]interface{}{
		"tempDir": tempDir,
	})

	// First attempt: Try xar extraction
	cmd := exec.Command("xar", "-xf", filePath, "-C", tempDir)
	err = cmd.Run()
	if err != nil {
		tflog.Debug(ctx, "xar extraction failed, attempting pkgutil", map[string]interface{}{
			"error": err.Error(),
		})

		// Second attempt: Try pkgutil expansion
		expandedDir := filepath.Join(tempDir, "expanded")
		cmd = exec.Command("pkgutil", "--expand", filePath, expandedDir)
		if err := cmd.Run(); err != nil {
			// Third attempt: Try using 7zip
			cmd = exec.Command("7z", "x", filePath, "-o"+tempDir)
			if err := cmd.Run(); err != nil {
				return "", "", nil, fmt.Errorf("failed to extract pkg file using any available method: %w", err)
			}
		}
	}

	// Log the extracted structure for debugging
	if err := listFiles(ctx, tempDir); err != nil {
		tflog.Warn(ctx, "Failed to list extracted files", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Look for PackageInfo in common locations
	packageInfoPaths := []string{
		filepath.Join(tempDir, "PackageInfo"),
		filepath.Join(tempDir, "expanded", "PackageInfo"),
		// Add more common paths if needed
	}

	var packageInfoData []byte
	var foundPath string
	for _, path := range packageInfoPaths {
		if data, err := os.ReadFile(path); err == nil {
			packageInfoData = data
			foundPath = path
			break
		}
	}

	if packageInfoData == nil {
		return "", "", nil, fmt.Errorf("PackageInfo not found in extracted pkg file")
	}

	tflog.Debug(ctx, "Found PackageInfo file", map[string]interface{}{
		"path": foundPath,
		"size": len(packageInfoData),
	})

	// Define the package info structure
	var pkgInfo struct {
		InstallLocation string `plist:"install-location"`
		Bundles         []struct {
			Path                       string `plist:"path"`
			Id                         string `plist:"id"`
			CFBundleShortVersionString string `plist:"CFBundleShortVersionString"`
			CFBundleVersion            string `plist:"CFBundleVersion"`
		} `plist:"bundle"`
	}

	// Decode the PackageInfo file
	decoder := plist.NewDecoder(bytes.NewReader(packageInfoData))
	if err := decoder.Decode(&pkgInfo); err != nil {
		return "", "", nil, fmt.Errorf("failed to decode PackageInfo plist: %w", err)
	}

	// Process bundles and populate included apps
	for _, bundle := range pkgInfo.Bundles {
		if bundle.Path != "" {
			includedApps = append(includedApps, MacOSIncludedAppResourceModel{
				BundleId:      types.StringValue(bundle.Id),
				BundleVersion: types.StringValue(bundle.CFBundleShortVersionString),
			})
		}
	}

	if len(includedApps) > 0 {
		bundleId = includedApps[0].BundleId.ValueString()
		bundleVersion = includedApps[0].BundleVersion.ValueString()
	} else {
		return "", "", nil, fmt.Errorf("no valid included apps found in PackageInfo")
	}

	tflog.Debug(ctx, "Successfully extracted metadata", map[string]interface{}{
		"bundleId":      bundleId,
		"bundleVersion": bundleVersion,
		"includedApps":  len(includedApps),
	})

	return bundleId, bundleVersion, includedApps, nil
}

// Utility function to list all files in a directory (for debugging)
func listFiles(ctx context.Context, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		tflog.Debug(ctx, fmt.Sprintf("Found file: %s", rel))
		return nil
	})
}
