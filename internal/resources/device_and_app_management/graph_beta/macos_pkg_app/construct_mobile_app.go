package graphBetaMacOSPKGApp

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// EncryptionInfo holds the encryption metadata for the file.
type EncryptionInfo struct {
	EncryptionKey        string `json:"encryptionKey"`
	FileDigest           string `json:"fileDigest"`
	FileDigestAlgorithm  string `json:"fileDigestAlgorithm"`
	InitializationVector string `json:"initializationVector"`
	Mac                  string `json:"mac"`
	MacKey               string `json:"macKey"`
	ProfileIdentifier    string `json:"profileIdentifier"`
}

// EncryptedFileAnalysis contains hex representations of key portions of the encrypted file.
type EncryptedFileAnalysis struct {
	FileLength       int64  `json:"fileLength"`
	HMACHex          string `json:"hmacHex"`
	IVHex            string `json:"ivHex"`
	CiphertextSample string `json:"ciphertextSample"`
	FullHeaderHex    string `json:"fullHeaderHex"`
}

// constructMobileAppContentFile maps the Terraform schema to the SDK model,
// encrypts the installer file, logs its hex details and uploads the encrypted file.
func constructMobileAppContentFile(ctx context.Context, filePath string) (graphmodels.MobileAppContentFileable, *EncryptionInfo, error) {
	tflog.Debug(ctx, fmt.Sprintf("Starting content file construction for file: %s", filePath), map[string]interface{}{"file_path": filePath})

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
	tflog.Debug(ctx, "Starting file encryption", map[string]interface{}{"file": filePath})
	encryptionInfo, err := encryptFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt file: %v", err)
	}
	tflog.Debug(ctx, "File encryption completed", map[string]interface{}{"encrypted_file": filePath + ".bin"})

	// Analyze the encrypted file and log details
	encryptedFilePath := filePath + ".bin"
	analysis, err := analyzeEncryptedFileHex(encryptedFilePath)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error analyzing encrypted file: %v", err))
	} else {
		tflog.Debug(ctx, "Encrypted file analysis",
			map[string]interface{}{
				"FileLength":       analysis.FileLength,
				"HMACHex":          analysis.HMACHex,
				"IVHex":            analysis.IVHex,
				"CiphertextSample": analysis.CiphertextSample,
				"FullHeaderHex":    analysis.FullHeaderHex,
			})
	}

	// Get encrypted file size
	encryptedFileInfo, err := os.Stat(encryptedFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get encrypted file info: %v", err)
	}
	encryptedSize := encryptedFileInfo.Size()
	contentFile.SetSizeEncrypted(&encryptedSize)
	contentFile.SetSizeEncryptedInBytes(&encryptedSize)

	tflog.Debug(ctx, "Intune Mobile App encrypted package content file construction completed",
		map[string]interface{}{
			"original_size":  size,
			"encrypted_size": encryptedSize,
			"file_name":      fileName,
		})

	return contentFile, encryptionInfo, nil
}

// encryptFile exactly replicates the PowerShell implementation of Intune file encryption
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

	// Create target file path
	targetFilePath := sourcePath + ".bin"

	// Open source file for reading and calculate SHA256 digest
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %v", err)
	}

	// Calculate SHA256 digest of source file
	sourceDigest := sha256.New()
	if _, err := io.Copy(sourceDigest, sourceFile); err != nil {
		sourceFile.Close()
		return nil, fmt.Errorf("failed to calculate source digest: %v", err)
	}
	sourceDigestValue := sourceDigest.Sum(nil)

	// Reset source file position to beginning
	if _, err := sourceFile.Seek(0, io.SeekStart); err != nil {
		sourceFile.Close()
		return nil, fmt.Errorf("failed to seek to beginning of source file: %v", err)
	}

	// Create AES cipher
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		sourceFile.Close()
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Generate random IV
	iv := make([]byte, aesCipher.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		sourceFile.Close()
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Create target file
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		sourceFile.Close()
		return nil, fmt.Errorf("failed to create target file: %v", err)
	}

	// Write placeholder for HMAC (32 bytes of zeros)
	hmacLength := 32 // HMAC-SHA256 is 32 bytes
	hmacPlaceholder := make([]byte, hmacLength)
	if _, err := targetFile.Write(hmacPlaceholder); err != nil {
		sourceFile.Close()
		targetFile.Close()
		return nil, fmt.Errorf("failed to write HMAC placeholder: %v", err)
	}

	// Write IV
	if _, err := targetFile.Write(iv); err != nil {
		sourceFile.Close()
		targetFile.Close()
		return nil, fmt.Errorf("failed to write IV: %v", err)
	}

	// Create CBC encrypter
	blockMode := cipher.NewCBCEncrypter(aesCipher, iv)

	// Buffer for reading source file
	blockSize := aesCipher.BlockSize()
	buffer := make([]byte, 4096) // Use 4KB buffer for reading
	hasMoreData := true
	var paddingNeeded bool
	var lastPartial []byte

	// Read and encrypt in chunks
	for hasMoreData {
		n, err := sourceFile.Read(buffer)
		if err != nil && err != io.EOF {
			sourceFile.Close()
			targetFile.Close()
			return nil, fmt.Errorf("error reading source file: %v", err)
		}

		if err == io.EOF {
			hasMoreData = false
			paddingNeeded = true
		}

		if n > 0 {
			// Append any leftover bytes from previous read
			var currentChunk []byte
			if lastPartial != nil {
				currentChunk = append(lastPartial, buffer[:n]...)
			} else {
				currentChunk = buffer[:n]
			}

			// Process complete blocks
			completeBlocks := len(currentChunk) / blockSize * blockSize
			if completeBlocks > 0 {
				// Encrypt complete blocks
				encryptedBlocks := make([]byte, completeBlocks)
				blockMode.CryptBlocks(encryptedBlocks, currentChunk[:completeBlocks])

				// Write encrypted blocks
				if _, err := targetFile.Write(encryptedBlocks); err != nil {
					sourceFile.Close()
					targetFile.Close()
					return nil, fmt.Errorf("failed to write encrypted blocks: %v", err)
				}

				// Save any remaining bytes for next iteration
				if len(currentChunk) > completeBlocks {
					lastPartial = currentChunk[completeBlocks:]
				} else {
					lastPartial = nil
				}
			} else {
				// Save all bytes for next iteration
				lastPartial = currentChunk
			}
		}

		if !hasMoreData && paddingNeeded {
			// Handle final padding
			var finalBlock []byte
			if lastPartial != nil {
				finalBlock = lastPartial
			} else {
				// If the file was a perfect multiple of block size, add a full padding block
				finalBlock = []byte{}
			}

			// Apply PKCS7 padding
			padding := blockSize - (len(finalBlock) % blockSize)
			paddedFinal := make([]byte, len(finalBlock)+padding)
			copy(paddedFinal, finalBlock)

			// Fill padding bytes with padding value
			for i := len(finalBlock); i < len(paddedFinal); i++ {
				paddedFinal[i] = byte(padding)
			}

			// Encrypt final padded block
			encryptedFinal := make([]byte, len(paddedFinal))
			blockMode.CryptBlocks(encryptedFinal, paddedFinal)

			// Write final encrypted block
			if _, err := targetFile.Write(encryptedFinal); err != nil {
				sourceFile.Close()
				targetFile.Close()
				return nil, fmt.Errorf("failed to write final encrypted block: %v", err)
			}
		}
	}

	// Close files
	sourceFile.Close()
	targetFile.Sync() // Ensure all data is flushed to disk
	targetFile.Close()

	// Calculate HMAC
	hmacFile, err := os.Open(targetFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for HMAC calculation: %v", err)
	}

	// Seek past the HMAC placeholder
	if _, err := hmacFile.Seek(int64(hmacLength), io.SeekStart); err != nil {
		hmacFile.Close()
		return nil, fmt.Errorf("failed to seek past HMAC placeholder: %v", err)
	}

	// Calculate HMAC over everything after the placeholder
	hmacCalculator := hmac.New(sha256.New, hmacKey)
	if _, err := io.Copy(hmacCalculator, hmacFile); err != nil {
		hmacFile.Close()
		return nil, fmt.Errorf("failed to calculate HMAC: %v", err)
	}
	hmacFile.Close()

	mac := hmacCalculator.Sum(nil)

	// Write HMAC to beginning of file
	hmacWriter, err := os.OpenFile(targetFilePath, os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for writing HMAC: %v", err)
	}

	if _, err := hmacWriter.Write(mac); err != nil {
		hmacWriter.Close()
		return nil, fmt.Errorf("failed to write HMAC: %v", err)
	}
	hmacWriter.Close()

	// Return encryption info
	return &EncryptionInfo{
		EncryptionKey:        base64.StdEncoding.EncodeToString(key),
		FileDigest:           base64.StdEncoding.EncodeToString(sourceDigestValue),
		FileDigestAlgorithm:  "SHA256",
		InitializationVector: base64.StdEncoding.EncodeToString(iv),
		Mac:                  base64.StdEncoding.EncodeToString(mac),
		MacKey:               base64.StdEncoding.EncodeToString(hmacKey),
		ProfileIdentifier:    "ProfileVersion1",
	}, nil
}

// analyzeEncryptedFileHex reads the encrypted file at the given path and returns details
// about its structure in hexadecimal format. It assumes the file layout is:
//   - Bytes 0-31: HMAC-SHA256 MAC,
//   - Bytes 32-47: AES-256-CBC Initialization Vector (IV),
//   - Bytes 48-end: Encrypted content.
//
// The returned FullHeaderHex is formatted as uppercase hex bytes separated by spaces.
func analyzeEncryptedFileHex(encryptedFilePath string) (*EncryptedFileAnalysis, error) {
	data, err := os.ReadFile(encryptedFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted file: %v", err)
	}
	fileLength := int64(len(data))
	if fileLength < 48 {
		return nil, fmt.Errorf("file too short to contain valid encryption information: %d bytes", fileLength)
	}

	// Extract HMAC (first 32 bytes)
	hmacBytes := data[0:32]
	hmacHex := strings.ToUpper(hex.EncodeToString(hmacBytes))

	// Extract IV (next 16 bytes, positions 32-47)
	ivBytes := data[32:48]
	ivHex := strings.ToUpper(hex.EncodeToString(ivBytes))

	// Extract a sample of ciphertext (first 16 bytes of ciphertext, positions 48-63)
	var ciphertextSample string
	if fileLength >= 64 {
		ciphertextSample = strings.ToUpper(hex.EncodeToString(data[48:64]))
	} else {
		ciphertextSample = "N/A"
	}

	// Build a full header hex dump of the first 64 bytes (or the entire file if shorter),
	// formatting each byte as uppercase hex and joining with spaces.
	headerLength := int(math.Min(64, float64(fileLength)))
	headerBytes := data[0:headerLength]
	var headerHexParts []string
	for _, b := range headerBytes {
		headerHexParts = append(headerHexParts, fmt.Sprintf("%02X", b))
	}
	fullHeaderHex := strings.Join(headerHexParts, " ")

	return &EncryptedFileAnalysis{
		FileLength:       fileLength,
		HMACHex:          hmacHex,
		IVHex:            ivHex,
		CiphertextSample: ciphertextSample,
		FullHeaderHex:    fullHeaderHex,
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
