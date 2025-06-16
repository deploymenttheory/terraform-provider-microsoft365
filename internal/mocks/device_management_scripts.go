package mocks

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/jarcoal/httpmock"
)

// ActivateDeviceShellScriptMocks sets up API responders for Device Shell Script endpoints.
func ActivateDeviceShellScriptMocks() {
	// Call individual mock activations for each operation type
	activateDeviceShellScriptBasicOperationMocks()
	activateDeviceShellScriptAssignmentMocks()
	activateDeviceShellScriptErrorScenarioMocks()
}

// activateDeviceShellScriptBasicOperationMocks sets up mocks for CRUD operations
func activateDeviceShellScriptBasicOperationMocks() {
	// Device Shell Scripts - List operation
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test macOS Shell Script 1",
						"description": "Test Description 1",
						"runAsAccount": "system",
						"fileName": "test-script-1.sh",
						"scriptContent": "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gV29ybGQn",
						"createdDateTime": "2023-01-01T12:00:00Z",
						"lastModifiedDateTime": "2023-01-01T12:00:00Z",
						"roleScopeTagIds": ["0"],
						"blockExecutionNotifications": false,
						"executionFrequency": null,
						"retryCount": 3
					},
					{
						"id": "00000000-0000-0000-0000-000000000002",
						"displayName": "Test macOS Shell Script 2",
						"description": "Test Description 2",
						"runAsAccount": "user",
						"fileName": "test-script-2.sh",
						"scriptContent": "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gVXNlcic=",
						"createdDateTime": "2023-01-02T12:00:00Z",
						"lastModifiedDateTime": "2023-01-02T12:00:00Z",
						"roleScopeTagIds": ["0"],
						"blockExecutionNotifications": true,
						"executionFrequency": "P1D",
						"retryCount": 5
					}
				]
			}`), nil
		})

	// Device Shell Scripts - Get individual resource
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[^/]+$`),
		func(req *http.Request) (*http.Response, error) {
			// Extract the ID from the URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts/$entity",
				"id": "`+id+`",
				"displayName": "Test macOS Shell Script",
				"description": "Test Description",
				"runAsAccount": "system",
				"fileName": "test-script.sh",
				"scriptContent": "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gV29ybGQn",
				"createdDateTime": "2023-01-01T12:00:00Z",
				"lastModifiedDateTime": "2023-01-01T12:00:00Z",
				"roleScopeTagIds": ["0"],
				"blockExecutionNotifications": false,
				"executionFrequency": null,
				"retryCount": 3
			}`), nil
		})

	// Device Shell Scripts - Get with expanded assignments
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[^/]+\?.*expand=assignments.*$`),
		func(req *http.Request) (*http.Response, error) {
			// Extract the ID from the URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]
			// Remove query parameters from ID
			if strings.Contains(id, "?") {
				id = strings.Split(id, "?")[0]
			}

			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts(assignments())/$entity",
				"id": "`+id+`",
				"displayName": "Test macOS Shell Script",
				"description": "Test Description",
				"runAsAccount": "system",
				"fileName": "test-script.sh",
				"scriptContent": "IyEvYmluL2Jhc2gKZWNobyAnSGVsbG8gV29ybGQn",
				"createdDateTime": "2023-01-01T12:00:00Z",
				"lastModifiedDateTime": "2023-01-01T12:00:00Z",
				"roleScopeTagIds": ["0"],
				"blockExecutionNotifications": false,
				"executionFrequency": null,
				"retryCount": 3,
				"assignments": {
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts('`+id+`')/assignments",
					"value": [
						{
							"id": "assignment-1",
							"target": {
								"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget"
							}
						}
					]
				}
			}`), nil
		})

	// Device Shell Scripts - Create operation
	httpmock.RegisterResponder("POST", `https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts/$entity",
				"id": "00000000-0000-0000-0000-000000000001",
				"displayName": "Test Created macOS Script",
				"description": "Created Test Description",
				"runAsAccount": "system",
				"fileName": "created-script.sh",
				"scriptContent": "IyEvYmluL2Jhc2gKZWNobyAnQ3JlYXRlZCBTY3JpcHQn",
				"createdDateTime": "2023-01-01T12:00:00Z",
				"lastModifiedDateTime": "2023-01-01T12:00:00Z",
				"roleScopeTagIds": ["0"],
				"blockExecutionNotifications": false,
				"executionFrequency": null,
				"retryCount": 3
			}`), nil
		})

	// Device Shell Scripts - Update operation
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[^/]+$`),
		func(req *http.Request) (*http.Response, error) {
			// Extract the ID from the URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts/$entity",
				"id": "`+id+`",
				"displayName": "Test Updated macOS Script",
				"description": "Updated Test Description",
				"runAsAccount": "user",
				"fileName": "updated-script.sh",
				"scriptContent": "IyEvYmluL2Jhc2gKZWNobyAnVXBkYXRlZCBTY3JpcHQn",
				"createdDateTime": "2023-01-01T12:00:00Z",
				"lastModifiedDateTime": "2023-01-02T12:00:00Z",
				"roleScopeTagIds": ["0"],
				"blockExecutionNotifications": true,
				"executionFrequency": "P7D",
				"retryCount": 5
			}`), nil
		})

	// Device Shell Scripts - Delete operation
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[^/]+$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ``), nil
		})
}

// activateDeviceShellScriptAssignmentMocks sets up mocks for assignment operations
func activateDeviceShellScriptAssignmentMocks() {
	// Device Shell Scripts - Assign operation
	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[^/]+/assign$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#Edm.String",
				"value": "Assignment successful"
			}`), nil
		})

	// Device Shell Scripts - Get assignments
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceShellScripts/[^/]+/assignments$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceShellScripts/assignments",
				"value": [
					{
						"id": "assignment-1",
						"target": {
							"@odata.type": "#microsoft.graph.allLicensedUsersAssignmentTarget"
						}
					}
				]
			}`), nil
		})
}

// activateDeviceShellScriptErrorScenarioMocks sets up mocks for error scenarios
func activateDeviceShellScriptErrorScenarioMocks() {
	// 404 Not Found for non-existent resources
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/99999999-9999-9999-9999-999999999999",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, `{
				"error": {
					"code": "NotFound",
					"message": "The specified device shell script was not found.",
					"innerError": {
						"date": "2023-01-01T12:00:00",
						"request-id": "test-request-id",
						"client-request-id": "test-client-request-id"
					}
				}
			}`), nil
		})

	// 400 Bad Request for invalid script content
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/invalid",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{
				"error": {
					"code": "BadRequest",
					"message": "Invalid script content provided.",
					"innerError": {
						"date": "2023-01-01T12:00:00",
						"request-id": "test-request-id",
						"client-request-id": "test-client-request-id"
					}
				}
			}`), nil
		})

	// 403 Forbidden for insufficient permissions
	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/beta/deviceManagement/deviceShellScripts/forbidden",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusForbidden, `{
				"error": {
					"code": "Forbidden",
					"message": "Insufficient privileges to complete the operation.",
					"innerError": {
						"date": "2023-01-01T12:00:00",
						"request-id": "test-request-id",
						"client-request-id": "test-client-request-id"
					}
				}
			}`), nil
		})
}

// MockDeviceShellScriptRequest is a helper function to register a custom mock for a specific endpoint
func MockDeviceShellScriptRequest(method, urlPattern string, statusCode int, responseBody string) {
	httpmock.RegisterResponder(method, urlPattern,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(statusCode, responseBody), nil
		})
}

// MockDeviceShellScriptRequestWithRegexp is a helper function to register a custom mock with a regexp pattern
func MockDeviceShellScriptRequestWithRegexp(method string, urlRegexp *regexp.Regexp, statusCode int, responseBody string) {
	httpmock.RegisterRegexpResponder(method, urlRegexp,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(statusCode, responseBody), nil
		})
}
