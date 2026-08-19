package androidManagedAppProtectionMocks

import (
	_ "embed"
	"net/http"
	"sync"

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
	state map[string]string
}

// RegisterMocks registers all standard HTTP mock responders.
func (m *AndroidManagedAppProtectionMock) RegisterMocks() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = make(map[string]string)

	httpmock.RegisterResponder("POST", baseURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, minimalResponse), nil
		},
	)

	httpmock.RegisterResponder("GET", baseURL+"/"+mockIDMin,
		httpmock.NewStringResponder(http.StatusOK, minimalResponse),
	)
	httpmock.RegisterResponder("PATCH", baseURL+"/"+mockIDMin,
		httpmock.NewStringResponder(http.StatusOK, minimalResponse),
	)
	httpmock.RegisterResponder("DELETE", baseURL+"/"+mockIDMin,
		httpmock.NewStringResponder(http.StatusNoContent, ""),
	)

	httpmock.RegisterResponder("GET", baseURL+"/"+mockIDMax,
		httpmock.NewStringResponder(http.StatusOK, maximalResponse),
	)
	httpmock.RegisterResponder("PATCH", baseURL+"/"+mockIDMax,
		httpmock.NewStringResponder(http.StatusOK, maximalResponse),
	)
	httpmock.RegisterResponder("DELETE", baseURL+"/"+mockIDMax,
		httpmock.NewStringResponder(http.StatusNoContent, ""),
	)
}

// RegisterErrorMocks registers HTTP mock responders that return errors.
func (m *AndroidManagedAppProtectionMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", baseURL,
		httpmock.NewStringResponder(http.StatusBadRequest, errorResponse),
	)
}

// CleanupMockState resets the mock state between tests.
func (m *AndroidManagedAppProtectionMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = make(map[string]string)
}
