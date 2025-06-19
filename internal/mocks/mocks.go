package mocks

import (
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jarcoal/httpmock"
)

// Mocks provides a centralized way to manage mock HTTP responses for testing.
type Mocks struct {
	authMocks *AuthenticationMocks
}

// NewMocks creates a new instance of Mocks, initializing all mock types.
func NewMocks() *Mocks {
	return &Mocks{
		authMocks: NewAuthenticationMocks(),
	}
}

// Activate activates all mock responders.
func (m *Mocks) Activate() {
	httpmock.Activate()
	m.authMocks.RegisterMocks()
	m.registerGraphMocks()
}

// DeactivateAndReset deactivates all mock responders and resets the mock state.
func (m *Mocks) DeactivateAndReset() {
	httpmock.DeactivateAndReset()
}

func TestName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

func TestsEntraLicesingGroupName() string {
	return "pptestusers"
}

var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider("test", true)()),
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider("test", false)()),
}

func (m *Mocks) registerGraphMocks() {
	// Mock for fetching user details
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/v1.0/users/testuser@example.com",
		func(req *http.Request) (*http.Response, error) {
			user := map[string]interface{}{
				"id":                "mock-user-id",
				"displayName":       "Test User",
				"userPrincipalName": "testuser@example.com",
			}
			return httpmock.NewJsonResponse(200, user)
		},
	)

	// Mock for group details
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/v1.0/groups/mock-group-id",
		func(req *http.Request) (*http.Response, error) {
			group := map[string]interface{}{
				"id":          "mock-group-id",
				"displayName": "Test Group",
			}
			return httpmock.NewJsonResponse(200, group)
		},
	)

	// Mock for beta endpoints, for example
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/me",
		func(req *http.Request) (*http.Response, error) {
			me := map[string]interface{}{
				"id":                "mock-user-id",
				"displayName":       "Test User (Beta)",
				"userPrincipalName": "testuser@example.com",
			}
			return httpmock.NewJsonResponse(200, me)
		},
	)
}
