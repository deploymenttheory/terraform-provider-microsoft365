package mocks

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"

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
