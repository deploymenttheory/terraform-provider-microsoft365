package mocks

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

// ProviderMock provides mock responses for provider authentication testing
type ProviderMock struct{}

// RegisterMocks registers authentication mock responses for various auth methods
func (p *ProviderMock) RegisterMocks() {
	// Mock Azure AD token endpoint for client_secret auth
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000001/oauth2/v2.0/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"access_token": "mock-access-token-client-secret",
			"token_type":   "Bearer",
			"expires_in":   3600,
			"scope":        "https://graph.microsoft.com/.default",
		}))

	// Mock device code flow endpoint
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/common/oauth2/v2.0/devicecode",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"device_code":      "device_code_mock",
			"user_code":        "ABC123",
			"verification_uri": "https://microsoft.com/devicelogin",
			"expires_in":       900,
			"interval":         5,
		}))

	// Mock instance discovery for all clouds
	httpmock.RegisterResponder("GET",
		"https://login.microsoftonline.com/common/discovery/instance",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"tenant_discovery_endpoint": "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration",
		}))

	// Mock OIDC configuration endpoint
	httpmock.RegisterResponder("GET",
		"https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"issuer":                   "https://login.microsoftonline.com/{tenantid}/v2.0",
			"authorization_endpoint":   "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			"token_endpoint":          "https://login.microsoftonline.com/common/oauth2/v2.0/token",
			"device_authorization_endpoint": "https://login.microsoftonline.com/common/oauth2/v2.0/devicecode",
		}))

	// Mock managed identity endpoint (for Azure VMs)
	httpmock.RegisterResponder("GET",
		"http://169.254.169.254/metadata/identity/oauth2/token",
		httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
			"access_token": "mock-managed-identity-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		}))

	// Mock workload identity token file content
	httpmock.RegisterResponder("GET",
		"file:///var/run/secrets/azure/tokens/azure-identity-token",
		httpmock.NewStringResponder(200, "mock-workload-identity-jwt-token"))
}

// RegisterErrorMocks registers authentication error responses for testing failure scenarios
func (p *ProviderMock) RegisterErrorMocks() {
	// Mock authentication failures
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/00000000-0000-0000-0000-000000000001/oauth2/v2.0/token",
		httpmock.NewJsonResponderOrPanic(401, map[string]interface{}{
			"error":             "invalid_client",
			"error_description": "AADSTS7000215: Invalid client secret is provided.",
		}))

	// Mock invalid tenant ID
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/invalid-tenant-id/oauth2/v2.0/token",
		httpmock.NewJsonResponderOrPanic(400, map[string]interface{}{
			"error":             "invalid_request",
			"error_description": "AADSTS90002: Tenant 'invalid-tenant-id' not found.",
		}))

	// Mock certificate authentication failure
	httpmock.RegisterResponder("POST",
		"https://login.microsoftonline.com/common/oauth2/v2.0/token",
		httpmock.NewJsonResponderOrPanic(401, map[string]interface{}{
			"error":             "invalid_client",
			"error_description": "AADSTS700027: Client assertion contains an invalid signature.",
		}))
}

// RegisterCloudSpecificMocks registers cloud-specific mock endpoints
func (p *ProviderMock) RegisterCloudSpecificMocks() {
	clouds := map[string]string{
		"gcc":     "https://login.microsoftonline.us",
		"gcchigh": "https://login.microsoftonline.us",
		"dod":     "https://login.microsoftonline.us", 
		"china":   "https://login.chinacloudapi.cn",
		"public":  "https://login.microsoftonline.com",
	}

	for cloud, loginUrl := range clouds {
		// Register token endpoint for each cloud
		httpmock.RegisterResponder("POST",
			loginUrl+"/common/oauth2/v2.0/token",
			httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
				"access_token": "mock-token-" + cloud,
				"token_type":   "Bearer",
				"expires_in":   3600,
			}))

		// Register discovery endpoint for each cloud
		httpmock.RegisterResponder("GET",
			loginUrl+"/common/discovery/instance",
			httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
				"tenant_discovery_endpoint": loginUrl + "/common/v2.0/.well-known/openid-configuration",
			}))
	}
}

func init() {
	// Register the provider mock with the global registry
	mocks.GlobalRegistry.Register("provider", &ProviderMock{})
}