package mocks

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

type WindowsUpdateDeploymentAudienceMembersMock struct{}

var (
	mockState = struct {
		sync.Mutex
		audiences map[string]map[string]any
	}{
		audiences: make(map[string]map[string]any),
	}
)

func (m *WindowsUpdateDeploymentAudienceMembersMock) RegisterMocks() {
	// Register GET for audience (to read members/exclusions)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)$`,
		m.getAudienceResponder())

	// Register GET for members collection
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)/members$`,
		m.getMembersResponder())

	// Register GET for exclusions collection
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)/exclusions$`,
		m.getExclusionsResponder())

	// Register POST for updateAudience action
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)/microsoft\.graph\.windowsUpdates\.updateAudience$`,
		m.updateAudienceResponder())

	// Register POST for creating audience (for dependencies)
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences$`,
		m.createAudienceResponder())

	// Register DELETE for audience (for cleanup)
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)$`,
		m.deleteAudienceResponder())
}

func (m *WindowsUpdateDeploymentAudienceMembersMock) createAudienceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		id := uuid.New().String()

		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_audience_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load response JSON"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse response JSON"}}`), nil
		}

		response["id"] = id
		response["members"] = []any{}
		response["exclusions"] = []any{}

		mockState.Lock()
		mockState.audiences[id] = response
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, response)
	}
}

func (m *WindowsUpdateDeploymentAudienceMembersMock) getAudienceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deploymentAudiences/")

		mockState.Lock()
		audience, exists := mockState.audiences[id]
		mockState.Unlock()

		if !exists {
			// For unit tests, if audience doesn't exist, create a minimal one
			// This handles the case where the audience was created but we're reading it
			audience = map[string]any{
				"@odata.type": "#microsoft.graph.windowsUpdates.deploymentAudience",
				"id":          id,
				"members":     []any{},
				"exclusions":  []any{},
			}
			mockState.Lock()
			mockState.audiences[id] = audience
			mockState.Unlock()
		}

		return httpmock.NewJsonResponse(200, audience)
	}
}

func (m *WindowsUpdateDeploymentAudienceMembersMock) updateAudienceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deploymentAudiences/")

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		mockState.Lock()
		audience, exists := mockState.audiences[id]
		if !exists {
			// For unit tests, auto-create the audience if it doesn't exist
			audience = map[string]any{
				"@odata.type": "#microsoft.graph.windowsUpdates.deploymentAudience",
				"id":          id,
				"members":     []any{},
				"exclusions":  []any{},
			}
			mockState.audiences[id] = audience
		}

		// Initialize members and exclusions if they don't exist
		if _, hasMembers := audience["members"]; !hasMembers {
			audience["members"] = []any{}
		}
		if _, hasExclusions := audience["exclusions"]; !hasExclusions {
			audience["exclusions"] = []any{}
		}

		// Make copies to avoid modifying the original slices
		membersSlice := audience["members"].([]any)
		members := make([]any, len(membersSlice))
		copy(members, membersSlice)

		exclusionsSlice := audience["exclusions"].([]any)
		exclusions := make([]any, len(exclusionsSlice))
		copy(exclusions, exclusionsSlice)

		// Process addMembers
		if addMembers, ok := requestBody["addMembers"].([]any); ok {
			for _, member := range addMembers {
				memberMap := member.(map[string]any)
				// Ensure @odata.type is set
				if _, hasODataType := memberMap["@odata.type"]; !hasODataType {
					memberMap["@odata.type"] = "#microsoft.graph.windowsUpdates.azureADDevice"
				}
				members = append(members, memberMap)
			}
		}

		// Process removeMembers
		if removeMembers, ok := requestBody["removeMembers"].([]any); ok {
			for _, removeMember := range removeMembers {
				removeMemberMap := removeMember.(map[string]any)
				removeMemberID := removeMemberMap["id"].(string)
				newMembers := []any{}
				for _, existingMember := range members {
					existingMemberMap := existingMember.(map[string]any)
					if existingMemberMap["id"].(string) != removeMemberID {
						newMembers = append(newMembers, existingMember)
					}
				}
				members = newMembers
			}
		}

		// Process addExclusions
		if addExclusions, ok := requestBody["addExclusions"].([]any); ok {
			for _, exclusion := range addExclusions {
				exclusionMap := exclusion.(map[string]any)
				// Ensure @odata.type is set
				if _, hasODataType := exclusionMap["@odata.type"]; !hasODataType {
					exclusionMap["@odata.type"] = "#microsoft.graph.windowsUpdates.azureADDevice"
				}
				exclusions = append(exclusions, exclusionMap)
			}
		}

		// Process removeExclusions
		if removeExclusions, ok := requestBody["removeExclusions"].([]any); ok {
			for _, removeExclusion := range removeExclusions {
				removeExclusionMap := removeExclusion.(map[string]any)
				removeExclusionID := removeExclusionMap["id"].(string)
				newExclusions := []any{}
				for _, existingExclusion := range exclusions {
					existingExclusionMap := existingExclusion.(map[string]any)
					if existingExclusionMap["id"].(string) != removeExclusionID {
						newExclusions = append(newExclusions, existingExclusion)
					}
				}
				exclusions = newExclusions
			}
		}

		audience["members"] = members
		audience["exclusions"] = exclusions
		mockState.audiences[id] = audience
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	}
}

func (m *WindowsUpdateDeploymentAudienceMembersMock) getMembersResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deploymentAudiences/")

		mockState.Lock()
		audience, exists := mockState.audiences[id]
		if !exists {
			// Auto-create for unit tests
			audience = map[string]any{
				"@odata.type": "#microsoft.graph.windowsUpdates.deploymentAudience",
				"id":          id,
				"members":     []any{},
				"exclusions":  []any{},
			}
			mockState.audiences[id] = audience
		}

		members, hasMem := audience["members"]
		if !hasMem || members == nil {
			members = []any{}
		}
		mockState.Unlock()

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/windows/updates/deploymentAudiences('" + id + "')/members",
			"value":          members,
		}

		return httpmock.NewJsonResponse(200, response)
	}
}

func (m *WindowsUpdateDeploymentAudienceMembersMock) getExclusionsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deploymentAudiences/")

		mockState.Lock()
		audience, exists := mockState.audiences[id]
		if !exists {
			// Auto-create for unit tests
			audience = map[string]any{
				"@odata.type": "#microsoft.graph.windowsUpdates.deploymentAudience",
				"id":          id,
				"members":     []any{},
				"exclusions":  []any{},
			}
			mockState.audiences[id] = audience
		}

		exclusions, hasExcl := audience["exclusions"]
		if !hasExcl || exclusions == nil {
			exclusions = []any{}
		}
		mockState.Unlock()

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/windows/updates/deploymentAudiences('" + id + "')/exclusions",
			"value":          exclusions,
		}

		return httpmock.NewJsonResponse(200, response)
	}
}

func (m *WindowsUpdateDeploymentAudienceMembersMock) deleteAudienceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deploymentAudiences/")

		mockState.Lock()
		_, exists := mockState.audiences[id]
		if exists {
			delete(mockState.audiences, id)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	}
}

func (m *WindowsUpdateDeploymentAudienceMembersMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)/microsoft\.graph\.windowsUpdates\.updateAudience$`,
		factories.ErrorResponse(400, "BadRequest", "Invalid request"))

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/error-id$`,
		factories.ErrorResponse(404, "ResourceNotFound", "Resource not found"))
}

func (m *WindowsUpdateDeploymentAudienceMembersMock) CleanupMockState() {
	mockState.Lock()
	mockState.audiences = make(map[string]map[string]any)
	mockState.Unlock()
}
