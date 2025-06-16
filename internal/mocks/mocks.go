package mocks

import (
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jarcoal/httpmock"
)

// TestName extracts the name of the calling test function.
func TestName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

// TestUnitTestProtoV6ProviderFactories provides a map of provider factories for unit tests.
var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// TestAccProtoV6ProviderFactories provides a map of provider factories for acceptance tests.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// ActivateMicrosoftGraphMocks activates all Microsoft Graph API mocks by domain.
func ActivateMicrosoftGraphMocks() {
	// Activate authentication and common endpoint mocks
	activateAuthenticationMocks()

	// Activate domain-specific mocks
	ActivateDeviceManagementMocks()
	ActivateDeviceAndAppManagementMocks()
	ActivateIdentityAndAccessMocks()
	ActivateM365AdminMocks()

	// Activate specific resource mocks
	ActivateDeviceShellScriptMocks()

	// Add more domain activations as needed
}

// MockMicrosoftGraphRequest is a helper function to register a custom mock for a specific endpoint
func MockMicrosoftGraphRequest(method, urlPattern string, statusCode int, responseBody string) {
	httpmock.RegisterResponder(method, urlPattern,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(statusCode, responseBody), nil
		})
}

// MockMicrosoftGraphRequestWithRegexp is a helper function to register a custom mock with a regexp pattern
func MockMicrosoftGraphRequestWithRegexp(method string, urlRegexp *regexp.Regexp, statusCode int, responseBody string) {
	httpmock.RegisterRegexpResponder(method, urlRegexp,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(statusCode, responseBody), nil
		})
}

// ProviderConfigMinimal returns a minimal valid provider config for unit tests.
func ProviderConfigMinimal() string {
	return `
provider "microsoft365" {
  cloud = "public"
  tenant_id = "5298c688-49cd-4aaa-9978-77aeb00e1000"
  auth_method = "client_secret"
  entra_id_options = {
    client_id = "5298c688-49cd-4aaa-9978-77aeb00e1000"
    client_secret = "5298c688-49cd-4aaa-9978-77aeb00e1000"
  }
  telemetry_optout = true
  debug_mode = false
}
`
}

// ProviderConfigMaximal returns a maximal valid provider config for provider-level or edge-case tests.
func ProviderConfigMaximal() string {
	return `
provider "microsoft365" {
  cloud = "public"
  tenant_id = "00000000-0000-0000-0000-000000000000"
  auth_method = "client_secret"
  entra_id_options = {
    client_id = "fake-client-id"
    client_secret = "fake-client-secret"
    client_certificate = "fake-cert"
    client_certificate_password = "fake-password"
    send_certificate_chain = false
    username = "fake-user"
    disable_instance_discovery = false
    additionally_allowed_tenants = ["*"]
    redirect_url = "http://localhost"
    federated_token_file_path = "/tmp/fake-token"
    managed_identity_id = "fake-mi-id"
    oidc_token_file_path = "/tmp/fake-oidc-token"
    ado_service_connection_id = "fake-ado-id"
  }
  client_options = {
    enable_headers_inspection = true
    enable_retry = true
    max_retries = 5
    retry_delay_seconds = 2
    enable_redirect = true
    max_redirects = 3
    enable_compression = true
    custom_user_agent = "test-agent"
    use_proxy = true
    proxy_url = "http://localhost:8888"
    proxy_username = "proxy-user"
    proxy_password = "proxy-pass"
    timeout_seconds = 60
    enable_chaos = true
    chaos_percentage = 50
    chaos_status_code = 500
    chaos_status_message = "Internal Server Error"
  }
  telemetry_optout = true
  debug_mode = true
}
`
}
