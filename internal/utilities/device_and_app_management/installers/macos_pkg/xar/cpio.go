package xar

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/cavaliergopher/cpio"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Extracts individual files from a cpio archive
func extractCpioFiles(ctx context.Context, metadata *PayloadMetadata, r io.Reader) error {
	tflog.Debug(ctx, "Extracting files from cpio archive", map[string]interface{}{
		"payload_path": metadata.Path,
	})

	cr := cpio.NewReader(r)

	for {
		hdr, err := cr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			tflog.Error(ctx, "Failed to read cpio archive", map[string]interface{}{
				"payload_path": metadata.Path,
				"error":        err.Error(),
			})
			return fmt.Errorf("failed to read cpio archive: %w", err)
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, cr); err != nil {
			return fmt.Errorf("failed to extract file %s: %w", hdr.Name, err)
		}
		metadata.ExtractedFiles[hdr.Name] = buf.Bytes()
	}

	tflog.Info(ctx, "Successfully extracted files from cpio", map[string]interface{}{
		"payload_path": metadata.Path,
		"total_files":  len(metadata.ExtractedFiles),
	})

	return nil
}
