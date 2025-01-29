package xar

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ExtractedContent represents content found and extracted from either XAR or Payload
type ExtractedContent struct {
	Path        string
	Content     []byte
	FromPayload bool
	PayloadPath string // Original payload path if FromPayload is true
}

// ExtractContent gets content from either XAR directly or from within a Payload
// Fix the extraction logic for compressed payloads
func ExtractContent(ctx context.Context, processedXar *ProcessedXAR, match FileMatch) (*ExtractedContent, error) {
	// First check if the file is directly in the XAR
	if file, ok := processedXar.Reader.File[match.ID]; ok {
		tflog.Debug(ctx, "Extracting file from XAR", map[string]interface{}{
			"path":       match.FullPath,
			"compressed": match.IsCompressed,
		})

		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer rc.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, rc); err != nil {
			return nil, fmt.Errorf("failed to read content: %w", err)
		}

		return &ExtractedContent{
			Path:        match.FullPath,
			Content:     buf.Bytes(),
			FromPayload: false,
		}, nil
	}

	// If not in XAR directly, check extracted payloads
	for _, metadata := range processedXar.Payloads {
		if extractedContent, found := metadata.ExtractedFiles[match.FullPath]; found {
			tflog.Debug(ctx, "Found file in extracted payload", map[string]interface{}{
				"file":         match.FullPath,
				"payload_path": metadata.Path,
			})

			return &ExtractedContent{
				Path:        match.FullPath,
				Content:     extractedContent,
				FromPayload: true,
				PayloadPath: metadata.Path,
			}, nil
		}
	}

	return nil, fmt.Errorf("file not found in XAR or extracted payloads: %s", match.FullPath)
}
