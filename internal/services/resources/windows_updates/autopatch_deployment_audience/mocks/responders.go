package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// extractAudienceID extracts the deployment audience ID from URL paths that may have
// sub-resource segments (e.g., /members, /exclusions, /microsoft.graph...).
func extractAudienceID(urlPath string) string {
	const marker = "deploymentAudiences/"
	idx := strings.Index(urlPath, marker)
	if idx < 0 {
		return ""
	}
	rest := urlPath[idx+len(marker):]
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return ""
}

type WindowsUpdateDeploymentAudienceMock struct{}

var (
	mockState = struct {
		sync.Mutex
		audiences map[string]map[string]any
	}{
		audiences: make(map[string]map[string]any),
	}
)

func (m *WindowsUpdateDeploymentAudienceMock) RegisterMocks() {
	m.registerCreateAudienceResponder()
	m.registerGetAudienceResponder()
	m.registerGetMembersResponder()
	m.registerGetExclusionsResponder()
	m.registerUpdateAudienceResponder()
	m.registerDeleteAudienceResponder()
}

func (m *WindowsUpdateDeploymentAudienceMock) registerCreateAudienceResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences$`,
		m.createAudienceResponder())
}

func (m *WindowsUpdateDeploymentAudienceMock) registerGetAudienceResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)$`,
		m.getAudienceResponder())
}

func (m *WindowsUpdateDeploymentAudienceMock) registerGetMembersResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)/members$`,
		m.getMembersResponder())
}

func (m *WindowsUpdateDeploymentAudienceMock) registerGetExclusionsResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)/exclusions$`,
		m.getExclusionsResponder())
}

func (m *WindowsUpdateDeploymentAudienceMock) registerUpdateAudienceResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)/microsoft\.graph\.windowsUpdates\.updateAudience$`,
		m.updateAudienceResponder())
}

func (m *WindowsUpdateDeploymentAudienceMock) registerDeleteAudienceResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)$`,
		m.deleteAudienceResponder())
}

func (m *WindowsUpdateDeploymentAudienceMock) createAudienceResponder() httpmock.Responder {
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

func (m *WindowsUpdateDeploymentAudienceMock) getAudienceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deploymentAudiences/")

		mockState.Lock()
		audience, exists := mockState.audiences[id]
		mockState.Unlock()

		if !exists {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/error_handling/not_found.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		return httpmock.NewJsonResponse(200, audience)
	}
}

func (m *WindowsUpdateDeploymentAudienceMock) getMembersResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := extractAudienceID(req.URL.Path)

		mockState.Lock()
		audience, exists := mockState.audiences[id]
		if !exists {
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

func (m *WindowsUpdateDeploymentAudienceMock) getExclusionsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := extractAudienceID(req.URL.Path)

		mockState.Lock()
		audience, exists := mockState.audiences[id]
		if !exists {
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

func (m *WindowsUpdateDeploymentAudienceMock) updateAudienceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := extractAudienceID(req.URL.Path)

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		mockState.Lock()
		audience, exists := mockState.audiences[id]
		if !exists {
			audience = map[string]any{
				"@odata.type": "#microsoft.graph.windowsUpdates.deploymentAudience",
				"id":          id,
				"members":     []any{},
				"exclusions":  []any{},
			}
			mockState.audiences[id] = audience
		}

		if _, hasMembers := audience["members"]; !hasMembers {
			audience["members"] = []any{}
		}
		if _, hasExclusions := audience["exclusions"]; !hasExclusions {
			audience["exclusions"] = []any{}
		}

		membersSlice := audience["members"].([]any)
		members := make([]any, len(membersSlice))
		copy(members, membersSlice)

		exclusionsSlice := audience["exclusions"].([]any)
		exclusions := make([]any, len(exclusionsSlice))
		copy(exclusions, exclusionsSlice)

		if addMembers, ok := requestBody["addMembers"].([]any); ok {
			for _, member := range addMembers {
				memberMap := member.(map[string]any)
				if _, hasODataType := memberMap["@odata.type"]; !hasODataType {
					memberMap["@odata.type"] = "#microsoft.graph.windowsUpdates.azureADDevice"
				}
				members = append(members, memberMap)
			}
		}

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

		if addExclusions, ok := requestBody["addExclusions"].([]any); ok {
			for _, exclusion := range addExclusions {
				exclusionMap := exclusion.(map[string]any)
				if _, hasODataType := exclusionMap["@odata.type"]; !hasODataType {
					exclusionMap["@odata.type"] = "#microsoft.graph.windowsUpdates.azureADDevice"
				}
				exclusions = append(exclusions, exclusionMap)
			}
		}

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

func (m *WindowsUpdateDeploymentAudienceMock) deleteAudienceResponder() httpmock.Responder {
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

func (m *WindowsUpdateDeploymentAudienceMock) RegisterErrorMocks() {
	m.registerCreateAudienceErrorResponder()
	m.registerGetAudienceErrorResponder()
	m.registerDeleteAudienceErrorResponder()
}

func (m *WindowsUpdateDeploymentAudienceMock) registerCreateAudienceErrorResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences$`,
		factories.ErrorResponse(400, "BadRequest", "Invalid request"))
}

func (m *WindowsUpdateDeploymentAudienceMock) registerGetAudienceErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/error-id$`,
		factories.ErrorResponse(404, "ResourceNotFound", "Resource not found"))
}

func (m *WindowsUpdateDeploymentAudienceMock) registerDeleteAudienceErrorResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Audience is in use"))
}

func (m *WindowsUpdateDeploymentAudienceMock) CleanupMockState() {
	mockState.Lock()
	mockState.audiences = make(map[string]map[string]any)
	mockState.Unlock()
}
