package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	autopatchGroups map[string]map[string]any
}

func init() {
	mockState.autopatchGroups = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
	})
	mocks.GlobalRegistry.Register("autopatch_groups", &AutopatchGroupsMock{})
}

type AutopatchGroupsMock struct{}

var _ mocks.MockRegistrar = (*AutopatchGroupsMock)(nil)

// getJSONFileForName determines which JSON file to load based on the autopatch group name
func getJSONFileForName(name string) string {
	// Map test names to JSON response files
	switch name {
	case "test", "unit-test":
		return "post_autopatch_groups_test_success.json"
	case "unit-test-autopatch-group":
		return "post_autopatch_groups_unittest_success.json"
	case "auto-patch-group":
		return "post_autopatch_groups_autopatchgroup_success.json"
	default:
		// Use test success file as generic fallback
		return "post_autopatch_groups_test_success.json"
	}
}

func (m *AutopatchGroupsMock) RegisterMocks() {
	mockState.Lock()
	mockState.autopatchGroups = make(map[string]map[string]any)
	mockState.Unlock()

	// Create autopatch group - POST /device/v2/autopatchGroups
	httpmock.RegisterResponder("POST", "https://services.autopatch.microsoft.com/device/v2/autopatchGroups", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Get the name from request
		name, ok := requestBody["name"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"name is required"}}`), nil
		}

		jsonFileName := getJSONFileForName(name)

		// Load JSON response from file
		responsesPath := filepath.Join("tests", "responses", constants.TfOperationCreate, jsonFileName)
		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		// Parse the JSON response
		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Generate a UUID for the new resource
		newID := uuid.New().String()
		responseObj["id"] = newID

		// Store in mock state
		mockState.Lock()
		mockState.autopatchGroups[newID] = responseObj
		mockState.Unlock()

		// Return the response
		respJSON, _ := json.Marshal(responseObj)
		return httpmock.NewBytesResponse(201, respJSON), nil
	})

	// Read autopatch group - GET /device/v2/autopatchGroups/{id}
	httpmock.RegisterResponder("GET", `=~^https://services\.autopatch\.microsoft\.com/device/v2/autopatchGroups/([0-9a-fA-F-]+)$`, func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		parts := req.URL.Path[len("/device/v2/autopatchGroups/"):]
		id := parts

		mockState.Lock()
		group, exists := mockState.autopatchGroups[id]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Autopatch group not found"}}`), nil
		}

		respJSON, _ := json.Marshal(group)
		return httpmock.NewBytesResponse(200, respJSON), nil
	})

	// Update autopatch group - PUT /device/v2/autopatchGroups
	httpmock.RegisterResponder("PUT", "https://services.autopatch.microsoft.com/device/v2/autopatchGroups", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id, ok := requestBody["id"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"id is required for update"}}`), nil
		}

		mockState.Lock()
		group, exists := mockState.autopatchGroups[id]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Autopatch group not found"}}`), nil
		}

		// Merge request body into existing group
		for key, value := range requestBody {
			group[key] = value
		}

		mockState.autopatchGroups[id] = group
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete autopatch group - DELETE /device/v2/autopatchGroups/{id}
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://services\.autopatch\.microsoft\.com/device/v2/autopatchGroups/([0-9a-fA-F-]+)$`, func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		parts := req.URL.Path[len("/device/v2/autopatchGroups/"):]
		id := parts

		mockState.Lock()
		_, exists := mockState.autopatchGroups[id]
		if exists {
			delete(mockState.autopatchGroups, id)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Autopatch group not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AutopatchGroupsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.autopatchGroups = make(map[string]map[string]any)
	mockState.Unlock()

	// Create - return error
	httpmock.RegisterResponder("POST", "https://services.autopatch.microsoft.com/device/v2/autopatchGroups",
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))

	// Read - return error
	httpmock.RegisterResponder("GET", `=~^https://services\.autopatch\.microsoft\.com/device/v2/autopatchGroups/([0-9a-fA-F-]+)$`,
		httpmock.NewStringResponder(404, `{"error":{"code":"NotFound","message":"Autopatch group not found"}}`))

	// Update - return error
	httpmock.RegisterResponder("PUT", "https://services.autopatch.microsoft.com/device/v2/autopatchGroups",
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))

	// Delete - return error
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://services\.autopatch\.microsoft\.com/device/v2/autopatchGroups/([0-9a-fA-F-]+)$`,
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Cannot delete autopatch group"}}`))
}

func (m *AutopatchGroupsMock) CleanupMockState() {
	mockState.Lock()
	mockState.autopatchGroups = make(map[string]map[string]any)
	mockState.Unlock()
}
