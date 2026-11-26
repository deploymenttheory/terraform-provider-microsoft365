package sharedConstructors

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var errUnexpectedAzureUploadStatus = errors.New("unexpected status")

const (
	blockSize                = 8 * 1024 * 1024  // 8 MiB chunks
	blockUploadTimeout       = 2 * time.Minute  // Permit slower 8 MiB body writes within the caller deadline
	blockListCommitTimeout   = time.Minute      // Keep the small block-list request independently bounded
	retryMaxAttempts         = 3                // Number of retry attempts per block
	retryBackoffAllowance    = 10 * time.Second // Allowance for SDK retry scheduling
	sasPolicyPropagationWait = 30 * time.Second // Azure stored access policy propagation limit
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

	tflog.Debug(ctx, "Starting Azure Storage upload", map[string]any{
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

		tflog.Debug(ctx, "Uploading block", map[string]any{
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

		// Each attempt may consume the full request timeout. Give all configured
		// attempts room to run, bounded by the caller's operation deadline.
		attempts := 0
		retryBudget := time.Duration(retryMaxAttempts)*blockUploadTimeout + time.Duration(retryMaxAttempts-1)*(sasPolicyPropagationWait+retryBackoffAllowance)
		retryBlockError := func(err error) *retry.RetryError {
			if ctx.Err() != nil || attempts >= retryMaxAttempts {
				return retry.NonRetryableError(err)
			}
			return retry.RetryableError(err)
		}
		err = retry.RetryContext(ctx, retryBudget, func() *retry.RetryError {
			attempts++
			uploadCtx, cancel := context.WithTimeout(ctx, blockUploadTimeout)
			defer cancel()

			req, err := http.NewRequestWithContext(uploadCtx, "PUT", blockURL, bytes.NewReader(buffer[:n]))
			if err != nil {
				err = azureStorageRequestError(err)
				tflog.Error(ctx, "Failed to create block upload request", map[string]any{
					"block_number": blockNum + 1,
					"error":        err.Error(),
				})
				return retry.NonRetryableError(err)
			}

			req.Header.Set("x-ms-blob-type", "BlockBlob")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				err = azureStorageRequestError(err)
				if errors.Is(err, context.DeadlineExceeded) {
					tflog.Error(ctx, "Block upload timed out", map[string]any{
						"block_number": blockNum + 1,
						"timeout":      blockUploadTimeout.String(),
						"error":        err.Error(),
					})
				} else {
					tflog.Error(ctx, "Failed to upload block", map[string]any{
						"block_number": blockNum + 1,
						"error":        err.Error(),
					})
				}
				return retryBlockError(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				body, _ := io.ReadAll(resp.Body)
				tflog.Error(ctx, "Unexpected status code", map[string]any{
					"block_number": blockNum + 1,
					"status_code":  resp.StatusCode,
					"response":     string(body),
				})
				statusError := fmt.Errorf("%w: %d - %s", errUnexpectedAzureUploadStatus, resp.StatusCode, string(body))
				// Azure may issue a SAS before its stored access policy has
				// propagated. Replay the same block after that documented window
				// instead of exhausting fast retries and renewing a fresh URI.
				// https://learn.microsoft.com/en-us/rest/api/storageservices/define-stored-access-policy
				missingPolicy := resp.StatusCode == http.StatusForbidden && strings.Contains(string(body), "SAS identifier cannot be found for specified signed identifier")
				if missingPolicy && attempts < retryMaxAttempts && ctx.Err() == nil {
					resp.Body.Close()
					tflog.Warn(ctx, "Waiting for Azure stored access policy propagation", map[string]any{
						"block_number": blockNum + 1,
						"wait":         sasPolicyPropagationWait.String(),
					})
					timer := time.NewTimer(sasPolicyPropagationWait)
					select {
					case <-ctx.Done():
						timer.Stop()
						return retry.NonRetryableError(fmt.Errorf("%w: %w", statusError, ctx.Err()))
					case <-timer.C:
					}
				}
				return retryBlockError(statusError)
			}

			tflog.Debug(ctx, "Block upload successful", map[string]any{
				"block_number": blockNum + 1,
				"size_mb":      float64(n) / 1024 / 1024,
			})

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to upload block %d after %d attempts: %w", blockNum+1, attempts, err)
		}

		blockList = append(blockList, blockID)
	}

	tflog.Debug(ctx, "File upload completed, committing block list", map[string]any{
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

	commitCtx, cancel := context.WithTimeout(ctx, blockListCommitTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(commitCtx, "PUT", commitURL, strings.NewReader(blockListXML))
	if err != nil {
		return fmt.Errorf("failed to create commit request: %w", azureStorageRequestError(err))
	}
	req.Header.Set("Content-Type", "application/xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = azureStorageRequestError(err)
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("commit request timed out after %s: %w", blockListCommitTimeout, err)
		}
		return fmt.Errorf("failed to commit blocks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to commit blocks: status %d - %s", resp.StatusCode, string(body))
	}

	tflog.Debug(ctx, "Azure Storage upload completed successfully", map[string]any{
		"file_path":      filePath,
		"total_size_mb":  float64(fileInfo.Size()) / 1024 / 1024,
		"total_blocks":   len(blockList),
		"elapsed_time":   time.Since(startTime).Round(time.Second).String(),
		"avg_speed_mbps": float64(fileInfo.Size()) / 1024 / 1024 / time.Since(startTime).Seconds(),
	})

	return nil
}

// URL errors contain the complete signed query. Keep the transport cause for
// cancellation/retry handling while excluding the SAS credential from diagnostics.
func azureStorageRequestError(err error) error {
	var requestError *url.Error
	if errors.As(err, &requestError) {
		return fmt.Errorf("azure storage %s failed: %w", requestError.Op, requestError.Err)
	}
	return err
}
