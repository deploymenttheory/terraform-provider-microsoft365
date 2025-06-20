package mocks

import (
	"context"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jarcoal/httpmock"
	abstractions "github.com/microsoft/kiota-abstractions-go"
)

// MockAuthProvider implements the required authentication interface
type MockAuthProvider struct{}

// AuthenticateRequest adds a mock authorization header to requests
func (m *MockAuthProvider) AuthenticateRequest(ctx context.Context, request *abstractions.RequestInformation, additionalAuthenticationContext map[string]interface{}) error {
	if request.Headers != nil {
		request.Headers.Add("Authorization", "Bearer mock-token")
	}
	return nil
}

// Mocks provides a centralized way to manage mock HTTP responses for testing.
type Mocks struct {
	authMocks *AuthenticationMocks
	Clients   *client.MockGraphClients
}

// NewMocks creates a new instance of Mocks, initializing all mock types.
func NewMocks() *Mocks {
	// Create a new HTTP client that will be used with httpmock
	httpClient := &http.Client{}

	return &Mocks{
		authMocks: NewAuthenticationMocks(),
		Clients:   client.NewMockGraphClients(httpClient),
	}
}

// GetMockClients returns the mock clients for use in tests
func (m *Mocks) GetMockClients() client.GraphClientInterface {
	return m.Clients
}

// Activate activates all mock responders.
func (m *Mocks) Activate() {
	httpmock.Activate()
	// Configure httpmock to use the same client that our mock clients use
	httpmock.ActivateNonDefault(http.DefaultClient)
	m.authMocks.RegisterMocks()
	m.RegisterMacOSPlatformScriptMocks()
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

var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider("test", false)()),
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider("test", true)()),
}
