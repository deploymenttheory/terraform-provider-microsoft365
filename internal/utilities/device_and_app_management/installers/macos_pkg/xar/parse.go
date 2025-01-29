package xar

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"howett.net/plist"
)

// XarFile represents a parsed XAR archive with its structure
type XarFile struct {
	Reader    *Reader
	Structure map[string]interface{}
}

// ParseXAR opens and processes a XAR file, returning both the reader and directory structure
func ParseXAR(ctx context.Context, filePath string) (*XarFile, error) {
	tflog.Debug(ctx, "Starting XAR file parsing", map[string]interface{}{
		"path": filePath,
	})

	// Open the XAR file
	reader, err := OpenReader(filePath)
	if err != nil {
		tflog.Error(ctx, "Failed to open XAR file", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to open XAR file: %w", err)
	}

	// Create directory structure
	structure := make(map[string]interface{})
	for _, file := range reader.File {
		tflog.Debug(ctx, "Processing file", map[string]interface{}{
			"name": file.Name,
			"type": file.Type,
		})
		addToStructure(structure, file.Name, file.Type == FileTypeDirectory)
	}

	tflog.Info(ctx, "Successfully parsed XAR file", map[string]interface{}{
		"total_files": len(reader.File),
	})

	return &XarFile{
		Reader:    reader,
		Structure: structure,
	}, nil
}

// addToStructure adds a file or directory to the structure map
func addToStructure(structure map[string]interface{}, filePath string, isDirectory bool) {
	parts := strings.Split(strings.TrimPrefix(filePath, "/"), "/")
	current := structure

	for i, part := range parts {
		if i == len(parts)-1 {
			if isDirectory {
				current[part] = make(map[string]interface{})
			} else {
				current[part] = nil
			}
		} else {
			if _, exists := current[part]; !exists {
				current[part] = make(map[string]interface{})
			}
			current = current[part].(map[string]interface{})
		}
	}
}

// ParseAsPlist attempts to parse content as a plist
func ParseAsPlist(ctx context.Context, content *ExtractedContent) (map[string]interface{}, error) {
	tflog.Debug(ctx, "Parsing plist content", map[string]interface{}{
		"path":         content.Path,
		"from_payload": content.FromPayload,
	})

	var data map[string]interface{}
	decoder := plist.NewDecoder(bytes.NewReader(content.Content))

	if err := decoder.Decode(&data); err != nil {
		tflog.Error(ctx, "Failed to parse plist", map[string]interface{}{
			"path":  content.Path,
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to parse plist: %w", err)
	}

	tflog.Debug(ctx, "Successfully parsed plist", map[string]interface{}{
		"path":       content.Path,
		"keys_found": len(data),
	})

	return data, nil
}
