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
	audienceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_deployment_audience/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	updatePolicies map[string]map[string]any
}

func init() {
	mockState.updatePolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_update_policy", &WindowsUpdatePolicyMock{})
}

type WindowsUpdatePolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdatePolicyMock)(nil)

func (m *WindowsUpdatePolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.updatePolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// Register audience mocks as well
	audienceMock := &audienceMocks.WindowsUpdateDeploymentAudienceMock{}
	audienceMock.RegisterMocks()

	// Create update policy - POST /admin/windows/updates/updatePolicies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/updatePolicies",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body to determine which response file to use
			var requestBody map[string]any
			json.NewDecoder(req.Body).Decode(&requestBody)

			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)
			
			// Determine which response file based on request content
			var responsesPath string
			hasRules := false
			hasSettings := false
			
			if rules, ok := requestBody["complianceChangeRules"].([]any); ok && len(rules) > 0 {
				hasRules = true
			}
			if settings, ok := requestBody["deploymentSettings"].(map[string]any); ok && settings != nil {
				hasSettings = true
			}

			// Choose appropriate response file
			if hasRules && hasSettings {
				responsesPath = filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_update_policy_01_full.json")
			} else {
				responsesPath = filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_update_policy_03_minimal.json")
			}

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
			}

			newID := uuid.New().String()
			responseObj["id"] = newID

			mockState.Lock()
			mockState.updatePolicies[newID] = responseObj
			mockState.Unlock()

			resp, err := httpmock.NewJsonResponse(201, responseObj)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})

	// Get update policy - GET /admin/windows/updates/updatePolicies/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			policyID := parts[len(parts)-1]

			mockState.Lock()
			policy, exists := mockState.updatePolicies[policyID]
			mockState.Unlock()

			if !exists {
				_, filename, _, _ := runtime.Caller(0)
				sourceDir := filepath.Dir(filename)
				responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationRead, "get_update_policy_success.json")

				jsonData, err := os.ReadFile(responsesPath)
				if err != nil {
					return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Update policy not found"}}`), nil
				}

				var responseObj map[string]any
				if err := json.Unmarshal(jsonData, &responseObj); err != nil {
					return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
				}

				policy = responseObj
			}

			resp, err := httpmock.NewJsonResponse(200, policy)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})

	// Update update policy - PATCH /admin/windows/updates/updatePolicies/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			policyID := parts[len(parts)-1]

			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)
			responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationUpdate, "patch_update_policy_success.json")

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
			}

			responseObj["id"] = policyID

			mockState.Lock()
			mockState.updatePolicies[policyID] = responseObj
			mockState.Unlock()

			resp, err := httpmock.NewJsonResponse(200, responseObj)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})

	// Delete update policy - DELETE /admin/windows/updates/updatePolicies/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			policyID := parts[len(parts)-1]

			mockState.Lock()
			delete(mockState.updatePolicies, policyID)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsUpdatePolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.updatePolicies = make(map[string]map[string]any)
	mockState.Unlock()

	// Register audience mocks as well
	audienceMock := &audienceMocks.WindowsUpdateDeploymentAudienceMock{}
	audienceMock.RegisterMocks()

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/updatePolicies",
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

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+$`,
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

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+$`,
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

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+$`,
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

func (m *WindowsUpdatePolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.updatePolicies = make(map[string]map[string]any)
	mockState.Unlock()
}
