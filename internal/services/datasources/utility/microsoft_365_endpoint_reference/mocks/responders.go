package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"

	"github.com/jarcoal/httpmock"
)

func init() {
	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":"endpoint not found"}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("microsoft_365_endpoints", &Microsoft365EndpointReferenceMock{})
}

// Microsoft365EndpointReferenceMock provides mock responses for Microsoft 365 endpoints operations
type Microsoft365EndpointReferenceMock struct{}

// Ensure Microsoft365EndpointReferenceMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*Microsoft365EndpointReferenceMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for Microsoft 365 endpoints operations
// This implements the MockRegistrar interface
func (m *Microsoft365EndpointReferenceMock) RegisterMocks() {
	// GET https://endpoints.office.com/endpoints/worldwide
	httpmock.RegisterResponder("GET", `=~^https://endpoints\.office\.com/endpoints/worldwide`,
		m.getEndpointsResponder("worldwide"))

	// GET https://endpoints.office.com/endpoints/USGOVDoD
	httpmock.RegisterResponder("GET", `=~^https://endpoints\.office\.com/endpoints/USGOVDoD`,
		m.getEndpointsResponder("USGOVDoD"))

	// GET https://endpoints.office.com/endpoints/USGOVGCCHigh
	httpmock.RegisterResponder("GET", `=~^https://endpoints\.office\.com/endpoints/USGOVGCCHigh`,
		m.getEndpointsResponder("USGOVGCCHigh"))

	// GET https://endpoints.office.com/endpoints/China
	httpmock.RegisterResponder("GET", `=~^https://endpoints\.office\.com/endpoints/China`,
		m.getEndpointsResponder("China"))
}

// RegisterErrorMocks sets up mock HTTP responders that return error responses
func (m *Microsoft365EndpointReferenceMock) RegisterErrorMocks() {
	// For this datasource, we don't need specific error mocks as validation is client-side
	// The API itself doesn't require authentication and rarely fails
}

// CleanupMockState clears the mock state for clean test runs
func (m *Microsoft365EndpointReferenceMock) CleanupMockState() {
	// This datasource is read-only, no state to clean up
}

// loadJSONResponse loads a JSON response from a file
func (m *Microsoft365EndpointReferenceMock) loadJSONResponse(filePath string) ([]map[string]any, error) {
	var response []map[string]any

	content, err := os.ReadFile(filePath)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(content, &response)
	return response, err
}

// getEndpointsResponder returns a responder function for Microsoft 365 endpoints API
func (m *Microsoft365EndpointReferenceMock) getEndpointsResponder(instance string) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Map instance name to file name
		var fileName string
		switch instance {
		case "worldwide":
			fileName = "get_endpoints_worldwide.json"
		case "USGOVDoD":
			fileName = "get_endpoints_usgov_dod.json"
		case "USGOVGCCHigh":
			fileName = "get_endpoints_usgov_gcchigh.json"
		case "China":
			fileName = "get_endpoints_china.json"
		default:
			return httpmock.NewStringResponse(404, `{"error":"invalid instance"}`), nil
		}

		// Load the response from file
		responseData, err := m.loadJSONResponse(filepath.Join("tests", "responses", fileName))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":"failed to load response"}`), nil
		}

		// Apply filters based on query parameters
		filteredData := m.applyFilters(responseData, req)

		return httpmock.NewJsonResponse(200, filteredData)
	}
}

// applyFilters applies client-side filtering to mimic server behavior
func (m *Microsoft365EndpointReferenceMock) applyFilters(data []map[string]any, req *http.Request) []map[string]any {
	// For the mock, we don't need to implement filtering since the datasource
	// does filtering client-side after fetching all data
	return data
}

// Helper functions to extract query parameters (not used in current implementation but kept for completeness)
func getQueryParam(req *http.Request, param string) string {
	return req.URL.Query().Get(param)
}

func hasQueryParam(req *http.Request, param string) bool {
	return req.URL.Query().Has(param)
}

// Helper to check if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
