package sharedConstructors

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	blockSize          = 8 * 1024 * 1024  // 8 MiB chunks
	blockUploadTimeout = 60 * time.Second // Individual block upload timeout
	retryMaxAttempts   = 3                // Number of retry attempts per block
	retryMinWait       = 5 * time.Second  // Minimum wait between retries
	retryMaxWait       = 10 * time.Second // Maximum wait between retries
)

// UploadToAzureStorage handles the chunked upload of large files to Azure Blob Storage.
// This implementation generates block IDs as six-digit, zero-padded strings (e.g., "000001")
// before base64 encoding them. This ensures that block ordering matches the expected
// lexicographical order, preventing file corruption during the commit phase.
//
// The function follows these steps:
//  1. Opens the source file and retrieves its metadata.
//  2. Divides the file into 8 MiB blocks.
//  3. For each block, generates a block ID as a six-digit string, converts it to base64, and uploads
//     the block using the required "x-ms-blob-type" header.
//  4. Constructs an XML block list from these base64-encoded block IDs and sends a commit request.
//  5. Reports progress and handles errors with context-aware timeouts and retry logic.
//
// Parameters:
//   - ctx: Context for timeout and cancellation control.
//   - sasUri: Azure Storage SAS URI with write permissions for the target blob.
//   - filePath: Path to the local file to be uploaded.
//
// Returns an error if the upload fails at any stage, with details about the failure.
func UploadToAzureStorage(ctx context.Context, sasUri string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// Define the block size to match PowerShell (8 MiB)
	blockSize := 8 * 1024 * 1024
	totalBlocks := int(math.Ceil(float64(fileInfo.Size()) / float64(blockSize)))
	blockList := []string{}
	buffer := make([]byte, blockSize)

	tflog.Debug(ctx, "Starting Azure Storage upload", map[string]interface{}{
		"file_path":          filePath,
		"total_size_mb":      float64(fileInfo.Size()) / 1024 / 1024,
		"total_blocks":       totalBlocks,
		"block_size_mb":      blockSize / 1024 / 1024,
		"retry_max_attempts": retryMaxAttempts,
		"block_timeout":      blockUploadTimeout.String(),
	})

	uploadedBytes := int64(0)
	startTime := time.Now()

	for blockNum := 0; blockNum < totalBlocks; blockNum++ {
		// Generate block ID as a six-digit, zero-padded string and then base64 encode it.
		// This matches the PowerShell equivilent logic (ToString("D6") in UTF8) and ensures correct ordering.
		blockIDString := fmt.Sprintf("%06d", blockNum)
		blockID := base64.StdEncoding.EncodeToString([]byte(blockIDString))

		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file block: %v", err)
		}

		uploadedBytes += int64(n)
		percentComplete := float64(uploadedBytes) / float64(fileInfo.Size()) * 100
		uploadSpeed := float64(uploadedBytes) / time.Since(startTime).Seconds() / 1024 / 1024 // MB/s

		tflog.Debug(ctx, "Uploading block", map[string]interface{}{
			"block_number":        blockNum + 1,
			"blocks_remaining":    totalBlocks - (blockNum + 1),
			"bytes_uploaded_mb":   float64(uploadedBytes) / 1024 / 1024,
			"percent_complete":    fmt.Sprintf("%.1f%%", percentComplete),
			"elapsed_time":        time.Since(startTime).Round(time.Second).String(),
			"upload_speed_mbps":   fmt.Sprintf("%.2f", uploadSpeed),
			"estimated_remaining": fmt.Sprintf("%.0fs", float64(fileInfo.Size()-uploadedBytes)/float64(uploadedBytes)*time.Since(startTime).Seconds()),
		})

		// Create the block URL with the SAS token and the new block ID
		blockURL := fmt.Sprintf("%s&comp=block&blockid=%s", sasUri, blockID)

		// Upload the block with retry logic
		err = retry.RetryContext(ctx, blockUploadTimeout, func() *retry.RetryError {
			uploadCtx, cancel := context.WithTimeout(ctx, blockUploadTimeout)
			defer cancel()

			req, err := http.NewRequestWithContext(uploadCtx, "PUT", blockURL, bytes.NewReader(buffer[:n]))
			if err != nil {
				tflog.Error(ctx, "Failed to create block upload request", map[string]interface{}{
					"block_number": blockNum + 1,
					"error":        err.Error(),
				})
				return retry.NonRetryableError(err)
			}

			req.Header.Set("x-ms-blob-type", "BlockBlob")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if err == context.DeadlineExceeded {
					tflog.Error(ctx, "Block upload timed out", map[string]interface{}{
						"block_number": blockNum + 1,
						"timeout":      blockUploadTimeout.String(),
						"error":        err.Error(),
					})
				} else {
					tflog.Error(ctx, "Failed to upload block", map[string]interface{}{
						"block_number": blockNum + 1,
						"error":        err.Error(),
					})
				}
				return retry.RetryableError(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				body, _ := io.ReadAll(resp.Body)
				tflog.Error(ctx, "Unexpected status code", map[string]interface{}{
					"block_number": blockNum + 1,
					"status_code":  resp.StatusCode,
					"response":     string(body),
				})
				return retry.RetryableError(fmt.Errorf("unexpected status: %d - %s", resp.StatusCode, string(body)))
			}

			tflog.Debug(ctx, "Block upload successful", map[string]interface{}{
				"block_number": blockNum + 1,
				"size_mb":      float64(n) / 1024 / 1024,
			})

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to upload block %d after %d attempts: %v", blockNum+1, retryMaxAttempts, err)
		}

		blockList = append(blockList, blockID)
	}

	tflog.Debug(ctx, "File upload completed, committing block list", map[string]interface{}{
		"total_blocks":  len(blockList),
		"total_size_mb": float64(fileInfo.Size()) / 1024 / 1024,
		"elapsed_time":  time.Since(startTime).Round(time.Second).String(),
	})

	// Build the block list XML using the base64-encoded block IDs
	var blockListElements []string
	for _, blockID := range blockList {
		blockListElements = append(blockListElements, fmt.Sprintf("<Latest>%s</Latest>", blockID))
	}
	blockListXML := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?><BlockList>%s</BlockList>`,
		strings.Join(blockListElements, ""))

	commitURL := fmt.Sprintf("%s&comp=blocklist", sasUri)

	commitCtx, cancel := context.WithTimeout(ctx, blockUploadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(commitCtx, "PUT", commitURL, strings.NewReader(blockListXML))
	if err != nil {
		return fmt.Errorf("failed to create commit request: %v", err)
	}
	req.Header.Set("Content-Type", "application/xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("commit request timed out after %s", blockUploadTimeout)
		}
		return fmt.Errorf("failed to commit blocks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to commit blocks: status %d - %s", resp.StatusCode, string(body))
	}

	tflog.Debug(ctx, "Azure Storage upload completed successfully", map[string]interface{}{
		"file_path":      filePath,
		"total_size_mb":  float64(fileInfo.Size()) / 1024 / 1024,
		"total_blocks":   len(blockList),
		"elapsed_time":   time.Since(startTime).Round(time.Second).String(),
		"avg_speed_mbps": float64(fileInfo.Size()) / 1024 / 1024 / time.Since(startTime).Seconds(),
	})

	return nil
}
