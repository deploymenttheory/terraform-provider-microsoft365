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
	users map[string]map[string]any
}

func init() {
	mockState.users = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("user", &UserMock{})
}

type UserMock struct{}

var _ mocks.MockRegistrar = (*UserMock)(nil)

func (m *UserMock) RegisterMocks() {
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
	mockState.Unlock()

	RegisterGetByObjectIdMock()
	RegisterListAndFilterMocks()
}

func (m *UserMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users`,
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges to complete the operation"}}`))
}

func (m *UserMock) CleanupMockState() {
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
	mockState.Unlock()
}

func RegisterGetByObjectIdMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			objectId := parts[len(parts)-1]

			switch objectId {
			case "11111111-1111-1111-1111-111111111111":
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_user_by_object_id.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			default:
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}
		})
}

func RegisterListAndFilterMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users(\?.*)?$`,
		func(req *http.Request) (*http.Response, error) {
			queryParams := req.URL.Query()
			filter := queryParams.Get("$filter")

			// List all users
			if filter == "" {
				return jsonFileResponse("../tests/responses/validate_get/get_users_all.json")
			}

			switch {
			case strings.Contains(filter, "displayName eq") && strings.Contains(filter, "DT-TEST-USER-001"):
				return jsonFileResponse("../tests/responses/validate_get/get_user_by_display_name.json")
			case strings.Contains(filter, "employeeId eq") && strings.Contains(filter, "EMP-0001"):
				return jsonFileResponse("../tests/responses/validate_get/get_user_by_employee_id.json")
			case strings.Contains(filter, "givenName eq") && strings.Contains(filter, "Test"):
				return jsonFileResponse("../tests/responses/validate_get/get_user_by_given_name.json")
			case strings.Contains(filter, "userPrincipalName eq") && strings.Contains(filter, "dt.test.user001@contoso.com"):
				return jsonFileResponse("../tests/responses/validate_get/get_user_by_upn.json")
			case strings.Contains(filter, "onPremisesImmutableId eq") && strings.Contains(filter, "IMMUTABLE-0001"):
				return jsonFileResponse("../tests/responses/validate_get/get_user_by_on_premises_immutable_id.json")
			case strings.Contains(filter, "onPremisesDistinguishedName eq"):
				return jsonFileResponse("../tests/responses/validate_get/get_user_by_on_premises_distinguished_name.json")
			case strings.Contains(filter, "accountEnabled eq true") && strings.Contains(filter, "userType eq 'Member'"):
				return jsonFileResponse("../tests/responses/validate_get/get_users_odata_filter.json")
			}

			// Default empty response for unmatched filters
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#users",
				"value":          []any{},
			})
		})
}

func jsonFileResponse(path string) (*http.Response, error) {
	jsonStr, _ := helpers.ParseJSONFile(path)
	var responseObj map[string]any
	json.Unmarshal([]byte(jsonStr), &responseObj)
	return httpmock.NewJsonResponse(200, responseObj)
}
