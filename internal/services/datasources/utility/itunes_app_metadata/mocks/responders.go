package mocks

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jarcoal/httpmock"
)

// ItunesAppMetadataMock provides mock responses for iTunes app metadata operations
type ItunesAppMetadataMock struct{}

// RegisterMocks registers HTTP mock responses for iTunes app metadata operations
func (m *ItunesAppMetadataMock) RegisterMocks() {
	// Register responder for iTunes Search API
	httpmock.RegisterResponder("GET", `=~^https://itunes.apple.com/search`,
		func(req *http.Request) (*http.Response, error) {
			// Parse the query parameters
			query := req.URL.Query()
			term := query.Get("term")
			// country := query.Get("country") - Not used currently but available for future use

			// Determine which mock response to use based on the search term
			var responseFile string
			switch {
			case strings.Contains(term, "firefox"):
				responseFile = filepath.Join("tests", "firefox_search", "get_itunes_search.json")
			case strings.Contains(term, "office"):
				responseFile = filepath.Join("tests", "office_search", "get_itunes_search.json")
			case strings.Contains(term, "error"):
				// Return an error response for error testing
				return httpmock.NewStringResponse(500, `{"errorMessage": "Internal Server Error"}`), nil
			case strings.Contains(term, "empty"):
				// Return an empty result set
				return httpmock.NewStringResponse(200, `{"resultCount": 0, "results": []}`), nil
			default:
				// Default to Firefox search if no match
				responseFile = filepath.Join("tests", "firefox_search", "get_itunes_search.json")
			}

			// Read the mock response file
			jsonData, err := os.ReadFile(responseFile)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error": "Failed to read mock response: %s"}`, err)), nil
			}

			// Return the mock response
			return httpmock.NewStringResponse(200, string(jsonData)), nil
		})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error": "Resource not found"}`))
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *ItunesAppMetadataMock) RegisterErrorMocks() {
	// Register responder for iTunes Search API that returns an error
	httpmock.RegisterResponder("GET", `=~^https://itunes.apple.com/search`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(500, `{"errorMessage": "Internal Server Error"}`), nil
		})
}
