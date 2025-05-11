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
