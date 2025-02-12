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

	for blockNum := 0; blockNum < totalBlocks; blockNum++ {
		// Create block ID
		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%07d", blockNum)))

		// Read file chunk
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read file block: %v", err)
		}

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
				return retry.RetryableError(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				return retry.RetryableError(fmt.Errorf("unexpected status: %d", resp.StatusCode))
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to upload block: %v", err)
		}

		blockList = append(blockList, blockID)
	}

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

	return nil
}
