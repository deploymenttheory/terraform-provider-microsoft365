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
	"strings"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
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

type intuneWinPackage struct {
	fileName          string
	unencryptedSize   int64
	encryptedSize     int64
	encryptedFilePath string
	encryptionInfo    *construct.EncryptionInfo
}

func extractIntuneWinPackage(packagePath string) (*intuneWinPackage, error) {
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

	return &intuneWinPackage{
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

func (installer *intuneWinPackage) contentFile() graphmodels.MobileAppContentFileable {
	contentFile := graphmodels.NewMobileAppContentFile()
	contentFile.SetName(&installer.fileName)
	contentFile.SetSize(&installer.unencryptedSize)
	contentFile.SetSizeEncrypted(&installer.encryptedSize)

	falseValue := false
	contentFile.SetIsDependency(&falseValue)
	contentFile.SetIsFrameworkFile(&falseValue)

	return contentFile
}

func (installer *intuneWinPackage) cleanup() error {
	return os.Remove(installer.encryptedFilePath)
}

func cleanupIntuneWinPackage(ctx context.Context, installer *intuneWinPackage) {
	if err := installer.cleanup(); err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to remove temporary encrypted Win32 content file: %v", err))
	}
}
