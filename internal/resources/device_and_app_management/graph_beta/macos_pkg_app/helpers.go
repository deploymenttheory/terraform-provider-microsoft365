package graphBetaMacOSPKGApp

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	download "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/common"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
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

// captureAppMetadata computes and sets the metadata for the installer file (size and checksums)
// used for evaluation for content version updates and other purposes
func captureAppMetadata(ctx context.Context, object *MacOSPKGAppResourceModel, installerSourcePath string) error {
	fileInfo, err := os.Stat(installerSourcePath)
	if err != nil {
		return fmt.Errorf("failed to get file information: %v", err)
	}

	if object.AppMetadata == nil {
		object.AppMetadata = &sharedmodels.MobileAppMetaDataResourceModel{}
	}

	currentSize := fileInfo.Size()
	object.AppMetadata.InstallerSizeInBytes = types.Int64Value(currentSize)

	tflog.Debug(ctx, "Computed installer file size", map[string]interface{}{
		"file_path":  installerSourcePath,
		"size_bytes": currentSize,
	})

	md5Hash := md5.New()
	sha256Hash := sha256.New()

	file, err := os.Open(installerSourcePath)
	if err != nil {
		return fmt.Errorf("failed to open file for checksums: %v", err)
	}
	defer file.Close()

	multiWriter := io.MultiWriter(md5Hash, sha256Hash)

	if _, err := io.Copy(multiWriter, file); err != nil {
		return fmt.Errorf("failed to read file for checksums: %v", err)
	}

	md5Checksum := hex.EncodeToString(md5Hash.Sum(nil))
	sha256Checksum := hex.EncodeToString(sha256Hash.Sum(nil))

	object.AppMetadata.InstallerMD5Checksum = types.StringValue(md5Checksum)
	object.AppMetadata.InstallerSHA256Checksum = types.StringValue(sha256Checksum)

	tflog.Debug(ctx, "Computed installer file checksums", map[string]interface{}{
		"file_path":       installerSourcePath,
		"md5_checksum":    md5Checksum,
		"sha256_checksum": sha256Checksum,
	})

	return nil
}

// evaluateIfContentVersionUpdateRequired determines if a new content version needs to be created and uploaded
// based on path changes and metadata differences
func evaluateIfContentVersionUpdateRequired(ctx context.Context, object *MacOSPKGAppResourceModel, state *MacOSPKGAppResourceModel, client *msgraphbetasdk.GraphServiceClient) (bool, string, error) {

	pathChanged := false

	if !object.MacOSPkgApp.InstallerFilePathSource.Equal(state.MacOSPkgApp.InstallerFilePathSource) {
		tflog.Debug(ctx, "Installer file path has changed", map[string]interface{}{
			"old_path": state.MacOSPkgApp.InstallerFilePathSource.ValueString(),
			"new_path": object.MacOSPkgApp.InstallerFilePathSource.ValueString(),
		})
		pathChanged = true
	}

	if !object.MacOSPkgApp.InstallerURLSource.Equal(state.MacOSPkgApp.InstallerURLSource) {
		tflog.Debug(ctx, "Installer URL has changed", map[string]interface{}{
			"old_url": state.MacOSPkgApp.InstallerURLSource.ValueString(),
			"new_url": object.MacOSPkgApp.InstallerURLSource.ValueString(),
		})
		pathChanged = true
	}

	// If neither installer source is provided, no upload is needed
	if (object.MacOSPkgApp.InstallerFilePathSource.IsNull() || object.MacOSPkgApp.InstallerFilePathSource.ValueString() == "") &&
		(object.MacOSPkgApp.InstallerURLSource.IsNull() || object.MacOSPkgApp.InstallerURLSource.ValueString() == "") {
		tflog.Debug(ctx, "No installer source provided, skipping content evaluation")
		return false, "", nil
	}

	metadataChanged := false
	existingContentVersion := ""

	if state.ContentVersion.IsNull() {
		tflog.Debug(ctx, "No content versions in state, upload required")
		return true, "", nil
	}

	// Check size differences if app metadata is available in both objects
	if object.AppMetadata != nil && state.AppMetadata != nil {
		if !object.AppMetadata.InstallerSizeInBytes.IsNull() && !state.AppMetadata.InstallerSizeInBytes.IsNull() {
			currentSize := object.AppMetadata.InstallerSizeInBytes.ValueInt64()
			previousSize := state.AppMetadata.InstallerSizeInBytes.ValueInt64()

			if currentSize != previousSize {
				tflog.Debug(ctx, "Installer size has changed", map[string]interface{}{
					"previous_size": previousSize,
					"current_size":  currentSize,
				})
				metadataChanged = true
			}
		}

		// Compare SHA256 checksum (preferred)
		if !object.AppMetadata.InstallerSHA256Checksum.IsNull() && !state.AppMetadata.InstallerSHA256Checksum.IsNull() {
			currentChecksum := object.AppMetadata.InstallerSHA256Checksum.ValueString()
			previousChecksum := state.AppMetadata.InstallerSHA256Checksum.ValueString()

			if currentChecksum != previousChecksum && currentChecksum != "" && previousChecksum != "" {
				tflog.Debug(ctx, "Installer SHA256 checksum has changed", map[string]interface{}{
					"previous_checksum": previousChecksum,
					"current_checksum":  currentChecksum,
				})
				metadataChanged = true
			}
		} else if !object.AppMetadata.InstallerMD5Checksum.IsNull() && !state.AppMetadata.InstallerMD5Checksum.IsNull() {
			currentChecksum := object.AppMetadata.InstallerMD5Checksum.ValueString()
			previousChecksum := state.AppMetadata.InstallerMD5Checksum.ValueString()

			if currentChecksum != previousChecksum && currentChecksum != "" && previousChecksum != "" {
				tflog.Debug(ctx, "Installer MD5 checksum has changed", map[string]interface{}{
					"previous_checksum": previousChecksum,
					"current_checksum":  currentChecksum,
				})
				metadataChanged = true
			}
		}
	} else {
		tflog.Debug(ctx, "Metadata is missing from one of the objects, assuming content needs update")
		metadataChanged = true
	}

	resource, err := client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Get(ctx, nil)

	if err == nil {
		macOSPkgApp, ok := resource.(graphmodels.MacOSPkgAppable)
		if ok && macOSPkgApp.GetCommittedContentVersion() != nil &&
			*macOSPkgApp.GetCommittedContentVersion() != "" {
			existingContentVersion = *macOSPkgApp.GetCommittedContentVersion()
			tflog.Debug(ctx, fmt.Sprintf("Found existing committed content version from API: %s", existingContentVersion))
		}
	}

	updateNeeded := pathChanged || metadataChanged

	tflog.Debug(ctx, "Content version update evaluation result", map[string]interface{}{
		"update_needed":            updateNeeded,
		"path_changed":             pathChanged,
		"metadata_changed":         metadataChanged,
		"existing_content_version": existingContentVersion,
	})

	return updateNeeded, existingContentVersion, nil
}
