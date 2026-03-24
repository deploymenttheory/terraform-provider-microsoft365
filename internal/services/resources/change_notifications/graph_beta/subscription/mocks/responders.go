package mocks

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	commonMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

var mockState struct {
	sync.Mutex
	byID map[string]map[string]any
}

func init() {
	mockState.byID = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(
		httpmock.NewStringResponder(
			404,
			`{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`,
		),
	)
	commonMocks.GlobalRegistry.Register("change_notifications_subscription", &SubscriptionMock{})
}

type SubscriptionMock struct{}

var _ commonMocks.MockRegistrar = (*SubscriptionMock)(nil)

func (m *SubscriptionMock) loadJSONResponse(filePath string) (map[string]any, error) {
	fullPath := filepath.Join("..", "tests", "responses", filePath)
	jsonContent, err := helpers.ParseJSONFile(fullPath)
	if err != nil {
		return nil, err
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (m *SubscriptionMock) RegisterMocks() {
	mockState.Lock()
	mockState.byID = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/subscriptions$`,
		func(req *http.Request) (*http.Response, error) {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			var payload map[string]any
			if err := json.Unmarshal(bodyBytes, &payload); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest"}}`), nil
			}

			if res, err := m.loadJSONResponse(
				"validate_create/post_subscription_error.json",
			); err == nil {
				if nurl, ok := payload["notificationUrl"].(string); ok &&
					strings.Contains(nurl, "invalid-webhook") {
					return httpmock.NewJsonResponse(400, res)
				}
			}

			base, err := m.loadJSONResponse("validate_create/post_subscription_success.json")
			if err != nil {
				return nil, err
			}
			for k, v := range payload {
				base[k] = v
			}
			mockState.Lock()
			id := base["id"].(string)
			clone := cloneMap(base)
			mockState.byID[id] = clone
			mockState.Unlock()
			return httpmock.NewJsonResponse(201, base)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/subscriptions/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			id := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
			mockState.Lock()
			data, ok := mockState.byID[id]
			mockState.Unlock()
			if !ok {
				if nf, err := m.loadJSONResponse(
					"validate_delete/get_subscription_not_found.json",
				); err == nil {
					return httpmock.NewJsonResponse(404, nf)
				}
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound"}}`), nil
			}
			return httpmock.NewJsonResponse(200, data)
		})

	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/subscriptions/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			id := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			var patch map[string]any
			if err := json.Unmarshal(bodyBytes, &patch); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest"}}`), nil
			}
			mockState.Lock()
			data, ok := mockState.byID[id]
			if ok {
				for k, v := range patch {
					data[k] = v
				}
				mockState.byID[id] = data
			}
			mockState.Unlock()
			if !ok {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound"}}`), nil
			}
			return httpmock.NewJsonResponse(200, data)
		})

	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/subscriptions/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			id := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
			mockState.Lock()
			delete(mockState.byID, id)
			mockState.Unlock()
			return httpmock.NewBytesResponse(204, nil), nil
		})
}

func cloneMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func (m *SubscriptionMock) RegisterErrorMocks() {}

func (m *SubscriptionMock) CleanupMockState() {
	mockState.Lock()
	mockState.byID = make(map[string]map[string]any)
	mockState.Unlock()
}
