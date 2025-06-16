package helpers

import (
	"context"
	"fmt"
	"os"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	download "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TempFileInfo holds information about a temporary file that needs to be cleaned up
type TempFileInfo struct {
	// Path to the temporary file
	FilePath string
	// Whether this is a downloaded file that should be cleaned up
	ShouldCleanup bool
}

// SetInstallerSourcePath inspects the resource model's app_installer and returns the actual installer file path
// along with information about whether it needs cleanup.
func SetInstallerSourcePath(ctx context.Context, metadataObj types.Object) (string, TempFileInfo, error) {
	var fileInfo TempFileInfo
	var metadata sharedmodels.MobileAppMetaDataResourceModel

	if !metadataObj.IsNull() {
		diags := metadataObj.As(ctx, &metadata, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			tflog.Warn(ctx, "Failed to parse app_installer, proceeding with empty struct", map[string]interface{}{
				"errors": diags.Errors(),
			})
		}
	}

	tflog.Debug(ctx, "setInstallerSourcePath input", map[string]interface{}{
		"installer_file_path_source": metadata.InstallerFilePathSource.String(),
		"installer_url_source":       metadata.InstallerURLSource.String(),
	})

	if metadata.InstallerFilePathSource.IsUnknown() && metadata.InstallerURLSource.IsUnknown() {
		tflog.Debug(ctx, "Installer sources are unknown during plan - skipping")
		return "", fileInfo, nil
	}

	// Prefer local file path
	if !metadata.InstallerFilePathSource.IsNull() && metadata.InstallerFilePathSource.ValueString() != "" {
		fileInfo.FilePath = metadata.InstallerFilePathSource.ValueString()
		fileInfo.ShouldCleanup = false
		tflog.Debug(ctx, fmt.Sprintf("Using local installer file: %s", fileInfo.FilePath))
		return fileInfo.FilePath, fileInfo, nil
	}

	// Otherwise try URL
	if !metadata.InstallerURLSource.IsNull() && metadata.InstallerURLSource.ValueString() != "" {
		url := metadata.InstallerURLSource.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Downloading installer from URL: %s", url))
		downloadedPath, err := download.DownloadFile(url)
		if err != nil {
			return "", fileInfo, fmt.Errorf("failed to download installer file from %s: %v", url, err)
		}
		fileInfo.FilePath = downloadedPath
		fileInfo.ShouldCleanup = true
		return fileInfo.FilePath, fileInfo, nil
	}

	return "", fileInfo, fmt.Errorf("installer file not provided; please supply either a local file path or a URL")
}

// CleanupTempFile removes a temporary file if it exists and should be cleaned up
func CleanupTempFile(ctx context.Context, fileInfo TempFileInfo) {
	if !fileInfo.ShouldCleanup || fileInfo.FilePath == "" {
		return
	}

	if err := os.Remove(fileInfo.FilePath); err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to remove temporary file %s: %v", fileInfo.FilePath, err))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Successfully removed temporary file: %s", fileInfo.FilePath))
	}
}
