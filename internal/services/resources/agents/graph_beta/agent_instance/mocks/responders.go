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
	agentInstances map[string]map[string]any
}

func init() {
	mockState.agentInstances = make(map[string]map[string]any)
	mocks.GlobalRegistry.Register("agent_instance", &AgentInstanceMock{})
}

type AgentInstanceMock struct{}

var _ mocks.MockRegistrar = (*AgentInstanceMock)(nil)

func (m *AgentInstanceMock) RegisterMocks() {
	mockState.Lock()
	mockState.agentInstances = make(map[string]map[string]any)
	mockState.Unlock()

	// Create agent instance - POST /agentRegistry/agentInstances
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/agentRegistry/agentInstances", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		displayName, ok := requestBody["displayName"].(string)
		if !ok || displayName == "" {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_agent_instance_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		newId := uuid.New().String()
		responseObj["id"] = newId
		responseObj["displayName"] = displayName
		responseObj["@odata.type"] = "#microsoft.graph.agentInstance"

		if ownerIds, ok := requestBody["ownerIds"].([]any); ok {
			ownerIdStrings := make([]string, len(ownerIds))
			for i, id := range ownerIds {
				if idStr, ok := id.(string); ok {
					ownerIdStrings[i] = idStr
				}
			}
			responseObj["ownerIds"] = ownerIdStrings
		}

		if originatingStore, ok := requestBody["originatingStore"].(string); ok {
			responseObj["originatingStore"] = originatingStore
		}

		if managedBy, ok := requestBody["managedBy"].(string); ok {
			responseObj["managedBy"] = managedBy
		}

		if url, ok := requestBody["url"].(string); ok {
			responseObj["url"] = url
		}

		if preferredTransport, ok := requestBody["preferredTransport"].(string); ok {
			responseObj["preferredTransport"] = preferredTransport
		}

		if additionalInterfaces, ok := requestBody["additionalInterfaces"].([]any); ok {
			responseObj["additionalInterfaces"] = additionalInterfaces
		}

		if signatures, ok := requestBody["signatures"].([]any); ok {
			responseObj["signatures"] = signatures
		}

		if agentCardManifest, ok := requestBody["agentCardManifest"].(map[string]any); ok {
			// Add generated id for the manifest if not present
			if _, hasId := agentCardManifest["id"]; !hasId {
				agentCardManifest["id"] = uuid.New().String()
			}
			responseObj["agentCardManifest"] = agentCardManifest
		}

		mockState.Lock()
		mockState.agentInstances[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get agent instance - GET /agentRegistry/agentInstances/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentInstances/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentInstanceId := parts[len(parts)-1]

		mockState.Lock()
		agentInstance, exists := mockState.agentInstances[agentInstanceId]
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		return httpmock.NewJsonResponse(200, agentInstance)
	})

	// Update agent instance - PATCH /agentRegistry/agentInstances/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentInstances/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentInstanceId := parts[len(parts)-1]

		mockState.Lock()
		agentInstance, exists := mockState.agentInstances[agentInstanceId]
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		for key, value := range requestBody {
			agentInstance[key] = value
		}

		agentInstance["id"] = agentInstanceId
		agentInstance["@odata.type"] = "#microsoft.graph.agentInstance"

		mockState.Lock()
		mockState.agentInstances[agentInstanceId] = agentInstance
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, agentInstance)
	})

	// Delete agent instance - DELETE /agentRegistry/agentInstances/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentInstances/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentInstanceId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.agentInstances[agentInstanceId]
		if exists {
			delete(mockState.agentInstances, agentInstanceId)
		}
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Get agent card manifest - GET /agentRegistry/agentCardManifests/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCardManifests/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		manifestId := parts[len(parts)-1]

		// Find the manifest by ID in stored agent instances
		mockState.Lock()
		defer mockState.Unlock()

		for _, agentInstance := range mockState.agentInstances {
			if manifest, ok := agentInstance["agentCardManifest"].(map[string]any); ok {
				if id, ok := manifest["id"].(string); ok && id == manifestId {
					manifest["@odata.type"] = "#microsoft.graph.agentCardManifest"
					return httpmock.NewJsonResponse(200, manifest)
				}
			}
		}

		errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
		return httpmock.NewStringResponse(404, errorResp), nil
	})

	// Get agent card manifest via agent instance - GET /agentRegistry/agentInstances/{id}/agentCardManifest
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentInstances/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/agentCardManifest$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		// Path is: /beta/agentRegistry/agentInstances/{id}/agentCardManifest
		agentInstanceId := parts[len(parts)-2]

		mockState.Lock()
		agentInstance, exists := mockState.agentInstances[agentInstanceId]
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		if manifest, ok := agentInstance["agentCardManifest"].(map[string]any); ok {
			manifest["@odata.type"] = "#microsoft.graph.agentCardManifest"
			return httpmock.NewJsonResponse(200, manifest)
		}

		errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
		return httpmock.NewStringResponse(404, errorResp), nil
	})

	// Update agent card manifest - PATCH /agentRegistry/agentCardManifests/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCardManifests/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		manifestId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
			return httpmock.NewStringResponse(400, errorResp), nil
		}

		// Find and update the manifest by ID in stored agent instances
		mockState.Lock()
		defer mockState.Unlock()

		for _, agentInstance := range mockState.agentInstances {
			if manifest, ok := agentInstance["agentCardManifest"].(map[string]any); ok {
				if id, ok := manifest["id"].(string); ok && id == manifestId {
					// Update manifest properties
					for key, value := range requestBody {
						manifest[key] = value
					}
					manifest["id"] = manifestId
					manifest["@odata.type"] = "#microsoft.graph.agentCardManifest"
					return httpmock.NewJsonResponse(200, manifest)
				}
			}
		}

		errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
		return httpmock.NewStringResponse(404, errorResp), nil
	})
}

func (m *AgentInstanceMock) RegisterErrorMocks() {
	errorBadRequest, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
	errorNotFound, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/agentRegistry/agentInstances",
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentInstances/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentInstances/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentInstances/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCardManifests/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCardManifests/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
}

func (m *AgentInstanceMock) CleanupMockState() {
	mockState.Lock()
	mockState.agentInstances = make(map[string]map[string]any)
	mockState.Unlock()
}
