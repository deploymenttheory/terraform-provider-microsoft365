package installers

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	xar "github.com/korylprince/goxar"
	"howett.net/plist"
)

// PackageInfo represents the structure of PackageInfo files in macOS packages
type PackageInfo struct {
	XMLName    xml.Name `xml:"pkg-info"`
	Identifier string   `xml:"identifier,attr"`
	Version    string   `xml:"version,attr"`
	Bundles    []Bundle `xml:"bundle"`
}

type Bundle struct {
	Path                       string `xml:"path,attr"`
	ID                         string `xml:"id,attr"`
	CFBundleShortVersionString string `xml:"CFBundleShortVersionString,attr"`
	CFBundleVersion            string `xml:"CFBundleVersion,attr"`
}

type InfoPlist struct {
	CFBundleIdentifier         string `plist:"CFBundleIdentifier"`
	CFBundleShortVersionString string `plist:"CFBundleShortVersionString"`
	CFBundleVersion            string `plist:"CFBundleVersion"`
}

type MacOSIncludedAppResourceModel struct {
	BundleId      types.String `tfsdk:"bundle_id"`
	BundleVersion types.String `tfsdk:"bundle_version"`
}

func ExtractmacOSPkgMetadata(ctx context.Context, filePath string) (string, string, []MacOSIncludedAppResourceModel, error) {
	tflog.Debug(ctx, "Starting PKG metadata extraction", map[string]interface{}{
		"filePath": filePath,
	})

	reader, err := xar.OpenReader(filePath)
	if err != nil {
		tflog.Error(ctx, "Failed to open pkg file", map[string]interface{}{
			"error": err.Error(),
		})
		return "", "", nil, fmt.Errorf("failed to open pkg file: %w", err)
	}
	defer reader.Close()

	tflog.Debug(ctx, "Searching for applications in package", map[string]interface{}{
		"totalFiles": len(reader.File),
	})

	// First, try to find Info.plist files
	var includedApps []MacOSIncludedAppResourceModel
	var primaryBundleId, primaryBundleVersion string

	// Walk through files looking for .app/Contents/Info.plist
	for _, file := range reader.File {
		cleanPath := path.Clean(strings.ReplaceAll(file.Name, "\\", "/"))
		parts := strings.Split(cleanPath, "/")

		tflog.Debug(ctx, "Examining file", map[string]interface{}{
			"path":  cleanPath,
			"parts": parts,
		})

		// Look for Info.plist in app bundles
		if len(parts) >= 4 {
			if strings.HasSuffix(parts[len(parts)-1], "Info.plist") &&
				strings.HasSuffix(parts[len(parts)-3], ".app") &&
				parts[len(parts)-2] == "Contents" {

				tflog.Debug(ctx, "Found potential Info.plist", map[string]interface{}{
					"path": cleanPath,
				})

				appInfo, err := extractInfoPlist(ctx, file)
				if err != nil {
					tflog.Warn(ctx, "Failed to extract Info.plist", map[string]interface{}{
						"path":  cleanPath,
						"error": err.Error(),
					})
					continue
				}

				if appInfo.CFBundleIdentifier != "" {
					app := MacOSIncludedAppResourceModel{
						BundleId:      types.StringValue(appInfo.CFBundleIdentifier),
						BundleVersion: types.StringValue(appInfo.CFBundleShortVersionString),
					}
					includedApps = append(includedApps, app)

					if primaryBundleId == "" {
						primaryBundleId = appInfo.CFBundleIdentifier
						primaryBundleVersion = appInfo.CFBundleShortVersionString
					}
				}
			}
		}
	}

	// If no Info.plist found, try PackageInfo
	if len(includedApps) == 0 {
		tflog.Debug(ctx, "No Info.plist found, checking PackageInfo")

		// Find PackageInfo file
		var packageInfoFile *xar.File
		for _, file := range reader.File {
			if strings.HasSuffix(file.Name, "PackageInfo") {
				packageInfoFile = file
				break
			}
		}

		if packageInfoFile != nil {
			pkgInfo, err := extractPackageInfo(ctx, packageInfoFile)
			if err != nil {
				tflog.Warn(ctx, "Failed to extract PackageInfo", map[string]interface{}{
					"error": err.Error(),
				})
			} else {
				// Process bundles from PackageInfo
				for _, bundle := range pkgInfo.Bundles {
					if bundle.ID != "" && bundle.CFBundleShortVersionString != "" {
						app := MacOSIncludedAppResourceModel{
							BundleId:      types.StringValue(bundle.ID),
							BundleVersion: types.StringValue(bundle.CFBundleShortVersionString),
						}
						includedApps = append(includedApps, app)

						if primaryBundleId == "" {
							primaryBundleId = bundle.ID
							primaryBundleVersion = bundle.CFBundleShortVersionString
						}
					}
				}

				// If still no bundles, use package identifier
				if len(includedApps) == 0 && pkgInfo.Identifier != "" {
					tflog.Debug(ctx, "Using package identifier as fallback", map[string]interface{}{
						"identifier": pkgInfo.Identifier,
					})

					includedApps = append(includedApps, MacOSIncludedAppResourceModel{
						BundleId:      types.StringValue(pkgInfo.Identifier),
						BundleVersion: types.StringValue(pkgInfo.Version),
					})
					primaryBundleId = pkgInfo.Identifier
					primaryBundleVersion = pkgInfo.Version
				}
			}
		}
	}

	if len(includedApps) == 0 {
		tflog.Error(ctx, "No applications found in package", nil)
		return "", "", nil, fmt.Errorf("no applications found in package")
	}

	tflog.Debug(ctx, "Successfully extracted metadata", map[string]interface{}{
		"primaryBundleId":      primaryBundleId,
		"primaryBundleVersion": primaryBundleVersion,
		"includedAppsCount":    len(includedApps),
	})

	return primaryBundleId, primaryBundleVersion, includedApps, nil
}

func extractInfoPlist(ctx context.Context, file *xar.File) (*InfoPlist, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open Info.plist: %w", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read Info.plist: %w", err)
	}

	tflog.Debug(ctx, "Read Info.plist content", map[string]interface{}{
		"fileName": file.Name,
		"size":     len(content),
	})

	var info InfoPlist
	decoder := plist.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode Info.plist: %w", err)
	}

	tflog.Debug(ctx, "Successfully parsed Info.plist", map[string]interface{}{
		"bundleId":           info.CFBundleIdentifier,
		"bundleVersion":      info.CFBundleVersion,
		"bundleShortVersion": info.CFBundleShortVersionString,
	})

	return &info, nil
}

func extractPackageInfo(ctx context.Context, file *xar.File) (*PackageInfo, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open PackageInfo: %w", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read PackageInfo: %w", err)
	}

	var pkgInfo PackageInfo
	if err := xml.Unmarshal(content, &pkgInfo); err != nil {
		return nil, fmt.Errorf("failed to parse PackageInfo XML: %w", err)
	}

	tflog.Debug(ctx, "Successfully parsed PackageInfo", map[string]interface{}{
		"identifier": pkgInfo.Identifier,
		"version":    pkgInfo.Version,
		"bundles":    len(pkgInfo.Bundles),
	})

	return &pkgInfo, nil
}
