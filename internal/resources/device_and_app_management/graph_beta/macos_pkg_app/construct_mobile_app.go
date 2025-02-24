package graphBetaMacOSPKGApp

import (
	"bytes"
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

// pkcs7Pad appends PKCS7 padding to data to make its length a multiple of blockSize.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padText...)
}

// encryptFile encrypts the source file using AES-CBC with PKCS7 padding.
// It writes a 32-byte HMAC placeholder, then the IV, then the encrypted data.
// After encrypting, it computes the HMAC over the IV and ciphertext, and writes it in the placeholder.
func encryptFile(sourcePath string) (*EncryptionInfo, error) {
	// Generate AES key (32 bytes for 256-bit key)
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}

	// Create AES cipher block
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}
	blockSize := aesCipher.BlockSize() // should be 16 bytes

	// Generate IV (same length as block size)
	iv := make([]byte, blockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Generate HMAC key (32 bytes)
	hmacKey := make([]byte, 32)
	if _, err := rand.Read(hmacKey); err != nil {
		return nil, fmt.Errorf("failed to generate HMAC key: %v", err)
	}

	// Read source file data
	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %v", err)
	}

	// Calculate SHA256 digest of the source data
	sourceHash := sha256.Sum256(sourceData)

	// Apply PKCS7 padding to the source data
	paddedData := pkcs7Pad(sourceData, blockSize)

	// Create encrypted file (same as sourcePath+".bin")
	encryptedPath := sourcePath + ".bin"
	encrypted, err := os.Create(encryptedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create encrypted file: %v", err)
	}
	// Ensure the file is closed later
	defer encrypted.Close()

	// Write a 32-byte placeholder for the HMAC (SHA256 produces 32 bytes)
	hmacPlaceholder := make([]byte, 32)
	if _, err := encrypted.Write(hmacPlaceholder); err != nil {
		return nil, fmt.Errorf("failed to write HMAC placeholder: %v", err)
	}

	// Write the IV
	if _, err := encrypted.Write(iv); err != nil {
		return nil, fmt.Errorf("failed to write IV: %v", err)
	}

	// Encrypt the padded data using AES-CBC
	mode := cipher.NewCBCEncrypter(aesCipher, iv)
	encryptedData := make([]byte, len(paddedData))
	mode.CryptBlocks(encryptedData, paddedData)

	// Write the encrypted data
	if _, err := encrypted.Write(encryptedData); err != nil {
		return nil, fmt.Errorf("failed to write encrypted data: %v", err)
	}

	// At this point the file layout is:
	// [0..31]       : Placeholder for HMAC (32 bytes)
	// [32..(32+blockSize-1)]: IV
	// [remaining]   : Encrypted data

	// Flush file writes (by closing and reopening for reading)
	if err := encrypted.Sync(); err != nil {
		return nil, fmt.Errorf("failed to sync encrypted file: %v", err)
	}

	// Reopen the file for reading the portion for HMAC calculation.
	// We need to read from offset 32 until end.
	encryptedForHMAC, err := os.Open(encryptedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open encrypted file for HMAC: %v", err)
	}
	defer encryptedForHMAC.Close()

	// Seek to offset 32 (skip the HMAC placeholder)
	if _, err := encryptedForHMAC.Seek(32, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek for HMAC calculation: %v", err)
	}

	// Compute HMAC over the IV and encrypted data
	hmacHash := hmac.New(sha256.New, hmacKey)
	if _, err := io.Copy(hmacHash, encryptedForHMAC); err != nil {
		return nil, fmt.Errorf("failed to calculate HMAC: %v", err)
	}
	macSum := hmacHash.Sum(nil)

	// Open the file again for writing the computed HMAC in place of the placeholder.
	encryptedForWrite, err := os.OpenFile(encryptedPath, os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open encrypted file for writing HMAC: %v", err)
	}
	defer encryptedForWrite.Close()

	// Write the computed HMAC at the beginning of the file.
	if _, err := encryptedForWrite.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start for writing HMAC: %v", err)
	}
	if _, err := encryptedForWrite.Write(macSum); err != nil {
		return nil, fmt.Errorf("failed to write HMAC: %v", err)
	}

	// Return the encryption metadata in the same format as the PowerShell script.
	return &EncryptionInfo{
		EncryptionKey:        base64.StdEncoding.EncodeToString(key),
		FileDigest:           base64.StdEncoding.EncodeToString(sourceHash[:]),
		FileDigestAlgorithm:  "SHA256",
		InitializationVector: base64.StdEncoding.EncodeToString(iv),
		Mac:                  base64.StdEncoding.EncodeToString(macSum),
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
