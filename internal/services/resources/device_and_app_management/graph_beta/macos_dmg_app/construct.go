package graphBetaMacOSDmgApp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	sharedConstructors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform configuration to a macOS DMG App request suitable for the Microsoft Graph API
func constructResource(ctx context.Context, data *MacOSDmgAppResourceModel, installerSourcePath string) (graphmodels.MacOSDmgAppable, error) {
	tflog.Debug(ctx, "Constructing macOS DMG app resource from Terraform configuration")

	requestBody := graphmodels.NewMacOSDmgApp()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, requestBody.SetPublisher)
	convert.FrameworkToGraphString(data.Developer, requestBody.SetDeveloper)
	convert.FrameworkToGraphString(data.Owner, requestBody.SetOwner)
	convert.FrameworkToGraphString(data.Notes, requestBody.SetNotes)
	convert.FrameworkToGraphString(data.InformationUrl, requestBody.SetInformationUrl)
	convert.FrameworkToGraphString(data.PrivacyInformationUrl, requestBody.SetPrivacyInformationUrl)
	convert.FrameworkToGraphBool(data.IsFeatured, requestBody.SetIsFeatured)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Handle app icon (either from file path or web source)
	if data.AppIcon != nil {
		largeIcon, tempFiles, err := sharedConstructors.ConstructMobileAppIcon(ctx, data.AppIcon)
		if err != nil {
			return nil, err
		}

		defer func() {
			for _, tempFile := range tempFiles {
				helpers.CleanupTempFile(ctx, tempFile)
			}
		}()

		requestBody.SetLargeIcon(largeIcon)
	}

	// For creating resources, we need the installer file to set filename
	if installerSourcePath != "" {
		if _, err := os.Stat(installerSourcePath); err != nil {
			return nil, fmt.Errorf("installer file not found at path %s: %w", installerSourcePath, err)
		}

		filename := filepath.Base(installerSourcePath)
		tflog.Debug(ctx, fmt.Sprintf("Using filename from installer path: %s", filename))
		convert.FrameworkToGraphString(types.StringValue(filename), requestBody.SetFileName)
	}

	// Handle macOS DMG app specific properties
	if data.MacOSDmgApp != nil {
		baseApp := requestBody

		if data.MacOSDmgApp.MinimumSupportedOperatingSystem != nil {
			minOS := graphmodels.NewMacOSMinimumOperatingSystem()
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
			convert.FrameworkToGraphBool(data.MacOSDmgApp.MinimumSupportedOperatingSystem.V150, minOS.SetV150)

			baseApp.SetMinimumSupportedOperatingSystem(minOS)
			tflog.Debug(ctx, "Mapped minimum supported operating system")
		}

		// Handle included apps from Terraform configuration
		if !data.MacOSDmgApp.IncludedApps.IsNull() && !data.MacOSDmgApp.IncludedApps.IsUnknown() {
			includedAppsElements := data.MacOSDmgApp.IncludedApps.Elements()
			var includedApps []graphmodels.MacOSIncludedAppable
			var primaryBundleId, primaryBundleVersion string
			var isPrimarySet bool

			for _, elem := range includedAppsElements {
				objVal, ok := elem.(types.Object)
				if !ok {
					tflog.Error(ctx, fmt.Sprintf("Expected object type for included app, got %T", elem))
					continue
				}

				attrs := objVal.Attributes()
				includedApp := graphmodels.NewMacOSIncludedApp()

				// Set the OData type for the included app
				odataType := "#microsoft.graph.macOSIncludedApp"
				includedApp.SetOdataType(&odataType)

				var currentBundleId, currentBundleVersion string

				if bundleId, exists := attrs["bundle_id"]; exists {
					if bundleIdStr, ok := bundleId.(types.String); ok && !bundleIdStr.IsNull() {
						convert.FrameworkToGraphString(bundleIdStr, includedApp.SetBundleId)
						currentBundleId = bundleIdStr.ValueString()
					}
				}

				if bundleVersion, exists := attrs["bundle_version"]; exists {
					if bundleVersionStr, ok := bundleVersion.(types.String); ok && !bundleVersionStr.IsNull() {
						convert.FrameworkToGraphString(bundleVersionStr, includedApp.SetBundleVersion)
						currentBundleVersion = bundleVersionStr.ValueString()
					}
				}

				// Set the first included app as the primary bundle
				if !isPrimarySet && currentBundleId != "" && currentBundleVersion != "" {
					primaryBundleId = currentBundleId
					primaryBundleVersion = currentBundleVersion
					isPrimarySet = true
					tflog.Debug(ctx, fmt.Sprintf("Set primary bundle: ID=%s, Version=%s", primaryBundleId, primaryBundleVersion))
				}

				includedApps = append(includedApps, includedApp)
			}

			// Set the primary bundle properties
			if isPrimarySet {
				baseApp.SetPrimaryBundleId(&primaryBundleId)
				baseApp.SetPrimaryBundleVersion(&primaryBundleVersion)
			}

			baseApp.SetIncludedApps(includedApps)
			tflog.Debug(ctx, fmt.Sprintf("Added %d included apps from Terraform configuration", len(includedApps)))
		}

		// Set ignore version detection
		convert.FrameworkToGraphBool(data.MacOSDmgApp.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Successfully constructed macOS DMG app resource")
	return requestBody, nil
}
