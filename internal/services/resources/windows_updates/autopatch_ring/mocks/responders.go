package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	rings map[string]map[string]any
}

func init() {
	mockState.rings = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_update_ring", &WindowsUpdateRingMock{})
}

type WindowsUpdateRingMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdateRingMock)(nil)

func (m *WindowsUpdateRingMock) RegisterMocks() {
	mockState.Lock()
	mockState.rings = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreateRingResponder()
	m.registerGetRingResponder()
	m.registerUpdateRingResponder()
	m.registerDeleteRingResponder()
}

func (m *WindowsUpdateRingMock) registerCreateRingResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F-]+/rings$`,
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)
			responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_ring.json")

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load response: %s"}}`, err.Error())), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse response: %s"}}`, err.Error())), nil
			}

			// Reflect mutable fields from the request into the response
			for _, field := range []string{"displayName", "description", "isPaused", "deferralInDays", "isHotpatchEnabled", "includedGroupAssignment", "excludedGroupAssignment"} {
				if val, ok := requestBody[field]; ok {
					responseObj[field] = val
				}
			}

			mockState.Lock()
			if id, ok := responseObj["id"].(string); ok {
				mockState.rings[id] = responseObj
			}
			mockState.Unlock()

			resp, err := httpmock.NewJsonResponse(200, responseObj)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateRingMock) registerGetRingResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F-]+/rings/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			ringId := parts[len(parts)-1]

			mockState.Lock()
			ring, exists := mockState.rings[ringId]
			mockState.Unlock()

			if !exists {
				_, filename, _, _ := runtime.Caller(0)
				sourceDir := filepath.Dir(filename)
				responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationRead, "get_ring.json")

				jsonData, err := os.ReadFile(responsesPath)
				if err != nil {
					return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Ring not found"}}`), nil
				}

				var responseObj map[string]any
				if err := json.Unmarshal(jsonData, &responseObj); err != nil {
					return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse response: %s"}}`, err.Error())), nil
				}
				ring = responseObj
			}

			resp, err := httpmock.NewJsonResponse(200, ring)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateRingMock) registerUpdateRingResponder() {
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F-]+/rings/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			ringId := parts[len(parts)-1]

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)
			responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationUpdate, "patch_ring.json")

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load response: %s"}}`, err.Error())), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse response: %s"}}`, err.Error())), nil
			}

			responseObj["id"] = ringId
			for _, field := range []string{"displayName", "description", "isPaused", "deferralInDays", "isHotpatchEnabled", "includedGroupAssignment", "excludedGroupAssignment"} {
				if val, ok := requestBody[field]; ok {
					responseObj[field] = val
				}
			}

			mockState.Lock()
			mockState.rings[ringId] = responseObj
			mockState.Unlock()

			resp, err := httpmock.NewJsonResponse(200, responseObj)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateRingMock) registerDeleteRingResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F-]+/rings/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			ringId := parts[len(parts)-1]

			mockState.Lock()
			delete(mockState.rings, ringId)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsUpdateRingMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.rings = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreateRingErrorResponder()
	m.registerGetRingErrorResponder()
}

func (m *WindowsUpdateRingMock) registerCreateRingErrorResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies`,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(403, map[string]any{
				"error": map[string]any{
					"code":    "Forbidden",
					"message": "Insufficient privileges to complete the operation.",
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateRingMock) registerGetRingErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies`,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(403, map[string]any{
				"error": map[string]any{
					"code":    "Forbidden",
					"message": "Insufficient privileges to complete the operation.",
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateRingMock) CleanupMockState() {
	mockState.Lock()
	mockState.rings = make(map[string]map[string]any)
	mockState.Unlock()
}
