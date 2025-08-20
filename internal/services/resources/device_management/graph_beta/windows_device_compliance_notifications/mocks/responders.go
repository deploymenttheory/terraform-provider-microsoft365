package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	templates         map[string]map[string]interface{}
	localizedMessages map[string]map[string]interface{}
}

func init() {
	mockState.templates = make(map[string]map[string]interface{})
	mockState.localizedMessages = make(map[string]map[string]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_device_compliance_notifications", &WindowsDeviceComplianceNotificationsMock{})
}

type WindowsDeviceComplianceNotificationsMock struct{}

var _ mocks.MockRegistrar = (*WindowsDeviceComplianceNotificationsMock)(nil)

func (m *WindowsDeviceComplianceNotificationsMock) RegisterMocks() {
	mockState.Lock()
	mockState.templates = make(map[string]map[string]interface{})
	mockState.localizedMessages = make(map[string]map[string]interface{})
	mockState.Unlock()

	// GET /deviceManagement/notificationMessageTemplates - List templates
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates", func(req *http.Request) (*http.Response, error) {
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_notifications_list.json")
		if err != nil {
			return httpmock.NewStringResponse(500, "Internal Server Error"), nil
		}
		
		var responseObj map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, "Internal Server Error"), nil
		}

		mockState.Lock()
		defer mockState.Unlock()
		
		if len(mockState.templates) == 0 {
			responseObj["value"] = []interface{}{}
		} else {
			list := make([]map[string]interface{}, 0, len(mockState.templates))
			for _, v := range mockState.templates {
				c := map[string]interface{}{}
				for k, vv := range v {
					c[k] = vv
				}
				list = append(list, c)
			}
			responseObj["value"] = list
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /deviceManagement/notificationMessageTemplates/{id} - Get specific template
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		
		mockState.Lock()
		template, ok := mockState.templates[id]
		mockState.Unlock()
		
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_device_compliance_notifications_not_found.json")
			var errObj map[string]interface{}
			json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_windows_device_compliance_notifications.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override template values with actual template values
		for k, v := range template {
			responseObj[k] = v
		}

		// Include localized messages if expand parameter is present
		if strings.Contains(req.URL.RawQuery, "expand") && strings.Contains(req.URL.RawQuery, "localizedNotificationMessages") {
			messages := []interface{}{}
			for messageId, message := range mockState.localizedMessages {
				if strings.HasPrefix(messageId, id+"_") {
					messages = append(messages, message)
				}
			}
			responseObj["localizedNotificationMessages"] = messages
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /deviceManagement/notificationMessageTemplates - Create template
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates", func(req *http.Request) (*http.Response, error) {
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_device_compliance_notifications_success.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override template values with request values
		responseObj["id"] = id
		responseObj["displayName"] = body["displayName"]
		if v, ok := body["brandingOptions"]; ok {
			responseObj["brandingOptions"] = v
		}
		if v, ok := body["defaultLocale"]; ok {
			responseObj["defaultLocale"] = v
		}
		if v, ok := body["roleScopeTagIds"]; ok {
			responseObj["roleScopeTagIds"] = v
		} else {
			responseObj["roleScopeTagIds"] = []string{"0"}
		}

		// Store in mock state
		mockState.Lock()
		mockState.templates[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// PATCH /deviceManagement/notificationMessageTemplates/{id} - Update template
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, ok := mockState.templates[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_device_compliance_notifications_not_found.json")
			var errObj map[string]interface{}
			json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_windows_device_compliance_notifications_success.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override with existing values
		for k, v := range existing {
			responseObj[k] = v
		}

		// Apply updates
		for k, v := range body {
			responseObj[k] = v
			existing[k] = v
		}

		// Update last modified time
		responseObj["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
		existing["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"

		mockState.templates[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /deviceManagement/notificationMessageTemplates/{id}/localizedNotificationMessages - Get localized messages
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+/localizedNotificationMessages$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		templateId := parts[len(parts)-2]
		
		mockState.Lock()
		defer mockState.Unlock()
		
		// Find all localized messages for this template
		messages := []interface{}{}
		for messageId, message := range mockState.localizedMessages {
			if strings.HasPrefix(messageId, templateId+"_") {
				messages = append(messages, message)
			}
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_localized_notification_messages.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["value"] = messages

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /deviceManagement/notificationMessageTemplates/{id}/localizedNotificationMessages - Create localized message
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+/localizedNotificationMessages$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		templateId := parts[len(parts)-2]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		messageId := templateId + "_" + body["locale"].(string)
		
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_localized_notification_message_success.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override with actual values
		responseObj["id"] = messageId
		responseObj["locale"] = body["locale"]
		responseObj["subject"] = body["subject"]
		responseObj["messageTemplate"] = body["messageTemplate"]
		responseObj["isDefault"] = body["isDefault"]

		mockState.Lock()
		mockState.localizedMessages[messageId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// PATCH /deviceManagement/notificationMessageTemplates/{templateId}/localizedNotificationMessages/{messageId} - Update localized message
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+/localizedNotificationMessages/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		messageId := parts[len(parts)-1]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, ok := mockState.localizedMessages[messageId]
		if ok {
			for k, v := range body {
				existing[k] = v
			}
			mockState.localizedMessages[messageId] = existing
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(200, ""), nil
	})

	// DELETE /deviceManagement/notificationMessageTemplates/{templateId}/localizedNotificationMessages/{messageId} - Delete localized message
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+/localizedNotificationMessages/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		messageId := parts[len(parts)-1]
		
		mockState.Lock()
		delete(mockState.localizedMessages, messageId)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// DELETE /deviceManagement/notificationMessageTemplates/{id} - Delete template
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		
		mockState.Lock()
		delete(mockState.templates, id)
		// Also delete associated localized messages
		for messageId := range mockState.localizedMessages {
			if strings.HasPrefix(messageId, id+"_") {
				delete(mockState.localizedMessages, messageId)
			}
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *WindowsDeviceComplianceNotificationsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.templates = make(map[string]map[string]interface{})
	mockState.localizedMessages = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Error response for creation
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_device_compliance_notifications_error.json")
		var errObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	// Error response for GET operations
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_device_compliance_notifications_not_found.json")
		var errObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *WindowsDeviceComplianceNotificationsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.templates {
		delete(mockState.templates, id)
	}
	for id := range mockState.localizedMessages {
		delete(mockState.localizedMessages, id)
	}
}