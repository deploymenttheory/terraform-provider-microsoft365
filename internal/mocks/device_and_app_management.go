package mocks

import (
	"net/http"
	"regexp"

	"github.com/jarcoal/httpmock"
)

// ActivateDeviceAndAppManagementMocks sets up mocks for Device and App Management endpoints.
func ActivateDeviceAndAppManagementMocks() {
	// Application Categories
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppCategories\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileAppCategories",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test Category 1",
						"lastModifiedDateTime": "2023-01-01T12:00:00Z"
					},
					{
						"id": "00000000-0000-0000-0000-000000000002",
						"displayName": "Test Category 2",
						"lastModifiedDateTime": "2023-01-01T12:00:00Z"
					}
				]
			}`), nil
		})

	// macOS PKG Apps
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps",
				"value": [
					{
						"@odata.type": "#microsoft.graph.macOSPkgApp",
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Test macOS App",
						"description": "Test Description",
						"publisher": "Test Publisher",
						"packageIdentifier": "com.test.app"
					}
				]
			}`), nil
		})

	// Cloud PC Device Images (v1.0)
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/v1\.0/deviceManagement/virtualEndpoint/deviceImages\?.*$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#deviceManagement/virtualEndpoint/deviceImages",
				"value": [
					{
						"id": "00000000-0000-0000-0000-000000000001",
						"displayName": "Windows 11 Enterprise",
						"operatingSystem": "Windows",
						"osVersion": "Windows 11 Enterprise"
					},
					{
						"id": "00000000-0000-0000-0000-000000000002",
						"displayName": "Windows 10 Enterprise",
						"operatingSystem": "Windows",
						"osVersion": "Windows 10 Enterprise"
					}
				]
			}`), nil
		})
}
