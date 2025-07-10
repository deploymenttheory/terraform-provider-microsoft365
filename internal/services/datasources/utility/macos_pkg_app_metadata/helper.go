package utilityMacOSPKGAppMetadata

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/common"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/macos_pkg/xar"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// bytesToMB converts bytes to megabytes as a float64 value
func bytesToMB(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}

// extractMetadataFromFile extracts metadata from a local PKG file
func extractMetadataFromFile(ctx context.Context, filePath string) (*xar.InstallerMetadata, []byte, []byte, error) {
	tflog.Debug(ctx, fmt.Sprintf("Starting metadata extraction from PKG file: %s", filePath))

	// Validate file exists and is accessible
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error accessing file: %v", err))
		if os.IsNotExist(err) {
			return nil, nil, nil, fmt.Errorf("file not found: %s", filePath)
		}
		return nil, nil, nil, fmt.Errorf("error accessing file: %w", err)
	}

	// Validate it's not a directory
	if fileInfo.IsDir() {
		tflog.Error(ctx, fmt.Sprintf("Path points to a directory, not a file: %s", filePath))
		return nil, nil, nil, fmt.Errorf("path points to a directory, not a file: %s", filePath)
	}

	// Open the file for checksum calculation
	file, err := os.Open(filePath)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error opening file: %v", err))
		return nil, nil, nil, fmt.Errorf("error opening file: %w", err)
	}

	// Calculate MD5 and SHA256 checksums
	md5Hash := md5.New()
	sha256Hash := sha256.New()
	multiWriter := io.MultiWriter(md5Hash, sha256Hash)

	if _, err := io.Copy(multiWriter, file); err != nil {
		file.Close()
		return nil, nil, nil, fmt.Errorf("error calculating checksums: %w", err)
	}

	// Get the checksums
	md5Checksum := md5Hash.Sum(nil)
	sha256Checksum := sha256Hash.Sum(nil)

	// Reset file position
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		file.Close()
		return nil, nil, nil, fmt.Errorf("error resetting file position: %w", err)
	}

	// Create a TempFileReader for XAR processing
	tfr, err := xar.NewTempFileReader(file, os.TempDir)
	file.Close() // Close the file as it's now copied to the temp file
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating temp file reader: %w", err)
	}
	defer tfr.Close()

	// Extract metadata
	metadata, err := xar.ExtractXARMetadata(tfr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error extracting XAR metadata: %w", err)
	}

	// Set the file size
	metadata.Size = fileInfo.Size()

	// Log the successful extraction
	sizeMB := bytesToMB(fileInfo.Size())
	tflog.Debug(ctx, fmt.Sprintf("Successfully extracted metadata - Bundle ID: '%s', Version: '%s', Size: %.2f MB",
		metadata.BundleIdentifier, metadata.Version, sizeMB))

	return metadata, md5Checksum, sha256Checksum, nil
}

// extractMetadataFromURL downloads a PKG file from a URL and extracts metadata
func extractMetadataFromURL(ctx context.Context, sourceURL string) (*xar.InstallerMetadata, []byte, []byte, error) {
	tflog.Debug(ctx, fmt.Sprintf("Downloading PKG file from URL: %s", sourceURL))

	downloadedPath, err := common.DownloadFile(sourceURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to download PKG file from %s: %w", sourceURL, err)
	}

	tempFile := helpers.TempFileInfo{
		FilePath:      downloadedPath,
		ShouldCleanup: true,
	}
	defer helpers.CleanupTempFile(ctx, tempFile)

	tflog.Debug(ctx, fmt.Sprintf("Successfully downloaded PKG file to temporary location: %s", downloadedPath))

	// Now process the downloaded file the same way as local files
	return extractMetadataFromFile(ctx, downloadedPath)
}
