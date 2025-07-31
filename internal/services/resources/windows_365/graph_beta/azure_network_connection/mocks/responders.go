package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	connections map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.connections = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// AzureNetworkConnectionMock provides mock responses for azure network connection operations
type AzureNetworkConnectionMock struct{}

// RegisterMocks registers HTTP mock responses for azure network connection operations
func (m *AzureNetworkConnectionMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.connections = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register GET for listing connections
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/onPremisesConnections",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			connections := make([]map[string]interface{}, 0, len(mockState.connections))
			for _, conn := range mockState.connections {
				connections = append(connections, conn)
			}
			mockState.Unlock()

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/virtualEndpoint/onPremisesConnections",
				"value":          connections,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual connection
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/onPremisesConnections/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			connectionId := urlParts[len(urlParts)-1]

			mockState.Lock()
			connectionData, exists := mockState.connections[connectionId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Connection not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, connectionData)
		})

	// Register POST for creating connection
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/onPremisesConnections",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate new connection ID
			connectionId := uuid.New().String()

			// Create connection data with all fields
			connectionData := map[string]interface{}{
				"id":                  connectionId,
				"displayName":         requestBody["displayName"],
				"connectionType":      requestBody["connectionType"], // Use correct field name
				"adDomainName":        requestBody["adDomainName"],
				"adDomainUsername":    requestBody["adDomainUsername"],
				"organizationalUnit":  requestBody["organizationalUnit"],
				"resourceGroupId":     requestBody["resourceGroupId"],
				"subnetId":            requestBody["subnetId"],
				"subscriptionId":      requestBody["subscriptionId"],
				"virtualNetworkId":    requestBody["virtualNetworkId"],
				"healthCheckStatus":   "passed",
				"managedBy":           "windows365",
				"inUse":               false,
			}

			// Store in mock state
			mockState.Lock()
			mockState.connections[connectionId] = connectionData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, connectionData)
		})

	// Register PATCH for updating connection
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/onPremisesConnections/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			connectionId := urlParts[len(urlParts)-1]

			mockState.Lock()
			connectionData, exists := mockState.connections[connectionId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Connection not found"}}`), nil
			}

			// Parse request body
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update connection data
			mockState.Lock()
			
			// For PATCH operations, we need to handle the case where optional fields 
			// are removed from the configuration (like going from maximal to minimal)
			// Check for specific field patterns to simulate real API behavior
			
			// If this looks like a minimal config update (no organizationalUnit in request)
			_, hasOrgUnit := requestBody["organizationalUnit"]
			if !hasOrgUnit {
				// Remove organizationalUnit from the stored state to simulate API clearing it
				delete(connectionData, "organizationalUnit")
			}
			
			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(connectionData, key)
				} else {
					connectionData[key] = value
				}
			}
			// Ensure the ID is preserved
			connectionData["id"] = connectionId
			mockState.connections[connectionId] = connectionData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, connectionData)
		})

	// Register DELETE for removing connection
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/onPremisesConnections/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			connectionId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.connections[connectionId]
			if exists {
				delete(mockState.connections, connectionId)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Connection not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register specific connection mocks for testing
	registerSpecificConnectionMocks()
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *AzureNetworkConnectionMock) RegisterErrorMocks() {
	// Register error response for creating connection with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/onPremisesConnections",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid resource group ID"))

	// Register error response for connection not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/onPremisesConnections/not-found-connection",
		factories.ErrorResponse(404, "ResourceNotFound", "Connection not found"))
}

// registerSpecificConnectionMocks registers mocks for specific test scenarios
func registerSpecificConnectionMocks() {
	// Minimal connection
	minimalConnectionId := "11111111-1111-1111-1111-111111111111"
	minimalConnectionData := map[string]interface{}{
		"id":                 minimalConnectionId,
		"displayName":        "Test Minimal Connection",
		"connectionType":     "hybridAzureADJoin",
		"adDomainName":       "example.local",
		"adDomainUsername":   "testuser",
		"resourceGroupId":    "/subscriptions/11111111-1111-1111-1111-111111111111/resourcegroups/test-rg",
		"subnetId":           "/subscriptions/11111111-1111-1111-1111-111111111111/resourcegroups/test-rg/providers/microsoft.network/virtualnetworks/test-vnet/subnets/test-subnet",
		"subscriptionId":     "11111111-1111-1111-1111-111111111111",
		"virtualNetworkId":   "/subscriptions/11111111-1111-1111-1111-111111111111/resourcegroups/test-rg/providers/microsoft.network/virtualnetworks/test-vnet",
		"healthCheckStatus":  "passed",
		"managedBy":          "windows365",
		"inUse":              false,
	}

	mockState.Lock()
	mockState.connections[minimalConnectionId] = minimalConnectionData
	mockState.Unlock()

	// Maximal connection
	maximalConnectionId := "22222222-2222-2222-2222-222222222222"
	maximalConnectionData := map[string]interface{}{
		"id":                 maximalConnectionId,
		"displayName":        "Test Maximal Connection",
		"connectionType":     "hybridAzureADJoin",
		"adDomainName":       "example.local",
		"adDomainUsername":   "testuser",
		"organizationalUnit": "OU=CloudPCs,DC=example,DC=local",
		"resourceGroupId":    "/subscriptions/22222222-2222-2222-2222-222222222222/resourcegroups/test-rg-maximal",
		"subnetId":           "/subscriptions/22222222-2222-2222-2222-222222222222/resourcegroups/test-rg-maximal/providers/microsoft.network/virtualnetworks/test-vnet-maximal/subnets/test-subnet-maximal",
		"subscriptionId":     "22222222-2222-2222-2222-222222222222",
		"virtualNetworkId":   "/subscriptions/22222222-2222-2222-2222-222222222222/resourcegroups/test-rg-maximal/providers/microsoft.network/virtualnetworks/test-vnet-maximal",
		"healthCheckStatus":  "passed",
		"managedBy":          "windows365",
		"inUse":              false,
	}

	mockState.Lock()
	mockState.connections[maximalConnectionId] = maximalConnectionData
	mockState.Unlock()

	// Register specific GET for these connections
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/onPremisesConnections/"+minimalConnectionId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			connectionData := mockState.connections[minimalConnectionId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, connectionData)
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/onPremisesConnections/"+maximalConnectionId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			connectionData := mockState.connections[maximalConnectionId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, connectionData)
		})
}