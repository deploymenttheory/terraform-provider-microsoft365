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
	contentApprovals map[string]map[string]any
}

func init() {
	mockState.contentApprovals = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_update_content_approval", &WindowsUpdateContentApprovalMock{})
}

type WindowsUpdateContentApprovalMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdateContentApprovalMock)(nil)

func (m *WindowsUpdateContentApprovalMock) RegisterMocks() {
	mockState.Lock()
	mockState.contentApprovals = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreateContentApprovalResponder()
	m.registerGetContentApprovalResponder()
	m.registerUpdateContentApprovalResponder()
	m.registerDeleteContentApprovalResponder()
}

func (m *WindowsUpdateContentApprovalMock) registerCreateContentApprovalResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+/complianceChanges$`, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		_, filename, _, _ := runtime.Caller(0)
		sourceDir := filepath.Dir(filename)

		var responsesPath string

		content, _ := requestBody["content"].(map[string]any)
		catalogEntry, _ := content["catalogEntry"].(map[string]any)
		catalogEntryId, _ := catalogEntry["id"].(string)
		odataType, _ := catalogEntry["@odata.type"].(string)

		isRevoked, _ := requestBody["isRevoked"].(bool)
		deploymentSettings, hasDeploymentSettings := requestBody["deploymentSettings"]

		if isRevoked {
			responsesPath = filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_content_approval_revoked.json")
		} else if !hasDeploymentSettings || deploymentSettings == nil {
			responsesPath = filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_content_approval_minimal.json")
		} else if strings.Contains(odataType, "qualityUpdate") {
			responsesPath = filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_content_approval_quality_update.json")
		} else {
			responsesPath = filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_content_approval_success.json")
		}

		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		if content, ok := responseObj["content"].(map[string]any); ok {
			if catalogEntry, ok := content["catalogEntry"].(map[string]any); ok {
				catalogEntry["id"] = catalogEntryId
			}
		}

		mockState.Lock()
		if id, ok := responseObj["id"].(string); ok {
			mockState.contentApprovals[id] = responseObj
		}
		mockState.Unlock()

		resp, err := httpmock.NewJsonResponse(201, responseObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *WindowsUpdateContentApprovalMock) registerGetContentApprovalResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+/complianceChanges/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		complianceChangeId := parts[len(parts)-1]

		mockState.Lock()
		approval, exists := mockState.contentApprovals[complianceChangeId]
		mockState.Unlock()

		if !exists {
			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)
			responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationRead, "get_content_approval_success.json")

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Content approval not found"}}`), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
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

func (m *WindowsUpdateContentApprovalMock) registerUpdateContentApprovalResponder() {
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+/complianceChanges/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		complianceChangeId := parts[len(parts)-1]

		_, filename, _, _ := runtime.Caller(0)
		sourceDir := filepath.Dir(filename)
		responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationUpdate, "patch_content_approval_success.json")

		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		mockState.Lock()
		mockState.contentApprovals[complianceChangeId] = responseObj
		mockState.Unlock()

		resp, err := httpmock.NewJsonResponse(200, responseObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *WindowsUpdateContentApprovalMock) registerDeleteContentApprovalResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies/[0-9a-fA-F-]+/complianceChanges/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		complianceChangeId := parts[len(parts)-1]

		mockState.Lock()
		delete(mockState.contentApprovals, complianceChangeId)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *WindowsUpdateContentApprovalMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.contentApprovals = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreateContentApprovalErrorResponder()
	m.registerGetContentApprovalErrorResponder()
	m.registerUpdateContentApprovalErrorResponder()
	m.registerDeleteContentApprovalErrorResponder()
}

func (m *WindowsUpdateContentApprovalMock) registerCreateContentApprovalErrorResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		resp, err := httpmock.NewJsonResponse(403, errorObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *WindowsUpdateContentApprovalMock) registerGetContentApprovalErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		resp, err := httpmock.NewJsonResponse(403, errorObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *WindowsUpdateContentApprovalMock) registerUpdateContentApprovalErrorResponder() {
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		resp, err := httpmock.NewJsonResponse(403, errorObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *WindowsUpdateContentApprovalMock) registerDeleteContentApprovalErrorResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatePolicies`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		resp, err := httpmock.NewJsonResponse(403, errorObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *WindowsUpdateContentApprovalMock) CleanupMockState() {
	mockState.Lock()
	mockState.contentApprovals = make(map[string]map[string]any)
	mockState.Unlock()
}
