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

// encryptFile encrypts the source file using AES-CBC for Intune file upload.
// The encrypted file structure is:
// - HMAC-SHA256 MAC (32 bytes) - calculated over the rest of the file
// - AES Initialization Vector (16 bytes)
// - Encrypted Content (variable length) - using AES-CBC
func encryptFile(sourcePath string) (*EncryptionInfo, error) {
	// Generate AES key (32 bytes for 256-bit key)
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}

	// Generate HMAC key (32 bytes)
	hmacKey := make([]byte, 32)
	if _, err := rand.Read(hmacKey); err != nil {
		return nil, fmt.Errorf("failed to generate HMAC key: %v", err)
	}

	// Create AES cipher block
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Generate random IV (16 bytes for AES)
	iv := make([]byte, aesCipher.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Read source file for hash calculation
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %v", err)
	}

	// Calculate SHA256 digest of the source file
	sourceHash := sha256.New()
	if _, err := io.Copy(sourceHash, sourceFile); err != nil {
		sourceFile.Close()
		return nil, fmt.Errorf("failed to calculate source file hash: %v", err)
	}
	fileDigest := sourceHash.Sum(nil)

	// Reset to beginning of file for encryption
	if _, err := sourceFile.Seek(0, io.SeekStart); err != nil {
		sourceFile.Close()
		return nil, fmt.Errorf("failed to seek source file: %v", err)
	}

	// Create encrypted file
	encryptedPath := sourcePath + ".bin"
	encryptedFile, err := os.Create(encryptedPath)
	if err != nil {
		sourceFile.Close()
		return nil, fmt.Errorf("failed to create encrypted file: %v", err)
	}

	// Reserve space for the HMAC (32 bytes)
	hmacPlaceholder := make([]byte, 32)
	if _, err := encryptedFile.Write(hmacPlaceholder); err != nil {
		sourceFile.Close()
		encryptedFile.Close()
		return nil, fmt.Errorf("failed to write HMAC placeholder: %v", err)
	}

	// Write the IV
	if _, err := encryptedFile.Write(iv); err != nil {
		sourceFile.Close()
		encryptedFile.Close()
		return nil, fmt.Errorf("failed to write IV: %v", err)
	}

	// Create CBC encrypter
	blockMode := cipher.NewCBCEncrypter(aesCipher, iv)

	// Create a streaming crypto writer for encryption
	// We use a custom writer that performs PKCS7 padding and AES-CBC encryption
	cryptoWriter := &cbcPaddingWriter{
		blockMode: blockMode,
		writer:    encryptedFile,
		blockSize: aesCipher.BlockSize(),
	}

	// Copy the source file to the crypto writer (encrypting in the process)
	if _, err := io.Copy(cryptoWriter, sourceFile); err != nil {
		sourceFile.Close()
		encryptedFile.Close()
		return nil, fmt.Errorf("failed to encrypt file: %v", err)
	}

	// Flush the final padded block
	if err := cryptoWriter.Close(); err != nil {
		sourceFile.Close()
		encryptedFile.Close()
		return nil, fmt.Errorf("failed to finalize encryption: %v", err)
	}

	sourceFile.Close()

	// Flush writes to ensure all data is on disk
	if err := encryptedFile.Sync(); err != nil {
		encryptedFile.Close()
		return nil, fmt.Errorf("failed to sync encrypted file: %v", err)
	}

	// Calculate HMAC over the IV and encrypted content
	// Reset file position to after the HMAC placeholder
	if _, err := encryptedFile.Seek(32, io.SeekStart); err != nil {
		encryptedFile.Close()
		return nil, fmt.Errorf("failed to seek for HMAC calculation: %v", err)
	}

	// Initialize HMAC calculator
	hmacCalculator := hmac.New(sha256.New, hmacKey)

	// Read from the file to calculate HMAC (IV + encrypted content)
	if _, err := io.Copy(hmacCalculator, encryptedFile); err != nil {
		encryptedFile.Close()
		return nil, fmt.Errorf("failed to calculate HMAC: %v", err)
	}

	mac := hmacCalculator.Sum(nil)

	// Write the HMAC at the beginning of the file
	if _, err := encryptedFile.Seek(0, io.SeekStart); err != nil {
		encryptedFile.Close()
		return nil, fmt.Errorf("failed to seek to start for writing HMAC: %v", err)
	}

	if _, err := encryptedFile.Write(mac); err != nil {
		encryptedFile.Close()
		return nil, fmt.Errorf("failed to write HMAC: %v", err)
	}

	encryptedFile.Close()

	// Return the encryption metadata
	return &EncryptionInfo{
		EncryptionKey:        base64.StdEncoding.EncodeToString(key),
		FileDigest:           base64.StdEncoding.EncodeToString(fileDigest),
		FileDigestAlgorithm:  "SHA256",
		InitializationVector: base64.StdEncoding.EncodeToString(iv),
		Mac:                  base64.StdEncoding.EncodeToString(mac),
		MacKey:               base64.StdEncoding.EncodeToString(hmacKey),
		ProfileIdentifier:    "ProfileVersion1",
	}, nil
}

// cbcPaddingWriter is a writer that performs PKCS7 padding and AES-CBC encryption
type cbcPaddingWriter struct {
	blockMode cipher.BlockMode
	writer    io.Writer
	buffer    []byte
	blockSize int
}

// Write encrypts and writes data
func (w *cbcPaddingWriter) Write(p []byte) (int, error) {
	// Add incoming data to buffer
	w.buffer = append(w.buffer, p...)

	// Process complete blocks
	blocksToProcess := len(w.buffer) / w.blockSize
	if blocksToProcess > 0 {
		// Get the blocks that can be processed
		toProcess := w.buffer[:blocksToProcess*w.blockSize]

		// Encrypt the blocks
		encrypted := make([]byte, len(toProcess))
		w.blockMode.CryptBlocks(encrypted, toProcess)

		// Write the encrypted blocks
		if _, err := w.writer.Write(encrypted); err != nil {
			return 0, err
		}

		// Keep remaining bytes in buffer
		w.buffer = w.buffer[blocksToProcess*w.blockSize:]
	}

	return len(p), nil
}

// Close adds padding to the last block and writes it
func (w *cbcPaddingWriter) Close() error {
	// Apply PKCS7 padding to the remaining data
	padding := w.blockSize - (len(w.buffer) % w.blockSize)
	padBytes := bytes.Repeat([]byte{byte(padding)}, padding)
	paddedData := append(w.buffer, padBytes...)

	// Encrypt the final padded block
	encrypted := make([]byte, len(paddedData))
	w.blockMode.CryptBlocks(encrypted, paddedData)

	// Write the final encrypted block
	_, err := w.writer.Write(encrypted)
	return err
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
