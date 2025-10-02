package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	agreements map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.agreements = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("conditional_access_terms_of_use", &ConditionalAccessTermsOfUseMock{})
}

// ConditionalAccessTermsOfUseMock provides mock responses for terms of use operations
type ConditionalAccessTermsOfUseMock struct{}

// Ensure ConditionalAccessTermsOfUseMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*ConditionalAccessTermsOfUseMock)(nil)

// RegisterMocks registers HTTP mock responses for terms of use operations
func (m *ConditionalAccessTermsOfUseMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.agreements = make(map[string]map[string]any)
	mockState.Unlock()

	// Register POST for creating agreements
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/v1.0/agreements",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Load the base response template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_agreement.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}

			// Update response with request data
			if displayName, ok := requestBody["displayName"].(string); ok {
				response["displayName"] = displayName
			}
			if isViewingBeforeAcceptanceRequired, ok := requestBody["isViewingBeforeAcceptanceRequired"].(bool); ok {
				response["isViewingBeforeAcceptanceRequired"] = isViewingBeforeAcceptanceRequired
			}
			if isPerDeviceAcceptanceRequired, ok := requestBody["isPerDeviceAcceptanceRequired"].(bool); ok {
				response["isPerDeviceAcceptanceRequired"] = isPerDeviceAcceptanceRequired
			}
			if userReacceptRequiredFrequency, ok := requestBody["userReacceptRequiredFrequency"].(string); ok {
				response["userReacceptRequiredFrequency"] = userReacceptRequiredFrequency
			} else {
				delete(response, "userReacceptRequiredFrequency")
			}
			if termsExpiration, ok := requestBody["termsExpiration"].(map[string]any); ok {
				// Preserve the exact format from the request to avoid normalization issues
				responseTermsExpiration := make(map[string]any)
				if startDateTime, exists := termsExpiration["startDateTime"]; exists {
					responseTermsExpiration["startDateTime"] = startDateTime
				}
				if frequency, exists := termsExpiration["frequency"]; exists {
					responseTermsExpiration["frequency"] = frequency
				}
				response["termsExpiration"] = responseTermsExpiration
			} else {
				delete(response, "termsExpiration")
			}
			if file, ok := requestBody["file"].(map[string]any); ok {
				response["file"] = file
			}

			// Store in mock state
			mockState.Lock()
			mockState.agreements[response["id"].(string)] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register GET for individual agreements
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/v1\.0/agreements/[a-fA-F0-9-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			agreement, exists := mockState.agreements[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_agreement_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Create response copy
			agreementCopy := make(map[string]any)
			for k, v := range agreement {
				agreementCopy[k] = v
			}

			return httpmock.NewJsonResponse(200, agreementCopy)
		})

	// Register PATCH for updating agreements
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/v1\.0/agreements/[a-fA-F0-9-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_agreement_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Load update template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/get_agreement_updated.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var updatedAgreement map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &updatedAgreement); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}

			mockState.Lock()
			agreement, exists := mockState.agreements[id]
			if !exists {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_agreement_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Start with existing data
			for k, v := range agreement {
				updatedAgreement[k] = v
			}

			// Apply updates from request body
			for k, v := range requestBody {
				updatedAgreement[k] = v
			}

			// Store updated state
			mockState.agreements[id] = updatedAgreement
			mockState.Unlock()

			return factories.SuccessResponse(200, updatedAgreement)(req)
		})

	// Register DELETE for agreements
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/v1\.0/agreements/[a-fA-F0-9-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.agreements[id]
			if exists {
				delete(mockState.agreements, id)
			}
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_agreement_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *ConditionalAccessTermsOfUseMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored agreements
	for id := range mockState.agreements {
		delete(mockState.agreements, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *ConditionalAccessTermsOfUseMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.agreements = make(map[string]map[string]any)
	mockState.Unlock()

	// Register error response for creating agreement with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/v1.0/agreements",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_agreement_error.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Register error response for agreement not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/v1\.0/agreements/[a-fA-F0-9-]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_agreement_not_found.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}
