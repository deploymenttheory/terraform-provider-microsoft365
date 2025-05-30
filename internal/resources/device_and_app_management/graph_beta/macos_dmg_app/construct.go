package graphBetaMacOSDmgApp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud/graph_beta/device_and_app_management"
	download "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/common"
	//utility "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/macos_dmg"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform configuration to a macOS DMG App request suitable for the Microsoft Graph API
func constructResource(ctx context.Context, data *MacOSDmgAppResourceModel, installerSourcePath string) (graphmodels.MacOSDmgAppable, error) {
	tflog.Debug(ctx, "Constructing macOS DMG app resource from Terraform configuration")

	requestBody := graphmodels.NewMacOSDmgApp()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetStringProperty(data.Publisher, requestBody.SetPublisher)
	constructors.SetStringProperty(data.Developer, requestBody.SetDeveloper)
	constructors.SetStringProperty(data.Owner, requestBody.SetOwner)
	constructors.SetStringProperty(data.Notes, requestBody.SetNotes)
	constructors.SetStringProperty(data.InformationUrl, requestBody.SetInformationUrl)
	constructors.SetStringProperty(data.PrivacyInformationUrl, requestBody.SetPrivacyInformationUrl)
	constructors.SetBoolProperty(data.IsFeatured, requestBody.SetIsFeatured)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
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
			requestBody.SetLargeIcon(largeIcon)
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
			requestBody.SetLargeIcon(largeIcon)
		}
	}

	// For creating resources, we need the installer file to set filename
	if installerSourcePath != "" {
		if _, err := os.Stat(installerSourcePath); err != nil {
			return nil, fmt.Errorf("installer file not found at path %s: %w", installerSourcePath, err)
		}

		filename := filepath.Base(installerSourcePath)
		tflog.Debug(ctx, fmt.Sprintf("Using filename from installer path: %s", filename))
		constructors.SetStringProperty(types.StringValue(filename), requestBody.SetFileName)
	}

	// Handle macOS DMG app specific properties
	if data.MacOSDmgApp != nil {
		baseApp := requestBody

		if data.MacOSDmgApp.MinimumSupportedOperatingSystem != nil {
			minOS := graphmodels.NewMacOSMinimumOperatingSystem()
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
			constructors.SetBoolProperty(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V150, minOS.SetV150)

			baseApp.SetMinimumSupportedOperatingSystem(minOS)
			tflog.Debug(ctx, "Mapped minimum supported operating system")
		}

	// 	// Extract and set metadata from DMG file if available
	// 	if installerSourcePath != "" {
	// 		err := extractAndSetDMGMetadata(ctx, baseApp, data, installerSourcePath)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("failed to extract DMG metadata: %w", err)
	// 		}
	// 	}

	// 	constructors.SetBoolProperty(data.MacOSDmgApp.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)

	// 	// Map included apps if provided
	// 	if !data.MacOSDmgApp.IncludedApps.IsNull() && !data.MacOSDmgApp.IncludedApps.IsUnknown() {
	// 		var includedAppsElements []MacOSIncludedAppResourceModel
	// 		diags := data.MacOSDmgApp.IncludedApps.ElementsAs(ctx, &includedAppsElements, false)
	// 		if diags.HasError() {
	// 			return nil, fmt.Errorf("error extracting included apps: %s", diags.Errors())
	// 		}

	// 		var includedApps []graphmodels.MacOSIncludedAppable
	// 		for _, app := range includedAppsElements {
	// 			includedApp := graphmodels.NewMacOSIncludedApp()
	// 			constructors.SetStringProperty(app.BundleId, includedApp.SetBundleId)
	// 			constructors.SetStringProperty(app.BundleVersion, includedApp.SetBundleVersion)
	// 			includedApps = append(includedApps, includedApp)
	// 		}
	// 		baseApp.SetIncludedApps(includedApps)
	// 		tflog.Debug(ctx, fmt.Sprintf("Mapped %d included apps", len(includedApps)))
	// 	}

	// 	// Set primary bundle properties if provided
	// 	constructors.SetStringProperty(data.MacOSDmgApp.PrimaryBundleId, baseApp.SetPrimaryBundleId)
	// 	constructors.SetStringProperty(data.MacOSDmgApp.PrimaryBundleVersion, baseApp.SetPrimaryBundleVersion)
	// }

	tflog.Debug(ctx, "Successfully constructed macOS DMG app resource")
	return requestBody, nil
}

// // extractAndSetDMGMetadata extracts metadata from DMG file and sets it on the base app
// // Note: DMG metadata extraction will reuse the same PKG utility since both are Apple package formats
// func extractAndSetDMGMetadata(ctx context.Context, baseApp graphmodels.MacOSDmgAppable, data *MacOSDmgAppResourceModel, installerSourcePath string) error {
// 	if !strings.HasSuffix(strings.ToLower(installerSourcePath), ".dmg") {
// 		return fmt.Errorf("file must be a .dmg file, got: %s", installerSourcePath)
// 	}

// 	// Extract fields from all Info.plist files in the DMG using the resolved path
// 	fields := []utility.Field{
// 		{Key: "CFBundleIdentifier", Required: true},
// 		{Key: "CFBundleShortVersionString", Required: true},
// 	}

// 	tflog.Debug(ctx, fmt.Sprintf("Extracting metadata from DMG file: %s", installerSourcePath))
// 	extractedFields, err := utility.ExtractFieldsFromDMGFile(
// 		ctx,
// 		installerSourcePath,
// 		"Info.plist",
// 		fields,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to extract metadata from DMG file: %w", err)
// 	}

// 	if len(extractedFields) == 0 {
// 		return fmt.Errorf("no Info.plist files found in the DMG installer at %s", installerSourcePath)
// 	}

// 	// First Info.plist becomes primary bundle
// 	primaryBundleId := extractedFields[0].Values["CFBundleIdentifier"]
// 	primaryBundleVersion := extractedFields[0].Values["CFBundleShortVersionString"]

// 	tflog.Debug(ctx, fmt.Sprintf("Setting primary bundle ID: %s, version: %s", primaryBundleId, primaryBundleVersion))
// 	constructors.SetStringProperty(types.StringValue(primaryBundleId), baseApp.SetPrimaryBundleId)
// 	constructors.SetStringProperty(types.StringValue(primaryBundleVersion), baseApp.SetPrimaryBundleVersion)

// 	// All entries are set as included apps (including primary)
// 	var includedApps []graphmodels.MacOSIncludedAppable
// 	for _, fields := range extractedFields {
// 		includedApp := graphmodels.NewMacOSIncludedApp()
// 		constructors.SetStringProperty(
// 			types.StringValue(fields.Values["CFBundleIdentifier"]),
// 			includedApp.SetBundleId,
// 		)
// 		constructors.SetStringProperty(
// 			types.StringValue(fields.Values["CFBundleShortVersionString"]),
// 			includedApp.SetBundleVersion,
// 		)
// 		includedApps = append(includedApps, includedApp)
// 	}

// 	baseApp.SetIncludedApps(includedApps)
// 	tflog.Debug(ctx, fmt.Sprintf("Added %d included apps from DMG metadata", len(includedApps)))

// 	return nil
// }
