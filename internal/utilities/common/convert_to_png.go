package common

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	// Register standard image formats
	_ "image/gif"
	_ "image/jpeg"
)

// ConvertToPNG converts an image from various formats (JPEG, GIF, etc.) to PNG.
// It accepts either a byte slice containing the image data or a file path.
// Returns the PNG image as a byte slice.
func ConvertToPNG(ctx context.Context, input any) ([]byte, error) {
	var img image.Image
	var err error
	var format string

	switch v := input.(type) {
	case []byte:
		// Input is a byte slice
		img, format, err = image.Decode(bytes.NewReader(v))
		if err != nil {
			return nil, fmt.Errorf("failed to decode image data: %v", err)
		}
	case string:
		// Input is a file path
		file, err := os.Open(v)
		if err != nil {
			return nil, fmt.Errorf("failed to open image file: %v", err)
		}
		defer file.Close()

		img, format, err = image.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image file: %v", err)
		}
	default:
		return nil, fmt.Errorf("input must be either a byte slice or a file path string")
	}

	tflog.Debug(ctx, fmt.Sprintf("Converting image from format %s to PNG", format))

	// Create a buffer to store the PNG image
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode image to PNG: %v", err)
	}

	return buf.Bytes(), nil
}

// ConvertFileToPNG converts an image file to PNG format and saves it to the specified output path.
// If outputPath is empty, it will replace the original file extension with .png.
func ConvertFileToPNG(ctx context.Context, inputPath string, outputPath string) (string, error) {
	if outputPath == "" {
		// Generate output path by replacing the extension with .png
		ext := filepath.Ext(inputPath)
		outputPath = strings.TrimSuffix(inputPath, ext) + ".png"
	}

	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	// Decode the image
	img, format, err := image.Decode(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Converting image file from format %s to PNG", format))

	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Encode as PNG
	if err := png.Encode(outputFile, img); err != nil {
		return "", fmt.Errorf("failed to encode image to PNG: %v", err)
	}

	return outputPath, nil
}

// ConvertURLImageToPNG downloads an image from a URL and converts it to PNG format.
// Returns the PNG image as a byte slice.
func ConvertURLImageToPNG(ctx context.Context, url string) ([]byte, error) {
	// Download the image
	imageData, err := DownloadImage(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %v", err)
	}

	// Convert to PNG
	pngData, err := ConvertToPNG(ctx, imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image to PNG: %v", err)
	}

	return pngData, nil
}

// IsPNG checks if the provided image data is already in PNG format.
func IsPNG(data []byte) bool {
	// Check for PNG signature
	return len(data) > 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n"
}

// GetImageFormat detects the format of an image from its bytes.
// Returns the format as a string (e.g., "png", "jpeg", "gif").
func GetImageFormat(data []byte) (string, error) {
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to detect image format: %v", err)
	}
	return format, nil
}

// SaveImageAsPNG saves image data as a PNG file.
func SaveImageAsPNG(data []byte, outputPath string) error {
	// Create the output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// If the data is already PNG, write it directly
	if IsPNG(data) {
		_, err = io.Copy(outputFile, bytes.NewReader(data))
		return err
	}

	// Otherwise, decode and re-encode as PNG
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode image data: %v", err)
	}

	// Encode as PNG
	if err := png.Encode(outputFile, img); err != nil {
		return fmt.Errorf("failed to encode image to PNG: %v", err)
	}

	return nil
}
