package androidManagedAppProtectionMocks

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

const (
	baseURL   = "https://graph.microsoft.com/beta/deviceAppManagement/androidManagedAppProtections"
	mockIDMin = "00000000-0000-0000-0000-000000000001"
	mockIDMax = "00000000-0000-0000-0000-000000000002"
)

//go:embed testdata/minimal_response.json
var minimalResponse string

//go:embed testdata/maximal_response.json
var maximalResponse string

//go:embed testdata/error_response.json
var errorResponse string

// AndroidManagedAppProtectionMock manages mock HTTP responses for unit tests.
type AndroidManagedAppProtectionMock struct {
	mu    sync.Mutex
	state map[string]map[string]any
}

// RegisterMocks registers all standard HTTP mock responders.
func (m *AndroidManagedAppProtectionMock) RegisterMocks() {
	m.mu.Lock()
	m.state = make(map[string]map[string]any)
	m.mu.Unlock()

	httpmock.RegisterResponder("POST", baseURL, func(req *http.Request) (*http.Response, error) {
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, fmt.Errorf("decode mock create request: %w", err)
		}
		fixture := minimalResponse
		if payload["printBlocked"] == true {
			fixture = maximalResponse
		}
		var stored map[string]any
		if err := json.Unmarshal([]byte(fixture), &stored); err != nil {
			return nil, fmt.Errorf("decode mock create fixture: %w", err)
		}
		for key, value := range payload {
			stored[key] = value
		}
		id := uuid.NewString()
		stored["id"] = id
		m.mu.Lock()
		defer m.mu.Unlock()
		m.state[id] = stored
		return newJSONResponse(http.StatusCreated, stored)
	})

	itemURL := "=~^" + regexp.QuoteMeta(baseURL) + "/[^/?]+(?:\\?.*)?$"
	httpmock.RegisterResponder("GET", itemURL, func(req *http.Request) (*http.Response, error) {
		m.mu.Lock()
		defer m.mu.Unlock()
		stored, exists := m.state[path.Base(req.URL.Path)]
		if !exists {
			return newJSONResponse(
				http.StatusNotFound,
				map[string]any{
					"error": map[string]string{
						"code":    "ResourceNotFound",
						"message": "Policy not found",
					},
				},
			)
		}
		return newJSONResponse(http.StatusOK, stored)
	})

	httpmock.RegisterResponder("PATCH", itemURL, func(req *http.Request) (*http.Response, error) {
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, fmt.Errorf("decode mock update request: %w", err)
		}
		m.mu.Lock()
		defer m.mu.Unlock()
		id := path.Base(req.URL.Path)
		stored, exists := m.state[id]
		if !exists {
			return newJSONResponse(
				http.StatusNotFound,
				map[string]any{
					"error": map[string]string{
						"code":    "ResourceNotFound",
						"message": "Policy not found",
					},
				},
			)
		}
		for key, value := range payload {
			stored[key] = value
		}
		stored["id"] = id
		return newJSONResponse(http.StatusOK, stored)
	})

	httpmock.RegisterResponder("DELETE", itemURL, func(req *http.Request) (*http.Response, error) {
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.state, path.Base(req.URL.Path))
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	})
}

// newJSONResponse returns a mock response with the content type required by the Graph SDK.
func newJSONResponse(status int, payload any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, payload)
	if err != nil {
		return nil, fmt.Errorf("encode mock response: %w", err)
	}
	return response, nil
}

// RegisterErrorMocks registers HTTP mock responders that return errors.
func (m *AndroidManagedAppProtectionMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", baseURL,
		httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, json.RawMessage(errorResponse)),
	)
}

// CleanupMockState resets the mock state between tests.
func (m *AndroidManagedAppProtectionMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = make(map[string]map[string]any)
}
