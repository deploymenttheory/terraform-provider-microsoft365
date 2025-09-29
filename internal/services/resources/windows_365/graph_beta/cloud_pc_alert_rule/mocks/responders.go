package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	cloudPcAlertRules map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.cloudPcAlertRules = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// CloudPcAlertRuleMock provides mock responses for Cloud PC alert rule operations
type CloudPcAlertRuleMock struct{}

// RegisterMocks registers HTTP mock responses for Cloud PC alert rule operations
func (m *CloudPcAlertRuleMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.cloudPcAlertRules = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Cloud PC alert rules
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/monitoring/alertRules",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			rules := make([]map[string]any, 0, len(mockState.cloudPcAlertRules))
			for _, rule := range mockState.cloudPcAlertRules {
				// Ensure @odata.type is present
				ruleCopy := make(map[string]any)
				for k, v := range rule {
					ruleCopy[k] = v
				}
				if _, hasODataType := ruleCopy["@odata.type"]; !hasODataType {
					ruleCopy["@odata.type"] = "#microsoft.graph.deviceManagement.alertRule"
				}

				rules = append(rules, ruleCopy)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/monitoring/alertRules",
				"value":          rules,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for retrieving a specific Cloud PC alert rule
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/monitoring/alertRules/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			rule, exists := mockState.cloudPcAlertRules[id]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Alert rule not found"}}`), nil
			}

			// Ensure @odata.type is present
			ruleCopy := make(map[string]any)
			for k, v := range rule {
				ruleCopy[k] = v
			}
			if _, hasODataType := ruleCopy["@odata.type"]; !hasODataType {
				ruleCopy["@odata.type"] = "#microsoft.graph.deviceManagement.alertRule"
			}

			return httpmock.NewJsonResponse(200, ruleCopy)
		})

	// Register POST for creating Cloud PC alert rules
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/monitoring/alertRules",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate a new ID
			id := uuid.New().String()
			requestBody["id"] = id
			requestBody["@odata.type"] = "#microsoft.graph.deviceManagement.alertRule"

			// Set default values if not provided
			if _, hasEnabled := requestBody["enabled"]; !hasEnabled {
				requestBody["enabled"] = true
			}
			if _, hasIsSystemRule := requestBody["isSystemRule"]; !hasIsSystemRule {
				requestBody["isSystemRule"] = false
			}

			// Store in mock state
			mockState.Lock()
			mockState.cloudPcAlertRules[id] = requestBody
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, requestBody)
		})

	// Register PATCH for updating Cloud PC alert rules
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/monitoring/alertRules/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			existingRule, exists := mockState.cloudPcAlertRules[id]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Alert rule not found"}}`), nil
			}

			// Update existing rule with new values
			for k, v := range requestBody {
				existingRule[k] = v
			}

			// Handle optional fields that should be cleared if not present in request
			// This simulates the real Microsoft Graph API behavior where optional fields
			// are cleared when not included in PATCH requests
			optionalFields := []string{"description"}
			for _, field := range optionalFields {
				if _, present := requestBody[field]; !present {
					// Remove the field to simulate clearing it
					delete(existingRule, field)
				}
			}

			existingRule["@odata.type"] = "#microsoft.graph.deviceManagement.alertRule"

			mockState.cloudPcAlertRules[id] = existingRule
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, existingRule)
		})

	// Register DELETE for deleting Cloud PC alert rules
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/monitoring/alertRules/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			_, exists := mockState.cloudPcAlertRules[id]
			if exists {
				delete(mockState.cloudPcAlertRules, id)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Alert rule not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *CloudPcAlertRuleMock) RegisterErrorMocks() {
	// Register GET that returns an error
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/monitoring/alertRules",
		httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Internal server error"}}`))

	// Register POST that returns an error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/monitoring/alertRules",
		httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid alert rule configuration"}}`))

	// Register PATCH that returns an error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/monitoring/alertRules/([^/]+)$`,
		httpmock.NewStringResponder(409, `{"error":{"code":"Conflict","message":"Alert rule configuration conflict"}}`))

	// Register DELETE that returns an error
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/monitoring/alertRules/([^/]+)$`,
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient permissions to delete alert rule"}}`))
}

// GetMockCloudPcAlertRule returns a mock Cloud PC alert rule for testing
func GetMockCloudPcAlertRule() map[string]any {
	return map[string]any{
		"id":                uuid.New().String(),
		"@odata.type":       "#microsoft.graph.deviceManagement.alertRule",
		"alertRuleTemplate": "cloudPcProvisionScenario",
		"description":       "Test alert rule for Cloud PC provisioning failures",
		"displayName":       "Test Cloud PC Alert Rule",
		"enabled":           true,
		"isSystemRule":      false,
		"severity":          "warning",
		"notificationChannels": []map[string]any{
			{
				"notificationChannelType": "portal",
				"notificationReceivers": []map[string]any{
					{
						"contactInformation": "admin@contoso.com",
						"locale":             "en-US",
					},
				},
			},
		},
		"threshold": map[string]any{
			"aggregation": "count",
			"operator":    "greaterOrEqual",
			"target":      1,
		},
		"conditions": []map[string]any{
			{
				"relationshipType":  "and",
				"conditionCategory": "provisionFailures",
				"aggregation":       "count",
				"operator":          "greaterOrEqual",
				"thresholdValue":    "1",
			},
		},
	}
}

// CreateMockCloudPcAlertRuleInState creates a mock alert rule in the mock state for testing
func CreateMockCloudPcAlertRuleInState(id string, rule map[string]any) {
	mockState.Lock()
	defer mockState.Unlock()

	if rule == nil {
		rule = GetMockCloudPcAlertRule()
	}
	rule["id"] = id
	rule["@odata.type"] = "#microsoft.graph.deviceManagement.alertRule"

	mockState.cloudPcAlertRules[id] = rule
}

// ClearMockState clears all mock state for testing
func ClearMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.cloudPcAlertRules = make(map[string]map[string]any)
}
