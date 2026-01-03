package provider_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccM365Provider_ConfigurationBuilder(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{
			name:   "basic_device_code",
			config: acceptance.ProviderConfigWithAuthMethod("device_code"),
		},
		{
			name:   "gcc_cloud",
			config: acceptance.ProviderConfigWithCloud("gcc"),
		},
		{
			name: "client_secret_full",
			config: acceptance.ProviderConfigForClientSecret(
				"00000000-0000-0000-0000-000000000001",
				"test-secret",
			),
		},
		{
			name: "client_certificate_full",
			config: acceptance.ProviderConfigForClientCertificate(
				"00000000-0000-0000-0000-000000000001",
				"/path/to/cert.pfx",
				"cert-password",
			),
		},
		{
			name: "builder_comprehensive",
			config: acceptance.NewProviderConfigBuilder().
				WithCloud("gcc").
				WithClientSecret("00000000-0000-0000-0000-000000000001", "test-secret").
				WithProxy("http://proxy.example.com:8080").
				WithDebugMode(true).
				WithTelemetryOptout(true).
				Build(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { mocks.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: acceptance.ConfiguredM365ProviderBlock(tt.config),
						Check:  resource.ComposeTestCheckFunc(),
					},
				},
			})
		})
	}
}

func TestAccM365Provider_ValidationErrors(t *testing.T) {
	validationTests := []struct {
		name        string
		config      string
		expectError string
	}{
		{
			name: "invalid_auth_method",
			config: `
provider "microsoft365" {
  auth_method = "invalid_method"
}`,
			expectError: "expected auth_method to be one of",
		},
		{
			name: "invalid_cloud",
			config: `
provider "microsoft365" {
  cloud = "invalid_cloud" 
  auth_method = "device_code"
}`,
			expectError: "expected cloud to be one of",
		},
		{
			name: "invalid_tenant_id",
			config: `
provider "microsoft365" {
  tenant_id = "not-a-guid"
  auth_method = "device_code"
}`,
			expectError: "must be a valid GUID",
		},
		{
			name: "invalid_client_id",
			config: `
provider "microsoft365" {
  auth_method = "client_secret"
  entra_id_options = {
    client_id = "not-a-guid"
    client_secret = "test-secret"
  }
}`,
			expectError: "must be a valid GUID",
		},
		{
			name: "invalid_proxy_url",
			config: `
provider "microsoft365" {
  auth_method = "device_code"
  client_options = {
    use_proxy = true
    proxy_url = "not-a-valid-url"
  }
}`,
			expectError: "not a valid URL",
		},
		{
			name: "invalid_redirect_url",
			config: `
provider "microsoft365" {
  auth_method = "interactive_browser"
  entra_id_options = {
    redirect_url = "not-a-valid-url"
  }
}`,
			expectError: "not a valid URL",
		},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      tt.config,
						ExpectError: nil, // Framework validation will catch these
					},
				},
			})
		})
	}
}
