package graphBetaApplications

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	//pkg "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"howett.net/plist"
)

func constructMacOSPkgAppResource(ctx context.Context, data *MacOSPkgAppResourceModel, baseApp graphmodels.MacOSPkgAppable) (graphmodels.MacOSPkgAppable, error) {
	if data.PackageInstallerFileSource.IsNull() || data.PackageInstallerFileSource.ValueString() == "" {
		return nil, fmt.Errorf("package_installer_file_source is required but not provided")
	}

	var includedApps []MacOSIncludedAppResourceModel
	var primaryBundleId, primaryBundleVersion string

	// Extract metadata from the provided package installer
	if !data.PackageInstallerFileSource.IsNull() {
		primaryBundleId, primaryBundleVersion, _, _, bundledApps, err := ExtractmacOSPkgMetadata(ctx, data.PackageInstallerFileSource.ValueString())
		if err != nil {
			return nil, fmt.Errorf("failed to extract metadata from the provided .pkg file at '%s': %w", data.PackageInstallerFileSource.ValueString(), err)
		}

		// Convert bundledApps to the appropriate type
		convertedApps := make([]MacOSIncludedAppResourceModel, len(bundledApps))
		for i, app := range bundledApps {
			convertedApps[i] = MacOSIncludedAppResourceModel{
				BundleId:      app.BundleId,
				BundleVersion: app.BundleVersion,
			}
		}

		// Set PrimaryBundleId and PrimaryBundleVersion if not explicitly set in data
		if data.PrimaryBundleId.IsNull() {
			data.PrimaryBundleId = types.StringValue(primaryBundleId)
		}
		if data.PrimaryBundleVersion.IsNull() {
			data.PrimaryBundleVersion = types.StringValue(primaryBundleVersion)
		}

		// Assign the converted apps to IncludedApps
		data.IncludedApps = convertedApps
	}

	// Deduplicate IncludedApps based on BundleId
	uniqueApps := make(map[string]MacOSIncludedAppResourceModel)
	for _, app := range includedApps {
		uniqueApps[app.BundleId.ValueString()] = app
	}
	includedApps = make([]MacOSIncludedAppResourceModel, 0, len(uniqueApps))
	for _, app := range uniqueApps {
		includedApps = append(includedApps, app)
	}

	// Set IncludedApps in the base app model
	graphIncludedApps := make([]graphmodels.MacOSIncludedAppable, 0, len(includedApps))
	for _, app := range includedApps {
		includedApp := graphmodels.NewMacOSIncludedApp()
		constructors.SetStringProperty(app.BundleId, includedApp.SetBundleId)
		constructors.SetStringProperty(app.BundleVersion, includedApp.SetBundleVersion)
		graphIncludedApps = append(graphIncludedApps, includedApp)
	}
	baseApp.SetIncludedApps(graphIncludedApps)

	// Set PrimaryBundleId and PrimaryBundleVersion
	constructors.SetStringProperty(types.StringValue(primaryBundleId), baseApp.SetPrimaryBundleId)
	constructors.SetStringProperty(types.StringValue(primaryBundleVersion), baseApp.SetPrimaryBundleVersion)

	// Set IgnoreVersionDetection
	constructors.SetBoolProperty(data.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)

	// Set MinimumSupportedOperatingSystem
	if data.MinimumSupportedOperatingSystem != nil {
		minOS := graphmodels.NewMacOSMinimumOperatingSystem()
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
		baseApp.SetMinimumSupportedOperatingSystem(minOS)
	}

	// Set PreInstallScript
	if data.PreInstallScript != nil {
		preScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PreInstallScript.ScriptContent, preScript.SetScriptContent)
		baseApp.SetPreInstallScript(preScript)
	}

	// Set PostInstallScript
	if data.PostInstallScript != nil {
		postScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PostInstallScript.ScriptContent, postScript.SetScriptContent)
		baseApp.SetPostInstallScript(postScript)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully constructed MacOSPkgApp resource with %d unique IncludedApps", len(includedApps)))
	return baseApp, nil
}

type InfoPlist struct {
	CFBundleIdentifier         string `plist:"CFBundleIdentifier"`
	CFBundleShortVersionString string `plist:"CFBundleShortVersionString"`
	CFBundleVersion            string `plist:"CFBundleVersion"`
	CFBundleName               string `plist:"CFBundleName"`
	LSMinimumSystemVersion     string `plist:"LSMinimumSystemVersion"`
}

func ExtractmacOSPkgMetadata(ctx context.Context, pkgPath string) (primaryBundleId string, primaryBundleVersion string, primaryBundleName string, minOSVersion string, bundledApps []MacOSIncludedAppResourceModel, err error) {
	tflog.Info(ctx, fmt.Sprintf("Starting metadata extraction from pkg file: %s", pkgPath))

	// Open the .pkg file
	file, err := os.Open(pkgPath)
	if err != nil {
		return "", "", "", "", nil, fmt.Errorf("failed to open pkg file: %w", err)
	}
	defer file.Close()

	// Read and validate XAR header
	header := make([]byte, 28)
	if _, err := file.Read(header); err != nil {
		return "", "", "", "", nil, fmt.Errorf("failed to read XAR header: %w", err)
	}

	// Verify magic number "xar!"
	if string(header[0:4]) != "xar!" {
		return "", "", "", "", nil, fmt.Errorf("invalid XAR format: incorrect magic number")
	}

	// Parse header fields (big endian)
	headerSize := uint16(header[4])<<8 | uint16(header[5])
	version := uint16(header[6])<<8 | uint16(header[7])
	tocCompressedLength := uint64(header[8])<<56 | uint64(header[9])<<48 | uint64(header[10])<<40 | uint64(header[11])<<32 |
		uint64(header[12])<<24 | uint64(header[13])<<16 | uint64(header[14])<<8 | uint64(header[15])

	tflog.Debug(ctx, fmt.Sprintf("XAR Header: size=%d, version=%d, TOC compressed length=%d",
		headerSize, version, tocCompressedLength))

	// Skip any remaining header bytes
	if headerSize > 28 {
		if _, err := file.Seek(int64(headerSize), io.SeekStart); err != nil {
			return "", "", "", "", nil, fmt.Errorf("failed to skip header padding: %w", err)
		}
	}

	// Read compressed TOC
	compressedTOC := make([]byte, tocCompressedLength)
	if _, err := file.Read(compressedTOC); err != nil {
		return "", "", "", "", nil, fmt.Errorf("failed to read compressed TOC: %w", err)
	}

	// Decompress TOC (zlib)
	tocReader, err := gzip.NewReader(bytes.NewReader(compressedTOC))
	if err != nil {
		return "", "", "", "", nil, fmt.Errorf("failed to create TOC decompressor: %w", err)
	}
	defer tocReader.Close()

	tocContent, err := io.ReadAll(tocReader)
	if err != nil {
		return "", "", "", "", nil, fmt.Errorf("failed to decompress TOC: %w", err)
	}

	tflog.Debug(ctx, "Successfully decompressed TOC XML")

	// Track unique bundle IDs to avoid duplicates
	uniqueBundleIds := make(map[string]bool)
	var apps []MacOSIncludedAppResourceModel

	// Parse TOC XML to find Info.plist files and their data
	type XarTOC struct {
		Files []struct {
			ID   string `xml:"id,attr"`
			Name string `xml:"name"`
			Type string `xml:"type"`
			Data struct {
				Offset   int64  `xml:"offset"`
				Length   int64  `xml:"length"`
				Size     int64  `xml:"size"`
				Encoding string `xml:"encoding>style"`
			} `xml:"data"`
		} `xml:"toc>file"`
	}

	var toc XarTOC
	if err := xml.Unmarshal(tocContent, &toc); err != nil {
		return "", "", "", "", nil, fmt.Errorf("failed to parse TOC XML: %w", err)
	}

	// Process each file in the archive
	for _, xarFile := range toc.Files {
		if strings.HasSuffix(xarFile.Name, "Info.plist") {
			tflog.Debug(ctx, fmt.Sprintf("Found Info.plist at offset %d", xarFile.Data.Offset))

			// Seek to file data
			if _, err := file.Seek(int64(headerSize)+int64(tocCompressedLength)+xarFile.Data.Offset, io.SeekStart); err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to seek to Info.plist data: %s", err))
				continue
			}

			// Read compressed data
			compressedData := make([]byte, xarFile.Data.Length)
			if _, err := file.Read(compressedData); err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to read Info.plist compressed data: %s", err))
				continue
			}

			// Decompress if needed
			var plistData []byte
			if xarFile.Data.Encoding == "application/x-gzip" {
				reader, err := gzip.NewReader(bytes.NewReader(compressedData))
				if err != nil {
					tflog.Warn(ctx, fmt.Sprintf("Failed to create Info.plist decompressor: %s", err))
					continue
				}
				plistData, err = io.ReadAll(reader)
				reader.Close()
				if err != nil {
					tflog.Warn(ctx, fmt.Sprintf("Failed to decompress Info.plist: %s", err))
					continue
				}
			} else {
				plistData = compressedData
			}

			// Parse Info.plist
			var info InfoPlist
			decoder := plist.NewDecoder(bytes.NewReader(plistData))
			if err := decoder.Decode(&info); err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to parse Info.plist: %s", err))
				continue
			}

			if info.CFBundleIdentifier != "" && !uniqueBundleIds[info.CFBundleIdentifier] {
				tflog.Info(ctx, fmt.Sprintf("Found application bundle: %s (version: %s)",
					info.CFBundleIdentifier, info.CFBundleShortVersionString))

				// Use CFBundleShortVersionString if available, fallback to CFBundleVersion
				version := info.CFBundleShortVersionString
				if version == "" {
					version = info.CFBundleVersion
				}

				app := MacOSIncludedAppResourceModel{
					BundleId:      types.StringValue(info.CFBundleIdentifier),
					BundleVersion: types.StringValue(version),
				}
				apps = append(apps, app)
				uniqueBundleIds[info.CFBundleIdentifier] = true

				// If this is the first app found, use it as primary
				if primaryBundleId == "" {
					primaryBundleId = info.CFBundleIdentifier
					primaryBundleVersion = version
					primaryBundleName = info.CFBundleName
					minOSVersion = info.LSMinimumSystemVersion
					tflog.Info(ctx, fmt.Sprintf("Set primary bundle: %s (version: %s, min OS: %s)",
						primaryBundleId, primaryBundleVersion, minOSVersion))
				}
			}
		}
	}

	if len(apps) == 0 {
		tflog.Warn(ctx, "No application bundles found in the pkg file")
		return "", "", "", "", nil, fmt.Errorf("no application bundles found in the pkg file")
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully extracted metadata for %d application bundles", len(apps)))
	return primaryBundleId, primaryBundleVersion, primaryBundleName, minOSVersion, apps, nil
}
