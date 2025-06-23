package mocks

import (
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestUnitTestProtoV6ProviderFactories provides a map of provider factories specifically
// configured for unit testing scenarios. The Microsoft365 provider is instantiated with
// testMode=true, which enables test-specific behavior such as:
//   - Using mock HTTP clients instead of real Microsoft Graph API calls
//   - Bypassing authentication requirements
//   - Using in-memory or stubbed data sources
//   - Faster execution without network dependencies
//
// This factory is intended for isolated unit tests that focus on testing provider logic,
// resource schemas, validation, and state management without making actual API calls to
// Microsoft365 services.
var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider("test", true)()),
}

// TestAccProtoV6ProviderFactories provides a map of provider factories specifically
// configured for acceptance testing scenarios. The Microsoft365 provider is instantiated with
// testMode=false, which enables production-like behavior such as:
//   - Making real HTTP calls to Microsoft Graph API endpoints
//   - Requiring valid authentication credentials
//   - Interacting with actual Microsoft365 tenant resources
//   - Full end-to-end integration testing
//
// This factory is intended for acceptance tests that verify the provider's ability to
// correctly manage real Microsoft365 resources through the Graph API. These tests
// require valid credentials and may create, modify, or delete actual resources in a
// test tenant.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider("test", false)()),
}

// Mocks provides a centralized way to manage mock HTTP responses for testing.
type Mocks struct {
	AuthMocks *AuthenticationMocks
	Clients   *client.MockGraphClients
}

// NewMocks creates a new instance of Mocks, initializing shared mock
// and set up the mock client with the default http client.
func NewMocks() *Mocks {
	return &Mocks{
		AuthMocks: NewAuthenticationMocks(),
		Clients:   client.NewMockGraphClients(http.DefaultClient),
	}
}

// TestName returns the name of the function that called it.
func TestName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}
