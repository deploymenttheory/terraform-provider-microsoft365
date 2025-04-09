package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
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

	// Set base properties
	constructors.SetStringProperty(data.Description, baseApp.SetDescription)
	constructors.SetStringProperty(data.Publisher, baseApp.SetPublisher)
	constructors.SetStringProperty(data.DisplayName, baseApp.SetDisplayName)
	constructors.SetStringProperty(data.InformationUrl, baseApp.SetInformationUrl)
	constructors.SetBoolProperty(data.IsFeatured, baseApp.SetIsFeatured)
	constructors.SetStringProperty(data.Owner, baseApp.SetOwner)
	constructors.SetStringProperty(data.Developer, baseApp.SetDeveloper)
	constructors.SetStringProperty(data.Notes, baseApp.SetNotes)
	constructors.SetStringProperty(data.PrivacyInformationUrl, baseApp.SetPrivacyInformationUrl)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, baseApp.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if !data.Categories.IsNull() {
		if err := constructors.SetObjectsFromStringSet(
			ctx,
			data.Categories,
			constructCategories,
			baseApp.SetCategories); err != nil {
			return nil, fmt.Errorf("failed to set categories: %s", err)
		}
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

			// Create temp file info for cleanup
			iconTempFile := TempFileInfo{
				FilePath:      downloadedPath,
				ShouldCleanup: true,
			}
			// Clean up the icon file when done with this function
			defer cleanupTempFile(ctx, iconTempFile)

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

	// Verify file exists and path is valid
	if _, err := os.Stat(installerSourcePath); err != nil {
		return nil, fmt.Errorf("installer file not found at path %s: %w", installerSourcePath, err)
	}

	// Set filename from the installer source path
	filename := filepath.Base(installerSourcePath)
	tflog.Debug(ctx, fmt.Sprintf("Using filename from installer path: %s", filename))
	constructors.SetStringProperty(types.StringValue(filename), baseApp.SetFileName)

	// Extract fields from all Info.plist files in the package using the resolved path
	fields := []utility.Field{
		{Key: "CFBundleIdentifier", Required: true},
		{Key: "CFBundleShortVersionString", Required: true},
	}

	tflog.Debug(ctx, fmt.Sprintf("Extracting metadata from PKG file: %s", installerSourcePath))
	extractedFields, err := utility.ExtractFieldsFromFiles(
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
	tflog.Debug(ctx, fmt.Sprintf("Added %d included apps from PKG metadata", len(includedApps)))

	constructors.SetBoolProperty(data.MacOSPkgApp.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)

	if data.MacOSPkgApp.MinimumSupportedOperatingSystem != nil {
		minOS := graphmodels.NewMacOSMinimumOperatingSystem()
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
		baseApp.SetMinimumSupportedOperatingSystem(minOS)
	}

	if data.MacOSPkgApp.PreInstallScript != nil {
		preScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.MacOSPkgApp.PreInstallScript.ScriptContent, preScript.SetScriptContent)
		baseApp.SetPreInstallScript(preScript)
	}

	if data.MacOSPkgApp.PostInstallScript != nil {
		postScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.MacOSPkgApp.PostInstallScript.ScriptContent, postScript.SetScriptContent)
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
