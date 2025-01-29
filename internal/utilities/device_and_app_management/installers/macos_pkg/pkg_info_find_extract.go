package utility

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/macos_pkg/xar"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Field defines a field to be extracted from file content
type Field struct {
	Key      string
	Required bool
}

// ExtractedFields represents fields extracted from a file
type ExtractedFields struct {
	FilePath string
	Values   map[string]string
}

// ExtractFieldsFromFiles finds and extracts specified fields from matching files
func ExtractFieldsFromFiles(ctx context.Context, filePath string, pattern string, fields []Field) ([]ExtractedFields, error) {
	// Layer 1: Parse XAR
	tflog.Debug(ctx, "Opening XAR file for field extraction", map[string]interface{}{
		"path":    filePath,
		"pattern": pattern,
	})

	xarFile, err := xar.ParseXAR(ctx, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XAR: %w", err)
	}

	// Layer 2: Process payloads
	processedXar, err := xar.ProcessPayloads(ctx, xarFile)
	if err != nil {
		return nil, fmt.Errorf("failed to process payloads: %w", err)
	}

	// Layer 3: Find matching files
	matches, err := xar.FindFilesByPattern(ctx, processedXar, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find files: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no files found matching pattern: %s", pattern)
	}

	var results []ExtractedFields

	// Layer 4: Extract content and parse fields from each match
	for _, match := range matches {
		tflog.Debug(ctx, "Processing file for field extraction", map[string]interface{}{
			"file": match.FullPath,
		})

		// Layer 4a: Extract content
		content, err := xar.ExtractContent(ctx, processedXar, match)
		if err != nil {
			tflog.Error(ctx, "Failed to extract content", map[string]interface{}{
				"file":  match.FullPath,
				"error": err.Error(),
			})
			continue
		}

		// Layer 4b: Parse as plist
		plistData, err := xar.ParseAsPlist(ctx, content)
		if err != nil {
			tflog.Error(ctx, "Failed to parse plist", map[string]interface{}{
				"file":  match.FullPath,
				"error": err.Error(),
			})
			continue
		}

		// Extract requested fields
		extractedFields := ExtractedFields{
			FilePath: match.FullPath,
			Values:   make(map[string]string),
		}

		missingRequired := false
		for _, field := range fields {
			if val, ok := plistData[field.Key].(string); ok {
				extractedFields.Values[field.Key] = val
				tflog.Debug(ctx, "Extracted field value", map[string]interface{}{
					"file":  match.FullPath,
					"field": field.Key,
					"value": val,
				})
			} else if field.Required {
				tflog.Error(ctx, "Required field missing", map[string]interface{}{
					"file":  match.FullPath,
					"field": field.Key,
				})
				missingRequired = true
				break
			}
		}

		if !missingRequired {
			results = append(results, extractedFields)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no valid fields extracted from matching files")
	}

	tflog.Info(ctx, "Completed field extraction", map[string]interface{}{
		"files_processed": len(matches),
		"files_extracted": len(results),
	})

	return results, nil
}
