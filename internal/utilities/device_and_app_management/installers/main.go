package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/extract_macos_installer_metadata/extract"
	"howett.net/plist"
)

// PkgInfo represents the structure of a PackageInfo plist file
type PkgInfo struct {
	Identifier      string   `plist:"identifier,omitempty"`
	Version         string   `plist:"version,omitempty"`
	InstallLocation string   `plist:"install-location,omitempty"`
	Bundles         []Bundle `plist:"bundle,omitempty"`
}

type Bundle struct {
	Path                       string `plist:"path"`
	ID                         string `plist:"id"`
	CFBundleShortVersionString string `plist:"CFBundleShortVersionString"`
	CFBundleVersion            string `plist:"CFBundleVersion"`
}

type ExtractedMetadata struct {
	PrimaryBundleId      string        `json:"primaryBundleId"`
	PrimaryBundleVersion string        `json:"primaryBundleVersion"`
	IncludedApps         []IncludedApp `json:"includedApps"`
	PackageInfo          PackageInfo   `json:"packageInfo"`
}

type IncludedApp struct {
	BundleId      string `json:"bundleId"`
	BundleVersion string `json:"bundleVersion"`
}

type PackageInfo struct {
	Identifier      string `json:"identifier"`
	Version         string `json:"version"`
	InstallLocation string `json:"installLocation"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	inputPath := flag.String("input", "", "Path to the .pkg file")
	outputPath := flag.String("output", "", "Path to save the JSON output (optional, defaults to stdout)")
	flag.Parse()

	if *inputPath == "" {
		log.Fatal("Error: -input flag is required")
	}

	// Create log file
	logFile, err := os.CreateTemp("", "pkginfo-*.log")
	if err != nil {
		log.Fatal("Failed to create log file:", err)
	}
	defer logFile.Close()

	// Set up logging
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// Extract metadata
	metadata, err := extractPkgMetadata(*inputPath)
	if err != nil {
		log.Printf("Failed to extract metadata: %v", err)
		log.Fatal("Error extracting metadata:", err)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		log.Printf("Failed to convert to JSON: %v", err)
		log.Fatal("Error converting to JSON:", err)
	}

	// Output results
	if *outputPath != "" {
		err = os.WriteFile(*outputPath, jsonData, 0644)
		if err != nil {
			log.Printf("Failed to write output file: %v", err)
			log.Fatal("Error writing output file:", err)
		}
		log.Printf("Metadata written to %s", *outputPath)
	} else {
		fmt.Println(string(jsonData))
	}
}

func extractPkgMetadata(filePath string) (*ExtractedMetadata, error) {
	log.Printf("Starting metadata extraction for: %s", filePath)

	// Open pkg file
	reader, err := extract.OpenPkg(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open pkg file: %w", err)
	}
	defer reader.Close()

	// Extract PackageInfo
	packageInfoData, err := reader.ExtractPackageInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to extract PackageInfo: %w", err)
	}

	log.Printf("Successfully read PackageInfo (%d bytes)", len(packageInfoData))

	// Parse PackageInfo as plist
	var pkgInfo PkgInfo
	decoder := plist.NewDecoder(bytes.NewReader(packageInfoData))
	if err := decoder.Decode(&pkgInfo); err != nil {
		return nil, fmt.Errorf("failed to decode PackageInfo plist: %w", err)
	}

	// Create metadata structure
	metadata := &ExtractedMetadata{
		IncludedApps: make([]IncludedApp, 0),
		PackageInfo: PackageInfo{
			Identifier:      pkgInfo.Identifier,
			Version:         pkgInfo.Version,
			InstallLocation: pkgInfo.InstallLocation,
		},
	}

	// Process bundles
	for i, bundle := range pkgInfo.Bundles {
		log.Printf("Processing bundle %d: Path=%s, ID=%s, Version=%s",
			i+1, bundle.Path, bundle.ID, bundle.CFBundleShortVersionString)

		if bundle.Path != "" && strings.Contains(bundle.Path, ".app") {
			metadata.IncludedApps = append(metadata.IncludedApps, IncludedApp{
				BundleId:      bundle.ID,
				BundleVersion: bundle.CFBundleShortVersionString,
			})
			log.Printf("Added bundle to included apps")
		} else {
			log.Printf("Skipped bundle (not an .app or empty path)")
		}
	}

	// If no bundles found, use package identifier
	if len(metadata.IncludedApps) == 0 && pkgInfo.Identifier != "" {
		log.Printf("No app bundles found, using package identifier as fallback")
		metadata.IncludedApps = append(metadata.IncludedApps, IncludedApp{
			BundleId:      pkgInfo.Identifier,
			BundleVersion: pkgInfo.Version,
		})
	}

	if len(metadata.IncludedApps) > 0 {
		metadata.PrimaryBundleId = metadata.IncludedApps[0].BundleId
		metadata.PrimaryBundleVersion = metadata.IncludedApps[0].BundleVersion
		log.Printf("Set primary bundle: ID=%s, Version=%s",
			metadata.PrimaryBundleId, metadata.PrimaryBundleVersion)
	} else {
		return nil, fmt.Errorf("no valid included apps found in package")
	}

	log.Printf("Successfully extracted metadata for %d apps", len(metadata.IncludedApps))
	return metadata, nil
}
