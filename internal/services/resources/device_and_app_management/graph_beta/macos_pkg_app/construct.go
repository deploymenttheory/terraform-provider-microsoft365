package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	download "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/common"
	utility "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/macos_pkg"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model, using the provided installer path
func constructResource(ctx context.Context, data *MacOSPKGAppResourceModel, installerSourcePath string) (graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	baseApp := graphmodels.NewMacOSPkgApp()

	convert.FrameworkToGraphString(data.Description, baseApp.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, baseApp.SetPublisher)
	convert.FrameworkToGraphString(data.DisplayName, baseApp.SetDisplayName)
	convert.FrameworkToGraphString(data.InformationUrl, baseApp.SetInformationUrl)
	convert.FrameworkToGraphBool(data.IsFeatured, baseApp.SetIsFeatured)
	convert.FrameworkToGraphString(data.Owner, baseApp.SetOwner)
	convert.FrameworkToGraphString(data.Developer, baseApp.SetDeveloper)
	convert.FrameworkToGraphString(data.Notes, baseApp.SetNotes)
	convert.FrameworkToGraphString(data.PrivacyInformationUrl, baseApp.SetPrivacyInformationUrl)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, baseApp.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Handle app icon (either from file path or web source)
	if data.AppIcon != nil {
		largeIcon := graphmodels.NewMimeContent()
		iconType := "image/png"
		largeIcon.SetTypeEscaped(&iconType)

		if !data.AppIcon.IconFilePathSource.IsNull() && data.AppIcon.IconFilePathSource.ValueString() != "" {
			iconPath := data.AppIcon.IconFilePathSource.ValueString()
			iconBytes, err := os.ReadFile(iconPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read PNG icon file from %s: %v", iconPath, err)
			}
			largeIcon.SetValue(iconBytes)
			baseApp.SetLargeIcon(largeIcon)
		} else if !data.AppIcon.IconURLSource.IsNull() && data.AppIcon.IconURLSource.ValueString() != "" {
			webSource := data.AppIcon.IconURLSource.ValueString()

			downloadedPath, err := download.DownloadFile(webSource)
			if err != nil {
				return nil, fmt.Errorf("failed to download icon file from %s: %v", webSource, err)
			}

			iconTempFile := helpers.TempFileInfo{
				FilePath:      downloadedPath,
				ShouldCleanup: true,
			}
			// Clean up the icon file when done with this function
			defer helpers.CleanupTempFile(ctx, iconTempFile)

			iconBytes, err := os.ReadFile(downloadedPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read downloaded PNG icon file from %s: %v", downloadedPath, err)
			}

			largeIcon.SetValue(iconBytes)
			baseApp.SetLargeIcon(largeIcon)
		}
	}

	// For creating resources, we need the installer file to extract metadata
	// Verify the installer path is provided and the file exists
	if installerSourcePath == "" {
		return nil, fmt.Errorf("installer source path is empty; a valid file path is required")
	}

	if _, err := os.Stat(installerSourcePath); err != nil {
		return nil, fmt.Errorf("installer file not found at path %s: %w", installerSourcePath, err)
	}

	filename := filepath.Base(installerSourcePath)
	tflog.Debug(ctx, fmt.Sprintf("Using filename from installer path: %s", filename))
	convert.FrameworkToGraphString(types.StringValue(filename), baseApp.SetFileName)

	// Extract fields from all Info.plist files in the package using the resolved path
	fields := []utility.Field{
		{Key: "CFBundleIdentifier", Required: true},
		{Key: "CFBundleShortVersionString", Required: true},
	}

	tflog.Debug(ctx, fmt.Sprintf("Extracting metadata from PKG file: %s", installerSourcePath))
	extractedFields, err := utility.ExtractFieldsFromPKGFile(
		ctx,
		installerSourcePath,
		"Info.plist",
		fields,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata from pkg file: %w", err)
	}

	if len(extractedFields) == 0 {
		return nil, fmt.Errorf("no Info.plist files found in the PKG installer at %s", installerSourcePath)
	}

	// First Info.plist becomes primary bundle
	primaryBundleId := extractedFields[0].Values["CFBundleIdentifier"]
	primaryBundleVersion := extractedFields[0].Values["CFBundleShortVersionString"]

	tflog.Debug(ctx, fmt.Sprintf("Setting primary bundle ID: %s, version: %s", primaryBundleId, primaryBundleVersion))
	convert.FrameworkToGraphString(types.StringValue(primaryBundleId), baseApp.SetPrimaryBundleId)
	convert.FrameworkToGraphString(types.StringValue(primaryBundleVersion), baseApp.SetPrimaryBundleVersion)

	// All entries are set as included apps (including primary)
	var includedApps []graphmodels.MacOSIncludedAppable
	for _, fields := range extractedFields {
		includedApp := graphmodels.NewMacOSIncludedApp()
		convert.FrameworkToGraphString(
			types.StringValue(fields.Values["CFBundleIdentifier"]),
			includedApp.SetBundleId,
		)
		convert.FrameworkToGraphString(
			types.StringValue(fields.Values["CFBundleShortVersionString"]),
			includedApp.SetBundleVersion,
		)
		includedApps = append(includedApps, includedApp)
	}

	baseApp.SetIncludedApps(includedApps)
	tflog.Debug(ctx, fmt.Sprintf("Added %d included apps from PKG metadata", len(includedApps)))

	convert.FrameworkToGraphBool(data.MacOSPkgApp.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)

	if data.MacOSPkgApp.MinimumSupportedOperatingSystem != nil {
		minOS := graphmodels.NewMacOSMinimumOperatingSystem()
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
		convert.FrameworkToGraphBool(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
		baseApp.SetMinimumSupportedOperatingSystem(minOS)
	}

	if data.MacOSPkgApp.PreInstallScript != nil {
		preScript := graphmodels.NewMacOSAppScript()
		convert.FrameworkToGraphString(data.MacOSPkgApp.PreInstallScript.ScriptContent, preScript.SetScriptContent)
		baseApp.SetPreInstallScript(preScript)
	}

	if data.MacOSPkgApp.PostInstallScript != nil {
		postScript := graphmodels.NewMacOSAppScript()
		convert.FrameworkToGraphString(data.MacOSPkgApp.PostInstallScript.ScriptContent, postScript.SetScriptContent)
		baseApp.SetPostInstallScript(postScript)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), baseApp); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return baseApp, nil
}
