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
	groupPolicyUploadedDefinitionFiles map[string]map[string]any
}

func init() {
	mockState.groupPolicyUploadedDefinitionFiles = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("group_policy_uploaded_definition_files", &GroupPolicyUploadedDefinitionFilesMock{})
}

type GroupPolicyUploadedDefinitionFilesMock struct{}

var _ mocks.MockRegistrar = (*GroupPolicyUploadedDefinitionFilesMock)(nil)

func (m *GroupPolicyUploadedDefinitionFilesMock) RegisterMocks() {
	mockState.Lock()
	mockState.groupPolicyUploadedDefinitionFiles = make(map[string]map[string]any)
	mockState.Unlock()

	// List all group policy uploaded definition files
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if len(mockState.groupPolicyUploadedDefinitionFiles) == 0 {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group_policy_uploaded_definition_files_list.json")
			var responseObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Return list of existing files
		list := make([]map[string]any, 0, len(mockState.groupPolicyUploadedDefinitionFiles))
		for _, v := range mockState.groupPolicyUploadedDefinitionFiles {
			c := map[string]any{}
			for k, vv := range v {
				c[k] = vv
			}
			list = append(list, c)
		}

		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyUploadedDefinitionFiles",
			"value":          list,
		})
	})

	// Get a specific group policy uploaded definition file
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		file, ok := mockState.groupPolicyUploadedDefinitionFiles[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_group_policy_uploaded_definition_files_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		// Check if expanded operations are requested
		var jsonTemplate string
		if strings.Contains(req.URL.RawQuery, "expand=groupPolicyOperations") {
			if file["status"] == "uploadFailed" {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group_policy_uploaded_definition_files_failed.json")
				jsonTemplate = jsonStr
			} else {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group_policy_uploaded_definition_files_with_operations.json")
				jsonTemplate = jsonStr
			}
		} else {
			// Get the appropriate template based on file configuration
			languageFiles, hasLanguageFiles := file["groupPolicyUploadedLanguageFiles"]
			if hasLanguageFiles && languageFiles != nil {
				if languageFilesSlice, ok := languageFiles.([]any); ok && len(languageFilesSlice) > 1 {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group_policy_uploaded_definition_files_multiple.json")
					jsonTemplate = jsonStr
				} else {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group_policy_uploaded_definition_files_single.json")
					jsonTemplate = jsonStr
				}
			} else {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group_policy_uploaded_definition_files_single.json")
				jsonTemplate = jsonStr
			}
		}

		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonTemplate), &responseObj)

		// Override template values with actual file values
		for k, v := range file {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// Create a new group policy uploaded definition file
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles", func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_group_policy_uploaded_definition_files_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		}

		id := uuid.New().String()

		// Load the success response template
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_group_policy_uploaded_definition_files_success.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Only include fields that were provided in the request
		responseObj["id"] = id
		if v, ok := body["fileName"]; ok {
			responseObj["fileName"] = v
		}
		if v, ok := body["defaultLanguageCode"]; ok {
			responseObj["defaultLanguageCode"] = v
		}

		// Content is not returned in the response
		responseObj["content"] = nil

		// Handle language files
		if v, ok := body["groupPolicyUploadedLanguageFiles"]; ok {
			responseObj["groupPolicyUploadedLanguageFiles"] = v

			// Update language codes based on language files
			if languageFiles, ok := v.([]any); ok && len(languageFiles) > 0 {
				languageCodes := []string{}
				for _, lf := range languageFiles {
					if lfMap, ok := lf.(map[string]any); ok {
						if lc, ok := lfMap["languageCode"].(string); ok {
							languageCodes = append(languageCodes, lc)
						}
					}
				}
				responseObj["languageCodes"] = languageCodes
			}
		}

		// Store in mock state
		mockState.Lock()
		mockState.groupPolicyUploadedDefinitionFiles[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Delete a group policy uploaded definition file
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		delete(mockState.groupPolicyUploadedDefinitionFiles, id)
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})

	// Remove a group policy uploaded definition file
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles/[^/]+/remove$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		mockState.Lock()
		delete(mockState.groupPolicyUploadedDefinitionFiles, id)
		mockState.Unlock()

		// Return empty success response
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *GroupPolicyUploadedDefinitionFilesMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.groupPolicyUploadedDefinitionFiles = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group_policy_uploaded_definition_files_list.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_group_policy_uploaded_definition_files_error.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_group_policy_uploaded_definition_files_not_found.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *GroupPolicyUploadedDefinitionFilesMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.groupPolicyUploadedDefinitionFiles {
		delete(mockState.groupPolicyUploadedDefinitionFiles, id)
	}
}
