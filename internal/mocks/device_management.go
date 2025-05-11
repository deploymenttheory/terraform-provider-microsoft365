package mocks

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/jarcoal/httpmock"
)

// ActivateDeviceManagementMocks sets up API responders for Device Management endpoints.
func ActivateDeviceManagementMocks() {
	// Call individual mock activations for each resource type
	activateDeviceManagementDeviceCategoryMocks()
	activateDeviceManagementRoleScopeTagMocks()
	activateDeviceManagementAssignmentFilterMocks()
	activateDeviceManagementWindowsPlatformScriptMocks()
	activateDeviceManagementWindowsFeatureUpdateProfileMocks()

	// Add more device management resource mock activations as needed
}

// activateDeviceManagementDeviceCategoryMocks sets up API responders for Device Category endpoints.
func activateDeviceManagementDeviceCategoryMocks() {
	// Device Categories - List operation
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceCategories\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			// In a real implementation, this might load from a test file
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCategories",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test Device Category 1",
						"description": "Test Description 1"
					},
					{
						"id": "00000000-0000-0000-0000-000000000002",
						"displayName": "Test Device Category 2",
						"description": "Test Description 2"
					}
				]
			}`), nil
		})

	// Device Categories - Get individual resource
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceCategories/[^/]+$`),
		func(req *http.Request) (*http.Response, error) {
			// Extract the ID from the URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCategories/$entity",
				"id": "`+id+`",
				"displayName": "Test Device Category",
				"description": "Test Description",
				"roleScopeTagIds": ["0"]
			}`), nil
		})

	// Device Categories - Create operation
	httpmock.RegisterResponder("POST", `https://graph.microsoft.com/beta/deviceManagement/deviceCategories`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCategories/$entity",
				"id": "00000000-0000-0000-0000-000000000003",
				"displayName": "Test Created Category",
				"description": "Created Test Description",
				"roleScopeTagIds": ["0"]
			}`), nil
		})

	// Device Categories - Update operation
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceCategories/[^/]+$`),
		func(req *http.Request) (*http.Response, error) {
			// Extract the ID from the URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceCategories/$entity",
				"id": "`+id+`",
				"displayName": "Test Updated Category",
				"description": "Updated Test Description",
				"roleScopeTagIds": ["0"]
			}`), nil
		})

	// Device Categories - Delete operation
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceCategories/[^/]+$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ``), nil
		})
}

// activateDeviceManagementRoleScopeTagMocks sets up API responders for Role Scope Tag endpoints.
func activateDeviceManagementRoleScopeTagMocks() {
	// Role Scope Tags - List operation
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/roleScopeTags\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/roleScopeTags",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test Scope Tag 1",
						"description": "Test Description 1"
					},
					{
						"id": "00000000-0000-0000-0000-000000000002",
						"displayName": "Test Scope Tag 2",
						"description": "Test Description 2"
					}
				]
			}`), nil
		})

	// Add other operations (GET individual, POST, PATCH, DELETE) as needed
}

// activateDeviceManagementAssignmentFilterMocks sets up API responders for Assignment Filter endpoints.
func activateDeviceManagementAssignmentFilterMocks() {
	// Assignment Filters - List operation
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/assignmentFilters\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/assignmentFilters",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test Assignment Filter 1",
						"description": "Test Description 1",
						"platform": "windows10AndLater"
					}
				]
			}`), nil
		})

	// Add other operations as needed
}

// activateDeviceManagementWindowsPlatformScriptMocks sets up API responders for Windows Platform Script endpoints.
func activateDeviceManagementWindowsPlatformScriptMocks() {
	// Windows Platform Scripts - List operation
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/deviceManagementScripts\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceManagementScripts",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test Windows Script",
						"description": "Test Script Description",
						"scriptContent": "UEsDBBQAAAgIAFmcKVaR8V40KwEAAEQCAAAHAAAAZm9vLnBzMZVRy27CMBA8R8p/WHGhSoRHBEEoVQ+9tKdWqppb5DibYuHYwXYKqeq/107SvlC/YK+emZ0dj1eSnXNSYTR6YgLHUwDFMvJ8X0YnHsCXqVpbr5xgzJ/rEgVrMt+vDIu5MSUPHppKt8iDSu+ggJlUWq1kBPe17rw9lZlFpWCrlAOJCvYOCwx6l6WulHagGBqEMpd87fLQRDoDa2CjRZUKuJnvO/N4rlyxFDW+T0yyV/AcmXK7ZWvQAC0YGxxTxNr4sLjgPZcTfEmI0A2EbSYTJ+XzGg8OdIldjBzS10KeOQgFxh/O35HjVzk1tJbjPxZ+bLw28F741o0JM9sM2z8FQXoWPwZLzVrXaW9ZYnrtN/+H7ULDxnAzCeZGfq7l6a4c7r9O+VHK0Oj54c4nUEsBAhQAFAAACAgAWZwpVpHxXjQrAQAARBIAAAcAAAAAAAAAAAAAALSBAAAAAGZvby5wczFQSwUGAAAAAAEAAQA1AAAAUAEAAAAA"
					}
				]
			}`), nil
		})

	// Add other operations as needed
}

// activateDeviceManagementWindowsFeatureUpdateProfileMocks sets up API responders for Windows Feature Update Profile endpoints.
func activateDeviceManagementWindowsFeatureUpdateProfileMocks() {
	// Windows Feature Update Profiles - List operation
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsFeatureUpdateProfiles",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test Feature Update Profile",
						"description": "Test Description",
						"featureUpdateVersion": "21H2"
					}
				]
			}`), nil
		})

	// Add other operations as needed
}
