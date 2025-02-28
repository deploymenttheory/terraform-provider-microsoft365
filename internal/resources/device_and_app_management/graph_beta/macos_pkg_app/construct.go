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

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *MacOSPKGAppResourceModel) (graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var baseApp graphmodels.MobileAppable

	baseApp = graphmodels.NewMacOSPkgApp()
	macOSApp := baseApp.(graphmodels.MacOSPkgAppable)

	// Set base properties
	constructors.SetStringProperty(data.Description, baseApp.SetDescription)
	constructors.SetStringProperty(data.Publisher, baseApp.SetPublisher)
	constructors.SetStringProperty(data.DisplayName, baseApp.SetDisplayName)
	constructors.SetStringProperty(data.InformationUrl, baseApp.SetInformationUrl)
	constructors.SetBoolProperty(data.IsFeatured, baseApp.SetIsFeatured)
	constructors.SetStringProperty(data.Owner, baseApp.SetOwner)
	constructors.SetStringProperty(data.Developer, baseApp.SetDeveloper)
	constructors.SetStringProperty(data.Notes, baseApp.SetNotes)

	if err := constructors.SetStringList(ctx, data.RoleScopeTagIds, baseApp.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Handle categories
	// if data.Categories != nil && len(data.Categories) > 0 {
	// 	// Extract display names from the categories
	// 	displayNames := make([]string, 0, len(data.Categories))
	// 	for _, category := range data.Categories {
	// 		if !category.DisplayName.IsNull() && !category.DisplayName.IsUnknown() {
	// 			displayNames = append(displayNames, category.DisplayName.ValueString())
	// 		}
	// 	}

	// 	// Build the categories using our helper function
	// 	if len(displayNames) > 0 {
	// 		categories := BuildCategoriesFromDisplayNames(displayNames)
	// 		baseApp.SetCategories(categories)
	// 	}
	// }

	// Handle app icon (either from file path or web source)
	if data.AppIcon != nil {
		largeIcon := graphmodels.NewMimeContent()
		iconType := "image/png"
		largeIcon.SetTypeEscaped(&iconType)

		// Get icon from file path
		if !data.AppIcon.IconFilePath.IsNull() && data.AppIcon.IconFilePath.ValueString() != "" {
			iconPath := data.AppIcon.IconFilePath.ValueString()
			iconBytes, err := os.ReadFile(iconPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read PNG icon file from %s: %v", iconPath, err)
			}
			largeIcon.SetValue(iconBytes)
			baseApp.SetLargeIcon(largeIcon)
		} else if !data.AppIcon.IconFileWebSource.IsNull() && data.AppIcon.IconFileWebSource.ValueString() != "" {
			// Get icon from web source
			webSource := data.AppIcon.IconFileWebSource.ValueString()

			// Download the file
			downloadedPath, err := download.DownloadFile(webSource)
			if err != nil {
				return nil, fmt.Errorf("failed to download icon file from %s: %v", webSource, err)
			}

			// Clean up temporary file when we're done
			defer func() {
				if err := os.Remove(downloadedPath); err != nil {
					tflog.Warn(ctx, fmt.Sprintf("Failed to clean up temporary icon file %s: %v", downloadedPath, err))
				}
			}()

			// Read the downloaded file
			iconBytes, err := os.ReadFile(downloadedPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read downloaded PNG icon file from %s: %v", downloadedPath, err)
			}

			largeIcon.SetValue(iconBytes)
			baseApp.SetLargeIcon(largeIcon)
		}
	}

	// Set MacOS PKG specific properties
	if data.MacOSPkgApp.PackageInstallerFileSource.IsNull() || data.MacOSPkgApp.PackageInstallerFileSource.ValueString() == "" {
		return nil, fmt.Errorf("package_installer_file_source is required but not provided")
	}

	filename := filepath.Base(data.MacOSPkgApp.PackageInstallerFileSource.ValueString())
	constructors.SetStringProperty(types.StringValue(filename), macOSApp.SetFileName)

	fields := []utility.Field{
		{Key: "CFBundleIdentifier", Required: true},
		{Key: "CFBundleShortVersionString", Required: true},
	}

	// Extract fields from all Info.plist files in the package
	extractedFields, err := utility.ExtractFieldsFromFiles(
		ctx,
		data.MacOSPkgApp.PackageInstallerFileSource.ValueString(),
		"Info.plist",
		fields,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata from pkg file: %w", err)
	}

	// First Info.plist becomes primary bundle
	primaryBundleId := extractedFields[0].Values["CFBundleIdentifier"]
	primaryBundleVersion := extractedFields[0].Values["CFBundleShortVersionString"]

	constructors.SetStringProperty(types.StringValue(primaryBundleId), macOSApp.SetPrimaryBundleId)
	constructors.SetStringProperty(types.StringValue(primaryBundleVersion), macOSApp.SetPrimaryBundleVersion)

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

	macOSApp.SetIncludedApps(includedApps)

	constructors.SetBoolProperty(data.MacOSPkgApp.IgnoreVersionDetection, macOSApp.SetIgnoreVersionDetection)

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
		macOSApp.SetMinimumSupportedOperatingSystem(minOS)
	}

	if data.MacOSPkgApp.PreInstallScript != nil {
		preScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.MacOSPkgApp.PreInstallScript.ScriptContent, preScript.SetScriptContent)
		macOSApp.SetPreInstallScript(preScript)
	}

	if data.MacOSPkgApp.PostInstallScript != nil {
		postScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.MacOSPkgApp.PostInstallScript.ScriptContent, postScript.SetScriptContent)
		macOSApp.SetPostInstallScript(postScript)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), macOSApp); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully constructed MacOSPkgApp resource with %d included apps", len(includedApps)))
	return macOSApp, nil
}

// BuildCategoriesFromDisplayNames creates an array of MobileAppCategoryable objects
// from an array of display names. The function maps each display name to its
// corresponding ID based on predefined mappings.
func BuildCategoriesFromDisplayNames(displayNames []string) []graphmodels.MobileAppCategoryable {
	// Define the mapping of display names to IDs
	categoryMapping := map[string]string{
		"Other apps":             "0720a99e-562b-4a77-83f0-9a7523fcf13e",
		"Books & Reference":      "f1fc9fe2-728d-4867-9a72-a61e18f8c606",
		"Data management":        "046e0b16-76ce-4b49-bf1b-1cc5bd94fb47",
		"Productivity":           "ed899483-3019-425e-a470-28e901b9790e",
		"Business":               "2b73ae71-12c8-49be-b462-3dae769ccd9d",
		"Development & Design":   "79bc98d4-7ddf-4841-9bc1-5c84a26d7ee8",
		"Photos & Media":         "5dcd7a90-0306-4f09-a75d-6b97a243f04e",
		"Collaboration & Social": "f79135dc-8e41-48c1-9a59-ab9a7259c38e",
		"Computer management":    "981deed8-6857-4e78-a50e-c3f61d312737",
	}

	// Create the array of category objects
	categories := make([]graphmodels.MobileAppCategoryable, 0, len(displayNames))

	// Process each display name
	for _, name := range displayNames {
		// Check if the display name exists in our mapping
		id, exists := categoryMapping[name]
		if !exists {
			// Skip invalid category names
			continue
		}

		// Create a new category object
		category := graphmodels.NewMobileAppCategory()

		// Set the display name
		displayNameCopy := name // Create a copy to avoid issues with loop variable capture
		category.SetDisplayName(&displayNameCopy)

		// Set the ID
		idCopy := id // Create a copy for the same reason
		category.SetId(&idCopy)

		// Add to the array
		categories = append(categories, category)
	}

	return categories
}
