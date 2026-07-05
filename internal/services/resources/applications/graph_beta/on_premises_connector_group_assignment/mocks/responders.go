package mocks

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

const (
	unitTestApplicationID      = "11111111-1111-1111-1111-111111111111"
	unitTestConnectorGroupID   = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	unitTestConnectorGroupID2  = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	connectorGroupODataPath    = "/onPremisesPublishingProfiles/applicationProxy/connectorGroups/"
	graphBetaConnectorGroupURL = "https://graph.microsoft.com/beta" + connectorGroupODataPath
)

var mockState struct {
	sync.Mutex
	assignments     map[string]string
	connectorGroups map[string]map[string]any
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("on_premises_connector_group_assignment", &OnPremisesConnectorGroupAssignmentMock{})
}

type OnPremisesConnectorGroupAssignmentMock struct{}

var _ mocks.MockRegistrar = (*OnPremisesConnectorGroupAssignmentMock)(nil)

func (m *OnPremisesConnectorGroupAssignmentMock) RegisterMocks() {
	resetMockState()

	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F-]+/connectorGroup/\$ref$`, func(req *http.Request) (*http.Response, error) {
		applicationID := applicationIDFromPath(req.URL.Path)
		if applicationID == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Application ID is required."}}`), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		odataID, ok := requestBody["@odata.id"].(string)
		if !ok || odataID == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"@odata.id is required."}}`), nil
		}

		connectorGroupID := connectorGroupIDFromODataID(odataID)
		if connectorGroupID == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"Request_BadRequest","message":"Invalid connector group reference."}}`), nil
		}

		mockState.Lock()
		_, connectorGroupExists := mockState.connectorGroups[connectorGroupID]
		if connectorGroupExists {
			mockState.assignments[applicationID] = connectorGroupID
		}
		mockState.Unlock()

		if !connectorGroupExists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Connector group not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F-]+/connectorGroup$`, func(req *http.Request) (*http.Response, error) {
		applicationID := applicationIDFromPath(req.URL.Path)

		mockState.Lock()
		connectorGroupID, assigned := mockState.assignments[applicationID]
		connectorGroup, exists := mockState.connectorGroups[connectorGroupID]
		mockState.Unlock()

		if !assigned || !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Connector group assignment not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, connectorGroup)
	})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F-]+/connectorGroup/\$ref$`, func(req *http.Request) (*http.Response, error) {
		applicationID := applicationIDFromPath(req.URL.Path)

		mockState.Lock()
		delete(mockState.assignments, applicationID)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *OnPremisesConnectorGroupAssignmentMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F-]+/connectorGroup/\$ref$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F-]+/connectorGroup$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Connector group assignment not found"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F-]+/connectorGroup/\$ref$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *OnPremisesConnectorGroupAssignmentMock) CleanupMockState() {
	resetMockState()
}

func (m *OnPremisesConnectorGroupAssignmentMock) AssignConnectorGroup(applicationID, connectorGroupID string) {
	mockState.Lock()
	defer mockState.Unlock()

	mockState.assignments[applicationID] = connectorGroupID
}

func (m *OnPremisesConnectorGroupAssignmentMock) AssignedConnectorGroup(applicationID string) (string, bool) {
	mockState.Lock()
	defer mockState.Unlock()

	connectorGroupID, ok := mockState.assignments[applicationID]
	return connectorGroupID, ok
}

func resetMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	mockState.assignments = make(map[string]string)
	mockState.connectorGroups = map[string]map[string]any{
		unitTestConnectorGroupID: {
			"@odata.context":      "https://graph.microsoft.com/beta/$metadata#onPremisesPublishingProfiles/applicationProxy/connectorGroups/$entity",
			"id":                  unitTestConnectorGroupID,
			"name":                "Unit Test Connector Group",
			"connectorGroupType":  "applicationProxy",
			"isDefault":           false,
			"region":              "nam",
			"members@odata.count": 0,
		},
		unitTestConnectorGroupID2: {
			"@odata.context":      "https://graph.microsoft.com/beta/$metadata#onPremisesPublishingProfiles/applicationProxy/connectorGroups/$entity",
			"id":                  unitTestConnectorGroupID2,
			"name":                "Unit Test Connector Group 2",
			"connectorGroupType":  "applicationProxy",
			"isDefault":           false,
			"region":              "eur",
			"members@odata.count": 0,
		},
	}
}

func applicationIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	for index, part := range parts {
		if part == "applications" && index+1 < len(parts) {
			return parts[index+1]
		}
	}
	return ""
}

func connectorGroupIDFromODataID(odataID string) string {
	escapedPath := regexp.QuoteMeta(connectorGroupODataPath)
	pattern := regexp.MustCompile(`(?i)` + escapedPath + `([0-9a-f-]{36})$`)
	matches := pattern.FindStringSubmatch(odataID)
	if len(matches) != 2 {
		return ""
	}
	return strings.ToLower(matches[1])
}
