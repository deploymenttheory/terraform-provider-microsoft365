package xar

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cavaliergopher/cpio"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// PayloadMetadata tracks compression info about a file
// and extracted content if applicable

type PayloadMetadata struct {
	ID             uint64
	Path           string
	Size           int64
	CompressedSize int64
	EncodingType   string
	IsCompressed   bool
	ExtractedFiles map[string][]byte // Stores extracted files from the payload
}

type ProcessedXAR struct {
	*XarFile
	Payloads map[uint64]PayloadMetadata
}

// ProcessPayloads analyzes all files in the XAR to identify compressed content
func ProcessPayloads(ctx context.Context, xarFile *XarFile) (*ProcessedXAR, error) {
	tflog.Debug(ctx, "Starting payload analysis")

	processed := &ProcessedXAR{
		XarFile:  xarFile,
		Payloads: make(map[uint64]PayloadMetadata),
	}

	for id, file := range xarFile.Reader.File {
		metadata := PayloadMetadata{
			ID:             id,
			Path:           file.Name,
			Size:           file.Size,
			EncodingType:   file.EncodingMimetype,
			IsCompressed:   file.EncodingMimetype != "application/octet-stream",
			ExtractedFiles: make(map[string][]byte), // Ensure map is initialized
		}

		// If it's the Payload file, extract it
		if strings.HasSuffix(file.Name, "Payload") {
			tflog.Debug(ctx, "Extracting payload contents", map[string]interface{}{
				"path": metadata.Path,
			})
			if err := extractPayloadContents(ctx, &metadata, file); err != nil {
				tflog.Warn(ctx, "Failed to extract payload contents", map[string]interface{}{
					"path":  metadata.Path,
					"error": err.Error(),
				})
			}
		}

		processed.Payloads[id] = metadata
	}

	tflog.Info(ctx, "Completed payload analysis", map[string]interface{}{
		"total_files":      len(processed.Payloads),
		"compressed_files": countCompressedFiles(processed.Payloads),
	})

	return processed, nil
}

// Extracts all files inside the Payload directory and decompresses it if needed
func extractPayloadContents(ctx context.Context, metadata *PayloadMetadata, file *File) error {
	tflog.Debug(ctx, "Extracting contents from Payload", map[string]interface{}{
		"path": metadata.Path,
	})

	payloadReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open Payload file: %w", err)
	}
	defer payloadReader.Close()

	// Read the first few bytes to determine the file type
	header := make([]byte, 6)
	_, err = io.ReadFull(payloadReader, header)
	if err != nil {
		return fmt.Errorf("failed to read Payload file header: %w", err)
	}

	tflog.Debug(ctx, "Raw payload header bytes", map[string]interface{}{
		"path":   metadata.Path,
		"header": fmt.Sprintf("%X", header), // Convert to hex
	})

	// Reset the reader to read the full content
	payloadBytes, err := io.ReadAll(payloadReader)
	if err != nil {
		return fmt.Errorf("failed to read full Payload content: %w", err)
	}

	if len(payloadBytes) < 10 {
		return fmt.Errorf("Payload too small to be valid")
	}

	// Detect compression format
	var reader io.Reader
	switch {
	case bytes.HasPrefix(header, []byte{0x1F, 0x8B}): // GZIP magic number
		tflog.Debug(ctx, "Detected gzip compression in Payload", map[string]interface{}{
			"path": metadata.Path,
		})
		gzr, err := gzip.NewReader(bytes.NewReader(payloadBytes))
		if err != nil {
			return fmt.Errorf("failed to initialize gzip reader: %w", err)
		}
		reader = gzr
	case bytes.HasPrefix(header, []byte{0x42, 0x5A, 0x68}): // BZIP2 magic number
		tflog.Debug(ctx, "Detected bzip2 compression in Payload", map[string]interface{}{
			"path": metadata.Path,
		})
		reader = bzip2.NewReader(bytes.NewReader(payloadBytes))
	case bytes.HasPrefix(header, []byte{0x78, 0x9C}) || bytes.HasPrefix(header, []byte{0x78, 0xDA}): // ZLIB magic numbers
		tflog.Debug(ctx, "Detected zlib compression in Payload", map[string]interface{}{
			"path": metadata.Path,
		})
		zr, err := zlib.NewReader(bytes.NewReader(payloadBytes))
		if err != nil {
			return fmt.Errorf("failed to initialize zlib reader: %w", err)
		}
		reader = zr
	default:
		tflog.Debug(ctx, "No compression detected in Payload, assuming raw cpio", map[string]interface{}{
			"path": metadata.Path,
		})
		reader = bytes.NewReader(payloadBytes)
	}

	// Ensure it's a valid CPIO archive
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, reader); err != nil {
		return fmt.Errorf("failed to read decompressed payload: %w", err)
	}
	if !bytes.HasPrefix(buf.Bytes(), []byte("070701")) {
		return fmt.Errorf("invalid CPIO header after extraction")
	}

	// Dump first 64 bytes of extracted buffer for debugging
	tflog.Debug(ctx, "First 64 bytes of extracted payload", map[string]interface{}{
		"path": metadata.Path,
		"data": fmt.Sprintf("%X", buf.Bytes()[:min(64, len(buf.Bytes()))]),
	})

	// Extract files from decompressed archive
	return extractCpioFiles(ctx, metadata, buf)
}

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
