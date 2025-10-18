package mocks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	roleDefinitions map[string]map[string]any
}

func init() {
	mockState.roleDefinitions = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("role_definitions", &RoleDefinitionsMock{})
}

type RoleDefinitionsMock struct{}

var _ mocks.MockRegistrar = (*RoleDefinitionsMock)(nil)

func (m *RoleDefinitionsMock) RegisterMocks() {
	mockState.Lock()
	mockState.roleDefinitions = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Get all role definitions - GET /roleManagement/directory/roleDefinitions
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/roleManagement/directory/roleDefinitions", func(req *http.Request) (*http.Response, error) {
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)

		// Handle different scenarios based on query parameters
		if filter := queryParams.Get("$filter"); filter != "" {
			if strings.Contains(filter, "isPrivileged eq true") {
				// Return filtered privileged roles
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_role_definitions_odata_filter.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			} else if strings.Contains(filter, "isBuiltIn eq true") {
				// Return built-in roles
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_role_definitions_odata_filter.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}
		}

		// Default: return all role definitions
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_role_definitions_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Get role definition by ID - GET /roleManagement/directory/roleDefinitions/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/roleManagement/directory/roleDefinitions/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		roleDefId := parts[len(parts)-1]

		// Return mock response for known IDs
		switch roleDefId {
		case "62e90394-69f5-4237-9190-012177145e10": // Global Administrator
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_role_definition_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		case "e8611ab8-c189-46e8-94e1-60213ab1f814": // Privileged Role Administrator
			responseObj := map[string]any{
				"id":             "e8611ab8-c189-46e8-94e1-60213ab1f814",
				"description":    "Can manage all aspects of privileged roles in Azure AD and Privileged Identity Management.",
				"displayName":    "Privileged Role Administrator",
				"isBuiltIn":      true,
				"isEnabled":      true,
				"isPrivileged":   true,
				"resourceScopes": []string{"/"},
				"templateId":     "e8611ab8-c189-46e8-94e1-60213ab1f814",
				"version":        "1",
				"rolePermissions": []map[string]any{
					{
						"allowedResourceActions": []string{
							"microsoft.directory/roleAssignments/allProperties/allTasks",
							"microsoft.directory/roleDefinitions/allProperties/allTasks",
							"microsoft.directory/scopedRoleMemberships/allProperties/allTasks",
						},
						"condition": nil,
					},
				},
			}
			return httpmock.NewJsonResponse(200, responseObj)
		case "194ae4cb-b126-40b2-bd5b-6091b380977d": // Security Administrator
			responseObj := map[string]any{
				"id":             "194ae4cb-b126-40b2-bd5b-6091b380977d",
				"description":    "Can read security information and reports, and manage configuration in Azure AD and Office 365.",
				"displayName":    "Security Administrator",
				"isBuiltIn":      true,
				"isEnabled":      true,
				"isPrivileged":   true,
				"resourceScopes": []string{"/"},
				"templateId":     "194ae4cb-b126-40b2-bd5b-6091b380977d",
				"version":        "1",
				"rolePermissions": []map[string]any{
					{
						"allowedResourceActions": []string{
							"microsoft.directory/identityProtection/allProperties/allTasks",
							"microsoft.directory/privilegedIdentityManagement/allProperties/read",
						},
						"condition": nil,
					},
				},
			}
			return httpmock.NewJsonResponse(200, responseObj)
		case "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3": // Application Administrator
			responseObj := map[string]any{
				"id":             "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3",
				"description":    "Can create and manage all aspects of app registrations and enterprise apps.",
				"displayName":    "Application Administrator",
				"isBuiltIn":      true,
				"isEnabled":      true,
				"isPrivileged":   false,
				"resourceScopes": []string{"/"},
				"templateId":     "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3",
				"version":        "1",
				"rolePermissions": []map[string]any{
					{
						"allowedResourceActions": []string{
							"microsoft.directory/applications/allProperties/allTasks",
							"microsoft.directory/servicePrincipals/allProperties/allTasks",
						},
						"condition": nil,
					},
				},
			}
			return httpmock.NewJsonResponse(200, responseObj)
		case "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9": // Conditional Access Administrator
			responseObj := map[string]any{
				"id":             "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9",
				"description":    "Can manage all aspects of Conditional Access.",
				"displayName":    "Conditional Access Administrator",
				"isBuiltIn":      true,
				"isEnabled":      true,
				"isPrivileged":   true,
				"resourceScopes": []string{"/"},
				"templateId":     "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9",
				"version":        "1",
				"rolePermissions": []map[string]any{
					{
						"allowedResourceActions": []string{
							"microsoft.directory/conditionalAccessPolicies/create",
							"microsoft.directory/conditionalAccessPolicies/delete",
							"microsoft.directory/conditionalAccessPolicies/standard/read",
						},
						"condition": nil,
					},
				},
			}
			return httpmock.NewJsonResponse(200, responseObj)
		default:
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Role definition not found"}}`), nil
		}
	})

	// 3. Handle OData queries with pagination simulation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/roleManagement/directory/roleDefinitions\?.*`, func(req *http.Request) (*http.Response, error) {
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)

		// Handle $count parameter
		if queryParams.Get("$count") == "true" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_role_definitions_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			responseObj["@odata.count"] = 2
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Handle $orderby parameter
		if orderBy := queryParams.Get("$orderby"); orderBy != "" && strings.Contains(orderBy, "displayName") {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_role_definitions_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Handle $select parameter
		if selectFields := queryParams.Get("$select"); selectFields != "" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_role_definitions_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Default OData response
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_role_definitions_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *RoleDefinitionsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.roleDefinitions = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/roleManagement/directory/roleDefinitions", func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/roleManagement/directory/roleDefinitions/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "NotFound",
				"message": "Role definition not found",
			},
		}
		return httpmock.NewJsonResponse(404, errorObj)
	})
}

func (m *RoleDefinitionsMock) CleanupMockState() {
	mockState.Lock()
	mockState.roleDefinitions = make(map[string]map[string]any)
	mockState.Unlock()
}
