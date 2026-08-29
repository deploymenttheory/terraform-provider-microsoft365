package graphBetaWin32App

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

const (
	intuneWinDetectionPath = "IntuneWinPackage/Metadata/Detection.xml"
	intuneWinContentsPath  = "IntuneWinPackage/Contents/"
	maxDetectionXMLSize    = 4 << 20
)

type intuneWinDetectionMetadata struct {
	XMLName                xml.Name                    `xml:"ApplicationInfo"`
	FileName               string                      `xml:"FileName"`
	UnencryptedContentSize int64                       `xml:"UnencryptedContentSize"`
	EncryptionInfo         intuneWinEncryptionMetadata `xml:"EncryptionInfo"`
}

type intuneWinEncryptionMetadata struct {
	EncryptionKey        string `xml:"EncryptionKey"`
	MacKey               string `xml:"MacKey"`
	InitializationVector string `xml:"InitializationVector"`
	Mac                  string `xml:"Mac"`
	FileDigest           string `xml:"FileDigest"`
	FileDigestAlgorithm  string `xml:"FileDigestAlgorithm"`
	ProfileIdentifier    string `xml:"ProfileIdentifier"`
}

type preparedWin32Content struct {
	fileName           string
	unencryptedSize    int64
	encryptedSize      int64
	encryptedFilePath  string
	encryptionInfo     *construct.EncryptionInfo
	temporaryDirectory string
}

func extractIntuneWinPackage(packagePath string) (*preparedWin32Content, error) {
	archive, err := zip.OpenReader(packagePath)
	if err != nil {
		return nil, fmt.Errorf("open .intunewin package: %w", err)
	}
	defer archive.Close()

	var detectionFile *zip.File
	for _, file := range archive.File {
		if file.Name == intuneWinDetectionPath {
			if detectionFile != nil {
				return nil, fmt.Errorf(".intunewin package contains multiple Detection.xml files")
			}
			detectionFile = file
		}
	}
	if detectionFile == nil {
		return nil, fmt.Errorf(".intunewin package does not contain %s", intuneWinDetectionPath)
	}
	if detectionFile.UncompressedSize64 > maxDetectionXMLSize {
		return nil, fmt.Errorf("Detection.xml exceeds the maximum supported size")
	}

	metadataReader, err := detectionFile.Open()
	if err != nil {
		return nil, fmt.Errorf("open Detection.xml: %w", err)
	}
	defer metadataReader.Close()

	var metadata intuneWinDetectionMetadata
	if err := xml.NewDecoder(io.LimitReader(metadataReader, maxDetectionXMLSize)).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("parse Detection.xml: %w", err)
	}
	if err := validateIntuneWinMetadata(&metadata); err != nil {
		return nil, err
	}

	contentPath := intuneWinContentsPath + metadata.FileName
	var encryptedContent *zip.File
	for _, file := range archive.File {
		if file.Name == contentPath {
			if encryptedContent != nil {
				return nil, fmt.Errorf(".intunewin package contains multiple encrypted content files named %q", metadata.FileName)
			}
			encryptedContent = file
		}
	}
	if encryptedContent == nil {
		return nil, fmt.Errorf(".intunewin package does not contain encrypted content %q", contentPath)
	}
	if encryptedContent.UncompressedSize64 == 0 || encryptedContent.UncompressedSize64 > math.MaxInt64 {
		return nil, fmt.Errorf("encrypted content has an invalid size")
	}

	contentReader, err := encryptedContent.Open()
	if err != nil {
		return nil, fmt.Errorf("open encrypted content: %w", err)
	}
	defer contentReader.Close()

	temporaryFile, err := os.CreateTemp("", "microsoft365-win32-content-*")
	if err != nil {
		return nil, fmt.Errorf("create temporary encrypted content file: %w", err)
	}

	encryptedSize, copyErr := io.Copy(temporaryFile, contentReader)
	closeErr := temporaryFile.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(temporaryFile.Name())
		if copyErr != nil {
			return nil, fmt.Errorf("extract encrypted content: %w", copyErr)
		}
		return nil, fmt.Errorf("close temporary encrypted content file: %w", closeErr)
	}
	if encryptedSize != int64(encryptedContent.UncompressedSize64) {
		_ = os.Remove(temporaryFile.Name())
		return nil, fmt.Errorf("encrypted content size does not match the package metadata")
	}

	return &preparedWin32Content{
		fileName:          metadata.FileName,
		unencryptedSize:   metadata.UnencryptedContentSize,
		encryptedSize:     encryptedSize,
		encryptedFilePath: temporaryFile.Name(),
		encryptionInfo: &construct.EncryptionInfo{
			EncryptionKey:        metadata.EncryptionInfo.EncryptionKey,
			FileDigest:           metadata.EncryptionInfo.FileDigest,
			FileDigestAlgorithm:  metadata.EncryptionInfo.FileDigestAlgorithm,
			InitializationVector: metadata.EncryptionInfo.InitializationVector,
			Mac:                  metadata.EncryptionInfo.Mac,
			MacKey:               metadata.EncryptionInfo.MacKey,
			ProfileIdentifier:    metadata.EncryptionInfo.ProfileIdentifier,
		},
	}, nil
}

func validateIntuneWinMetadata(metadata *intuneWinDetectionMetadata) error {
	if metadata.FileName == "" || path.Base(metadata.FileName) != metadata.FileName || strings.Contains(metadata.FileName, "\\") || metadata.FileName == "." || metadata.FileName == ".." {
		return fmt.Errorf("Detection.xml contains an invalid encrypted content file name")
	}
	if metadata.UnencryptedContentSize <= 0 {
		return fmt.Errorf("Detection.xml contains an invalid unencrypted content size")
	}

	fields := []struct {
		name string
		data string
		size int
	}{
		{name: "encryption key", data: metadata.EncryptionInfo.EncryptionKey, size: 32},
		{name: "MAC key", data: metadata.EncryptionInfo.MacKey, size: 32},
		{name: "initialization vector", data: metadata.EncryptionInfo.InitializationVector, size: 16},
		{name: "MAC", data: metadata.EncryptionInfo.Mac, size: 32},
		{name: "file digest", data: metadata.EncryptionInfo.FileDigest, size: 32},
	}
	for _, field := range fields {
		decoded, err := base64.StdEncoding.DecodeString(field.data)
		if err != nil || len(decoded) != field.size {
			return fmt.Errorf("Detection.xml contains an invalid %s", field.name)
		}
	}
	if !strings.EqualFold(metadata.EncryptionInfo.FileDigestAlgorithm, "SHA256") {
		return fmt.Errorf("Detection.xml contains an unsupported file digest algorithm %q", metadata.EncryptionInfo.FileDigestAlgorithm)
	}
	if metadata.EncryptionInfo.ProfileIdentifier == "" {
		return fmt.Errorf("Detection.xml does not contain an encryption profile identifier")
	}

	return nil
}

func (installer *preparedWin32Content) contentFile() graphmodels.MobileAppContentFileable {
	contentFile := graphmodels.NewMobileAppContentFile()
	contentFile.SetName(&installer.fileName)
	contentFile.SetSize(&installer.unencryptedSize)
	contentFile.SetSizeEncrypted(&installer.encryptedSize)

	falseValue := false
	contentFile.SetIsDependency(&falseValue)
	contentFile.SetIsFrameworkFile(&falseValue)

	return contentFile
}

func (installer *preparedWin32Content) cleanup() error {
	if installer.temporaryDirectory != "" {
		return os.RemoveAll(installer.temporaryDirectory)
	}
	return os.Remove(installer.encryptedFilePath)
}

func cleanupWin32Content(ctx context.Context, installer *preparedWin32Content) {
	if err := installer.cleanup(); err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to remove temporary encrypted Win32 content file: %v", err))
	}
}

func hasInstallerSource(data *Win32LobAppResourceModel) bool {
	return !data.AppInstaller.IsNull() || !data.AppInstallerZip.IsNull()
}

func installerSourceChanged(plan, state *Win32LobAppResourceModel) bool {
	return !plan.AppInstaller.Equal(state.AppInstaller) || !plan.AppInstallerZip.Equal(state.AppInstallerZip)
}

func selectInstallerSource(data *Win32LobAppResourceModel) (types.Object, bool, error) {
	if data.AppInstaller.IsUnknown() || data.AppInstallerZip.IsUnknown() {
		return types.Object{}, false, fmt.Errorf("installer source must be known before upload")
	}
	if !data.AppInstaller.IsNull() && !data.AppInstallerZip.IsNull() {
		return types.Object{}, false, fmt.Errorf("app_installer and app_installer_zip are mutually exclusive")
	}
	source, plainZip := data.AppInstaller, false
	if !data.AppInstallerZip.IsNull() {
		source, plainZip = data.AppInstallerZip, true
	}
	if source.IsNull() {
		return source, plainZip, fmt.Errorf("creating or uploading an application requires app_installer or app_installer_zip")
	}
	var metadata sharedmodels.MobileAppMetaDataResourceModel
	if diags := source.As(context.Background(), &metadata, basetypes.ObjectAsOptions{}); diags.HasError() {
		return source, plainZip, fmt.Errorf("invalid installer source: %v", diags.Errors())
	}
	if metadata.InstallerFilePathSource.IsUnknown() || metadata.InstallerURLSource.IsUnknown() {
		return source, plainZip, fmt.Errorf("installer path or URL must be known before upload")
	}
	local, remote := metadata.InstallerFilePathSource, metadata.InstallerURLSource
	if local.IsNull() == remote.IsNull() || (!local.IsNull() && strings.TrimSpace(local.ValueString()) == "") || (!remote.IsNull() && strings.TrimSpace(remote.ValueString()) == "") {
		return source, plainZip, fmt.Errorf("set exactly one nonempty installer_file_path_source or installer_url_source")
	}
	return source, plainZip, nil
}

func prepareInstaller(ctx context.Context, data *Win32LobAppResourceModel) (*preparedWin32Content, error) {
	source, plainZip, err := selectInstallerSource(data)
	if err != nil {
		return nil, err
	}
	sourcePath, temporaryFile, err := helpers.SetInstallerSourcePath(ctx, source)
	if err != nil {
		return nil, err
	}
	defer helpers.CleanupTempFile(ctx, temporaryFile)
	if plainZip {
		return prepareZipContent(ctx, sourcePath)
	}
	content, err := extractIntuneWinPackage(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("app_installer requires a prepackaged .intunewin file: %w; for an unencrypted ZIP (including one renamed .intunewin), use app_installer_zip", err)
	}
	return content, nil
}

// prepareZipContent encrypts a private copy of a plain installer ZIP. It never
// writes beside the source or falls back from a malformed Content Prep package.
func prepareZipContent(ctx context.Context, sourcePath string) (*preparedWin32Content, error) {
	archive, err := zip.OpenReader(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("open unencrypted installer ZIP: %w", err)
	}
	defer archive.Close()
	hasFile := false
	for _, entry := range archive.File {
		name := strings.ReplaceAll(entry.Name, "\\", "/")
		clean := path.Clean(name)
		if strings.EqualFold(clean, "IntuneWinPackage") || strings.HasPrefix(strings.ToLower(clean), "intunewinpackage/") {
			return nil, fmt.Errorf("app_installer_zip cannot contain a Content Prep Tool package; use app_installer for prepackaged .intunewin files")
		}
		if entry.Flags&1 != 0 {
			return nil, fmt.Errorf("app_installer_zip requires an unencrypted ZIP; entry %q is encrypted", entry.Name)
		}
		if strings.HasPrefix(clean, "/") || clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, ":") || entry.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("installer ZIP contains unsafe entry %q", entry.Name)
		}
		if !entry.FileInfo().IsDir() {
			hasFile = true
		}
	}
	if !hasFile {
		return nil, fmt.Errorf("installer ZIP contains no files")
	}
	temporaryDirectory, err := os.MkdirTemp("", "microsoft365-win32-zip-*")
	if err != nil {
		return nil, err
	}
	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(temporaryDirectory)
		}
	}()
	// Graph uses an .intunewin content filename even when the input is a plain ZIP.
	base := filepath.Base(sourcePath)
	fileName := strings.TrimSuffix(base, filepath.Ext(base)) + ".intunewin"
	copiedPath := filepath.Join(temporaryDirectory, fileName)
	source, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer source.Close()
	copied, err := os.OpenFile(copiedPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	_, copyErr := io.Copy(copied, source)
	closeErr := copied.Close()
	if copyErr != nil {
		return nil, copyErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	file, encryptionInfo, err := construct.EncryptMobileAppAndConstructFileContentMetadata(ctx, copiedPath)
	if err != nil {
		return nil, err
	}
	if err := os.Remove(copiedPath); err != nil {
		return nil, err
	}
	success = true
	return &preparedWin32Content{
		fileName:           fileName,
		unencryptedSize:    *file.GetSize(),
		encryptedSize:      *file.GetSizeEncrypted(),
		encryptedFilePath:  copiedPath + ".bin",
		encryptionInfo:     encryptionInfo,
		temporaryDirectory: temporaryDirectory,
	}, nil
}
