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
	ipApplicationSegments map[string]map[string]any
}

func init() {
	mockState.ipApplicationSegments = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("ip_application_segment", &IpApplicationSegmentMock{})
}

type IpApplicationSegmentMock struct{}

var _ mocks.MockRegistrar = (*IpApplicationSegmentMock)(nil)

func (m *IpApplicationSegmentMock) RegisterMocks() {
	mockState.Lock()
	mockState.ipApplicationSegments = make(map[string]map[string]any)
	mockState.Unlock()

	// Create IP application segment - POST /applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments$`, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Generate a UUID for the new resource
		newId := uuid.New().String()

		// Load the template response
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_ip_application_segment_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load response"}}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse response"}}`), nil
		}

		// Update response with request data
		responseObj["id"] = newId
		if destinationHost, ok := requestBody["destinationHost"]; ok {
			responseObj["destinationHost"] = destinationHost
		}
		if destinationType, ok := requestBody["destinationType"]; ok {
			responseObj["destinationType"] = destinationType
		}
		if ports, ok := requestBody["ports"]; ok {
			responseObj["ports"] = ports
		}
		if protocol, ok := requestBody["protocol"]; ok {
			responseObj["protocol"] = protocol
		}

		// Store in mock state
		mockState.Lock()
		mockState.ipApplicationSegments[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get IP application segment - GET /applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		segmentId := parts[len(parts)-1]

		mockState.Lock()
		segment, exists := mockState.ipApplicationSegments[segmentId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Return the stored segment data which includes all the request data
		return httpmock.NewJsonResponse(200, segment)
	})

	// Update IP application segment - PATCH /applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		segmentId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		segment, exists := mockState.ipApplicationSegments[segmentId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update fields from request
		for key, value := range requestBody {
			segment[key] = value
		}
		mockState.ipApplicationSegments[segmentId] = segment
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete IP application segment - DELETE /applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		segmentId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.ipApplicationSegments[segmentId]
		if exists {
			delete(mockState.ipApplicationSegments, segmentId)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *IpApplicationSegmentMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *IpApplicationSegmentMock) CleanupMockState() {
	mockState.Lock()
	mockState.ipApplicationSegments = make(map[string]map[string]any)
	mockState.Unlock()
}

