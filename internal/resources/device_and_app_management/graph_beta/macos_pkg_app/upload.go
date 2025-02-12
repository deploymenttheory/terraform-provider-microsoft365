package graphBetaMacOSPKGApp

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

// uploadToAzureStorage handles chunked upload of large files to Azure Storage
func uploadToAzureStorage(ctx context.Context, sasUri string, filePath string) error {
	const blockSize = 4 * 1024 * 1024 // 4 MiB chunks

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	totalBlocks := int(math.Ceil(float64(fileInfo.Size()) / float64(blockSize)))
	blockList := []string{}
	buffer := make([]byte, blockSize)

	tflog.Debug(ctx, "Starting Azure Storage upload", map[string]interface{}{
		"file_path":     filePath,
		"total_size":    fileInfo.Size(),
		"total_blocks":  totalBlocks,
		"block_size_mb": blockSize / 1024 / 1024,
	})

	uploadedBytes := int64(0)
	startTime := time.Now()

	for blockNum := 0; blockNum < totalBlocks; blockNum++ {
		// Create block ID
		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%07d", blockNum)))

		// Read file chunk
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file block: %v", err)
		}

		uploadedBytes += int64(n)
		percentComplete := float64(uploadedBytes) / float64(fileInfo.Size()) * 100

		tflog.Debug(ctx, "Uploading block", map[string]interface{}{
			"block_number":     blockNum + 1,
			"blocks_remaining": totalBlocks - (blockNum + 1),
			"bytes_uploaded":   uploadedBytes,
			"percent_complete": fmt.Sprintf("%.1f%%", percentComplete),
			"elapsed_time":     time.Since(startTime).Round(time.Second).String(),
		})

		// Create block URL with SAS token
		blockURL := fmt.Sprintf("%s&comp=block&blockid=%s", sasUri, blockID)

		// Upload block with retry logic
		err = retry.RetryContext(ctx, 30*time.Second, func() *retry.RetryError {
			req, err := http.NewRequestWithContext(ctx, "PUT", blockURL, bytes.NewReader(buffer[:n]))
			if err != nil {
				return retry.NonRetryableError(err)
			}

			req.Header.Set("x-ms-blob-type", "BlockBlob")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				tflog.Error(ctx, "Failed to upload block, retrying", map[string]interface{}{
					"block_number": blockNum + 1,
					"error":        err.Error(),
				})
				return retry.RetryableError(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				tflog.Error(ctx, "Unexpected status code, retrying", map[string]interface{}{
					"block_number": blockNum + 1,
					"status_code":  resp.StatusCode,
				})
				return retry.RetryableError(fmt.Errorf("unexpected status: %d", resp.StatusCode))
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to upload block: %v", err)
		}

		blockList = append(blockList, blockID)
	}

	tflog.Debug(ctx, "File upload completed, committing block list", map[string]interface{}{
		"total_blocks":  len(blockList),
		"total_size_mb": float64(fileInfo.Size()) / 1024 / 1024,
		"elapsed_time":  time.Since(startTime).Round(time.Second).String(),
	})

	// Commit block list
	blockListXML := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?><BlockList>%s</BlockList>`,
		strings.Join(blockList, ""))

	commitURL := fmt.Sprintf("%s&comp=blocklist", sasUri)

	req, err := http.NewRequestWithContext(ctx, "PUT", commitURL, strings.NewReader(blockListXML))
	if err != nil {
		return fmt.Errorf("failed to create commit request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to commit blocks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to commit blocks: unexpected status %d", resp.StatusCode)
	}

	tflog.Debug(ctx, "Azure Storage upload completed successfully", map[string]interface{}{
		"file_path":     filePath,
		"total_size_mb": float64(fileInfo.Size()) / 1024 / 1024,
		"total_blocks":  len(blockList),
		"elapsed_time":  time.Since(startTime).Round(time.Second).String(),
	})

	return nil
}
