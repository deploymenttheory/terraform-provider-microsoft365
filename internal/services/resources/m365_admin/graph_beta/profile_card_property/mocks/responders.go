package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

var mockState struct {
	sync.Mutex
	profileCardProperties map[string]map[string]any
}

const (
	profileCardPropertiesURL   = "https://graph.microsoft.com/beta/admin/people/profileCardProperties"
	profileCardPropertyItemURL = `=~^https://graph\.microsoft\.com/beta/admin/people/profileCardProperties/[^/?]+$`
	existingPropertyName       = "customAttribute9"
	notFoundBody               = `{"error":{"code":"notFound","message":"Not Found","innerError":{"peopleAdminErrorCode":"peopleAdminItemNotFound"}}}`
	conflictBody               = `{"error":{"code":"conflict","message":"A profile card property with this directoryPropertyName already exists."}}`
	badRequestBody             = `{"error":{"code":"badRequest","message":"Invalid request"}}`
	forbiddenBody              = `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`
)

func init() {
	mockState.profileCardProperties = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, notFoundBody))
	mocks.GlobalRegistry.Register("m365_admin_profile_card_property", &ProfileCardPropertyMock{})
}

type ProfileCardPropertyMock struct{}

var _ mocks.MockRegistrar = (*ProfileCardPropertyMock)(nil)

// loadResponse reads a JSON response fixture from tests/responses/validate_create
func loadResponse(fileName string) (map[string]any, error) {
	jsonData, err := os.ReadFile(filepath.Join("tests", "responses", "validate_create", fileName))
	if err != nil {
		return nil, fmt.Errorf("read fixture %s: %w", fileName, err)
	}
	var responseObj map[string]any
	if err := json.Unmarshal(jsonData, &responseObj); err != nil {
		return nil, fmt.Errorf("parse fixture %s: %w", fileName, err)
	}
	return responseObj, nil
}

// propertyNameFromPath returns the directoryPropertyName key from an item URL
func propertyNameFromPath(req *http.Request) string {
	parts := strings.Split(req.URL.Path, "/")
	name, err := url.PathUnescape(parts[len(parts)-1])
	if err != nil {
		return parts[len(parts)-1]
	}
	return name
}

// normalizeAnnotations mirrors Graph, which returns an empty localizations collection for annotations without localizations
func normalizeAnnotations(value any) []any {
	annotations, ok := value.([]any)
	if !ok {
		return []any{}
	}
	for _, item := range annotations {
		annotation, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if localizations, ok := annotation["localizations"].([]any); !ok || localizations == nil {
			annotation["localizations"] = []any{}
		}
	}
	return annotations
}

func (m *ProfileCardPropertyMock) RegisterMocks() {
	mockState.Lock()
	mockState.profileCardProperties = make(map[string]map[string]any)
	mockState.Unlock()

	// Seed a pre-existing property with lower-case casing so case-sensitive import is exercised
	if existing, err := loadResponse("get_profile_card_property_existing.json"); err == nil {
		mockState.Lock()
		mockState.profileCardProperties[existingPropertyName] = existing
		mockState.Unlock()
	}

	// Create profile card property - POST /admin/people/profileCardProperties
	httpmock.RegisterResponder("POST", profileCardPropertiesURL, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, badRequestBody), nil
		}

		name, ok := requestBody["directoryPropertyName"].(string)
		if !ok || name == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"badRequest","message":"directoryPropertyName is required"}}`), nil
		}
		if value, ok := requestBody["annotations"]; ok && value == nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"badRequest","message":"annotations cannot be null"}}`), nil
		}

		responseObj, err := loadResponse("post_profile_card_property_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"%s"}}`, err.Error())), nil
		}

		// Graph stores the key exactly as sent and defaults isVisible to true
		responseObj["directoryPropertyName"] = name
		if isVisible, ok := requestBody["isVisible"].(bool); ok {
			responseObj["isVisible"] = isVisible
		}
		if annotations, ok := requestBody["annotations"]; ok {
			responseObj["annotations"] = normalizeAnnotations(annotations)
		}

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.profileCardProperties[name]; exists {
			return httpmock.NewStringResponse(409, conflictBody), nil
		}
		mockState.profileCardProperties[name] = responseObj

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get profile card property - GET /admin/people/profileCardProperties/{directoryPropertyName} (case-sensitive key)
	httpmock.RegisterResponder("GET", profileCardPropertyItemURL, func(req *http.Request) (*http.Response, error) {
		name := propertyNameFromPath(req)

		mockState.Lock()
		property, exists := mockState.profileCardProperties[name]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		return httpmock.NewJsonResponse(200, property)
	})

	// Update profile card property - PATCH /admin/people/profileCardProperties/{directoryPropertyName}
	httpmock.RegisterResponder("PATCH", profileCardPropertyItemURL, func(req *http.Request) (*http.Response, error) {
		name := propertyNameFromPath(req)

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, badRequestBody), nil
		}
		if value, ok := requestBody["annotations"]; ok && value == nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"badRequest","message":"annotations cannot be null"}}`), nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		property, exists := mockState.profileCardProperties[name]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		// Properties left out of a PATCH keep their value; an empty annotations collection clears them
		if isVisible, ok := requestBody["isVisible"].(bool); ok {
			property["isVisible"] = isVisible
		}
		if annotations, ok := requestBody["annotations"]; ok {
			property["annotations"] = normalizeAnnotations(annotations)
		}
		mockState.profileCardProperties[name] = property

		return httpmock.NewJsonResponse(200, property)
	})

	// Delete profile card property - DELETE /admin/people/profileCardProperties/{directoryPropertyName}
	httpmock.RegisterResponder("DELETE", profileCardPropertyItemURL, func(req *http.Request) (*http.Response, error) {
		name := propertyNameFromPath(req)

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.profileCardProperties[name]; exists {
			delete(mockState.profileCardProperties, name)
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewStringResponse(404, notFoundBody), nil
	})
}

func (m *ProfileCardPropertyMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", profileCardPropertiesURL, httpmock.NewStringResponder(400, badRequestBody))
	httpmock.RegisterResponder("GET", profileCardPropertyItemURL, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterResponder("PATCH", profileCardPropertyItemURL, httpmock.NewStringResponder(400, badRequestBody))
	httpmock.RegisterResponder("DELETE", profileCardPropertyItemURL, httpmock.NewStringResponder(400, badRequestBody))
}

// RegisterForbiddenMocks registers 403 responses, as returned when the caller lacks PeopleSettings.ReadWrite.All
func (m *ProfileCardPropertyMock) RegisterForbiddenMocks() {
	httpmock.RegisterResponder("POST", profileCardPropertiesURL, httpmock.NewStringResponder(403, forbiddenBody))
	httpmock.RegisterResponder("GET", profileCardPropertyItemURL, httpmock.NewStringResponder(403, forbiddenBody))
	httpmock.RegisterResponder("PATCH", profileCardPropertyItemURL, httpmock.NewStringResponder(403, forbiddenBody))
	httpmock.RegisterResponder("DELETE", profileCardPropertyItemURL, httpmock.NewStringResponder(403, forbiddenBody))
}

func (m *ProfileCardPropertyMock) CleanupMockState() {
	mockState.Lock()
	mockState.profileCardProperties = make(map[string]map[string]any)
	mockState.Unlock()
}
