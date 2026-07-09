package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	privateNetworks map[string]map[string]any
}

func init() {
	mockState.privateNetworks = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("network_private_network", &PrivateNetworkMock{})
}

type PrivateNetworkMock struct{}

var _ mocks.MockRegistrar = (*PrivateNetworkMock)(nil)

func (m *PrivateNetworkMock) RegisterMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkaccess/privateNetworks",
		m.createPrivateNetworkResponder())
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkaccess/privateNetworks/([^/]+)$`,
		m.getPrivateNetworkResponder())
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/networkaccess/privateNetworks/([^/]+)$`,
		m.updatePrivateNetworkResponder())
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/networkaccess/privateNetworks/([^/]+)$`,
		m.deletePrivateNetworkResponder())
}

func (m *PrivateNetworkMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkaccess/privateNetworks",
		func(req *http.Request) (*http.Response, error) {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_private_network_error.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			return httpmock.NewStringResponse(400, jsonContent), nil
		})
}

func (m *PrivateNetworkMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.privateNetworks = make(map[string]map[string]any)
}

func (m *PrivateNetworkMock) createPrivateNetworkResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_private_network.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}
		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
		}

		id := "00000000-0000-0000-0000-000000000101"
		response["id"] = id
		mergePrivateNetworkRequest(response, requestBody)

		mockState.Lock()
		mockState.privateNetworks[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

func (m *PrivateNetworkMock) getPrivateNetworkResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkaccess/privateNetworks/")

		mockState.Lock()
		privateNetwork, exists := mockState.privateNetworks[id]
		mockState.Unlock()

		if !exists {
			switch {
			case strings.Contains(id, "minimal"):
				return jsonFixtureResponse(req, 200, filepath.Join("..", "tests", "responses", "validate_create", "get_private_network_minimal.json"))
			case strings.Contains(id, "maximal"):
				return jsonFixtureResponse(req, 200, filepath.Join("..", "tests", "responses", "validate_create", "get_private_network_maximal.json"))
			default:
				return jsonFixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_private_network_not_found.json"))
			}
		}

		privateNetworkCopy := make(map[string]any)
		for key, value := range privateNetwork {
			privateNetworkCopy[key] = value
		}
		if _, exists := privateNetworkCopy["@odata.context"]; !exists {
			privateNetworkCopy["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#networkAccess/privateNetworks/$entity"
		}

		return factories.SuccessResponse(200, privateNetworkCopy)(req)
	}
}

func (m *PrivateNetworkMock) updatePrivateNetworkResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkaccess/privateNetworks/")

		mockState.Lock()
		privateNetwork, exists := mockState.privateNetworks[id]
		mockState.Unlock()
		if !exists {
			return jsonFixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_private_network_not_found.json"))
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		updated := make(map[string]any)
		for key, value := range privateNetwork {
			updated[key] = value
		}
		mergePrivateNetworkRequest(updated, requestBody)

		mockState.Lock()
		mockState.privateNetworks[id] = updated
		mockState.Unlock()

		return factories.SuccessResponse(200, updated)(req)
	}
}

func (m *PrivateNetworkMock) deletePrivateNetworkResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkaccess/privateNetworks/")

		mockState.Lock()
		_, exists := mockState.privateNetworks[id]
		if exists {
			delete(mockState.privateNetworks, id)
		}
		mockState.Unlock()

		if !exists {
			return jsonFixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_private_network_not_found.json"))
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

func mergePrivateNetworkRequest(response map[string]any, requestBody map[string]any) {
	response["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#networkAccess/privateNetworks/$entity"
	for _, key := range []string{"name", "appIds", "networkIdentifications"} {
		if value, ok := requestBody[key]; ok {
			response[key] = value
		}
	}
}

func jsonFixtureResponse(req *http.Request, status int, path string) (*http.Response, error) {
	jsonContent, err := helpers.ParseJSONFile(path)
	if err != nil {
		return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
	}

	return factories.SuccessResponse(status, response)(req)
}
