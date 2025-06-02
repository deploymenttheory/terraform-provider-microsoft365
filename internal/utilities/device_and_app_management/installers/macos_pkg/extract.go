package utility

import (
	"context"
	"fmt"
	"os"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/macos_pkg/xar"
)

// Field defines a field to be extracted
type Field struct {
	Key      string
	Required bool
}

// ExtractedFields represents fields extracted from a file
type ExtractedFields struct {
	FilePath string
	Values   map[string]string
}

func ExtractFieldsFromPKGFile(ctx context.Context, filePath string, pattern string, fields []Field) ([]ExtractedFields, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	tfr, err := xar.NewTempFileReader(file, os.TempDir)
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	defer tfr.Close()

	metadata, err := xar.ExtractXARMetadata(tfr)
	if err != nil {
		return nil, fmt.Errorf("extracting XAR metadata: %w", err)
	}

	var results []ExtractedFields

	// Always include primary bundle as an included app
	results = append(results, ExtractedFields{
		FilePath: "primary",
		Values: map[string]string{
			"CFBundleIdentifier":         metadata.BundleIdentifier,
			"CFBundleShortVersionString": metadata.Version,
		},
	})

	// Add all bundles found in package
	for _, bundle := range metadata.IncludedBundles {
		if bundle.BundleID != "" && bundle.Version != "" {
			results = append(results, ExtractedFields{
				FilePath: bundle.Path,
				Values: map[string]string{
					"CFBundleIdentifier":         bundle.BundleID,
					"CFBundleShortVersionString": bundle.Version,
				},
			})
		}
	}

	return results, nil
}
