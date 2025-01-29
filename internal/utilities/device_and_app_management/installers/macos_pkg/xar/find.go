package xar

import (
	"context"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// FileMatch contains all information about a matched file
type FileMatch struct {
	ID             uint64
	FullPath       string
	Size           int64
	Type           FileType
	IsCompressed   bool
	EncodingType   string
	CompressedSize int64
}

// FindFilesByPattern searches recursively through the processed XAR archive for files matching the pattern
func FindFilesByPattern(ctx context.Context, processedXar *ProcessedXAR, pattern string) ([]FileMatch, error) {
	tflog.Debug(ctx, "Starting file search", map[string]interface{}{
		"pattern": pattern,
	})

	var matches []FileMatch

	// Search through all files in the archive
	for id, file := range processedXar.Reader.File {
		// Check if the file matches the pattern
		match, err := filepath.Match(pattern, filepath.Base(file.Name))
		if err != nil {
			tflog.Warn(ctx, "Pattern matching error", map[string]interface{}{
				"file":    file.Name,
				"pattern": pattern,
				"error":   err.Error(),
			})
			continue
		}

		if match {
			// Get payload metadata if available
			payloadMeta := processedXar.Payloads[id]

			fileMatch := FileMatch{
				ID:             id,
				FullPath:       file.Name,
				Size:           file.Size,
				Type:           file.Type,
				IsCompressed:   payloadMeta.IsCompressed,
				EncodingType:   payloadMeta.EncodingType,
				CompressedSize: payloadMeta.CompressedSize,
			}

			tflog.Debug(ctx, "Found matching file", map[string]interface{}{
				"path":       fileMatch.FullPath,
				"size":       fileMatch.Size,
				"compressed": fileMatch.IsCompressed,
				"encoding":   fileMatch.EncodingType,
			})

			matches = append(matches, fileMatch)
		}
	}

	tflog.Info(ctx, "Completed file search", map[string]interface{}{
		"pattern": pattern,
		"matches": len(matches),
	})

	return matches, nil
}
