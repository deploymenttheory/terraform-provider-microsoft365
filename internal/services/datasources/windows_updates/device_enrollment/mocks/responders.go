package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jarcoal/httpmock"
)

// RegisterGetDeviceByIdSuccessMock registers a mock for successful device lookup by ID
func RegisterGetDeviceByIdSuccessMock() {
	responseFile := getResponseFilePath("get_device_by_id_success.json")
	responseData, err := os.ReadFile(responseFile)
	if err != nil {
		panic("Failed to load mock response file: " + err.Error())
	}

	var responseObj map[string]interface{}
	if err := json.Unmarshal(responseData, &responseObj); err != nil {
		panic("Failed to parse mock response JSON: " + err.Error())
	}

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, responseObj)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})
}

// RegisterListAllDevicesSuccessMock registers a mock for successful list all devices
func RegisterListAllDevicesSuccessMock() {
	responseFile := getResponseFilePath("list_all_devices_success.json")
	responseData, err := os.ReadFile(responseFile)
	if err != nil {
		panic("Failed to load mock response file: " + err.Error())
	}

	var responseObj map[string]interface{}
	if err := json.Unmarshal(responseData, &responseObj); err != nil {
		panic("Failed to parse mock response JSON: " + err.Error())
	}

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets`, func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, responseObj)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})
}

// RegisterListDevicesWithFilterSuccessMock registers a mock for list devices with OData filter
func RegisterListDevicesWithFilterSuccessMock() {
	responseFile := getResponseFilePath("list_devices_with_filter_success.json")
	responseData, err := os.ReadFile(responseFile)
	if err != nil {
		panic("Failed to load mock response file: " + err.Error())
	}

	var responseObj map[string]interface{}
	if err := json.Unmarshal(responseData, &responseObj); err != nil {
		panic("Failed to parse mock response JSON: " + err.Error())
	}

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets\?.*filter=`, func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, responseObj)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})
}

// RegisterGetDeviceByIdErrorMock registers a mock for device not found error
func RegisterGetDeviceByIdErrorMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		errorResponse := map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "NotFound",
				"message": "Device not found",
			},
		}
		resp, err := httpmock.NewJsonResponse(404, errorResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})
}

// RegisterGetDeviceWithRegistrationErrorMock registers a mock for device with registration errors
func RegisterGetDeviceWithRegistrationErrorMock() {
	responseFile := getResponseFilePath("get_device_with_error_success.json")
	responseData, err := os.ReadFile(responseFile)
	if err != nil {
		panic("Failed to load mock response file: " + err.Error())
	}

	var responseObj map[string]interface{}
	if err := json.Unmarshal(responseData, &responseObj); err != nil {
		panic("Failed to parse mock response JSON: " + err.Error())
	}

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+$`, func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, responseObj)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})
}

// RegisterGetDeviceByNameSuccessMock registers mocks for device name lookup (managed devices + enrollment)
func RegisterGetDeviceByNameSuccessMock() {
	// Mock the managed devices API call
	managedDevicesResponse := map[string]interface{}{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/managedDevices",
		"value": []map[string]interface{}{
			{
				"id":               "managed-device-id-001",
				"deviceName":       "TEST-DEVICE-001",
				"azureADDeviceId":  "fb95f07d-9e73-411d-99ab-7eca3a5122b1",
				"operatingSystem":  "Windows",
				"osVersion":        "10.0.19045",
				"complianceState":  "compliant",
				"managementAgent":  "mdm",
			},
		},
	}

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices`, func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, managedDevicesResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})

	// Mock the enrollment status API call
	RegisterGetDeviceByIdSuccessMock()
}

// getResponseFilePath returns the absolute path to a response file
func getResponseFilePath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	mockDir := filepath.Dir(currentFile)
	return filepath.Join(mockDir, "..", "tests", "responses", "validate_get", filename)
}
