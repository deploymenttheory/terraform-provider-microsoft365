package graphBetaMacOSPKGApp

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// EncryptionInfo holds the encryption metadata for the file
type EncryptionInfo struct {
	EncryptionKey        string `json:"encryptionKey"`
	FileDigest           string `json:"fileDigest"`
	FileDigestAlgorithm  string `json:"fileDigestAlgorithm"`
	InitializationVector string `json:"initializationVector"`
	Mac                  string `json:"mac"`
	MacKey               string `json:"macKey"`
	ProfileIdentifier    string `json:"profileIdentifier"`
}

// constructMobileAppContentFile creates and initializes the content file process
func constructMobileAppContentFile(ctx context.Context, filePath string) (graphmodels.MobileAppContentFileable, *EncryptionInfo, error) {
	tflog.Debug(ctx, "Starting content file construction", map[string]interface{}{
		"file_path": filePath,
	})

	// Create initial content file object
	contentFile := graphmodels.NewMobileAppContentFile()

	// Set basic properties
	fileName := filepath.Base(filePath)
	contentFile.SetName(&fileName)

	// Get file info for size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get file info: %s", err)
	}

	// Set size properties
	size := fileInfo.Size()
	contentFile.SetSize(&size)
	contentFile.SetSizeInBytes(&size)

	// Set required flags
	falseValue := false
	contentFile.SetIsDependency(&falseValue)
	contentFile.SetIsFrameworkFile(&falseValue)

	// Encrypt the file
	tflog.Debug(ctx, "Starting file encryption")
	encryptionInfo, err := encryptFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt file: %v", err)
	}
	tflog.Debug(ctx, "File encryption completed")

	// Get encrypted file size
	encryptedFileInfo, err := os.Stat(filePath + ".bin")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get encrypted file info: %v", err)
	}

	encryptedSize := encryptedFileInfo.Size()
	contentFile.SetSizeEncrypted(&encryptedSize)
	contentFile.SetSizeEncryptedInBytes(&encryptedSize)

	tflog.Debug(ctx, "Intune Mobile App encrypted package Content file construction completed", map[string]interface{}{
		"original_size":  size,
		"encrypted_size": encryptedSize,
		"file_name":      fileName,
	})

	return contentFile, encryptionInfo, nil
}

// encryptFile encrypts the source file using AES encryption for Intune upload
func encryptFile(sourcePath string) (*EncryptionInfo, error) {
	// Generate AES key
	key := make([]byte, 32) // 256-bit key
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Generate HMAC key
	hmacKey := make([]byte, 32)
	if _, err := rand.Read(hmacKey); err != nil {
		return nil, fmt.Errorf("failed to generate HMAC key: %v", err)
	}

	// Read source file
	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %v", err)
	}

	// Calculate SHA256 of source
	sourceHash := sha256.Sum256(sourceData)

	// Generate IV
	iv := make([]byte, aesCipher.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Create encrypted file path
	encryptedPath := sourcePath + ".bin"

	// Create encrypted file
	encrypted, err := os.Create(encryptedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create encrypted file: %v", err)
	}
	defer encrypted.Close()

	// Initialize HMAC
	mac := hmac.New(sha256.New, hmacKey)

	// Write placeholder for HMAC (will be filled later)
	hmacSize := int64(mac.Size())
	if _, err := encrypted.Write(make([]byte, hmacSize)); err != nil {
		return nil, fmt.Errorf("failed to write HMAC placeholder: %v", err)
	}

	// Write IV
	if _, err := encrypted.Write(iv); err != nil {
		return nil, fmt.Errorf("failed to write IV: %v", err)
	}

	// Create AES encryptor
	stream := cipher.NewCTR(aesCipher, iv)
	writer := &cipher.StreamWriter{S: stream, W: encrypted}

	// Write encrypted data
	if _, err := writer.Write(sourceData); err != nil {
		return nil, fmt.Errorf("failed to write encrypted data: %v", err)
	}

	// Calculate HMAC
	if _, err := encrypted.Seek(hmacSize, 0); err != nil {
		return nil, fmt.Errorf("failed to seek for HMAC calculation: %v", err)
	}

	if _, err := io.Copy(mac, encrypted); err != nil {
		return nil, fmt.Errorf("failed to calculate HMAC: %v", err)
	}

	// Write HMAC at the beginning
	if _, err := encrypted.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek to start: %v", err)
	}

	if _, err := encrypted.Write(mac.Sum(nil)); err != nil {
		return nil, fmt.Errorf("failed to write HMAC: %v", err)
	}

	return &EncryptionInfo{
		EncryptionKey:        base64.StdEncoding.EncodeToString(key),
		FileDigest:           base64.StdEncoding.EncodeToString(sourceHash[:]),
		FileDigestAlgorithm:  "SHA256",
		InitializationVector: base64.StdEncoding.EncodeToString(iv),
		Mac:                  base64.StdEncoding.EncodeToString(mac.Sum(nil)),
		MacKey:               base64.StdEncoding.EncodeToString(hmacKey),
		ProfileIdentifier:    "ProfileVersion1",
	}, nil
}

// NewEncryptedIntuneMobileAppFileUpload creates a commit request for an encrypted Intune application file
// using the provided encryption information
func NewEncryptedIntuneMobileAppFileUpload(encryptionInfo *EncryptionInfo) (deviceappmanagement.MobileAppsItemGraphMacOSPkgAppContentVersionsItemFilesItemCommitPostRequestBodyable, error) {
	// Create the Graph API commit request body
	requestBody := deviceappmanagement.NewMobileAppsItemGraphMacOSPkgAppContentVersionsItemFilesItemCommitPostRequestBody()

	// Create and configure the encryption info
	fileEncryptionInfo := graphmodels.NewFileEncryptionInfo()

	// Decode base64 values to byte arrays
	encKey, err := base64.StdEncoding.DecodeString(encryptionInfo.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %v", err)
	}

	digest, err := base64.StdEncoding.DecodeString(encryptionInfo.FileDigest)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file digest: %v", err)
	}

	iv, err := base64.StdEncoding.DecodeString(encryptionInfo.InitializationVector)
	if err != nil {
		return nil, fmt.Errorf("failed to decode initialization vector: %v", err)
	}

	macBytes, err := base64.StdEncoding.DecodeString(encryptionInfo.Mac)
	if err != nil {
		return nil, fmt.Errorf("failed to decode mac: %v", err)
	}

	macKeyBytes, err := base64.StdEncoding.DecodeString(encryptionInfo.MacKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode mac key: %v", err)
	}

	// Set all the encryption info values
	fileEncryptionInfo.SetEncryptionKey(encKey)
	fileEncryptionInfo.SetFileDigest(digest)
	fileEncryptionInfo.SetFileDigestAlgorithm(&encryptionInfo.FileDigestAlgorithm)
	fileEncryptionInfo.SetInitializationVector(iv)
	fileEncryptionInfo.SetMac(macBytes)
	fileEncryptionInfo.SetMacKey(macKeyBytes)
	fileEncryptionInfo.SetProfileIdentifier(&encryptionInfo.ProfileIdentifier)

	// Set the encryption info on the request body
	requestBody.SetFileEncryptionInfo(fileEncryptionInfo)

	return requestBody, nil
}
