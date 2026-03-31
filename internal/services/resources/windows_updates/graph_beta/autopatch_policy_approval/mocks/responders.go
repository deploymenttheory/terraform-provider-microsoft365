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
	approvals map[string]map[string]any
}

func init() {
	mockState.approvals = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_update_policy_approval", &WindowsUpdatePolicyApprovalMock{})
}

type WindowsUpdatePolicyApprovalMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdatePolicyApprovalMock)(nil)

func (m *WindowsUpdatePolicyApprovalMock) RegisterMocks() {
	mockState.Lock()
	mockState.approvals = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreatePolicyApprovalResponder()
	m.registerGetPolicyApprovalResponder()
	m.registerUpdatePolicyApprovalResponder()
	m.registerDeletePolicyApprovalResponder()
}

func (m *WindowsUpdatePolicyApprovalMock) registerCreatePolicyApprovalResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F-]+/approvals$`,
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)

			status, _ := requestBody["status"].(string)

			var responsesPath string
			if status == "suspended" {
				responsesPath = filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_policy_approval_suspended.json")
			} else {
				responsesPath = filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_policy_approval_approved.json")
			}

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load response: %s"}}`, err.Error())), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse response: %s"}}`, err.Error())), nil
			}

			// Reflect the requested catalogEntryId and status into the response
			if catalogEntryId, ok := requestBody["catalogEntryId"].(string); ok {
				responseObj["catalogEntryId"] = catalogEntryId
			}
			if status != "" {
				responseObj["status"] = status
			}

			mockState.Lock()
			if id, ok := responseObj["id"].(string); ok {
				mockState.approvals[id] = responseObj
			}
			mockState.Unlock()

			resp, err := httpmock.NewJsonResponse(200, responseObj)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdatePolicyApprovalMock) registerGetPolicyApprovalResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F-]+/approvals/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			approvalId := parts[len(parts)-1]

			mockState.Lock()
			approval, exists := mockState.approvals[approvalId]
			mockState.Unlock()

			if !exists {
				_, filename, _, _ := runtime.Caller(0)
				sourceDir := filepath.Dir(filename)
				responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationRead, "get_policy_approval_approved.json")

				jsonData, err := os.ReadFile(responsesPath)
				if err != nil {
					return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Policy approval not found"}}`), nil
				}

				var responseObj map[string]any
				if err := json.Unmarshal(jsonData, &responseObj); err != nil {
					return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse response: %s"}}`, err.Error())), nil
				}
				approval = responseObj
			}

			resp, err := httpmock.NewJsonResponse(200, approval)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdatePolicyApprovalMock) registerUpdatePolicyApprovalResponder() {
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F-]+/approvals/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			approvalId := parts[len(parts)-1]

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)
			responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationUpdate, "patch_policy_approval_suspended.json")

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load response: %s"}}`, err.Error())), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse response: %s"}}`, err.Error())), nil
			}

			responseObj["id"] = approvalId
			if catalogEntryId, ok := requestBody["catalogEntryId"].(string); ok {
				responseObj["catalogEntryId"] = catalogEntryId
			}
			if status, ok := requestBody["status"].(string); ok {
				responseObj["status"] = status
			}

			mockState.Lock()
			mockState.approvals[approvalId] = responseObj
			mockState.Unlock()

			resp, err := httpmock.NewJsonResponse(200, responseObj)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdatePolicyApprovalMock) registerDeletePolicyApprovalResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/policies/[0-9a-fA-F-]+/approvals/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			approvalId := parts[len(parts)-1]

			mockState.Lock()
			delete(mockState.approvals, approvalId)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsUpdatePolicyApprovalMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.approvals = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreatePolicyApprovalErrorResponder()
	m.registerGetPolicyApprovalErrorResponder()
}

func (m *WindowsUpdatePolicyApprovalMock) registerCreatePolicyApprovalErrorResponder() {
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

func (m *WindowsUpdatePolicyApprovalMock) registerGetPolicyApprovalErrorResponder() {
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

func (m *WindowsUpdatePolicyApprovalMock) CleanupMockState() {
	mockState.Lock()
	mockState.approvals = make(map[string]map[string]any)
	mockState.Unlock()
}
