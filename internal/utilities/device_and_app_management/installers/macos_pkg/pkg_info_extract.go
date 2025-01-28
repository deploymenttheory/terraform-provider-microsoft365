package installers

// import (
// 	"bytes"
// 	"compress/gzip"
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"

// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// 	"howett.net/plist"
// )

// type InfoPlist struct {
// 	CFBundleIdentifier         string `plist:"CFBundleIdentifier"`
// 	CFBundleShortVersionString string `plist:"CFBundleShortVersionString"`
// 	CFBundleVersion            string `plist:"CFBundleVersion"`
// 	CFBundleName               string `plist:"CFBundleName"`
// 	LSMinimumSystemVersion     string `plist:"LSMinimumSystemVersion"`
// }

// func ExtractmacOSPkgMetadata(ctx context.Context, filePath string) (string, string, string, string, []MacOSIncludedAppResourceModel, error) {
// 	tflog.Debug(ctx, "Starting PKG metadata extraction", map[string]interface{}{
// 		"filePath": filePath,
// 	})

// 	// Open the .pkg file
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return "", "", "", "", nil, fmt.Errorf("failed to open pkg file: %w", err)
// 	}
// 	defer file.Close()

// 	// Skip xar header (first 28 bytes)
// 	_, err = file.Seek(28, 0)
// 	if err != nil {
// 		return "", "", "", "", nil, fmt.Errorf("failed to skip xar header: %w", err)
// 	}

// 	// Create a gzip reader for the first extraction
// 	gzr1, err := gzip.NewReader(file)
// 	if err != nil {
// 		return "", "", "", "", nil, fmt.Errorf("failed to create first gzip reader: %w", err)
// 	}
// 	defer gzr1.Close()

// 	tflog.Debug(ctx, "Reading first gzip layer", nil)

// 	// Read until we find the Payload file
// 	buffer := make([]byte, 4096)
// 	var payloadBuf bytes.Buffer
// 	foundPayload := false

// 	for {
// 		n, err := gzr1.Read(buffer)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return "", "", "", "", nil, fmt.Errorf("error reading first gzip content: %w", err)
// 		}

// 		chunk := buffer[:n]
// 		if !foundPayload && bytes.Contains(chunk, []byte("Payload")) {
// 			tflog.Debug(ctx, "Found Payload file", nil)
// 			foundPayload = true
// 		}
// 		if foundPayload {
// 			payloadBuf.Write(chunk)
// 		}
// 	}

// 	if payloadBuf.Len() == 0 {
// 		return "", "", "", "", nil, fmt.Errorf("no Payload file found")
// 	}

// 	// Create a gzip reader for the Payload content
// 	gzr2, err := gzip.NewReader(&payloadBuf)
// 	if err != nil {
// 		return "", "", "", "", nil, fmt.Errorf("failed to create second gzip reader: %w", err)
// 	}
// 	defer gzr2.Close()

// 	tflog.Debug(ctx, "Reading second gzip layer", nil)

// 	// Read the payload content looking for .app/Contents/Info.plist
// 	var includedApps []MacOSIncludedAppResourceModel

// 	buffer = make([]byte, 4096)
// 	var plistBuf bytes.Buffer
// 	foundPlist := false

// 	for {
// 		n, err := gzr2.Read(buffer)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return "", "", "", "", nil, fmt.Errorf("error reading second gzip content: %w", err)
// 		}

// 		chunk := buffer[:n]
// 		if !foundPlist && bytes.Contains(chunk, []byte(".app/Contents/Info.plist")) {
// 			tflog.Debug(ctx, "Found Info.plist file", nil)
// 			foundPlist = true
// 		}
// 		if foundPlist {
// 			plistBuf.Write(chunk)
// 			// Look for end of plist
// 			if bytes.Contains(chunk, []byte("</plist>")) {
// 				break
// 			}
// 		}
// 	}

// 	if plistBuf.Len() == 0 {
// 		return "", "", "", "", nil, fmt.Errorf("no Info.plist file found")
// 	}

// 	// Extract just the plist content
// 	plistContent := plistBuf.Bytes()
// 	startIdx := bytes.Index(plistContent, []byte("<?xml"))
// 	endIdx := bytes.Index(plistContent, []byte("</plist>"))
// 	if startIdx == -1 || endIdx == -1 {
// 		return "", "", "", "", nil, fmt.Errorf("invalid plist format")
// 	}
// 	plistContent = plistContent[startIdx : endIdx+8] // +8 for </plist>

// 	// Parse the plist
// 	var info InfoPlist
// 	decoder := plist.NewDecoder(bytes.NewReader(plistContent))
// 	if err := decoder.Decode(&info); err != nil {
// 		return "", "", "", "", nil, fmt.Errorf("failed to decode Info.plist: %w", err)
// 	}

// 	if info.CFBundleIdentifier == "" {
// 		return "", "", "", "", nil, fmt.Errorf("no bundle identifier found")
// 	}

// 	tflog.Debug(ctx, "Found bundle info", map[string]interface{}{
// 		"bundleId": info.CFBundleIdentifier,
// 		"version":  info.CFBundleShortVersionString,
// 	})

// 	app := MacOSIncludedAppResourceModel{
// 		BundleId:      types.StringValue(info.CFBundleIdentifier),
// 		BundleVersion: types.StringValue(info.CFBundleShortVersionString),
// 	}
// 	includedApps = append(includedApps, app)

// 	return info.CFBundleIdentifier,
// 		info.CFBundleShortVersionString,
// 		info.CFBundleName,
// 		info.LSMinimumSystemVersion,
// 		includedApps,
// 		nil
// }
