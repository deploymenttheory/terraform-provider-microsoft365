package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

var mockState struct {
	sync.Mutex
	roleDefinitions map[string]map[string]any
}

const (
	roleDefinitionsURL      = "https://graph.microsoft.com/beta/roleManagement/directory/roleDefinitions"
	roleDefinitionItemURL   = `=~^https://graph\.microsoft\.com/beta/roleManagement/directory/roleDefinitions/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	resourceActionsURL      = `=~^https://graph\.microsoft\.com/beta/roleManagement/directory/resourceNamespaces/microsoft\.directory/resourceActions`
	displayNamePrefix       = "unit-test-role-definition-"
	builtInRoleDefinitionID = "62e90394-69f5-4237-9190-012177145e10"
	invalidLengthBody       = `{"error":{"code":"Request_BadRequest","message":"Invalid value specified for property '%s' of resource 'RoleDefinition'."}}`
	notFoundBody            = `{"error":{"code":"Request_ResourceNotFound","message":"Resource does not exist or one of its queried reference-property objects are not present."}}`
)

// directoryResourceActions are the microsoft.directory actions the mock namespace exposes for validation
var directoryResourceActions = []string{
	"microsoft.directory/applications/standard/read",
	"microsoft.directory/groups/standard/read",
	"microsoft.directory/devices/standard/read",
	"microsoft.directory/bitlockerKeys/metadata/read",
}

func init() {
	mockState.roleDefinitions = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, notFoundBody))
	mocks.GlobalRegistry.Register("identity_and_access_role_definition", &RoleDefinitionMock{})
}

type RoleDefinitionMock struct{}

var _ mocks.MockRegistrar = (*RoleDefinitionMock)(nil)

// getJSONFileForDisplayName determines which JSON file to load based on the role's display name,
// e.g. "unit-test-role-definition-minimal" loads post_role_definition_minimal_success.json
func getJSONFileForDisplayName(displayName string) string {
	if !strings.HasPrefix(displayName, displayNamePrefix) {
		return ""
	}
	scenario := strings.TrimPrefix(displayName, displayNamePrefix)
	return fmt.Sprintf("post_role_definition_%s_success.json", scenario)
}

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

// nullRolePermissionsResponse mirrors Graph, which rejects a role definition without permissions.
func nullRolePermissionsResponse(request map[string]any) *http.Response {
	if value, ok := request["rolePermissions"]; ok && value == nil {
		resp, _ := httpmock.NewJsonResponse(400, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": "rolePermissions cannot be null."}})
		return resp
	}
	return nil
}

func (m *RoleDefinitionMock) RegisterMocks() {
	mockState.Lock()
	mockState.roleDefinitions = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerMockResourceActions()

	// Create role definition - POST /roleManagement/directory/roleDefinitions
	httpmock.RegisterResponder("POST", roleDefinitionsURL, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}
		if resp := nullRolePermissionsResponse(requestBody); resp != nil {
			return resp, nil
		}

		displayName, ok := requestBody["displayName"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"displayName is required"}}`), nil
		}

		jsonFileName := getJSONFileForDisplayName(displayName)
		if jsonFileName == "" {
			return httpmock.NewStringResponse(400, fmt.Sprintf(`{"error":{"code":"BadRequest","message":"Unable to determine JSON file for displayName: %s"}}`, displayName)), nil
		}

		responseObj, err := loadResponse(jsonFileName)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		// A supplied templateId becomes the role ID, as in Microsoft Entra
		newId := uuid.New().String()
		if templateId, ok := requestBody["templateId"].(string); ok && templateId != "" {
			newId = templateId
		}
		responseObj["id"] = newId
		responseObj["templateId"] = newId

		mockState.Lock()
		mockState.roleDefinitions[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get role definition - GET /roleManagement/directory/roleDefinitions/{id}
	httpmock.RegisterResponder("GET", roleDefinitionItemURL, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		roleDefinitionId := parts[len(parts)-1]

		if roleDefinitionId == builtInRoleDefinitionID {
			responseObj, err := loadResponse("get_role_definition_built_in.json")
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"%s"}}`, err.Error())), nil
			}
			return httpmock.NewJsonResponse(200, responseObj)
		}

		mockState.Lock()
		roleDefinition, exists := mockState.roleDefinitions[roleDefinitionId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		return httpmock.NewJsonResponse(200, roleDefinition)
	})

	// Update role definition - PATCH /roleManagement/directory/roleDefinitions/{id}
	httpmock.RegisterResponder("PATCH", roleDefinitionItemURL, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		roleDefinitionId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}
		if resp := nullRolePermissionsResponse(requestBody); resp != nil {
			return resp, nil
		}
		if _, ok := requestBody["templateId"]; ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"templateId cannot be updated."}}`), nil
		}
		// Graph rejects an empty description or version (InvalidLength); clearing needs an explicit null
		for _, field := range []string{"description", "version"} {
			if value, ok := requestBody[field]; ok && value == "" {
				return httpmock.NewStringResponse(400, fmt.Sprintf(invalidLengthBody, field)), nil
			}
		}

		mockState.Lock()
		defer mockState.Unlock()

		roleDefinition, exists := mockState.roleDefinitions[roleDefinitionId]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		for key, value := range requestBody {
			roleDefinition[key] = value
		}
		// Properties left out of a PATCH keep their value; description and version are cleared only when sent as null
		mockState.roleDefinitions[roleDefinitionId] = roleDefinition

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete role definition - DELETE /roleManagement/directory/roleDefinitions/{id} (permanent)
	httpmock.RegisterResponder("DELETE", roleDefinitionItemURL, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		roleDefinitionId := parts[len(parts)-1]

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.roleDefinitions[roleDefinitionId]; exists {
			delete(mockState.roleDefinitions, roleDefinitionId)
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewStringResponse(404, notFoundBody), nil
	})
}

func (m *RoleDefinitionMock) RegisterErrorMocks() {
	m.registerMockResourceActions()

	// Error scenarios for testing
	httpmock.RegisterResponder("POST", roleDefinitionsURL, httpmock.NewStringResponder(400, `{"error":{"code":"Request_BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", roleDefinitionItemURL, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterResponder("PATCH", roleDefinitionItemURL, httpmock.NewStringResponder(400, `{"error":{"code":"Request_BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", roleDefinitionItemURL, httpmock.NewStringResponder(400, `{"error":{"code":"Request_BadRequest","message":"Invalid request"}}`))
}

func (m *RoleDefinitionMock) CleanupMockState() {
	mockState.Lock()
	mockState.roleDefinitions = make(map[string]map[string]any)
	mockState.Unlock()
}

// registerMockResourceActions registers the microsoft.directory resource action list used by request validation.
// The first page links to a second page so pagination is exercised.
func (m *RoleDefinitionMock) registerMockResourceActions() {
	httpmock.RegisterResponder("GET", resourceActionsURL, func(req *http.Request) (*http.Response, error) {
		page := directoryResourceActions[:2]
		response := map[string]any{
			"@odata.context":  "https://graph.microsoft.com/beta/$metadata#roleManagement/directory/resourceNamespaces('microsoft.directory')/resourceActions",
			"@odata.nextLink": "https://graph.microsoft.com/beta/roleManagement/directory/resourceNamespaces/microsoft.directory/resourceActions?$skiptoken=page2",
		}
		if strings.Contains(req.URL.RawQuery, "skiptoken=page2") {
			page = directoryResourceActions[2:]
			delete(response, "@odata.nextLink")
		}

		values := make([]map[string]any, 0, len(page))
		for _, name := range page {
			values = append(values, map[string]any{"id": strings.ReplaceAll(name, "/", "-"), "name": name, "isPrivileged": false})
		}
		response["value"] = values

		return httpmock.NewJsonResponse(200, response)
	})
}
