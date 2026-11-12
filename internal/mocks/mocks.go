package mocks

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	providercore "github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider/core"
)

// NewProviderFactories creates protocol factories with only the supplied test registrations.
// Uses the real provider schema and configuration in unit or acceptance mode.
func NewProviderFactories(
	unitTestMode bool,
	registry providercore.Registry,
) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"microsoft365": func() (tfprotov6.ProviderServer, error) {
			server, err := providerserver.NewProtocol6WithError(
				providercore.NewWithRegistry("test", unitTestMode, registry),
			)()
			if err != nil {
				return nil, fmt.Errorf("failed to create isolated provider server: %w", err)
			}
			return server, nil
		},
	}
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
