package common

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DownloadFile downloads a file from a URL and saves it to a temporary directory with security validation.
// It follows redirects up to a maximum of 10 times. The filename is determined in the following order:
// 1. From Content-Disposition header if present and valid
// 2. From the final URL after redirects if valid
// 3. Falls back to a timestamp-based name if both sources are invalid
//
// The function implements several security measures:
// - Validates filenames to prevent directory traversal
// - Restricts filenames to alphanumeric characters, dots, hyphens, underscores, and spaces
// - Replaces '%' characters with '_'
// - Ensures final path remains within temporary directory
// - Verifies path safety after normalization
//
// Returns the path to the downloaded file and any error encountered. Possible errors include:
// - Failed to create temporary file
// - Too many redirects
// - Download failure
// - Write failure
// - Invalid filename
// - Path traversal attempt
// - Rename failure
// func DownloadFile(url string) (string, error) {
// 	tmpFile, err := os.CreateTemp("", "downloaded-*")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create temporary file: %v", err)
// 	}
// 	defer tmpFile.Close()

// 	client := &http.Client{
// 		CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 			if len(via) >= 10 {
// 				return fmt.Errorf("too many redirects when attempting to download file from %s", url)
// 			}
// 			return nil
// 		},
// 	}

// 	resp, err := client.Get(url)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to download file from %s: %v", url, err)
// 	}
// 	defer resp.Body.Close()

// 	_, err = io.Copy(tmpFile, resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to write to temporary file: %v", err)
// 	}

// 	var finalFileName string
// 	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
// 	filename := params["filename"]
// 	if err == nil && filename != "" {
// 		var sanitizeErr error
// 		finalFileName, sanitizeErr = sanitizeFileName(filename)
// 		if sanitizeErr != nil {
// 			log.Printf("[WARN] Failed to sanitize filename from Content-Disposition: %v", sanitizeErr)
// 		}
// 	}

// 	if finalFileName == "" {
// 		finalURL := resp.Request.URL.String()
// 		urlFileName := filepath.Base(finalURL)
// 		if validName, err := sanitizeFileName(urlFileName); err == nil {
// 			finalFileName = validName
// 		} else {
// 			finalFileName = fmt.Sprintf("downloaded-file-%d.png", time.Now().Unix())
// 		}
// 	}

// 	finalPath := filepath.Join(os.TempDir(), finalFileName)
// 	if !strings.HasPrefix(filepath.Clean(finalPath), os.TempDir()) {
// 		return "", fmt.Errorf("security error: final path '%s' would be outside temporary directory", finalPath)
// 	}

// 	err = os.Rename(tmpFile.Name(), finalPath)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to rename temporary file to final destination: %v", err)
// 	}

// 	log.Printf("[INFO] File downloaded to: %s", finalPath)
// 	return finalPath, nil
// }

// // sanitizeFileName cleans and validates a filename for secure file operations
// func sanitizeFileName(name string) (string, error) {
// 	name = filepath.Base(name)
// 	if strings.Contains(name, "..") {
// 		return "", fmt.Errorf("invalid filename: contains parent directory reference")
// 	}

// 	unescaped, unescapeErr := url.PathUnescape(name)
// 	if unescapeErr != nil {
// 		unescaped = name
// 	}

// 	cleaned := strings.TrimSpace(unescaped)
// 	if cleaned == "" || cleaned == "." {
// 		return "", fmt.Errorf("invalid filename: empty or invalid after cleaning")
// 	}

// 	return cleaned, nil
// }

// DownloadFile downloads a file from a URL and saves it to a temporary directory with security validation.
// It follows redirects up to a maximum of 10 times. The filename is determined in the following order:
// 1. From Content-Disposition header if present and valid
// 2. From the final URL after redirects if valid
// 3. Falls back to a timestamp-based name if both sources are invalid
//
// The function implements several security measures:
// - Validates filenames to prevent directory traversal
// - Restricts filenames to alphanumeric characters, dots, hyphens, underscores, and spaces
// - Ensures final path remains within temporary directory
// - Verifies path safety after normalization
//
// Returns the path to the downloaded file and any error encountered.
func DownloadFile(sourceURL string) (string, error) {
	// Create temporary file with a placeholder name
	tmpFile, err := os.CreateTemp("", "download-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	tmpPath := tmpFile.Name()

	// Close the temp file but keep the path for writing
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file: %v", err)
	}

	// Configure the HTTP client with redirect handling
	client := &http.Client{
		Timeout: 5 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	// Make the HTTP request
	req, err := http.NewRequest("GET", sourceURL, nil)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to create request for URL %s: %v", sourceURL, err)
	}

	// Add a User-Agent header to avoid 403 errors from some servers
	req.Header.Set("User-Agent", "Mozilla/5.0 Terraform-Microsoft365-Provider")

	resp, err := client.Do(req)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to download file from %s: %v", sourceURL, err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to download file, server returned status code: %d", resp.StatusCode)
	}

	// Create a file for writing the downloaded content
	out, err := os.Create(tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to open temporary file for writing: %v", err)
	}

	// Download the file content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		out.Close() // Close the file before attempting to remove it
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to write downloaded content to file: %v", err)
	}

	// Ensure all data is written to disk
	if err := out.Sync(); err != nil {
		out.Close() // Close the file before attempting to remove it
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to flush file to disk: %v", err)
	}

	// IMPORTANT: Close the file BEFORE renaming on Windows
	// Windows won't allow renaming a file that has an open handle
	out.Close()

	// Determine the final filename using a layered approach
	var finalFileName string

	// Try to get filename from Content-Disposition header
	if disposition := resp.Header.Get("Content-Disposition"); disposition != "" {
		if _, params, err := parseContentDisposition(disposition); err == nil {
			if filename, ok := params["filename"]; ok && filename != "" {
				if safe, safeFileName := sanitizeFileName(filename); safe {
					finalFileName = safeFileName
				}
			}
		}
	}

	// If Content-Disposition didn't yield a valid filename, try the URL path
	if finalFileName == "" {
		urlPath := resp.Request.URL.Path
		urlFileName := filepath.Base(urlPath)
		if safe, safeFileName := sanitizeFileName(urlFileName); safe {
			finalFileName = safeFileName
		}
	}

	// Last resort: use a timestamp-based name
	if finalFileName == "" {
		finalFileName = fmt.Sprintf("download-%d", time.Now().Unix())

		// Try to determine extension from Content-Type
		if contentType := resp.Header.Get("Content-Type"); contentType != "" {
			if ext, ok := mimeTypeToExtension(contentType); ok {
				finalFileName += ext
			}
		}
	}

	// Construct the final path in the temporary directory
	finalPath := filepath.Join(os.TempDir(), finalFileName)

	// Verify that the final path is still within the temporary directory (security check)
	if !strings.HasPrefix(filepath.Clean(finalPath), filepath.Clean(os.TempDir())) {
		os.Remove(tmpPath)
		return "", fmt.Errorf("security error: path traversal attempt detected in filename")
	}

	// If a file already exists at the destination path, use a unique name
	if _, err := os.Stat(finalPath); err == nil {
		timestamp := time.Now().UnixNano()
		ext := filepath.Ext(finalPath)
		base := strings.TrimSuffix(finalPath, ext)
		finalPath = fmt.Sprintf("%s-%d%s", base, timestamp, ext)
	}

	// On Windows, files with spaces can cause issues, especially with some tools
	// Use a safer filename by replacing spaces with underscores
	finalPath = strings.ReplaceAll(finalPath, " ", "_")

	// Move the temporary file to its final location
	// Try up to 3 times with short delays between attempts to handle possible file locking issues
	var renameErr error
	for retries := 0; retries < 3; retries++ {
		renameErr = os.Rename(tmpPath, finalPath)
		if renameErr == nil {
			break
		}

		// If rename failed, wait a bit before trying again
		time.Sleep(100 * time.Millisecond)
	}

	if renameErr != nil {
		// If rename still fails after retries, try copy and delete as fallback
		// This is slower but more reliable, especially on Windows
		srcFile, err := os.Open(tmpPath)
		if err == nil {
			defer srcFile.Close()

			destFile, err := os.Create(finalPath)
			if err == nil {
				defer destFile.Close()

				_, err = io.Copy(destFile, srcFile)
				if err == nil {
					destFile.Close() // Close explicitly before continuing
					srcFile.Close()  // Close explicitly before deletion

					// Now try to delete the temp file, but don't fail if we can't
					os.Remove(tmpPath)
					return finalPath, nil
				}
				destFile.Close()
			}
			srcFile.Close()
		}

		// If copy failed too, clean up and return error from original rename attempt
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to move downloaded file to final destination: %v", renameErr)
	}

	return finalPath, nil
}

// sanitizeFileName validates and cleans a filename for secure file operations.
// Returns a boolean indicating if the name is safe and the sanitized filename.
func sanitizeFileName(name string) (bool, string) {
	// Strip directory parts
	name = filepath.Base(name)

	// Attempt to URL-unescape if it seems to be URL-encoded
	if strings.Contains(name, "%") {
		if unescaped, err := url.PathUnescape(name); err == nil {
			name = unescaped
		}
	}

	// Remove any characters that aren't alphanumeric, dots, hyphens, underscores, or spaces
	safePattern := regexp.MustCompile(`[^a-zA-Z0-9.\-_ ]`)
	cleaned := safePattern.ReplaceAllString(name, "")

	// Trim leading/trailing spaces
	cleaned = strings.TrimSpace(cleaned)

	// Validate the cleaned name
	if cleaned == "" || cleaned == "." || cleaned == ".." || strings.Contains(cleaned, "..") {
		return false, ""
	}

	return true, cleaned
}

// parseContentDisposition parses a Content-Disposition header value
// Returns the disposition type and a map of parameters
func parseContentDisposition(header string) (string, map[string]string, error) {
	parts := strings.Split(header, ";")
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("invalid Content-Disposition header")
	}

	dispositon := strings.TrimSpace(parts[0])
	params := make(map[string]string)

	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		// Remove quotes if present
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		params[key] = value
	}

	return dispositon, params, nil
}

// mimeTypeToExtension maps common MIME types to file extensions
func mimeTypeToExtension(mimeType string) (string, bool) {
	// Extract primary type
	parts := strings.Split(mimeType, ";")
	primaryType := strings.TrimSpace(parts[0])

	extensionMap := map[string]string{
		"application/octet-stream":                      ".bin",
		"application/x-msdownload":                      ".exe",
		"application/x-apple-diskimage":                 ".dmg",
		"application/zip":                               ".zip",
		"application/x-zip-compressed":                  ".zip",
		"application/x-tar":                             ".tar",
		"application/x-gzip":                            ".gz",
		"application/x-bzip2":                           ".bz2",
		"application/x-rar-compressed":                  ".rar",
		"application/vnd.microsoft.portable-executable": ".exe",
		"application/pkix-cert":                         ".cer",
		"application/x-x509-ca-cert":                    ".crt",
		"application/x-pem-file":                        ".pem",
		"application/x-pkcs12":                          ".p12",
		"application/x-msi":                             ".msi",
		"application/x-ms-wim":                          ".wim",
		"application/x-ms-application":                  ".application",
		"application/x-ms-installer":                    ".msi",
		"application/pkg":                               ".pkg",
		"application/x-itunes-pkg":                      ".pkg",
		"application/vnd.apple.installer+xml":           ".pkg",
		"application/xml":                               ".xml",
		"application/json":                              ".json",
		"text/plain":                                    ".txt",
		"text/html":                                     ".html",
		"text/xml":                                      ".xml",
		"text/css":                                      ".css",
		"text/javascript":                               ".js",
		"image/jpeg":                                    ".jpg",
		"image/png":                                     ".png",
		"image/gif":                                     ".gif",
		"image/svg+xml":                                 ".svg",
		"image/bmp":                                     ".bmp",
		"image/webp":                                    ".webp",
		"image/tiff":                                    ".tif",
		"image/x-icon":                                  ".ico",
		"audio/mpeg":                                    ".mp3",
		"audio/wav":                                     ".wav",
		"audio/ogg":                                     ".ogg",
		"video/mp4":                                     ".mp4",
		"video/mpeg":                                    ".mpeg",
		"video/webm":                                    ".webm",
	}

	if ext, ok := extensionMap[primaryType]; ok {
		return ext, true
	}

	return "", false
}
