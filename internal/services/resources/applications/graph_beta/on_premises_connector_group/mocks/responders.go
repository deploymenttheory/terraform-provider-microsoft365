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
	connectorGroups map[string]map[string]any
}

func init() {
	mockState.connectorGroups = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("on_premises_connector_group", &OnPremisesConnectorGroupMock{})
}

type OnPremisesConnectorGroupMock struct{}

var _ mocks.MockRegistrar = (*OnPremisesConnectorGroupMock)(nil)

func (m *OnPremisesConnectorGroupMock) RegisterMocks() {
	mockState.Lock()
	mockState.connectorGroups = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups$`, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}
		if _, ok := requestBody["name"]; !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Property name is required."}}`), nil
		}
		if _, ok := requestBody["connectorGroupType"]; ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Property connectorGroupType is read-only."}}`), nil
		}
		if _, ok := requestBody["isDefault"]; ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Property isDefault is read-only."}}`), nil
		}

		newID := uuid.New().String()

		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_connector_group_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load response"}}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse response"}}`), nil
		}

		responseObj["id"] = newID
		responseObj["name"] = requestBody["name"]
		if region, ok := requestBody["region"]; ok {
			responseObj["region"] = region
		}

		mockState.Lock()
		mockState.connectorGroups[newID] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		connectorGroupID := parts[len(parts)-1]

		mockState.Lock()
		connectorGroup, exists := mockState.connectorGroups[connectorGroupID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, connectorGroup)
	})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		connectorGroupID := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}
		if _, ok := requestBody["connectorGroupType"]; ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Property connectorGroupType is read-only."}}`), nil
		}
		if _, ok := requestBody["isDefault"]; ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Property isDefault is read-only."}}`), nil
		}

		mockState.Lock()
		connectorGroup, exists := mockState.connectorGroups[connectorGroupID]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}
		for key, value := range requestBody {
			connectorGroup[key] = value
		}
		mockState.connectorGroups[connectorGroupID] = connectorGroup
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, connectorGroup)
	})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		connectorGroupID := parts[len(parts)-1]

		mockState.Lock()
		connectorGroup, exists := mockState.connectorGroups[connectorGroupID]
		if exists && connectorGroup["isDefault"] == true {
			mockState.Unlock()
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Default connector group cannot be deleted."}}`), nil
		}
		if exists {
			delete(mockState.connectorGroups, connectorGroupID)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *OnPremisesConnectorGroupMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/[0-9a-fA-F-]+$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/[0-9a-fA-F-]+$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups/[0-9a-fA-F-]+$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *OnPremisesConnectorGroupMock) CleanupMockState() {
	mockState.Lock()
	mockState.connectorGroups = make(map[string]map[string]any)
	mockState.Unlock()
}
