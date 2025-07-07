package sharedConstructors

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	download "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/common"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructMobileAppIcon handles the app icon processing, including format detection and conversion to PNG
// Returns a MimeContent object and any temporary files that need to be cleaned up
func ConstructMobileAppIcon(ctx context.Context, appIcon *sharedmodels.MobileAppIconResourceModel) (graphmodels.MimeContentable, []helpers.TempFileInfo, error) {
	largeIcon := graphmodels.NewMimeContent()
	iconType := "image/png"
	largeIcon.SetTypeEscaped(&iconType)

	var tempFiles []helpers.TempFileInfo

	if !appIcon.IconFilePathSource.IsNull() && appIcon.IconFilePathSource.ValueString() != "" {
		iconPath := appIcon.IconFilePathSource.ValueString()

		// Check if the icon is already PNG and if not, convert it
		iconBytes, err := os.ReadFile(iconPath)
		if err != nil {
			return nil, tempFiles, fmt.Errorf("failed to read icon file from %s: %v", iconPath, err)
		}

		if !download.IsPNG(iconBytes) {
			tflog.Debug(ctx, fmt.Sprintf("Converting icon from non-PNG format to PNG: %s", iconPath))

			pngBytes, err := download.ConvertToPNG(ctx, iconPath)
			if err != nil {
				return nil, tempFiles, fmt.Errorf("failed to convert icon to PNG format: %v", err)
			}

			tempDir := os.TempDir()
			tempPngPath := filepath.Join(tempDir, "converted_icon.png")

			if err := os.WriteFile(tempPngPath, pngBytes, 0644); err != nil {
				return nil, tempFiles, fmt.Errorf("failed to write converted PNG icon: %v", err)
			}

			// Add the temp file to the cleanup list
			tempFiles = append(tempFiles, helpers.TempFileInfo{
				FilePath:      tempPngPath,
				ShouldCleanup: true,
			})

			iconBytes = pngBytes
			tflog.Debug(ctx, "Successfully converted icon to PNG format")
		}

		largeIcon.SetValue(iconBytes)
	} else if !appIcon.IconURLSource.IsNull() && appIcon.IconURLSource.ValueString() != "" {
		webSource := appIcon.IconURLSource.ValueString()

		downloadedPath, err := download.DownloadFile(webSource)
		if err != nil {
			return nil, tempFiles, fmt.Errorf("failed to download icon file from %s: %v", webSource, err)
		}

		tempFiles = append(tempFiles, helpers.TempFileInfo{
			FilePath:      downloadedPath,
			ShouldCleanup: true,
		})

		iconBytes, err := os.ReadFile(downloadedPath)
		if err != nil {
			return nil, tempFiles, fmt.Errorf("failed to read downloaded icon file from %s: %v", downloadedPath, err)
		}

		if !download.IsPNG(iconBytes) {
			tflog.Debug(ctx, fmt.Sprintf("Converting downloaded icon from non-PNG format to PNG: %s", webSource))

			pngBytes, err := download.ConvertToPNG(ctx, downloadedPath)
			if err != nil {
				return nil, tempFiles, fmt.Errorf("failed to convert downloaded icon to PNG format: %v", err)
			}

			tempDir := os.TempDir()
			tempPngPath := filepath.Join(tempDir, "converted_icon.png")

			if err := os.WriteFile(tempPngPath, pngBytes, 0644); err != nil {
				return nil, tempFiles, fmt.Errorf("failed to write converted PNG icon: %v", err)
			}

			tempFiles = append(tempFiles, helpers.TempFileInfo{
				FilePath:      tempPngPath,
				ShouldCleanup: true,
			})

			iconBytes = pngBytes
			tflog.Debug(ctx, "Successfully converted downloaded icon to PNG format")
		}

		largeIcon.SetValue(iconBytes)
	}

	return largeIcon, tempFiles, nil
}
