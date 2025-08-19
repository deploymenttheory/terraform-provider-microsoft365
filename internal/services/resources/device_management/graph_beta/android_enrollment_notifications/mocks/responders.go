package mocks

import (
	"encoding/json"
	"fmt"
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
	enrollmentConfigs map[string]map[string]interface{}
	templates         map[string]map[string]interface{}
	localizedMessages map[string]map[string]interface{}
	assignments       map[string][]interface{}
}

func init() {
	mockState.enrollmentConfigs = make(map[string]map[string]interface{})
	mockState.templates = make(map[string]map[string]interface{})
	mockState.localizedMessages = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("android_enrollment_notifications", &AndroidEnrollmentNotificationsMock{})
}

type AndroidEnrollmentNotificationsMock struct{}

var _ mocks.MockRegistrar = (*AndroidEnrollmentNotificationsMock)(nil)

func (m *AndroidEnrollmentNotificationsMock) RegisterMocks() {
	mockState.Lock()
	mockState.enrollmentConfigs = make(map[string]map[string]interface{})
	mockState.templates = make(map[string]map[string]interface{})
	mockState.localizedMessages = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	mockState.Unlock()

	// Register basic group mocks for assignment validation
	m.registerGroupMocks()

	// GET /deviceManagement/deviceEnrollmentConfigurations - List configurations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_get/get_android_enrollment_notifications_list.json")
		if err != nil {
			return httpmock.NewStringResponse(500, "Internal Server Error"), nil
		}
		
		var responseObj map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, "Internal Server Error"), nil
		}

		mockState.Lock()
		defer mockState.Unlock()
		
		if len(mockState.enrollmentConfigs) == 0 {
			responseObj["value"] = []interface{}{}
		} else {
			list := make([]map[string]interface{}, 0, len(mockState.enrollmentConfigs))
			for _, v := range mockState.enrollmentConfigs {
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

	// GET /deviceManagement/deviceEnrollmentConfigurations/{id} - Get specific configuration
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		mockState.Lock()
		config, ok := mockState.enrollmentConfigs[id]
		mockState.Unlock()
		if !ok {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_android_enrollment_notifications_not_found.json")
			var errObj map[string]interface{}
			json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_android_enrollment_notifications.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override template values with actual config values
		for k, v := range config {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /deviceManagement/deviceEnrollmentConfigurations - Create configuration
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := uuid.New().String()

		// Generate notification template IDs based on notification templates array
		notificationTemplates := []string{}
		if templates, ok := body["notificationTemplates"].([]interface{}); ok {
			for _, template := range templates {
				templateStr := template.(string)
				templateGuid := uuid.New().String()
				
				if templateStr == "email" {
					templateId := "Email_" + templateGuid
					notificationTemplates = append(notificationTemplates, templateId)
					
					// Store email template
					mockState.Lock()
					mockState.templates[templateGuid] = map[string]interface{}{
						"id":           templateGuid,
						"displayName":  "Email Template",
						"brandingOptions": "none",
						"localizedNotificationMessages": []interface{}{},
					}
					mockState.Unlock()
				}
				if templateStr == "push" {
					templateId := "Push_" + templateGuid
					notificationTemplates = append(notificationTemplates, templateId)
					
					// Store push template
					mockState.Lock()
					mockState.templates[templateGuid] = map[string]interface{}{
						"id":           templateGuid,
						"displayName":  "Push Template",
						"brandingOptions": "none",
						"localizedNotificationMessages": []interface{}{},
					}
					mockState.Unlock()
				}
			}
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_android_enrollment_notifications_success.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override template values with request values
		responseObj["id"] = id
		responseObj["displayName"] = body["displayName"]
		responseObj["description"] = body["description"]
		responseObj["platformType"] = body["platformType"]
		responseObj["defaultLocale"] = body["defaultLocale"]
		responseObj["notificationTemplates"] = notificationTemplates

		if v, ok := body["roleScopeTagIds"]; ok {
			responseObj["roleScopeTagIds"] = v
		} else {
			responseObj["roleScopeTagIds"] = []string{"0"}
		}

		// Store in mock state
		mockState.Lock()
		mockState.enrollmentConfigs[id] = responseObj
		mockState.assignments[id] = []interface{}{}
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// PATCH /deviceManagement/deviceEnrollmentConfigurations/{id} - Update configuration
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, ok := mockState.enrollmentConfigs[id]
		if !ok {
			mockState.Unlock()
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_android_enrollment_notifications_not_found.json")
			var errObj map[string]interface{}
			json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_update/patch_android_enrollment_notifications_success.json")
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

		mockState.enrollmentConfigs[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// GET /deviceManagement/deviceEnrollmentConfigurations/{id}/assignments - Get assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+/assignments$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		
		mockState.Lock()
		storedAssignments, ok := mockState.assignments[id]
		mockState.Unlock()
		
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_android_enrollment_notifications_assignments.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)
		
		if !ok || len(storedAssignments) == 0 {
			responseObj["value"] = []interface{}{}
		} else {
			// Transform stored assignments (from POST body) to Graph API assignment format
			graphAssignments := make([]interface{}, 0, len(storedAssignments))
			for i, assignment := range storedAssignments {
				if assignmentMap, ok := assignment.(map[string]interface{}); ok {
					assignmentId := fmt.Sprintf("%s_assignment_%d", id, i)
					
					graphAssignment := map[string]interface{}{
						"@odata.type": "#microsoft.graph.enrollmentConfigurationAssignment",
						"id":          assignmentId,
						"target": map[string]interface{}{
							"@odata.type":                               "#microsoft.graph.groupAssignmentTarget",
							"deviceAndAppManagementAssignmentFilterId":  nil,
							"deviceAndAppManagementAssignmentFilterType": "none",
							"groupId": assignmentMap["group_id"],
						},
					}
					graphAssignments = append(graphAssignments, graphAssignment)
				}
			}
			responseObj["value"] = graphAssignments
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// POST /deviceManagement/deviceEnrollmentConfigurations/{id}/assign - Assign configuration
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+/assign$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		// The SDK sends assignments as "enrollmentConfigurationAssignments"
		if enrollmentAssignments, ok := body["enrollmentConfigurationAssignments"].([]interface{}); ok {
			// Convert Graph SDK assignment format to terraform assignment format for storage
			terraformAssignments := make([]interface{}, 0, len(enrollmentAssignments))
			for _, assignment := range enrollmentAssignments {
				if assignmentMap, ok := assignment.(map[string]interface{}); ok {
					if target, ok := assignmentMap["target"].(map[string]interface{}); ok {
						terraformAssignment := map[string]interface{}{
							"type": "groupAssignmentTarget",
						}
						if groupId, exists := target["groupId"]; exists {
							terraformAssignment["group_id"] = groupId
						}
						terraformAssignments = append(terraformAssignments, terraformAssignment)
					}
				}
			}
			mockState.assignments[id] = terraformAssignments
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// GET /deviceManagement/notificationMessageTemplates/{id} - Get template
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		
		mockState.Lock()
		template, ok := mockState.templates[id]
		mockState.Unlock()
		
		if !ok {
			return httpmock.NewJsonResponse(404, map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "NotFound",
					"message": "Notification message template not found",
				},
			})
		}

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_notification_template.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override with actual template values
		for k, v := range template {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// PATCH /deviceManagement/notificationMessageTemplates/{id} - Update template (branding options)
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		var body map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		template, ok := mockState.templates[id]
		if ok {
			for k, v := range body {
				template[k] = v
			}
			mockState.templates[id] = template
		}
		mockState.Unlock()

		return httpmock.NewStringResponse(200, ""), nil
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

	// DELETE /deviceManagement/deviceEnrollmentConfigurations/{id} - Delete configuration
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		
		mockState.Lock()
		delete(mockState.enrollmentConfigs, id)
		delete(mockState.assignments, id)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AndroidEnrollmentNotificationsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.enrollmentConfigs = make(map[string]map[string]interface{})
	mockState.templates = make(map[string]map[string]interface{})
	mockState.localizedMessages = make(map[string]map[string]interface{})
	mockState.assignments = make(map[string][]interface{})
	mockState.Unlock()

	// Register basic group mocks for assignment validation (successful for error tests)
	m.registerGroupMocks()

	// Error response for creation
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_android_enrollment_notifications_error.json")
		var errObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(400, errObj)
	})

	// Error response for GET operations
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_android_enrollment_notifications_not_found.json")
		var errObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(404, errObj)
	})
}

func (m *AndroidEnrollmentNotificationsMock) registerGroupMocks() {
	// GET /groups/{id} - Get group (for assignment validation)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]
		
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group.json")
		var responseObj map[string]interface{}
		json.Unmarshal([]byte(jsonStr), &responseObj)
		responseObj["id"] = id

		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *AndroidEnrollmentNotificationsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.enrollmentConfigs {
		delete(mockState.enrollmentConfigs, id)
	}
	for id := range mockState.templates {
		delete(mockState.templates, id)
	}
	for id := range mockState.localizedMessages {
		delete(mockState.localizedMessages, id)
	}
	for id := range mockState.assignments {
		delete(mockState.assignments, id)
	}
}