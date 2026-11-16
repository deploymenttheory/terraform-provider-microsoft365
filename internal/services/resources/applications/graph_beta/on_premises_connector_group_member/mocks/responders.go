package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	// UnitTestConnectorID is the connector the JSON responses describe.
	UnitTestConnectorID = "11111111-1111-1111-1111-111111111111"
	// UnitTestDefaultConnectorGroupID is the tenant default connector group.
	UnitTestDefaultConnectorGroupID = "dddddddd-dddd-dddd-dddd-dddddddddddd"
	// UnitTestConnectorGroupAID is the group the minimal config assigns.
	UnitTestConnectorGroupAID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	// UnitTestConnectorGroupBID is the group the updated config assigns.
	UnitTestConnectorGroupBID = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"

	connectorsURL      = `^https://graph\.microsoft\.com/beta/onPremisesPublishingProfiles/applicationProxy/connectors/([0-9a-fA-F-]+)`
	connectorGroupsURL = "https://graph.microsoft.com/beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups"
	connectorGroupPath = "/onPremisesPublishingProfiles/applicationProxy/connectorGroups/"
)

var mockState struct {
	sync.Mutex
	connectors      map[string]map[string]any
	connectorGroups map[string]map[string]any
	memberships     map[string]string
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("on_premises_connector_group_member", &OnPremisesConnectorGroupMemberMock{})
}

type OnPremisesConnectorGroupMemberMock struct{}

var _ mocks.MockRegistrar = (*OnPremisesConnectorGroupMemberMock)(nil)

// RegisterMocks registers stateful responders for the connector and connector group endpoints.
func (m *OnPremisesConnectorGroupMemberMock) RegisterMocks() {
	resetMockState()

	// Get connector - GET /onPremisesPublishingProfiles/applicationProxy/connectors/{id}
	httpmock.RegisterResponder("GET", "=~"+connectorsURL+`$`, func(req *http.Request) (*http.Response, error) {
		connectorID := connectorIDFromPath(req.URL.Path)

		mockState.Lock()
		connector, exists := mockState.connectors[connectorID]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Connector not found"}}`), nil
		}

		return httpmock.NewJsonResponse(200, connector)
	})

	// List connector memberships - GET /onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf
	httpmock.RegisterResponder("GET", "=~"+connectorsURL+`/memberOf$`, func(req *http.Request) (*http.Response, error) {
		connectorID := connectorIDFromPath(req.URL.Path)

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.connectors[connectorID]; !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Connector not found"}}`), nil
		}

		value := []any{}
		if group, ok := mockState.connectorGroups[mockState.memberships[connectorID]]; ok {
			value = append(value, group)
		}

		return httpmock.NewJsonResponse(200, map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#onPremisesPublishingProfiles('applicationProxy')/connectors('" + connectorID + "')/memberOf",
			"value":          value,
		})
	})

	// Assign connector to group - POST /onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf/$ref
	httpmock.RegisterResponder("POST", "=~"+connectorsURL+`/memberOf/\$ref$`, func(req *http.Request) (*http.Response, error) {
		connectorID := connectorIDFromPath(req.URL.Path)

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		odataID, ok := requestBody["@odata.id"].(string)
		if !ok || odataID == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"@odata.id is required"}}`), nil
		}

		connectorGroupID := connectorGroupIDFromODataID(odataID)
		if connectorGroupID == "" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid connector group reference"}}`), nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.connectors[connectorID]; !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Connector not found"}}`), nil
		}
		if _, exists := mockState.connectorGroups[connectorGroupID]; !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Connector group not found"}}`), nil
		}

		mockState.memberships[connectorID] = connectorGroupID

		return httpmock.NewStringResponse(204, ""), nil
	})

	// List connector groups - GET /onPremisesPublishingProfiles/applicationProxy/connectorGroups
	httpmock.RegisterResponder("GET", connectorGroupsURL, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		value := []any{}
		for _, id := range []string{UnitTestDefaultConnectorGroupID, UnitTestConnectorGroupAID, UnitTestConnectorGroupBID} {
			if group, ok := mockState.connectorGroups[id]; ok {
				value = append(value, group)
			}
		}

		return httpmock.NewJsonResponse(200, map[string]any{"value": value})
	})
}

// RegisterErrorMocks registers responders that reject the membership write.
func (m *OnPremisesConnectorGroupMemberMock) RegisterErrorMocks() {
	resetMockState()

	errorBody, err := os.ReadFile(filepath.Join("tests", "responses", "validate_create", "post_connector_member_of_error.json"))
	if err != nil {
		errorBody = []byte(`{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`)
	}

	httpmock.RegisterResponder("POST", "=~"+connectorsURL+`/memberOf/\$ref$`, httpmock.NewBytesResponder(403, errorBody))
	notFound := `{"error":{"code":"Request_ResourceNotFound","message":"Connector not found"}}`
	httpmock.RegisterResponder("GET", "=~"+connectorsURL+`$`, httpmock.NewStringResponder(404, notFound))
	httpmock.RegisterResponder("GET", "=~"+connectorsURL+`/memberOf$`, httpmock.NewStringResponder(404, notFound))
}

// CleanupMockState resets the in-memory tenant.
func (m *OnPremisesConnectorGroupMemberMock) CleanupMockState() {
	resetMockState()
}

// MoveConnector simulates a change made outside Terraform.
func (m *OnPremisesConnectorGroupMemberMock) MoveConnector(connectorID, connectorGroupID string) {
	mockState.Lock()
	defer mockState.Unlock()

	mockState.memberships[connectorID] = connectorGroupID
}

// CurrentConnectorGroup returns the group the connector is in.
func (m *OnPremisesConnectorGroupMemberMock) CurrentConnectorGroup(connectorID string) string {
	mockState.Lock()
	defer mockState.Unlock()

	return mockState.memberships[connectorID]
}

func resetMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	mockState.connectors = make(map[string]map[string]any)
	mockState.connectorGroups = make(map[string]map[string]any)
	mockState.memberships = make(map[string]string)

	if connector, err := loadJSONResponse("get_connector_success.json"); err == nil {
		delete(connector, "@odata.context")
		mockState.connectors[UnitTestConnectorID] = connector
	}

	if groups, err := loadJSONResponse("get_connector_groups_success.json"); err == nil {
		if values, ok := groups["value"].([]any); ok {
			for _, value := range values {
				if group, ok := value.(map[string]any); ok {
					if id, ok := group["id"].(string); ok {
						mockState.connectorGroups[id] = group
					}
				}
			}
		}
	}

	// New connectors are placed in the default connector group.
	mockState.memberships[UnitTestConnectorID] = UnitTestDefaultConnectorGroupID
}

func loadJSONResponse(fileName string) (map[string]any, error) {
	content, err := os.ReadFile(filepath.Join("tests", "responses", "validate_get", fileName))
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", fileName, err)
	}

	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", fileName, err)
	}

	return result, nil
}

func connectorIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	for index, part := range parts {
		if part == "connectors" && index+1 < len(parts) {
			return strings.ToLower(parts[index+1])
		}
	}
	return ""
}

func connectorGroupIDFromODataID(odataID string) string {
	pattern := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(connectorGroupPath) + `([0-9a-f-]{36})$`)
	matches := pattern.FindStringSubmatch(odataID)
	if len(matches) != 2 {
		return ""
	}
	return strings.ToLower(matches[1])
}
