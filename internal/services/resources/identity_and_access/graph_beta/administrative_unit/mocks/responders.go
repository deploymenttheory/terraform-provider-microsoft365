package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	administrativeUnits        map[string]map[string]any
	deletedAdministrativeUnits map[string]map[string]any
}

func init() {
	mockState.administrativeUnits = make(map[string]map[string]any)
	mockState.deletedAdministrativeUnits = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("administrative_unit", &AdministrativeUnitMock{})
}

type AdministrativeUnitMock struct{}

var _ mocks.MockRegistrar = (*AdministrativeUnitMock)(nil)

// getJSONFileForDisplayName determines which JSON file to load based on the administrative unit's display name
func getJSONFileForDisplayName(displayName string) string {
	displayNameLower := strings.ToLower(displayName)

	switch {
	case strings.Contains(displayNameLower, "user-based"):
		return "post_administrative_unit_au001_success.json"
	case strings.Contains(displayNameLower, "group-based"):
		return "post_administrative_unit_au002_success.json"
	case strings.Contains(displayNameLower, "mixed"):
		return "post_administrative_unit_au003_success.json"
	case strings.Contains(displayNameLower, "dynamic"):
		return "post_administrative_unit_au004_success.json"
	default:
		return "post_administrative_unit_au001_success.json"
	}
}

// loadJSONResponse loads a JSON response from the tests/responses/validate_create directory
func loadJSONResponse(filename string) (map[string]any, error) {
	responsePath := filepath.Join("tests", "responses", "validate_create", filename)

	data, err := os.ReadFile(responsePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read response file %s: %w", responsePath, err)
	}

	var response map[string]any
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

// RegisterMocks registers all mock responders for administrative unit operations
func (m *AdministrativeUnitMock) RegisterMocks() {
	// POST /administrativeUnits - Create
	httpmock.RegisterResponder(
		"POST",
		"https://graph.microsoft.com/beta/administrativeUnits",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			displayName, ok := requestBody["displayName"].(string)
			if !ok || displayName == "" {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"displayName is required"}}`), nil
			}

			// Load appropriate JSON response
			jsonFile := getJSONFileForDisplayName(displayName)
			response, err := loadJSONResponse(jsonFile)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalError","message":"%s"}}`, err.Error())), nil
			}

			// Generate a new ID for the administrative unit
			newID := uuid.New().String()
			response["id"] = newID
			response["displayName"] = displayName

			// Update with request body properties
			if description, ok := requestBody["description"].(string); ok {
				response["description"] = description
			}
			if isMemberManagementRestricted, ok := requestBody["isMemberManagementRestricted"].(bool); ok {
				response["isMemberManagementRestricted"] = isMemberManagementRestricted
			}
			if membershipRule, ok := requestBody["membershipRule"].(string); ok {
				response["membershipRule"] = membershipRule
			}
			if membershipRuleProcessingState, ok := requestBody["membershipRuleProcessingState"].(string); ok {
				response["membershipRuleProcessingState"] = membershipRuleProcessingState
			}
			if membershipType, ok := requestBody["membershipType"].(string); ok {
				response["membershipType"] = membershipType
			}
			if visibility, ok := requestBody["visibility"].(string); ok {
				response["visibility"] = visibility
			}

			// Store in mock state
			mockState.Lock()
			mockState.administrativeUnits[newID] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		},
	)

	// GET /administrativeUnits/{id} - Read
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/administrativeUnits/([0-9a-fA-F-]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			unit, exists := mockState.administrativeUnits[id]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, unit)
		},
	)

	// PATCH /administrativeUnits/{id} - Update
	httpmock.RegisterResponder(
		"PATCH",
		`=~^https://graph\.microsoft\.com/beta/administrativeUnits/([0-9a-fA-F-]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			unit, exists := mockState.administrativeUnits[id]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`), nil
			}

			var updateBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&updateBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update the administrative unit
			mockState.Lock()
			for key, value := range updateBody {
				unit[key] = value
			}
			mockState.administrativeUnits[id] = unit
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		},
	)

	// DELETE /administrativeUnits/{id} - Soft Delete
	httpmock.RegisterResponder(
		"DELETE",
		`=~^https://graph\.microsoft\.com/beta/administrativeUnits/([0-9a-fA-F-]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			unit, exists := mockState.administrativeUnits[id]
			if exists {
				mockState.deletedAdministrativeUnits[id] = unit
				delete(mockState.administrativeUnits, id)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		},
	)

	// GET /directory/deletedItems/{id} - Get deleted administrative unit
	// Used for soft delete verification (polling until resource appears in deleted items)
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/directory/deletedItems/([0-9a-fA-F-]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			deletedUnit, exists := mockState.deletedAdministrativeUnits[id]
			mockState.Unlock()

			if exists {
				return httpmock.NewJsonResponse(200, deletedUnit)
			}

			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`), nil
		},
	)

	// DELETE /directory/deletedItems/{id} - Permanent delete from deleted items
	httpmock.RegisterResponder(
		"DELETE",
		`=~^https://graph\.microsoft\.com/beta/directory/deletedItems/([0-9a-fA-F-]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			_, exists := mockState.deletedAdministrativeUnits[id]
			if exists {
				delete(mockState.deletedAdministrativeUnits, id)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		},
	)
}

// RegisterErrorMocks registers mock responders that return errors
func (m *AdministrativeUnitMock) RegisterErrorMocks() {
	// POST - Create error
	httpmock.RegisterResponder(
		"POST",
		"https://graph.microsoft.com/beta/administrativeUnits",
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`),
	)

	// GET - Read error
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/beta/administrativeUnits/([0-9a-fA-F-]+)$`,
		httpmock.NewStringResponder(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`),
	)

	// PATCH - Update error
	httpmock.RegisterResponder(
		"PATCH",
		`=~^https://graph\.microsoft\.com/beta/administrativeUnits/([0-9a-fA-F-]+)$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`),
	)

	// DELETE - Delete error
	httpmock.RegisterResponder(
		"DELETE",
		`=~^https://graph\.microsoft\.com/beta/administrativeUnits/([0-9a-fA-F-]+)$`,
		httpmock.NewStringResponder(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`),
	)
}

// CleanupMockState clears the mock state
func (m *AdministrativeUnitMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.administrativeUnits = make(map[string]map[string]any)
	mockState.deletedAdministrativeUnits = make(map[string]map[string]any)
}
