package utilities

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
)

// TestBase64Encode tests the Base64Encode function.
func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		expectErr bool
	}{
		{
			name:      "Normal string",
			input:     "Hello, World!",
			want:      "SGVsbG8sIFdvcmxkIQ==",
			expectErr: false,
		},
		{
			name:      "Empty string",
			input:     "",
			want:      "",
			expectErr: true,
		},
		{
			name:      "Special characters",
			input:     "GoLang@2024!",
			want:      "R29MYW5nQDIwMjQh",
			expectErr: false,
		},
		{
			name:      "Unicode characters",
			input:     "こんにちは",
			want:      "44GT44KT44Gr44Gh44Gv",
			expectErr: false,
		},
		{
			name:      "Long string",
			input:     "The quick brown fox jumps over the lazy dog",
			want:      "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZw==",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64Encode(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("Base64Encode() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("Base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStringToInt tests the StringToInt function.
func TestStringToInt(t *testing.T) {
	mapping := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	tests := []struct {
		name      string
		input     string
		mapping   map[string]int
		want      int
		expectErr bool
	}{
		{
			name:      "Valid key 'one'",
			input:     "one",
			mapping:   mapping,
			want:      1,
			expectErr: false,
		},
		{
			name:      "Valid key 'two'",
			input:     "two",
			mapping:   mapping,
			want:      2,
			expectErr: false,
		},
		{
			name:      "Invalid key 'four'",
			input:     "four",
			mapping:   mapping,
			want:      -1,
			expectErr: true,
		},
		{
			name:      "Empty string",
			input:     "",
			mapping:   mapping,
			want:      -1,
			expectErr: true,
		},
		{
			name:      "Case sensitivity",
			input:     "One",
			mapping:   mapping,
			want:      -1,
			expectErr: true,
		},
		{
			name:      "Numeric string",
			input:     "123",
			mapping:   mapping,
			want:      -1,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToInt(tt.input, tt.mapping)
			if (err != nil) != tt.expectErr {
				t.Errorf("StringToInt() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBoolPtr tests the BoolPtr function.
func TestBoolPtr(t *testing.T) {
	tests := []struct {
		name  string
		input bool
		want  *bool
	}{
		{
			name:  "True value",
			input: true,
			want:  BoolPtr(true),
		},
		{
			name:  "False value",
			input: false,
			want:  BoolPtr(false),
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := BoolPtr(tt.input)
			if got == nil || *got != *tt.want {
				t.Errorf("BoolPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test StringPtr utility function
func TestStringPtr(t *testing.T) {

	// Test case 1: Valid string
	t.Run("Non-empty string", func(t *testing.T) {
		input := "test"
		expected := "test"
		result := StringPtr(input)

		// Check if result is not nil
		if result == nil {
			t.Fatalf("Expected non-nil result for input %s", input)
		}

		// Check if the value pointed to by result matches expected value
		if *result != expected {
			t.Errorf("Expected %s, got %s", expected, *result)
		}
	})

	// Test case 2: Empty string
	t.Run("Empty string", func(t *testing.T) {
		input := ""
		expected := ""
		result := StringPtr(input)

		// Check if result is not nil
		if result == nil {
			t.Fatalf("Expected non-nil result for input empty string")
		}

		// Check if the value pointed to by result matches expected value
		if *result != expected {
			t.Errorf("Expected empty string, got %s", *result)
		}
	})

	// Test case 3: String with special characters
	t.Run("String with special characters", func(t *testing.T) {
		input := "!@#$%^&*()_+"
		expected := "!@#$%^&*()_+"
		result := StringPtr(input)

		// Check if result is not nil
		if result == nil {
			t.Fatalf("Expected non-nil result for input %s", input)
		}

		// Check if the value pointed to by result matches expected value
		if *result != expected {
			t.Errorf("Expected %s, got %s", expected, *result)
		}
	})
}

// TestDownloadImage tests the DownloadImage function.
func TestDownloadImage(t *testing.T) {
	// Define a sample image data
	sampleImageData := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header bytes

	// Variables to handle multiple redirects
	var redirectCount int32

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/valid-image.png":
			w.WriteHeader(http.StatusOK)
			w.Write(sampleImageData)
		case "/not-found.png":
			w.WriteHeader(http.StatusNotFound)
		case "/redirect":
			// Increment redirect count atomically
			current := atomic.AddInt32(&redirectCount, 1)
			if current > 10 {
				// Exceeded redirect limit
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// Redirect back to /redirect to simulate a loop
			http.Redirect(w, r, "/redirect", http.StatusFound)
		case "/too-many-redirects":
			// Redirect to /redirect, which in turn may redirect back
			http.Redirect(w, r, "/redirect", http.StatusFound)
		default:
			w.WriteHeader(http.StatusOK)
			w.Write(sampleImageData)
		}
	}))
	defer server.Close()

	tests := []struct {
		name      string
		url       string
		want      []byte
		expectErr bool
	}{
		{
			name:      "Valid image download",
			url:       server.URL + "/valid-image.png",
			want:      sampleImageData,
			expectErr: false,
		},
		{
			name:      "Image not found (404)",
			url:       server.URL + "/not-found.png",
			want:      nil,
			expectErr: true,
		},
		{
			name:      "Redirect to valid image",
			url:       server.URL + "/redirect",
			want:      nil, // Since /redirect keeps redirecting, the final response is 500
			expectErr: true,
		},
		{
			name:      "Too many redirects",
			url:       server.URL + "/too-many-redirects",
			want:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadImage(tt.url)
			if (err != nil) != tt.expectErr {
				t.Errorf("DownloadImage() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !compareByteSlices(got, tt.want) {
				t.Errorf("DownloadImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDownloadImage_InvalidURL tests DownloadImage with an invalid URL format.
func TestDownloadImage_InvalidURL(t *testing.T) {
	invalidURL := "://invalid-url"

	_, err := DownloadImage(invalidURL)
	if err == nil {
		t.Errorf("DownloadImage() expected error for invalid URL, got nil")
	}
}

// TestDownloadImage_EmptyURL tests DownloadImage with an empty URL.
func TestDownloadImage_EmptyURL(t *testing.T) {
	emptyURL := ""

	_, err := DownloadImage(emptyURL)
	if err == nil {
		t.Errorf("DownloadImage() expected error for empty URL, got nil")
	}
}

// TestDownloadImage_NonImageContent tests DownloadImage with non-image content.
func TestDownloadImage_NonImageContent(t *testing.T) {
	// Start a local HTTP server that returns plain text
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is not an image.")
	}))
	defer server.Close()

	url := server.URL + "/not-an-image"

	// Attempt to download as image; depending on implementation, it might not error
	data, err := DownloadImage(url)
	if err != nil {
		t.Errorf("DownloadImage() unexpected error: %v", err)
	}

	// Since the function does not validate image content, it should return the data
	expectedData := []byte("This is not an image.\n")
	if !compareByteSlices(data, expectedData) {
		t.Errorf("DownloadImage() = %v, want %v", data, expectedData)
	}
}

// Helper function to compare two byte slices.
func compareByteSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

// TestIsDebugMode tests the IsDebugMode function with various environment variable settings.
func TestIsDebugMode(t *testing.T) {
	// Save the original environment variable value to restore it later
	originalValue, originalExists := os.LookupEnv("M365_DEBUG_MODE")

	// Clean up function to restore the original environment after tests
	defer func() {
		if originalExists {
			os.Setenv("M365_DEBUG_MODE", originalValue)
		} else {
			os.Unsetenv("M365_DEBUG_MODE")
		}
	}()

	tests := []struct {
		name     string
		envValue string
		setEnv   bool
		want     bool
	}{
		{
			name:     "Debug mode enabled with 'true'",
			envValue: "true",
			setEnv:   true,
			want:     true,
		},
		{
			name:     "Debug mode enabled with 'TRUE'",
			envValue: "TRUE",
			setEnv:   true,
			want:     true,
		},
		{
			name:     "Debug mode enabled with 'True'",
			envValue: "True",
			setEnv:   true,
			want:     true,
		},
		{
			name:     "Debug mode disabled with 'false'",
			envValue: "false",
			setEnv:   true,
			want:     false,
		},
		{
			name:     "Debug mode disabled with empty string",
			envValue: "",
			setEnv:   true,
			want:     false,
		},
		{
			name:     "Debug mode disabled with invalid value",
			envValue: "invalid",
			setEnv:   true,
			want:     false,
		},
		{
			name:     "Environment variable not set",
			envValue: "",
			setEnv:   false,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset environment variable based on test case
			if tt.setEnv {
				os.Setenv("M365_DEBUG_MODE", tt.envValue)
			} else {
				os.Unsetenv("M365_DEBUG_MODE")
			}

			got := IsDebugMode()
			if got != tt.want {
				t.Errorf("IsDebugMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
