package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	owners       map[string][]map[string]any // key: applicationId, value: list of owner objects
	users        map[string]map[string]any   // key: userId, value: user data
	servicePrins map[string]map[string]any   // key: spId, value: sp data
}

func init() {
	mockState.owners = make(map[string][]map[string]any)
	mockState.users = make(map[string]map[string]any)
	mockState.servicePrins = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("application_owner", &ApplicationOwnerMock{})
}

type ApplicationOwnerMock struct{}

var _ mocks.MockRegistrar = (*ApplicationOwnerMock)(nil)

func (m *ApplicationOwnerMock) RegisterMocks() {
	mockState.Lock()
	mockState.owners = make(map[string][]map[string]any)
	mockState.users = make(map[string]map[string]any)
	mockState.servicePrins = make(map[string]map[string]any)

	// Seed mock users
	mockState.users["user-11111111-1111-1111-1111-111111111111"] = map[string]any{
		"@odata.type": "#microsoft.graph.user",
		"id":          "user-11111111-1111-1111-1111-111111111111",
		"displayName": "Test User Owner",
	}
	mockState.users["user-22222222-2222-2222-2222-222222222222"] = map[string]any{
		"@odata.type": "#microsoft.graph.user",
		"id":          "user-22222222-2222-2222-2222-222222222222",
		"displayName": "Test User Owner 2",
	}

	// Seed mock service principals
	mockState.servicePrins["sp-11111111-1111-1111-1111-111111111111"] = map[string]any{
		"@odata.type": "#microsoft.graph.servicePrincipal",
		"id":          "sp-11111111-1111-1111-1111-111111111111",
		"displayName": "Test Service Principal Owner",
	}

	// Seed mock applications
	mockState.owners["11111111-1111-1111-1111-111111111111"] = []map[string]any{}
	mockState.owners["22222222-2222-2222-2222-222222222222"] = []map[string]any{}
	mockState.Unlock()

	// Get application owners - GET /applications/{applicationId}/owners
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/owners$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		applicationId := parts[3]

		mockState.Lock()
		owners, exists := mockState.owners[applicationId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Application not found"}}`), nil
		}

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
			"value":          owners,
		}

		return httpmock.NewJsonResponse(200, response)
	})

	// Add application owner - POST /applications/{applicationId}/owners/$ref
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/owners/\$ref$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		applicationId := parts[3]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		odataId, ok := requestBody["@odata.id"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"@odata.id is required"}}`), nil
		}

		// Extract owner ID from @odata.id
		var ownerId string
		var ownerData map[string]any
		if strings.Contains(odataId, "/users/") {
			ownerId = odataId[strings.LastIndex(odataId, "/")+1:]
			mockState.Lock()
			ownerData = mockState.users[ownerId]
			mockState.Unlock()
		} else if strings.Contains(odataId, "/servicePrincipals/") {
			ownerId = odataId[strings.LastIndex(odataId, "/")+1:]
			mockState.Lock()
			ownerData = mockState.servicePrins[ownerId]
			mockState.Unlock()
		}

		if ownerData == nil {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Owner object not found"}}`), nil
		}

		mockState.Lock()
		if mockState.owners[applicationId] == nil {
			mockState.owners[applicationId] = []map[string]any{}
		}
		mockState.owners[applicationId] = append(mockState.owners[applicationId], ownerData)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete application owner - DELETE /applications/{applicationId}/owners/{ownerId}/$ref
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F-]+/owners/[0-9a-fA-F-]+/\$ref$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		applicationId := parts[3]
		ownerIdToRemove := parts[5]

		mockState.Lock()
		defer mockState.Unlock()

		if owners, exists := mockState.owners[applicationId]; exists {
			for i, owner := range owners {
				if ownerIdFromMap, ok := owner["id"].(string); ok && ownerIdFromMap == ownerIdToRemove {
					mockState.owners[applicationId] = append(owners[:i], owners[i+1:]...)
					return httpmock.NewStringResponse(204, ""), nil
				}
			}
		}

		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Owner not found"}}`), nil
	})

	// Get user - GET /users/{userId}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		userId := parts[3]

		mockState.Lock()
		user, exists := mockState.users[userId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, user)
	})

	// Get service principal - GET /servicePrincipals/{spId}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		spId := parts[3]

		mockState.Lock()
		sp, exists := mockState.servicePrins[spId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, sp)
	})

	// Get application - GET /applications/{applicationId}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		applicationId := parts[3]

		mockState.Lock()
		_, exists := mockState.owners[applicationId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Application not found"}}`), nil
		}

		response := map[string]any{
			"id":          applicationId,
			"displayName": "Test Application",
		}

		return httpmock.NewJsonResponse(200, response)
	})
}

func (m *ApplicationOwnerMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/owners$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/owners/\$ref$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F-]+/owners/[0-9a-fA-F-]+/\$ref$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *ApplicationOwnerMock) CleanupMockState() {
	mockState.Lock()
	mockState.owners = make(map[string][]map[string]any)
	mockState.users = make(map[string]map[string]any)
	mockState.servicePrins = make(map[string]map[string]any)
	mockState.Unlock()
}
