package mocks

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/jarcoal/httpmock"
)

//go:embed ios_mobile_app_configuration_by_id.json
var iosMobileAppConfigurationByIDResponse []byte

//go:embed ios_mobile_app_configuration_with_xml.json
var iosMobileAppConfigurationWithXMLResponse []byte

// Removed unused embedded files - we're using dynamic responses instead

// IOSMobileAppConfigurationDataSourceMock provides mock responses for data source operations
type IOSMobileAppConfigurationDataSourceMock struct{}

// RegisterMocks registers HTTP mock responses for iOS mobile app configuration data source operations
func (m *IOSMobileAppConfigurationDataSourceMock) RegisterMocks() {
	// Get by ID - standard config
	httpmock.RegisterResponder(
		"GET",
		`https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000001`,
		httpmock.NewBytesResponder(http.StatusOK, iosMobileAppConfigurationByIDResponse).
			HeaderSet(http.Header{"Content-Type": []string{"application/json"}}),
	)

	// Get by ID - config with XML
	httpmock.RegisterResponder(
		"GET",
		`https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000004`,
		httpmock.NewBytesResponder(http.StatusOK, iosMobileAppConfigurationWithXMLResponse).
			HeaderSet(http.Header{"Content-Type": []string{"application/json"}}),
	)

	// Get assignments for any config by ID
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph\.microsoft\.com/v1\.0/deviceAppManagement/mobileAppConfigurations/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract the config ID from the URL
			parts := strings.Split(req.URL.Path, "/")
			configId := parts[len(parts)-2]

			// Return assignments based on config ID
			if configId == "00000000-0000-0000-0000-000000000001" {
				return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
					"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#deviceAppManagement/mobileAppConfigurations('00000000-0000-0000-0000-000000000001')/assignments",
					"value": []map[string]interface{}{
						{
							"@odata.type": "#microsoft.graph.managedDeviceMobileAppConfigurationAssignment",
							"id":          "00000000-0000-0000-0000-000000000002",
							"target": map[string]interface{}{
								"@odata.type": "#microsoft.graph.groupAssignmentTarget",
								"groupId":     "00000000-0000-0000-0000-000000000003",
							},
						},
					},
				})
			}

			// Default: empty assignments
			return httpmock.NewJsonResponse(http.StatusOK, map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#deviceAppManagement/mobileAppConfigurations",
				"value":          []interface{}{},
			})
		},
	)

	// List all configurations
	httpmock.RegisterResponder(
		"GET",
		`https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations`,
		func(req *http.Request) (*http.Response, error) {
			// Since we don't use OData filtering in the implementation,
			// we'll simulate the behavior based on the test expectation.
			// The tests will filter manually in code.

			// For testing display name searches, return different responses
			// based on test expectations

			// Test with "iOS App Config Test" expects single result
			// Test with "Non-existent Config" expects empty
			// Test with "Duplicate Config" expects multiple

			// We'll return the full list and let the code filter
			resp := httpmock.NewBytesResponse(
				http.StatusOK,
				iosMobileAppConfigurationListForTests,
			)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		},
	)
}

// Create a comprehensive list response for testing
var iosMobileAppConfigurationListForTests = []byte(`{
  "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#deviceAppManagement/mobileAppConfigurations",
  "value": [
    {
      "@odata.type": "#microsoft.graph.iosMobileAppConfiguration",
      "id": "00000000-0000-0000-0000-000000000001",
      "displayName": "iOS App Config Test",
      "description": "Test iOS app configuration",
      "createdDateTime": "2023-01-01T00:00:00Z",
      "lastModifiedDateTime": "2023-01-02T00:00:00Z",
      "version": 1,
      "targetedMobileApps": [
        "com.example.app1",
        "com.example.app2"
      ],
      "settings": [
        {
          "@odata.type": "#microsoft.graph.appConfigurationSettingItem",
          "appConfigKey": "serverUrl",
          "appConfigKeyType": "stringType",
          "appConfigKeyValue": "https://api.example.com"
        },
        {
          "@odata.type": "#microsoft.graph.appConfigurationSettingItem",
          "appConfigKey": "syncInterval",
          "appConfigKeyType": "integerType",
          "appConfigKeyValue": "300"
        }
      ]
    },
    {
      "@odata.type": "#microsoft.graph.iosMobileAppConfiguration",
      "id": "dup-00000000-0000-0000-0000-000000000001",
      "displayName": "Duplicate Config",
      "description": "First duplicate config",
      "createdDateTime": "2023-01-01T00:00:00Z",
      "lastModifiedDateTime": "2023-01-02T00:00:00Z",
      "version": 1,
      "targetedMobileApps": [],
      "settings": []
    },
    {
      "@odata.type": "#microsoft.graph.iosMobileAppConfiguration",
      "id": "dup-00000000-0000-0000-0000-000000000002",
      "displayName": "Duplicate Config",
      "description": "Second duplicate config",
      "createdDateTime": "2023-01-01T00:00:00Z",
      "lastModifiedDateTime": "2023-01-02T00:00:00Z",
      "version": 1,
      "targetedMobileApps": [],
      "settings": []
    },
    {
      "@odata.type": "#microsoft.graph.androidManagedAppConfiguration",
      "id": "android-00000000-0000-0000-0000-000000000001",
      "displayName": "Android Config",
      "description": "Should be filtered out",
      "createdDateTime": "2023-01-01T00:00:00Z",
      "lastModifiedDateTime": "2023-01-02T00:00:00Z",
      "version": 1
    }
  ]
}`)

// RegisterErrorMocks registers error responses for testing error handling
func (m *IOSMobileAppConfigurationDataSourceMock) RegisterErrorMocks() {
	// Register 404 for all requests
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations`,
		httpmock.NewStringResponder(
			http.StatusNotFound,
			`{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`,
		),
	)
}
