package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	appRoleAssignments map[string]map[string]any
}

func init() {
	mockState.appRoleAssignments = make(map[string]map[string]any)
	mocks.GlobalRegistry.Register("service_principal_app_role_assigned_to", &ServicePrincipalAppRoleAssignedToMock{})
}

type ServicePrincipalAppRoleAssignedToMock struct{}

var _ mocks.MockRegistrar = (*ServicePrincipalAppRoleAssignedToMock)(nil)

func (m *ServicePrincipalAppRoleAssignedToMock) RegisterMocks() {
	mockState.Lock()
	mockState.appRoleAssignments = make(map[string]map[string]any)
	mockState.Unlock()

	// Create app role assignment - POST /servicePrincipals/{id}/appRoleAssignedTo
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/appRoleAssignedTo$`, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Generate a UUID for the new resource
		newId := uuid.New().String()

		// Load the template response
		responseObj, err := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_app_role_assignment_success.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load response"}}`), nil
		}

		// Update response with request data and generated ID
		responseObj["id"] = newId
		if principalId, ok := requestBody["principalId"]; ok {
			responseObj["principalId"] = principalId
		}
		if resourceId, ok := requestBody["resourceId"]; ok {
			responseObj["resourceId"] = resourceId
		}
		if appRoleId, ok := requestBody["appRoleId"]; ok {
			responseObj["appRoleId"] = appRoleId
		}
		responseObj["creationTimestamp"] = time.Now().UTC().Format(time.RFC3339)

		// Store in mock state
		mockState.Lock()
		mockState.appRoleAssignments[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get app role assignment - GET /servicePrincipals/{id}/appRoleAssignedTo/{appRoleAssignmentId}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/appRoleAssignedTo/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		assignmentId := parts[len(parts)-1]

		mockState.Lock()
		assignment, exists := mockState.appRoleAssignments[assignmentId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, assignment)
	})

	// Delete app role assignment - DELETE /servicePrincipals/{id}/appRoleAssignedTo/{appRoleAssignmentId}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/appRoleAssignedTo/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		assignmentId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.appRoleAssignments[assignmentId]
		if exists {
			delete(mockState.appRoleAssignments, assignmentId)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *ServicePrincipalAppRoleAssignedToMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/appRoleAssignedTo$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/appRoleAssignedTo/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/appRoleAssignedTo/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *ServicePrincipalAppRoleAssignedToMock) CleanupMockState() {
	mockState.Lock()
	mockState.appRoleAssignments = make(map[string]map[string]any)
	mockState.Unlock()
}
