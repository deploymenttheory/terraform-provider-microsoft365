package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"
	"os"

	download "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/common"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TempFileInfo holds information about a temporary file that needs to be cleaned up
type TempFileInfo struct {
	// Path to the temporary file
	FilePath string
	// Whether this is a downloaded file that should be cleaned up
	ShouldCleanup bool
}

// setInstallerSourcePath inspects the resource model and returns the actual installer file path
// along with information about whether it needs cleanup.
func setInstallerSourcePath(ctx context.Context, model *MacOSPkgAppResourceModel) (string, TempFileInfo, error) {
	var fileInfo TempFileInfo

	// Prefer the local installer file if provided
	if !model.InstallerFilePathSource.IsNull() && model.InstallerFilePathSource.ValueString() != "" {
		fileInfo.FilePath = model.InstallerFilePathSource.ValueString()
		fileInfo.ShouldCleanup = false
		tflog.Debug(ctx, fmt.Sprintf("Using local installer file: %s", fileInfo.FilePath))
	} else if !model.InstallerURLSource.IsNull() && model.InstallerURLSource.ValueString() != "" {
		// Download the installer file if a URL is provided
		tflog.Debug(ctx, fmt.Sprintf("Downloading installer file from URL: %s", model.InstallerURLSource.ValueString()))
		downloadedPath, err := download.DownloadFile(model.InstallerURLSource.ValueString())
		if err != nil {
			return "", fileInfo, fmt.Errorf("failed to download installer file from %s: %v", model.InstallerURLSource.ValueString(), err)
		}
		fileInfo.FilePath = downloadedPath
		fileInfo.ShouldCleanup = true
		tflog.Debug(ctx, fmt.Sprintf("Downloaded installer file to: %s", fileInfo.FilePath))
	} else {
		return "", fileInfo, fmt.Errorf("installer file not provided; please supply either a local file path or a URL")
	}

	return fileInfo.FilePath, fileInfo, nil
}

// cleanupTempFile removes a temporary file if it exists and should be cleaned up
func cleanupTempFile(ctx context.Context, fileInfo TempFileInfo) {
	if !fileInfo.ShouldCleanup || fileInfo.FilePath == "" {
		return
	}

	if err := os.Remove(fileInfo.FilePath); err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to remove temporary file %s: %v", fileInfo.FilePath, err))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Successfully removed temporary file: %s", fileInfo.FilePath))
	}
}
