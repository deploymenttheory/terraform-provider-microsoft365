package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

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
