package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	namedLocations map[string]map[string]interface{}
}

func init() {
	mockState.namedLocations = make(map[string]map[string]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("named_location", &NamedLocationMock{})
}

type NamedLocationMock struct{}

var _ mocks.MockRegistrar = (*NamedLocationMock)(nil)

func (m *NamedLocationMock) RegisterMocks() {
	mockState.Lock()
	mockState.namedLocations = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Create named location - POST /identity/conditionalAccess/namedLocations
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/namedLocations", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Generate a UUID for the new resource
		newId := uuid.New().String()

		// Determine response based on @odata.type
		var jsonStr string
		var err error
		
		if odataType, ok := requestBody["@odata.type"].(string); ok {
			switch odataType {
			case "#microsoft.graph.ipNamedLocation":
				jsonStr, err = helpers.ParseJSONFile("../tests/responses/validate_create/post_ip_named_location_success.json")
			case "#microsoft.graph.countryNamedLocation":
				jsonStr, err = helpers.ParseJSONFile("../tests/responses/validate_create/post_country_named_location_success.json")
			default:
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid @odata.type"}}`), nil
			}
		} else {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Missing @odata.type"}}`), nil
		}

		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load response"}}`), nil
		}

		var responseObj map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse response"}}`), nil
		}

		// Update response with request data
		responseObj["id"] = newId
		if displayName, ok := requestBody["displayName"]; ok {
			responseObj["displayName"] = displayName
		}
		
		// Update response based on @odata.type
		if odataType, ok := requestBody["@odata.type"].(string); ok {
			responseObj["@odata.type"] = odataType
			
			switch odataType {
			case "#microsoft.graph.ipNamedLocation":
				if isTrusted, ok := requestBody["isTrusted"]; ok {
					responseObj["isTrusted"] = isTrusted
				}
				if ipRanges, ok := requestBody["ipRanges"]; ok {
					responseObj["ipRanges"] = ipRanges
				}
			case "#microsoft.graph.countryNamedLocation":
				if countryLookupMethod, ok := requestBody["countryLookupMethod"]; ok {
					responseObj["countryLookupMethod"] = countryLookupMethod
				}
				if includeUnknownCountriesAndRegions, ok := requestBody["includeUnknownCountriesAndRegions"]; ok {
					responseObj["includeUnknownCountriesAndRegions"] = includeUnknownCountriesAndRegions
				}
				if countriesAndRegions, ok := requestBody["countriesAndRegions"]; ok {
					responseObj["countriesAndRegions"] = countriesAndRegions
				}
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.namedLocations[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get named location - GET /identity/conditionalAccess/namedLocations/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/namedLocations/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		locationId := parts[len(parts)-1]

		mockState.Lock()
		location, exists := mockState.namedLocations[locationId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Return the stored location data which includes all the request data
		return httpmock.NewJsonResponse(200, location)
	})

	// Update named location - PATCH /identity/conditionalAccess/namedLocations/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/namedLocations/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		locationId := parts[len(parts)-1]

		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		location, exists := mockState.namedLocations[locationId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update fields from request
		for key, value := range requestBody {
			location[key] = value
		}
		location["modifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.namedLocations[locationId] = location
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete named location - DELETE /identity/conditionalAccess/namedLocations/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/namedLocations/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		locationId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.namedLocations[locationId]
		if exists {
			delete(mockState.namedLocations, locationId)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *NamedLocationMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/namedLocations", httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/namedLocations/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/namedLocations/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/namedLocations/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *NamedLocationMock) CleanupMockState() {
	mockState.Lock()
	mockState.namedLocations = make(map[string]map[string]interface{})
	mockState.Unlock()
}