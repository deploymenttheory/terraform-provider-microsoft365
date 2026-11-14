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

const (
	collectionURL = "https://graph.microsoft.com/beta/identityGovernance/entitlementManagement/connectedOrganizations"
	itemURLRegex  = `=~^https://graph\.microsoft\.com/beta/identityGovernance/entitlementManagement/connectedOrganizations/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`

	// unitTestDisplayNamePrefix prefixes every unit test display name; the remainder selects the response fixture.
	unitTestDisplayNamePrefix = "unit-test-connected-organization-"

	notFoundBody   = `{"error":{"code":"ConnectedOrganizationNotFound","message":"The connected organization was not found."}}`
	forbiddenBody  = `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`
	badRequestBody = `{"error":{"code":"BadRequest","message":"Invalid request"}}`
)

var mockState struct {
	sync.Mutex
	connectedOrganizations map[string]map[string]any
}

func init() {
	mockState.connectedOrganizations = make(map[string]map[string]any)
	mocks.GlobalRegistry.Register("identity_and_access_connected_organization", &ConnectedOrganizationMock{})
}

// ConnectedOrganizationMock provides stateful httpmock responders for connected organizations
type ConnectedOrganizationMock struct{}

var _ mocks.MockRegistrar = (*ConnectedOrganizationMock)(nil)

// getJSONFileForDisplayName determines which JSON file to load based on the connected organization's display name.
// "unit-test-connected-organization-azure-ad-tenant" loads "post_connected_organization_azure_ad_tenant_success.json".
func getJSONFileForDisplayName(displayName string) string {
	if !strings.HasPrefix(displayName, unitTestDisplayNamePrefix) {
		return ""
	}

	scenario := strings.ReplaceAll(strings.TrimPrefix(displayName, unitTestDisplayNamePrefix), "-", "_")
	if scenario == "" {
		return ""
	}

	return fmt.Sprintf("post_connected_organization_%s_success.json", scenario)
}

// badRequest builds a Graph-style 400 response.
func badRequest(code, message string) *http.Response {
	resp, _ := httpmock.NewJsonResponse(400, map[string]any{"error": map[string]any{"code": code, "message": message}})
	return resp
}

// nullFieldResponse mirrors Graph, which rejects identity sources sent as null instead of being omitted.
func nullFieldResponse(request map[string]any) *http.Response {
	if value, ok := request["identitySources"]; ok && value == nil {
		return badRequest("BadRequest", "The server could not process the request because it is malformed or incorrect.")
	}
	return nil
}

func (m *ConnectedOrganizationMock) RegisterMocks() {
	mockState.Lock()
	mockState.connectedOrganizations = make(map[string]map[string]any)
	mockState.Unlock()

	// Create connected organization - POST /identityGovernance/entitlementManagement/connectedOrganizations
	httpmock.RegisterResponder("POST", collectionURL, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return badRequest("BadRequest", "Invalid request body"), nil
		}
		if resp := nullFieldResponse(requestBody); resp != nil {
			return resp, nil
		}

		sources, _ := requestBody["identitySources"].([]any)
		if len(sources) == 0 {
			return badRequest("MissingIdentitySource", "Identity Source info is missing."), nil
		}

		displayName, ok := requestBody["displayName"].(string)
		if !ok {
			return badRequest("BadRequest", "displayName is required"), nil
		}

		jsonFileName := getJSONFileForDisplayName(displayName)
		if jsonFileName == "" {
			return badRequest("BadRequest", fmt.Sprintf("Unable to determine JSON file for displayName: %s", displayName)), nil
		}

		responsesPath := filepath.Join("tests", "responses", "validate_create", jsonFileName)
		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		newId := uuid.New().String()
		responseObj["id"] = newId

		mockState.Lock()
		mockState.connectedOrganizations[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get connected organization - GET /identityGovernance/entitlementManagement/connectedOrganizations/{id}
	httpmock.RegisterResponder("GET", itemURLRegex, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		mockState.Lock()
		connectedOrganization, exists := mockState.connectedOrganizations[id]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		return httpmock.NewJsonResponse(200, connectedOrganization)
	})

	// Update connected organization - PATCH /identityGovernance/entitlementManagement/connectedOrganizations/{id}
	httpmock.RegisterResponder("PATCH", itemURLRegex, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return badRequest("BadRequest", "Invalid request body"), nil
		}
		if resp := nullFieldResponse(requestBody); resp != nil {
			return resp, nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		connectedOrganization, exists := mockState.connectedOrganizations[id]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		for key, value := range requestBody {
			connectedOrganization[key] = value
		}
		connectedOrganization["modifiedDateTime"] = "2026-01-02T00:00:00Z"
		mockState.connectedOrganizations[id] = connectedOrganization

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete connected organization - DELETE /identityGovernance/entitlementManagement/connectedOrganizations/{id}
	// Graph deletes connected organizations permanently.
	httpmock.RegisterResponder("DELETE", itemURLRegex, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.connectedOrganizations[id]; exists {
			delete(mockState.connectedOrganizations, id)
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewStringResponse(404, notFoundBody), nil
	})
}

func (m *ConnectedOrganizationMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", collectionURL, httpmock.NewStringResponder(403, forbiddenBody))
	httpmock.RegisterResponder("GET", itemURLRegex, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterResponder("PATCH", itemURLRegex, httpmock.NewStringResponder(400, badRequestBody))
	httpmock.RegisterResponder("DELETE", itemURLRegex, httpmock.NewStringResponder(400, badRequestBody))
}

func (m *ConnectedOrganizationMock) CleanupMockState() {
	mockState.Lock()
	mockState.connectedOrganizations = make(map[string]map[string]any)
	mockState.Unlock()
}
