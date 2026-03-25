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
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("directory_role", &DirectoryRoleMock{})
}

type DirectoryRoleMock struct{}

var _ mocks.MockRegistrar = (*DirectoryRoleMock)(nil)

func (m *DirectoryRoleMock) RegisterMocks() {
	RegisterGetByIDMock()
	RegisterListMock()
}

func (m *DirectoryRoleMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directoryRoles`,
		httpmock.NewStringResponder(403, `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation"}}`))
}

func (m *DirectoryRoleMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
}

func RegisterGetByIDMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directoryRoles/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			roleID := parts[len(parts)-1]

			switch roleID {
			case "aaaaaaaa-0001-0000-0000-000000000000":
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_directory_role_by_id.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			default:
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Directory role not found"}}`), nil
			}
		})
}

func RegisterListMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directoryRoles$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_directory_roles_all.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		})
}
