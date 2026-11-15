package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	// InvalidPartnerId is a partner ID the mock rejects the way Graph rejects unknown partners.
	InvalidPartnerId = "00000000-0000-0000-0000-000000000400"
	// ForbiddenPartnerId is a partner ID the mock answers with a permission error.
	ForbiddenPartnerId = "00000000-0000-0000-0000-000000000403"

	microsoftDefenderForEndpointPartnerId = "fc780465-2017-40d4-a0c5-307022471b92"
	connectorsUrl                         = "https://graph.microsoft.com/beta/deviceManagement/mobileThreatDefenseConnectors"
	connectorByIdPattern                  = `=~^https://graph\.microsoft\.com/beta/deviceManagement/mobileThreatDefenseConnectors/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	notFoundBody                          = `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`
)

var mockState struct {
	sync.Mutex
	connectors map[string]map[string]any
}

func init() {
	mockState.connectors = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, notFoundBody))
	mocks.GlobalRegistry.Register("mobile_threat_defense_connector", &MobileThreatDefenseConnectorMock{})
}

type MobileThreatDefenseConnectorMock struct{}

var _ mocks.MockRegistrar = (*MobileThreatDefenseConnectorMock)(nil)

// derivedFields maps each read-only *DeviceBlockedOnMissingPartnerData field to the toggle Graph derives it from.
var derivedFields = map[string]string{
	"androidDeviceBlockedOnMissingPartnerData": "androidEnabled",
	"iosDeviceBlockedOnMissingPartnerData":     "iosEnabled",
	"macDeviceBlockedOnMissingPartnerData":     "macEnabled",
	"windowsDeviceBlockedOnMissingPartnerData": "windowsEnabled",
}

// serviceOwnedFields are set by the service and ignored when a client writes them.
var serviceOwnedFields = []string{
	"lastHeartbeatDateTime", "partnerState", "partnerDisplayName",
	"partnerSupportsIosCertificateSync", "partnerSupportsIosMam", "partnerSupportsAndroidMam", "partnerSupportsWindowsMam",
	"partnerMinimumSupportedIosVersion", "partnerMinimumSupportedAndroidVersion", "partnerMinimumSupportedMacVersion", "partnerMinimumSupportedWindowsVersion",
}

// loadBaseConnector loads the connector Graph returns for a newly created partner with default settings.
func loadBaseConnector() (map[string]any, error) {
	responsesPath := filepath.Join("tests", "responses", "validate_create", "post_mobile_threat_defense_connector_success.json")
	jsonData, err := os.ReadFile(responsesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON response file: %w", err)
	}

	var responseObj map[string]any
	if err := json.Unmarshal(jsonData, &responseObj); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return responseObj, nil
}

// applyWrite copies writable fields from the request body onto the connector and recomputes the derived fields,
// mirroring how Graph ignores read-only fields and partner-specific settings the partner does not support.
func applyWrite(id string, connector map[string]any, requestBody map[string]any) {
	for key, value := range requestBody {
		if key == "id" || strings.HasPrefix(key, "@odata") {
			continue
		}
		if _, derived := derivedFields[key]; derived {
			continue
		}
		skip := false
		for _, owned := range serviceOwnedFields {
			if key == owned {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		connector[key] = value
	}

	if id != microsoftDefenderForEndpointPartnerId {
		connector["microsoftDefenderForEndpointAttachEnabled"] = false
		connector["launchMobileThreatDefensePartnerOnSetupEnabled"] = false
	}

	for derived, source := range derivedFields {
		enabled, _ := connector[source].(bool)
		connector[derived] = enabled
	}
}

func graphError(status int, code, message string) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": code, "message": message}})
	if err != nil {
		return nil, fmt.Errorf("failed to build mock error response: %w", err)
	}
	return resp, nil
}

func connectorIdFromRequest(req *http.Request) string {
	parts := strings.Split(req.URL.Path, "/")
	return parts[len(parts)-1]
}

func (m *MobileThreatDefenseConnectorMock) RegisterMocks() {
	mockState.Lock()
	mockState.connectors = make(map[string]map[string]any)
	mockState.Unlock()

	// Create connector - POST /deviceManagement/mobileThreatDefenseConnectors
	// Graph treats the POST as an upsert keyed on the partner ID: the posted settings replace any existing connector.
	httpmock.RegisterResponder("POST", connectorsUrl, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return graphError(400, "BadRequest", "Invalid request body")
		}

		id, _ := requestBody["id"].(string)
		switch id {
		case "":
			return graphError(400, "InternalError", "An error has occurred")
		case InvalidPartnerId:
			return graphError(400, "BadRequest", fmt.Sprintf("Partner Id %s is not a valid partner.", id))
		case ForbiddenPartnerId:
			return graphError(403, "Forbidden", "Application is not authorized to perform this operation.")
		}

		connector, err := loadBaseConnector()
		if err != nil {
			return graphError(500, "InternalServerError", err.Error())
		}
		connector["id"] = id
		applyWrite(id, connector, requestBody)

		mockState.Lock()
		mockState.connectors[id] = connector
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, connector)
	})

	// Get connector - GET /deviceManagement/mobileThreatDefenseConnectors/{id}
	httpmock.RegisterResponder("GET", connectorByIdPattern, func(req *http.Request) (*http.Response, error) {
		id := connectorIdFromRequest(req)

		mockState.Lock()
		connector, exists := mockState.connectors[id]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		return httpmock.NewJsonResponse(200, connector)
	})

	// Update connector - PATCH /deviceManagement/mobileThreatDefenseConnectors/{id}
	httpmock.RegisterResponder("PATCH", connectorByIdPattern, func(req *http.Request) (*http.Response, error) {
		id := connectorIdFromRequest(req)

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return graphError(400, "BadRequest", "Invalid request body")
		}

		mockState.Lock()
		defer mockState.Unlock()

		connector, exists := mockState.connectors[id]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		applyWrite(id, connector, requestBody)
		mockState.connectors[id] = connector

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete connector - DELETE /deviceManagement/mobileThreatDefenseConnectors/{id}
	// Graph returns success whether or not the partner has a connector.
	httpmock.RegisterResponder("DELETE", connectorByIdPattern, func(req *http.Request) (*http.Response, error) {
		id := connectorIdFromRequest(req)

		mockState.Lock()
		delete(mockState.connectors, id)
		mockState.Unlock()

		return httpmock.NewStringResponse(200, ""), nil
	})
}

func (m *MobileThreatDefenseConnectorMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", connectorsUrl, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", connectorByIdPattern, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterResponder("PATCH", connectorByIdPattern, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", connectorByIdPattern, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

// RemoveConnector deletes a connector from the mock state to simulate an out-of-band deletion.
func (m *MobileThreatDefenseConnectorMock) RemoveConnector(id string) {
	mockState.Lock()
	delete(mockState.connectors, id)
	mockState.Unlock()
}

func (m *MobileThreatDefenseConnectorMock) CleanupMockState() {
	mockState.Lock()
	mockState.connectors = make(map[string]map[string]any)
	mockState.Unlock()
}
