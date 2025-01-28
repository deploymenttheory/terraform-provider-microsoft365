package pkg_metadata

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/korylprince/goxar"
)

// PkgInfoXML represents the XML structure of the pkg-info file
type PkgInfoXML struct {
	XMLName         xml.Name `xml:"pkg-info"`
	Identifier      string   `xml:"identifier,attr"`
	Version         string   `xml:"version,attr"`
	InstallLocation string   `xml:"install-location,attr"`
	Bundles         []Bundle `xml:"bundle"`
}

type Bundle struct {
	Path                       string `xml:"path,attr"`
	ID                         string `xml:"id,attr"`
	CFBundleShortVersionString string `xml:"CFBundleShortVersionString,attr"`
	CFBundleVersion            string `xml:"CFBundleVersion,attr"`
}

// PackageInfo represents the output structure for package information
type PackageInfo struct {
	Identifier      string `json:"identifier"`
	Version         string `json:"version"`
	InstallLocation string `json:"installLocation"`
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

// ExtractMetadata extracts metadata from a macOS .pkg file
func ExtractMetadata(filePath string) (*ExtractedMetadata, error) {
	// Open the XAR archive
	reader, err := goxar.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open pkg file: %w", err)
	}
	defer reader.Close()

	// Find PackageInfo file
	var packageInfoFile *goxar.File
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "PackageInfo") {
			packageInfoFile = file
			break
		}
	}

	if packageInfoFile == nil {
		return nil, fmt.Errorf("PackageInfo not found in archive")
	}

	// Open and read PackageInfo
	rc, err := packageInfoFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open PackageInfo: %w", err)
	}
	defer rc.Close()

	// Read PackageInfo content
	packageInfoData, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read PackageInfo: %w", err)
	}

	// Parse PackageInfo XML
	var pkgInfoXML PkgInfoXML
	if err := xml.Unmarshal(packageInfoData, &pkgInfoXML); err != nil {
		return nil, fmt.Errorf("failed to parse PackageInfo XML: %w", err)
	}

	// Create metadata structure
	metadata := &ExtractedMetadata{
		IncludedApps: make([]IncludedApp, 0),
		PackageInfo: PackageInfo{
			Identifier:      pkgInfoXML.Identifier,
			Version:         pkgInfoXML.Version,
			InstallLocation: pkgInfoXML.InstallLocation,
		},
	}

	// Process bundles
	for _, bundle := range pkgInfoXML.Bundles {
		if bundle.Path != "" && strings.Contains(bundle.Path, ".app") {
			metadata.IncludedApps = append(metadata.IncludedApps, IncludedApp{
				BundleId:      bundle.ID,
				BundleVersion: bundle.CFBundleShortVersionString,
			})
		}
	}

	// If no bundles found, use package identifier
	if len(metadata.IncludedApps) == 0 && pkgInfoXML.Identifier != "" {
		metadata.IncludedApps = append(metadata.IncludedApps, IncludedApp{
			BundleId:      pkgInfoXML.Identifier,
			BundleVersion: pkgInfoXML.Version,
		})
	}

	if len(metadata.IncludedApps) > 0 {
		metadata.PrimaryBundleId = metadata.IncludedApps[0].BundleId
		metadata.PrimaryBundleVersion = metadata.IncludedApps[0].BundleVersion
	} else {
		return nil, fmt.Errorf("no valid included apps found in package")
	}

	return metadata, nil
}
