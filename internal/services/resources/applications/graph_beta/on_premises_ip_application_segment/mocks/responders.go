package mocks

import (
	"encoding/json"
	"net"
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
	mocks.GlobalRegistry.Register("on_premises_ip_application_segment", &OnPremisesIpApplicationSegmentMock{})
}

type OnPremisesIpApplicationSegmentMock struct{}

var _ mocks.MockRegistrar = (*OnPremisesIpApplicationSegmentMock)(nil)

func (m *OnPremisesIpApplicationSegmentMock) RegisterMocks() {
	mockState.Lock()
	mockState.ipApplicationSegments = make(map[string]map[string]any)
	mockState.Unlock()

	// Create IP application segment - POST /applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments$`, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}
		if requestBody["destinationType"] == "ipAddress" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"InvalidJson_BadRequest","message":"Valid JSON content expected."}}`), nil
		}
		if response := normalizeSegmentRequest(requestBody); response != nil {
			return response, nil
		}
		if requestBody["port"] != float64(0) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"port must be 0 for ip application segments"}}`), nil
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
		if port, ok := requestBody["port"]; ok {
			responseObj["port"] = port
		}
		if protocol, ok := requestBody["protocol"]; ok {
			responseObj["protocol"] = normalizeProtocolResponse(protocol)
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
		if requestBody["destinationType"] == "ipAddress" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"InvalidJson_BadRequest","message":"Valid JSON content expected."}}`), nil
		}
		if response := normalizeSegmentRequest(requestBody); response != nil {
			return response, nil
		}
		if requestBody["port"] != float64(0) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"port must be 0 for ip application segments"}}`), nil
		}

		mockState.Lock()
		segment, exists := mockState.ipApplicationSegments[segmentId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update fields from request
		for key, value := range requestBody {
			if key == "protocol" {
				value = normalizeProtocolResponse(value)
			}
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

// normalizeSegmentRequest mirrors the real beta endpoint's observed behavior
// for destination types this resource does not send verbatim. The real
// endpoint infers the stored destination type from the host format (CIDR hosts
// become "ipRangeCidr", single IPs become "ip", FQDNs become "fqdn"); the
// resource's ValidateConfig prevents such mismatches from being sent for
// "ipRange". It mutates requestBody in place and returns a non-nil error
// response when the real endpoint would reject the request:
//   - "ipRange" requires destinationHost as start and end addresses separated
//     by "..", e.g. "192.168.1.1..192.168.1.10"; CIDR hosts are normalized to
//     destinationType "ipRangeCidr" and other forms return 400
//     DestinationHost_InvalidIP.
//   - "dnsSuffix" is accepted, but Graph discards the segment's ports
//     (returned as []) and protocol (returned as "0").
func normalizeSegmentRequest(requestBody map[string]any) *http.Response {
	if requestBody["destinationType"] == "ipRange" {
		host, _ := requestBody["destinationHost"].(string)
		if _, _, err := net.ParseCIDR(host); err == nil {
			requestBody["destinationType"] = "ipRangeCidr"
		} else if !isValidIpRangeHost(host) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"DestinationHost_InvalidIP","message":"IP address invalid"}}`)
		}
	}

	if requestBody["destinationType"] == "dnsSuffix" {
		requestBody["ports"] = []any{}
		requestBody["protocol"] = "0"
	}

	return nil
}

func isValidIpRangeHost(host string) bool {
	parts := strings.Split(host, "..")
	if len(parts) != 2 {
		return false
	}
	for _, part := range parts {
		if net.ParseIP(part) == nil {
			return false
		}
	}

	return true
}

func normalizeProtocolResponse(value any) any {
	protocol, ok := value.(string)
	if !ok {
		return value
	}

	protocols := strings.Split(protocol, ",")
	for i := range protocols {
		protocols[i] = strings.TrimSpace(protocols[i])
	}

	return strings.Join(protocols, ",")
}

func (m *OnPremisesIpApplicationSegmentMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/onPremisesPublishing/segmentsConfiguration/microsoft\.graph\.ipSegmentConfiguration/applicationSegments/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *OnPremisesIpApplicationSegmentMock) CleanupMockState() {
	mockState.Lock()
	mockState.ipApplicationSegments = make(map[string]map[string]any)
	mockState.Unlock()
}
