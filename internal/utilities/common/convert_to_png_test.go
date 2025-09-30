package common

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createTestImage creates a valid image in the specified format for testing
func createTestImage(format string, width, height int) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with some colors to make it a valid image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := color.RGBA{uint8(x % 256), uint8(y % 256), 128, 255}
			img.Set(x, y, c)
		}
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err := jpeg.Encode(&buf, img, nil)
		return buf.Bytes(), err
	case "png":
		err := png.Encode(&buf, img)
		return buf.Bytes(), err
	case "gif":
		err := gif.Encode(&buf, img, nil)
		return buf.Bytes(), err
	default:
		return nil, nil
	}
}

func TestConvertToPNG(t *testing.T) {
	ctx := context.Background()

	// Create valid test images
	jpegData, err := createTestImage("jpeg", 10, 10)
	if err != nil {
		t.Fatalf("Failed to create test JPEG: %v", err)
	}

	pngData, err := createTestImage("png", 10, 10)
	if err != nil {
		t.Fatalf("Failed to create test PNG: %v", err)
	}

	gifData, err := createTestImage("gif", 10, 10)
	if err != nil {
		t.Fatalf("Failed to create test GIF: %v", err)
	}

	tests := []struct {
		name      string
		input     any
		expectErr bool
	}{
		{
			name:      "JPEG bytes to PNG",
			input:     jpegData,
			expectErr: false,
		},
		{
			name:      "PNG bytes to PNG",
			input:     pngData,
			expectErr: false,
		},
		{
			name:      "GIF bytes to PNG",
			input:     gifData,
			expectErr: false,
		},
		{
			name:      "invalid bytes",
			input:     []byte{0x00, 0x01, 0x02, 0x03},
			expectErr: true,
		},
		{
			name:      "empty bytes",
			input:     []byte{},
			expectErr: true,
		},
		{
			name:      "invalid input type",
			input:     123,
			expectErr: true,
		},
		{
			name:      "nil input",
			input:     nil,
			expectErr: true,
		},
	}

	// Add file path test
	tempDir, err := os.MkdirTemp("", "png_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test JPEG file
	jpegFile := filepath.Join(tempDir, "test.jpg")
	if err := os.WriteFile(jpegFile, jpegData, 0644); err != nil {
		t.Fatalf("Failed to write test JPEG file: %v", err)
	}

	tests = append(tests, []struct {
		name      string
		input     any
		expectErr bool
	}{
		{
			name:      "JPEG file path to PNG",
			input:     jpegFile,
			expectErr: false,
		},
		{
			name:      "non-existent file path",
			input:     "/non/existent/file.jpg",
			expectErr: true,
		},
		{
			name:      "directory path instead of file",
			input:     tempDir,
			expectErr: true,
		},
	}...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToPNG(ctx, tt.input)

			if tt.expectErr {
				if err == nil {
					t.Errorf("ConvertToPNG() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ConvertToPNG() error = %v", err)
				return
			}

			// Check if result is PNG
			if !IsPNG(result) {
				t.Errorf("ConvertToPNG() result is not PNG")
			}

			if len(result) == 0 {
				t.Errorf("ConvertToPNG() returned empty result")
			}
		})
	}
}

func TestIsPNG(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "Valid PNG",
			data:     []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00},
			expected: true,
		},
		{
			name:     "Not PNG",
			data:     []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46},
			expected: false,
		},
		{
			name:     "Empty data",
			data:     []byte{},
			expected: false,
		},
		{
			name:     "Too short",
			data:     []byte{0x89, 0x50, 0x4E, 0x47},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPNG(tt.data)
			if result != tt.expected {
				t.Errorf("IsPNG() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetImageFormat(t *testing.T) {
	// Create valid test images
	jpegData, err := createTestImage("jpeg", 5, 5)
	if err != nil {
		t.Fatalf("Failed to create test JPEG: %v", err)
	}

	pngData, err := createTestImage("png", 5, 5)
	if err != nil {
		t.Fatalf("Failed to create test PNG: %v", err)
	}

	gifData, err := createTestImage("gif", 5, 5)
	if err != nil {
		t.Fatalf("Failed to create test GIF: %v", err)
	}

	tests := []struct {
		name           string
		data           []byte
		expectedFormat string
		expectErr      bool
	}{
		{
			name:           "JPEG format detection",
			data:           jpegData,
			expectedFormat: "jpeg",
			expectErr:      false,
		},
		{
			name:           "PNG format detection",
			data:           pngData,
			expectedFormat: "png",
			expectErr:      false,
		},
		{
			name:           "GIF format detection",
			data:           gifData,
			expectedFormat: "gif",
			expectErr:      false,
		},
		{
			name:      "Invalid data",
			data:      []byte{0x00, 0x01, 0x02, 0x03},
			expectErr: true,
		},
		{
			name:      "Empty data",
			data:      []byte{},
			expectErr: true,
		},
		{
			name:      "Nil data",
			data:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, err := GetImageFormat(tt.data)

			if tt.expectErr {
				if err == nil {
					t.Errorf("GetImageFormat() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetImageFormat() error = %v", err)
				return
			}

			if format != tt.expectedFormat {
				t.Errorf("GetImageFormat() = %s, want %s", format, tt.expectedFormat)
			}
		})
	}
}

func TestSaveImageAsPNG(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "image_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a simple PNG for testing
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
		0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	outputPath := filepath.Join(tempDir, "test.png")

	// Test saving PNG data
	err = SaveImageAsPNG(pngData, outputPath)
	if err != nil {
		t.Fatalf("SaveImageAsPNG() error = %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("SaveImageAsPNG() did not create file")
	}

	// Read the file back and check if it's still a PNG
	savedData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if !IsPNG(savedData) {
		t.Errorf("Saved file is not a PNG")
	}
}

func TestConvertFileToPNG(t *testing.T) {
	ctx := context.Background()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "png_convert_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test images
	jpegData, err := createTestImage("jpeg", 8, 8)
	if err != nil {
		t.Fatalf("Failed to create test JPEG: %v", err)
	}

	gifData, err := createTestImage("gif", 8, 8)
	if err != nil {
		t.Fatalf("Failed to create test GIF: %v", err)
	}

	// Write test files
	jpegFile := filepath.Join(tempDir, "test.jpg")
	if err := os.WriteFile(jpegFile, jpegData, 0644); err != nil {
		t.Fatalf("Failed to write JPEG file: %v", err)
	}

	gifFile := filepath.Join(tempDir, "test.gif")
	if err := os.WriteFile(gifFile, gifData, 0644); err != nil {
		t.Fatalf("Failed to write GIF file: %v", err)
	}

	tests := []struct {
		name        string
		inputPath   string
		outputPath  string
		expectErr   bool
		checkOutput func(outputPath string, t *testing.T)
	}{
		{
			name:       "JPEG to PNG with auto output path",
			inputPath:  jpegFile,
			outputPath: "",
			expectErr:  false,
			checkOutput: func(outputPath string, t *testing.T) {
				if !strings.HasSuffix(outputPath, ".png") {
					t.Errorf("Output path should end with .png, got %s", outputPath)
				}
				data, err := os.ReadFile(outputPath)
				if err != nil {
					t.Errorf("Failed to read output file: %v", err)
				}
				if !IsPNG(data) {
					t.Errorf("Output file is not a PNG")
				}
			},
		},
		{
			name:       "GIF to PNG with custom output path",
			inputPath:  gifFile,
			outputPath: filepath.Join(tempDir, "custom.png"),
			expectErr:  false,
			checkOutput: func(outputPath string, t *testing.T) {
				if !strings.HasSuffix(outputPath, "custom.png") {
					t.Errorf("Output path should end with custom.png, got %s", outputPath)
				}
			},
		},
		{
			name:       "non-existent input file",
			inputPath:  filepath.Join(tempDir, "nonexistent.jpg"),
			outputPath: "",
			expectErr:  true,
		},
		{
			name:       "invalid image file",
			inputPath:  jpegFile,
			outputPath: "",
			expectErr:  false, // Should succeed since we have valid JPEG
		},
	}

	// Create an invalid image file
	invalidFile := filepath.Join(tempDir, "invalid.jpg")
	if err := os.WriteFile(invalidFile, []byte("not an image"), 0644); err != nil {
		t.Fatalf("Failed to write invalid file: %v", err)
	}

	tests = append(tests, struct {
		name        string
		inputPath   string
		outputPath  string
		expectErr   bool
		checkOutput func(outputPath string, t *testing.T)
	}{
		name:       "invalid image data",
		inputPath:  invalidFile,
		outputPath: "",
		expectErr:  true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertFileToPNG(ctx, tt.inputPath, tt.outputPath)

			if tt.expectErr {
				if err == nil {
					t.Errorf("ConvertFileToPNG() expected error but got nil")
					if result != "" {
						os.Remove(result) // Cleanup
					}
				}
				return
			}

			if err != nil {
				t.Errorf("ConvertFileToPNG() error = %v", err)
				return
			}

			if result == "" {
				t.Errorf("ConvertFileToPNG() returned empty path")
				return
			}

			// Check file exists
			if _, err := os.Stat(result); os.IsNotExist(err) {
				t.Errorf("ConvertFileToPNG() did not create output file")
				return
			}

			// Run custom check if provided
			if tt.checkOutput != nil {
				tt.checkOutput(result, t)
			}

			// Cleanup
			os.Remove(result)
		})
	}
}

func TestConvertURLImageToPNG(t *testing.T) {
	ctx := context.Background()

	// This function depends on DownloadImage which we can't easily test without
	// a real server or mocking the entire HTTP infrastructure.
	// We'll test the error case with an invalid URL
	_, err := ConvertURLImageToPNG(ctx, "not-a-valid-url")
	if err == nil {
		t.Errorf("ConvertURLImageToPNG() expected error for invalid URL but got nil")
	}

	// Test with empty URL
	_, err = ConvertURLImageToPNG(ctx, "")
	if err == nil {
		t.Errorf("ConvertURLImageToPNG() expected error for empty URL but got nil")
	}
}

func TestSaveImageAsPNG_ExtendedTests(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "save_png_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test images
	jpegData, err := createTestImage("jpeg", 6, 6)
	if err != nil {
		t.Fatalf("Failed to create test JPEG: %v", err)
	}

	pngData, err := createTestImage("png", 6, 6)
	if err != nil {
		t.Fatalf("Failed to create test PNG: %v", err)
	}

	tests := []struct {
		name      string
		data      []byte
		path      string
		expectErr bool
	}{
		{
			name:      "Save JPEG as PNG",
			data:      jpegData,
			path:      filepath.Join(tempDir, "jpeg_to_png.png"),
			expectErr: false,
		},
		{
			name:      "Save PNG as PNG (direct copy)",
			data:      pngData,
			path:      filepath.Join(tempDir, "png_to_png.png"),
			expectErr: false,
		},
		{
			name:      "Save to invalid path",
			data:      pngData,
			path:      "/invalid/path/that/does/not/exist/file.png",
			expectErr: true,
		},
		{
			name:      "Save invalid image data",
			data:      []byte{0x00, 0x01, 0x02, 0x03},
			path:      filepath.Join(tempDir, "invalid.png"),
			expectErr: true,
		},
		{
			name:      "Save empty data",
			data:      []byte{},
			path:      filepath.Join(tempDir, "empty.png"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveImageAsPNG(tt.data, tt.path)

			if tt.expectErr {
				if err == nil {
					t.Errorf("SaveImageAsPNG() expected error but got nil")
					os.Remove(tt.path) // Cleanup if file was created
				}
				return
			}

			if err != nil {
				t.Errorf("SaveImageAsPNG() error = %v", err)
				return
			}

			// Check file exists
			if _, err := os.Stat(tt.path); os.IsNotExist(err) {
				t.Errorf("SaveImageAsPNG() did not create file")
				return
			}

			// Read file back and check if it's PNG
			savedData, err := os.ReadFile(tt.path)
			if err != nil {
				t.Errorf("Failed to read saved file: %v", err)
				return
			}

			if !IsPNG(savedData) {
				t.Errorf("Saved file is not a PNG")
			}

			// Cleanup
			os.Remove(tt.path)
		})
	}
}
