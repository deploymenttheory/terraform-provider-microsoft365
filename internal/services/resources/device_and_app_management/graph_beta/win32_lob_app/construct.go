package graphBetaWin32LobApp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	sharedConstructors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	helpersCrud "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *Win32LobAppResourceModel, installerSourcePath string) (graphmodels.Win32LobAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWin32LobApp()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, requestBody.SetPublisher)
	convert.FrameworkToGraphString(data.FileName, requestBody.SetFileName)
	convert.FrameworkToGraphString(data.InstallCommandLine, requestBody.SetInstallCommandLine)
	convert.FrameworkToGraphString(data.UninstallCommandLine, requestBody.SetUninstallCommandLine)
	convert.FrameworkToGraphString(data.SetupFilePath, requestBody.SetSetupFilePath)
	convert.FrameworkToGraphString(data.CommittedContentVersion, requestBody.SetCommittedContentVersion)
	convert.FrameworkToGraphString(data.DisplayVersion, requestBody.SetDisplayVersion)
	convert.FrameworkToGraphString(data.Developer, requestBody.SetDeveloper)
	convert.FrameworkToGraphString(data.InformationUrl, requestBody.SetInformationUrl)
	convert.FrameworkToGraphString(data.PrivacyInformationUrl, requestBody.SetPrivacyInformationUrl)
	convert.FrameworkToGraphString(data.Notes, requestBody.SetNotes)
	convert.FrameworkToGraphString(data.Owner, requestBody.SetOwner)
	convert.FrameworkToGraphString(data.MinimumSupportedWindowsRelease, requestBody.SetMinimumSupportedWindowsRelease)
	convert.FrameworkToGraphBool(data.AllowAvailableUninstall, requestBody.SetAllowAvailableUninstall)
	convert.FrameworkToGraphBool(data.IsFeatured, requestBody.SetIsFeatured)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Handle applicable architectures
	if err := convert.FrameworkToGraphBitmaskEnumFromSet(ctx, data.AllowedArchitectures,
		graphmodels.ParseWindowsArchitecture, requestBody.SetAllowedArchitectures); err != nil {
		tflog.Warn(ctx, "Failed to set applicable architectures", map[string]any{
			"error": err.Error(),
		})
	}

	// Handle app icon (either from file path or web source)
	if data.AppIcon != nil {
		largeIcon, tempFiles, err := sharedConstructors.ConstructMobileAppIcon(ctx, data.AppIcon)
		if err != nil {
			return nil, err
		}

		defer func() {
			for _, tempFile := range tempFiles {
				helpersCrud.CleanupTempFile(ctx, tempFile)
			}
		}()

		requestBody.SetLargeIcon(largeIcon)
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
	convert.FrameworkToGraphString(types.StringValue(filename), requestBody.SetFileName)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
