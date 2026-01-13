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
	cloudPcDeviceImages map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.cloudPcDeviceImages = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// CloudPcDeviceImageMock provides mock responses for Cloud PC device image operations
type CloudPcDeviceImageMock struct{}

// RegisterMocks registers HTTP mock responses for Cloud PC device image operations
func (m *CloudPcDeviceImageMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.cloudPcDeviceImages = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Cloud PC device images
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/deviceImages",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			images := make([]map[string]any, 0, len(mockState.cloudPcDeviceImages))
			for _, image := range mockState.cloudPcDeviceImages {
				// Ensure @odata.type is present
				imageCopy := make(map[string]any)
				for k, v := range image {
					imageCopy[k] = v
				}
				if _, hasODataType := imageCopy["@odata.type"]; !hasODataType {
					imageCopy["@odata.type"] = "#microsoft.graph.cloudPcDeviceImage"
				}

				images = append(images, imageCopy)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/virtualEndpoint/deviceImages",
				"value":          images,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual Cloud PC device image
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/virtualEndpoint/deviceImages/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			image, exists := mockState.cloudPcDeviceImages[id]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"The specified Cloud PC device image was not found"}}`), nil
			}

			// Create response copy
			imageCopy := make(map[string]any)
			for k, v := range image {
				imageCopy[k] = v
			}
			if _, hasODataType := imageCopy["@odata.type"]; !hasODataType {
				imageCopy["@odata.type"] = "#microsoft.graph.cloudPcDeviceImage"
			}

			return httpmock.NewJsonResponse(200, imageCopy)
		})

	// Register POST for creating Cloud PC device image
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/deviceImages",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate a unique ID for the new image
			id := uuid.New().String()

			// Create the image object with required fields
			image := map[string]any{
				"@odata.type":           "#microsoft.graph.cloudPcDeviceImage",
				"id":                    id,
				"displayName":           requestBody["displayName"],
				"version":               requestBody["version"],
				"sourceImageResourceId": requestBody["sourceImageResourceId"],
				"operatingSystem":       "Windows 10 Enterprise",
				"osBuildNumber":         "19045",
				"osVersionNumber":       "10.0.19045.3930",
				"status":                "ready",
				"osStatus":              "supported",
				"lastModifiedDateTime":  "2024-01-01T00:00:00Z",
			}

			// Store in mock state
			mockState.Lock()
			mockState.cloudPcDeviceImages[id] = image
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, image)
		})

	// Register PATCH for updating Cloud PC device image
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/virtualEndpoint/deviceImages/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			image, exists := mockState.cloudPcDeviceImages[id]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"The specified Cloud PC device image was not found"}}`), nil
			}

			// Update the image with new values
			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(image, key)
				} else {
					image[key] = value
				}
			}
			image["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"

			// Ensure @odata.type is present
			if _, hasODataType := image["@odata.type"]; !hasODataType {
				image["@odata.type"] = "#microsoft.graph.cloudPcDeviceImage"
			}

			mockState.cloudPcDeviceImages[id] = image
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, image)
		})

	// Register DELETE for Cloud PC device image
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/virtualEndpoint/deviceImages/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.cloudPcDeviceImages[id]
			if exists {
				delete(mockState.cloudPcDeviceImages, id)
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"The specified Cloud PC device image was not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *CloudPcDeviceImageMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.cloudPcDeviceImages = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Cloud PC device images (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/deviceImages",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/virtualEndpoint/deviceImages",
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating Cloud PC device image with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/deviceImages",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid display name"))

	// Register error response for Cloud PC device image not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/deviceImages/not-found-image",
		factories.ErrorResponse(404, "ResourceNotFound", "Cloud PC device image not found"))
}
