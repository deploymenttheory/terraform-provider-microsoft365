package mocks

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	enrollmentNotifications map[string]map[string]interface{}
}

func init() {
	mockState.enrollmentNotifications = make(map[string]map[string]interface{})
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("android_enrollment_notifications", &AndroidEnrollmentNotificationsMock{})
}

type AndroidEnrollmentNotificationsMock struct{}

var _ mocks.MockRegistrar = (*AndroidEnrollmentNotificationsMock)(nil)

func (m *AndroidEnrollmentNotificationsMock) RegisterMocks() {
	mockState.Lock()
	mockState.enrollmentNotifications = make(map[string]map[string]interface{})
	mockState.Unlock()

	// GET /deviceManagement/deviceEnrollmentConfigurations (list)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if len(mockState.enrollmentNotifications) == 0 {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceEnrollmentConfigurations",
				"value":          []interface{}{},
			})
		}

		list := make([]map[string]interface{}, 0, len(mockState.enrollmentNotifications))
		for _, v := range mockState.enrollmentNotifications {
			list = append(list, v)
		}

		return httpmock.NewJsonResponse(200, map[string]interface{}{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceEnrollmentConfigurations",
			"value":          list,
		})
	})

	// POST /deviceManagement/deviceEnrollmentConfigurations (create)
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/create_android_enrollment_notifications.json")
		var responseObj map[string]interface{}
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Generate a new ID for the created resource
		newID := uuid.New().String()
		responseObj["id"] = newID

		// Parse request body to get the actual values
		var requestBody map[string]interface{}
		_ = json.NewDecoder(req.Body).Decode(&requestBody)

		// Update response with request values
		if displayName, ok := requestBody["displayName"].(string); ok {
			responseObj["displayName"] = displayName
		}
		if description, ok := requestBody["description"].(string); ok {
			responseObj["description"] = description
		}

		// Store in mock state
		mockState.enrollmentNotifications[newID] = responseObj

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// GET /deviceManagement/deviceEnrollmentConfigurations/{id} (read)
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/([^/]+)$`), func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		if resource, exists := mockState.enrollmentNotifications[id]; exists {
			return httpmock.NewJsonResponse(200, resource)
		}

		return httpmock.NewJsonResponse(404, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "NotFound",
				"message": "Resource not found",
			},
		})
	})

	// PATCH /deviceManagement/deviceEnrollmentConfigurations/{id} (update)
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/([^/]+)$`), func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		if resource, exists := mockState.enrollmentNotifications[id]; exists {
			// Parse request body to update resource
			var requestBody map[string]interface{}
			_ = json.NewDecoder(req.Body).Decode(&requestBody)

			// Update the resource with new values
			for key, value := range requestBody {
				resource[key] = value
			}

			// Update version
			if version, ok := resource["version"].(float64); ok {
				resource["version"] = version + 1
			}

			return httpmock.NewJsonResponse(200, resource)
		}

		return httpmock.NewJsonResponse(404, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "NotFound",
				"message": "Resource not found",
			},
		})
	})

	// DELETE /deviceManagement/deviceEnrollmentConfigurations/{id} (delete)
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/([^/]+)$`), func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		if _, exists := mockState.enrollmentNotifications[id]; exists {
			delete(mockState.enrollmentNotifications, id)
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewJsonResponse(404, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "NotFound",
				"message": "Resource not found",
			},
		})
	})

	// POST /deviceManagement/deviceEnrollmentConfigurations/{id}/assign (assign)
	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/([^/]+)/assign$`), func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-2]

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.enrollmentNotifications[id]; exists {
			return httpmock.NewStringResponse(200, ""), nil
		}

		return httpmock.NewJsonResponse(404, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "NotFound",
				"message": "Resource not found",
			},
		})
	})

	// PATCH /deviceManagement/notificationMessageTemplates/{id} (branding options)
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/([^/]+)$`), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, ""), nil
	})

	// POST /deviceManagement/notificationMessageTemplates/{id}/localizedNotificationMessages (localized messages)
	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/notificationMessageTemplates/([^/]+)/localizedNotificationMessages$`), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(201, map[string]interface{}{
			"id":              uuid.New().String(),
			"locale":          "en-US",
			"subject":         "Test Subject",
			"messageTemplate": "Test Message",
			"isDefault":       true,
		})
	})
}

func (m *AndroidEnrollmentNotificationsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.enrollmentNotifications = make(map[string]map[string]interface{})
	mockState.Unlock()

	// All operations return errors
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(500, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "InternalServerError",
				"message": "Error creating resource",
			},
		})
	})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/([^/]+)$`), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(500, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "InternalServerError",
				"message": "Error reading resource",
			},
		})
	})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/([^/]+)$`), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(500, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "InternalServerError",
				"message": "Error updating resource",
			},
		})
	})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/deviceEnrollmentConfigurations/([^/]+)$`), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(500, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "InternalServerError",
				"message": "Error deleting resource",
			},
		})
	})
}

func (m *AndroidEnrollmentNotificationsMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.enrollmentNotifications = make(map[string]map[string]interface{})
}