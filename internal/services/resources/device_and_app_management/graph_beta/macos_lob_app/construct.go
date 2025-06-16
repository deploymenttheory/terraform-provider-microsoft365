package graphBetaMacOSLobApp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	download "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model, using the provided installer path
func constructResource(ctx context.Context, data *MacOSLobAppResourceModel, installerSourcePath string) (graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	baseApp := graphmodels.NewMacOSLobApp()

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

	// For creating resources, we need the installer file to set filename
	if installerSourcePath != "" {
		if _, err := os.Stat(installerSourcePath); err != nil {
			return nil, fmt.Errorf("installer file not found at path %s: %w", installerSourcePath, err)
		}

		filename := filepath.Base(installerSourcePath)
		tflog.Debug(ctx, fmt.Sprintf("Using filename from installer path: %s", filename))
		constructors.SetStringProperty(types.StringValue(filename), baseApp.SetFileName)
	}

	// Set macOS LOB app specific properties
	if data.MacOSLobApp != nil {
		constructors.SetStringProperty(data.MacOSLobApp.BundleId, baseApp.SetBundleId)
		constructors.SetStringProperty(data.MacOSLobApp.BuildNumber, baseApp.SetBuildNumber)
		constructors.SetStringProperty(data.MacOSLobApp.VersionNumber, baseApp.SetVersionNumber)
		constructors.SetBoolProperty(data.MacOSLobApp.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)
		constructors.SetBoolProperty(data.MacOSLobApp.InstallAsManaged, baseApp.SetInstallAsManaged)

		// Set child apps if provided
		if len(data.MacOSLobApp.ChildApps) > 0 {
			var childApps []graphmodels.MacOSLobChildAppable
			for _, childApp := range data.MacOSLobApp.ChildApps {
				childAppModel := graphmodels.NewMacOSLobChildApp()
				constructors.SetStringProperty(childApp.BundleId, childAppModel.SetBundleId)
				constructors.SetStringProperty(childApp.BuildNumber, childAppModel.SetBuildNumber)
				constructors.SetStringProperty(childApp.VersionNumber, childAppModel.SetVersionNumber)
				childApps = append(childApps, childAppModel)
			}
			baseApp.SetChildApps(childApps)
			tflog.Debug(ctx, fmt.Sprintf("Added %d child apps", len(childApps)))
		}

		// Set minimum supported operating system
		if data.MacOSLobApp.MinimumSupportedOperatingSystem != nil {
			minOS := graphmodels.NewMacOSMinimumOperatingSystem()
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
			constructors.SetBoolProperty(data.MacOSLobApp.MinimumSupportedOperatingSystem.V150, minOS.SetV150)
			baseApp.SetMinimumSupportedOperatingSystem(minOS)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), baseApp); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return baseApp, nil
}
