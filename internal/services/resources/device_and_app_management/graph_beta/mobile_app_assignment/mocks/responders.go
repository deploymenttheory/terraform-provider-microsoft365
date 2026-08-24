package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of assignments for consistent responses
var mockState struct {
	sync.Mutex
	assignments map[string]map[string]any
	// requests records every assignment payload received, so tests can assert on what the
	// provider actually sent rather than only on what ended up in Terraform state.
	requests []map[string]any
}

func init() {
	mockState.assignments = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("mobile_app_assignment", &MobileAppAssignmentMock{})
}

// MobileAppAssignmentMock provides mock responses for mobile app assignment operations
type MobileAppAssignmentMock struct{}

// isRemovableRejected mirrors the Intune service behaviour: the service rejects the
// isRemovable setting for any intent other than "required", responding with
// HTTP 400 "IsRemovable setting is only supported for Required intent."
//
// Returns the service error response when the payload is invalid, or nil when it is accepted.
func isRemovableRejected(assignment map[string]any) *http.Response {
	settings, ok := assignment["settings"].(map[string]any)
	if !ok {
		return nil
	}

	if _, present := settings["isRemovable"]; !present {
		return nil
	}

	intent, _ := assignment["intent"].(string)
	if intent == "required" {
		return nil
	}

	return httpmock.NewStringResponse(400,
		`{"error":{"code":"BadRequest","message":"IsRemovable setting is only supported for Required intent."}}`)
}

// RegisterMocks registers HTTP mock responses for mobile app assignment operations
func (m *MobileAppAssignmentMock) RegisterMocks() {
	mockState.Lock()
	mockState.assignments = make(map[string]map[string]any)
	mockState.requests = nil
	mockState.Unlock()

	// POST /deviceAppManagement/mobileApps/{mobileAppId}/assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			var assignment map[string]any
			if err := json.NewDecoder(req.Body).Decode(&assignment); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			mockState.requests = append(mockState.requests, assignment)
			mockState.Unlock()

			if errResp := isRemovableRejected(assignment); errResp != nil {
				return errResp, nil
			}

			if assignment["id"] == nil {
				assignment["id"] = uuid.New().String()
			}
			assignment["@odata.type"] = "#microsoft.graph.mobileAppAssignment"

			assignmentId := assignment["id"].(string)
			mockState.Lock()
			mockState.assignments[assignmentId] = assignment
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, assignment)
		})

	// GET /deviceAppManagement/mobileApps/{mobileAppId}/assignments/{assignmentId}
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/[^/]+/assignments/[^/?]+(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			assignmentId := urlParts[len(urlParts)-1]

			mockState.Lock()
			assignment, exists := mockState.assignments[assignmentId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Assignment not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, assignment)
		})

	// GET /deviceAppManagement/mobileApps/{mobileAppId}/assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/[^/]+/assignments(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			assignments := make([]map[string]any, 0, len(mockState.assignments))
			for _, assignment := range mockState.assignments {
				assignments = append(assignments, assignment)
			}

			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps('00000000-0000-0000-0000-000000000001')/assignments",
				"value":          assignments,
			})
		})

	// PATCH /deviceAppManagement/mobileApps/{mobileAppId}/assignments/{assignmentId}
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/[^/]+/assignments/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			assignmentId := urlParts[len(urlParts)-1]

			mockState.Lock()
			assignment, exists := mockState.assignments[assignmentId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Assignment not found"}}`), nil
			}

			var update map[string]any
			if err := json.NewDecoder(req.Body).Decode(&update); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			merged := make(map[string]any, len(assignment))
			for k, v := range assignment {
				merged[k] = v
			}
			for k, v := range update {
				merged[k] = v
			}

			if errResp := isRemovableRejected(merged); errResp != nil {
				return errResp, nil
			}

			mockState.Lock()
			mockState.assignments[assignmentId] = merged
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, merged)
		})

	// DELETE /deviceAppManagement/mobileApps/{mobileAppId}/assignments/{assignmentId}
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/[^/]+/assignments/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			assignmentId := urlParts[len(urlParts)-1]

			mockState.Lock()
			delete(mockState.assignments, assignmentId)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// CleanupMockState clears all stored assignments from the mock state
func (m *MobileAppAssignmentMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	for id := range mockState.assignments {
		delete(mockState.assignments, id)
	}
	mockState.requests = nil
}

// SettingsSent returns the settings object of the most recent assignment request whose
// settings contain the given block name, and whether such a request was seen. It lets tests
// assert on the payload the provider actually sent.
func (m *MobileAppAssignmentMock) SettingsSent(block string) (map[string]any, bool) {
	mockState.Lock()
	defer mockState.Unlock()

	for i := len(mockState.requests) - 1; i >= 0; i-- {
		settings, ok := mockState.requests[i]["settings"].(map[string]any)
		if !ok {
			continue
		}
		if odataType, ok := settings["@odata.type"].(string); ok && strings.Contains(strings.ToLower(odataType), strings.ToLower(block)) {
			return settings, true
		}
	}

	return nil, false
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *MobileAppAssignmentMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/[^/]+/assignments$`,
		factories.ErrorResponse(400, "BadRequest", "Error creating mobile app assignment"))
}
