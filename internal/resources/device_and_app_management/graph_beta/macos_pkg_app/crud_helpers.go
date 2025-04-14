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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

// setInstallerSourcePath inspects the resource model's app_metadata and returns the actual installer file path
// along with information about whether it needs cleanup.
func setInstallerSourcePath(ctx context.Context, metadataObj types.Object) (string, TempFileInfo, error) {
	var fileInfo TempFileInfo
	var metadata sharedmodels.MobileAppMetaDataResourceModel

	// Parse metadata object if available
	if !metadataObj.IsNull() {
		diags := metadataObj.As(ctx, &metadata, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			tflog.Warn(ctx, "Failed to parse app_metadata, proceeding with empty struct", map[string]interface{}{
				"errors": diags.Errors(),
			})
		}
	}

	tflog.Debug(ctx, "setInstallerSourcePath input", map[string]interface{}{
		"installer_file_path_source": metadata.InstallerFilePathSource.String(),
		"installer_url_source":       metadata.InstallerURLSource.String(),
	})

	// Skip if both are unknown (planning phase)
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

	// Nothing provided
	return "", fileInfo, fmt.Errorf("installer file not provided; please supply either a local file path or a URL")
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

// GetAppMetadata computes and sets the metadata for the installer file (size and checksums)
// used for evaluation for content version updates and other purposes
func GetAppMetadata(
	ctx context.Context,
	installerSourcePath string,
	existingMetadata *sharedmodels.MobileAppMetaDataResourceModel,
) (*sharedmodels.MobileAppMetaDataResourceModel, error) {
	fileInfo, err := os.Stat(installerSourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file information: %v", err)
	}

	currentSize := fileInfo.Size()

	tflog.Debug(ctx, "Computed installer file size", map[string]interface{}{
		"file_path":  installerSourcePath,
		"size_bytes": currentSize,
	})

	md5Hash := md5.New()
	sha256Hash := sha256.New()

	file, err := os.Open(installerSourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for checksums: %v", err)
	}
	defer file.Close()

	multiWriter := io.MultiWriter(md5Hash, sha256Hash)

	if _, err := io.Copy(multiWriter, file); err != nil {
		return nil, fmt.Errorf("failed to read file for checksums: %v", err)
	}

	md5Checksum := hex.EncodeToString(md5Hash.Sum(nil))
	sha256Checksum := hex.EncodeToString(sha256Hash.Sum(nil))

	tflog.Debug(ctx, "Computed installer file checksums", map[string]interface{}{
		"file_path":       installerSourcePath,
		"md5_checksum":    md5Checksum,
		"sha256_checksum": sha256Checksum,
	})

	return &sharedmodels.MobileAppMetaDataResourceModel{
		InstallerFilePathSource: existingMetadata.InstallerFilePathSource,
		InstallerURLSource:      existingMetadata.InstallerURLSource,
		InstallerSizeInBytes:    types.Int64Value(currentSize),
		InstallerMD5Checksum:    types.StringValue(md5Checksum),
		InstallerSHA256Checksum: types.StringValue(sha256Checksum),
	}, nil
}

// evaluateIfContentVersionUpdateRequired determines if a new content version needs to be created and uploaded
// based on source path changes and metadata differences
func evaluateIfContentVersionUpdateRequired(ctx context.Context, object *MacOSPKGAppResourceModel, state *MacOSPKGAppResourceModel, client *msgraphbetasdk.GraphServiceClient) (bool, string, error) {

	pathChanged := false
	metadataChanged := false
	existingContentVersion := ""

	// Check installer source path changes
	var objectMetadata, stateMetadata sharedmodels.MobileAppMetaDataResourceModel

	if !object.AppMetadata.IsNull() {
		diags := object.AppMetadata.As(ctx, &objectMetadata, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return false, "", fmt.Errorf("failed to parse object metadata: %v", diags)
		}
	}

	if !state.AppMetadata.IsNull() {
		diags := state.AppMetadata.As(ctx, &stateMetadata, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return false, "", fmt.Errorf("failed to parse state metadata: %v", diags)
		}
	}

	// Path changes
	if !objectMetadata.InstallerFilePathSource.Equal(stateMetadata.InstallerFilePathSource) {
		tflog.Debug(ctx, "Installer file path has changed", map[string]interface{}{
			"old_path": stateMetadata.InstallerFilePathSource.ValueString(),
			"new_path": objectMetadata.InstallerFilePathSource.ValueString(),
		})
		pathChanged = true
	}

	if !objectMetadata.InstallerURLSource.Equal(stateMetadata.InstallerURLSource) {
		tflog.Debug(ctx, "Installer URL has changed", map[string]interface{}{
			"old_url": stateMetadata.InstallerURLSource.ValueString(),
			"new_url": objectMetadata.InstallerURLSource.ValueString(),
		})
		pathChanged = true
	}

	// If no installer source is provided at all → skip content evaluation
	if (objectMetadata.InstallerFilePathSource.IsNull() || objectMetadata.InstallerFilePathSource.ValueString() == "") &&
		(objectMetadata.InstallerURLSource.IsNull() || objectMetadata.InstallerURLSource.ValueString() == "") {
		tflog.Debug(ctx, "No installer source provided, skipping content evaluation")
		return false, "", nil
	}

	// Check metadata changes
	if !object.AppMetadata.IsNull() && !state.AppMetadata.IsNull() {
		if objectMetadata.InstallerSizeInBytes.ValueInt64() != stateMetadata.InstallerSizeInBytes.ValueInt64() {
			tflog.Debug(ctx, "Installer size has changed", map[string]interface{}{
				"previous_size": stateMetadata.InstallerSizeInBytes.ValueInt64(),
				"current_size":  objectMetadata.InstallerSizeInBytes.ValueInt64(),
			})
			metadataChanged = true
		}

		if objectMetadata.InstallerSHA256Checksum.ValueString() != stateMetadata.InstallerSHA256Checksum.ValueString() {
			tflog.Debug(ctx, "Installer SHA256 checksum has changed", map[string]interface{}{
				"previous_checksum": stateMetadata.InstallerSHA256Checksum.ValueString(),
				"current_checksum":  objectMetadata.InstallerSHA256Checksum.ValueString(),
			})
			metadataChanged = true
		}
	} else {
		tflog.Debug(ctx, "Metadata is missing from one of the objects, assuming content needs update")
		metadataChanged = true
	}

	// Retrieve existing committed content version (from API)
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
