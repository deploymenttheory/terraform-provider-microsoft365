package xar

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ParseXarDirectoryStructure parses a XAR file and returns its directory structure.
// It logs the parsing process for debugging.
func ParseXarDirectoryStructure(ctx context.Context, xarFilePath string) (map[string]interface{}, error) {
	tflog.Debug(ctx, "Opening XAR file", map[string]interface{}{
		"path": xarFilePath,
	})

	reader, err := OpenReader(xarFilePath)
	if err != nil {
		tflog.Debug(ctx, "Failed to open XAR file", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to open XAR file: %w", err)
	}
	defer func() {
		tflog.Debug(ctx, "Closing XAR file")
		_ = reader // No explicit close needed but logging for completeness
	}()

	tflog.Debug(ctx, "Reading directory structure")

	dirStructure := make(map[string]interface{})

	for _, file := range reader.File {
		tflog.Debug(ctx, "Processing file", map[string]interface{}{
			"name": file.Name,
			"type": file.Type,
		})
		addFileToStructure(dirStructure, file.Name, file.Type == FileTypeDirectory)
	}

	tflog.Debug(ctx, "Finished reading directory structure")

	return dirStructure, nil
}

// Helper function to add a file or directory to the structure map
func addFileToStructure(structure map[string]interface{}, filePath string, isDirectory bool) {
	parts := strings.Split(filePath, "/")
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
