package graphBetaMacOSLobApp

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

// constructResource maps the Terraform schema to the SDK model, using the provided installer path
func constructResource(ctx context.Context, data *MacOSLobAppResourceModel, installerSourcePath string) (graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewMacOSLobApp()

	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, requestBody.SetPublisher)
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.InformationUrl, requestBody.SetInformationUrl)
	convert.FrameworkToGraphBool(data.IsFeatured, requestBody.SetIsFeatured)
	convert.FrameworkToGraphString(data.Owner, requestBody.SetOwner)
	convert.FrameworkToGraphString(data.Developer, requestBody.SetDeveloper)
	convert.FrameworkToGraphString(data.Notes, requestBody.SetNotes)
	convert.FrameworkToGraphString(data.PrivacyInformationUrl, requestBody.SetPrivacyInformationUrl)

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

	// Set macOS LOB app specific properties
	if data.MacOSLobApp != nil {
		convert.FrameworkToGraphString(data.MacOSLobApp.BundleId, requestBody.SetBundleId)
		convert.FrameworkToGraphString(data.MacOSLobApp.BuildNumber, requestBody.SetBuildNumber)
		convert.FrameworkToGraphString(data.MacOSLobApp.VersionNumber, requestBody.SetVersionNumber)
		convert.FrameworkToGraphBool(data.MacOSLobApp.IgnoreVersionDetection, requestBody.SetIgnoreVersionDetection)
		convert.FrameworkToGraphBool(data.MacOSLobApp.InstallAsManaged, requestBody.SetInstallAsManaged)

		// Set child apps if provided
		if len(data.MacOSLobApp.ChildApps) > 0 {
			var childApps []graphmodels.MacOSLobChildAppable
			for _, childApp := range data.MacOSLobApp.ChildApps {
				childAppModel := graphmodels.NewMacOSLobChildApp()
				convert.FrameworkToGraphString(childApp.BundleId, childAppModel.SetBundleId)
				convert.FrameworkToGraphString(childApp.BuildNumber, childAppModel.SetBuildNumber)
				convert.FrameworkToGraphString(childApp.VersionNumber, childAppModel.SetVersionNumber)
				childApps = append(childApps, childAppModel)
			}
			requestBody.SetChildApps(childApps)
			tflog.Debug(ctx, fmt.Sprintf("Added %d child apps", len(childApps)))
		}

		// Set minimum supported operating system
		if data.MacOSLobApp.MinimumSupportedOperatingSystem != nil {
			minOS := graphmodels.NewMacOSMinimumOperatingSystem()
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
			convert.FrameworkToGraphBool(data.MacOSLobApp.MinimumSupportedOperatingSystem.V150, minOS.SetV150)
			requestBody.SetMinimumSupportedOperatingSystem(minOS)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
