package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	subscribedSkus map[string]map[string]any
}

func init() {
	mockState.subscribedSkus = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("subscribed_skus", &SubscribedSkusMock{})
}

type SubscribedSkusMock struct{}

var _ mocks.MockRegistrar = (*SubscribedSkusMock)(nil)

func (m *SubscribedSkusMock) RegisterMocks() {
	mockState.Lock()
	mockState.subscribedSkus = make(map[string]map[string]any)
	mockState.Unlock()

	RegisterGetByIdMock()
	RegisterListMock()
}

func (m *SubscribedSkusMock) RegisterErrorMocks() {
	RegisterErrorMocks()
}

func (m *SubscribedSkusMock) CleanupMockState() {
	mockState.Lock()
	mockState.subscribedSkus = make(map[string]map[string]any)
	mockState.Unlock()
}

func RegisterGetByIdMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/v1\.0/subscribedSkus/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			skuId := parts[len(parts)-1]

			switch skuId {
			case "48a80680-7326-48cd-9935-b556b81d3a4e_c7df2760-2c81-4ef7-b578-5b5392b571df":
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_subscribed_sku_by_id.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			default:
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Subscribed SKU not found"}}`), nil
			}
		})
}

func RegisterListMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/v1\.0/subscribedSkus`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/list_subscribed_skus.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		})
}

func RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/v1\.0/subscribedSkus`,
		httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Internal server error"}}`))
}
