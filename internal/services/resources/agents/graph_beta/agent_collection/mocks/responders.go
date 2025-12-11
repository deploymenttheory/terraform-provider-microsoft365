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
	agentCollections map[string]map[string]any
}

func init() {
	mockState.agentCollections = make(map[string]map[string]any)
	mocks.GlobalRegistry.Register("agent_collection", &AgentCollectionMock{})
}

type AgentCollectionMock struct{}

var _ mocks.MockRegistrar = (*AgentCollectionMock)(nil)

func (m *AgentCollectionMock) RegisterMocks() {
	mockState.Lock()
	mockState.agentCollections = make(map[string]map[string]any)
	mockState.Unlock()

	// Create agent collection - POST /agentRegistry/agentCollections
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/agentRegistry/agentCollections", func(req *http.Request) (*http.Response, error) {
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

		// Select appropriate JSON file based on displayName
		var jsonFile string
		switch displayName {
		case "Unit Test Agent Collection Minimal":
			jsonFile = "../tests/responses/validate_create/post_agent_collection_minimal.json"
		case "Unit Test Agent Collection Maximal":
			jsonFile = "../tests/responses/validate_create/post_agent_collection_maximal.json"
		case "Unit Test Agent Collection Update Minimal":
			jsonFile = "../tests/responses/validate_create/post_agent_collection_update_minimal.json"
		case "Unit Test Agent Collection Update Maximal":
			jsonFile = "../tests/responses/validate_create/post_agent_collection_update_maximal.json"
		default:
			jsonFile = "../tests/responses/validate_create/post_agent_collection_minimal.json"
		}

		jsonContent, err := helpers.ParseJSONFile(jsonFile)
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load JSON response file"}}`), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		newId := uuid.New().String()
		responseObj["id"] = newId

		mockState.Lock()
		mockState.agentCollections[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// List agent collections - GET /agentRegistry/agentCollections
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/agentRegistry/agentCollections", func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		collections := make([]map[string]any, 0, len(mockState.agentCollections))
		for _, collection := range mockState.agentCollections {
			collections = append(collections, collection)
		}
		mockState.Unlock()

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#agentRegistry/agentCollections",
			"value":          collections,
		}

		return httpmock.NewJsonResponse(200, response)
	})

	// Get agent collection - GET /agentRegistry/agentCollections/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentCollectionId := parts[len(parts)-1]

		mockState.Lock()
		agentCollection, exists := mockState.agentCollections[agentCollectionId]
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		return httpmock.NewJsonResponse(200, agentCollection)
	})

	// Update agent collection - PATCH /agentRegistry/agentCollections/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentCollectionId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.agentCollections[agentCollectionId]
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

		// Select appropriate JSON file based on displayName
		displayName, _ := requestBody["displayName"].(string)
		var jsonFile string
		switch displayName {
		case "Unit Test Agent Collection Minimal":
			jsonFile = "../tests/responses/validate_update/patch_agent_collection_minimal.json"
		case "Unit Test Agent Collection Maximal":
			jsonFile = "../tests/responses/validate_update/patch_agent_collection_maximal.json"
		case "Unit Test Agent Collection Update Minimal":
			jsonFile = "../tests/responses/validate_update/patch_agent_collection_update_minimal.json"
		case "Unit Test Agent Collection Update Maximal":
			jsonFile = "../tests/responses/validate_update/patch_agent_collection_update_maximal.json"
		default:
			jsonFile = "../tests/responses/validate_update/patch_agent_collection_minimal.json"
		}

		jsonContent, err := helpers.ParseJSONFile(jsonFile)
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load JSON response file"}}`), nil
		}

		var updatedCollection map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &updatedCollection); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		updatedCollection["id"] = agentCollectionId

		mockState.Lock()
		mockState.agentCollections[agentCollectionId] = updatedCollection
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete agent collection - DELETE /agentRegistry/agentCollections/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		agentCollectionId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.agentCollections[agentCollectionId]
		if exists {
			delete(mockState.agentCollections, agentCollectionId)
		}
		mockState.Unlock()

		if !exists {
			errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
			return httpmock.NewStringResponse(404, errorResp), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AgentCollectionMock) RegisterErrorMocks() {
	errorBadRequest, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
	errorNotFound, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/agentRegistry/agentCollections",
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/agentRegistry/agentCollections",
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		httpmock.NewStringResponder(400, errorBadRequest))
}

func (m *AgentCollectionMock) CleanupMockState() {
	mockState.Lock()
	mockState.agentCollections = make(map[string]map[string]any)
	mockState.Unlock()
}
