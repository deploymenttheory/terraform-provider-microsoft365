package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccM365Provider_AuthMethods(t *testing.T) {
	validAuthMethods := []string{
		"azure_developer_cli",
		"client_secret",
		"client_certificate", 
		"interactive_browser",
		"device_code",
		"workload_identity",
		"managed_identity",
		"oidc",
		"oidc_github", 
		"oidc_azure_devops",
	}

	for _, method := range validAuthMethods {
		t.Run(method, func(t *testing.T) {
			config := mocks.LoadTerraformTemplateFile("provider_auth_method.tf", map[string]string{
				"AuthMethod": method,
			})

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { mocks.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: acceptance.ConfiguredM365ProviderBlock(config),
						Check:  resource.ComposeTestCheckFunc(),
					},
				},
			})
		})
	}
}

func TestAccM365Provider_InvalidAuthMethod(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "microsoft365" {
  auth_method = "invalid_method"
}

# Try to use a data source to trigger provider initialization
data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type = "display_name"
  filter_value = "Test Script"
}
`,
				ExpectError: regexp.MustCompile("Invalid Attribute Value|expected auth_method to be one of|value must be one of|authentication method.*invalid"),
			},
		},
	})
}

func TestAccM365Provider_CloudEnvironments(t *testing.T) {
	validClouds := []string{
		"public",
		"dod",
		"gcc", 
		"gcchigh",
		"china",
		"ex",
		"rx",
	}

	for _, cloud := range validClouds {
		t.Run(cloud, func(t *testing.T) {
			config := mocks.LoadTerraformTemplateFile("provider_cloud.tf", map[string]string{
				"Cloud": cloud,
			})

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { mocks.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: acceptance.ConfiguredM365ProviderBlock(config),
						Check:  resource.ComposeTestCheckFunc(),
					},
				},
			})
		})
	}
}

func TestAccM365Provider_InvalidCloud(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "microsoft365" {
  cloud = "invalid_cloud"
  auth_method = "device_code"
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type = "display_name"
  filter_value = "Test Script"
}
`,
				ExpectError: regexp.MustCompile("expected cloud to be one of|value must be one of|Invalid Attribute Value"),
			},
		},
	})
}

func TestAccM365Provider_EnvVarPrecedence(t *testing.T) {
	// Test that environment variables take precedence over HCL configuration
	envVars := map[string]string{
		"M365_CLOUD":           "gcc",
		"M365_AUTH_METHOD":     "device_code",
		"M365_CLIENT_ID":       "env-client-id",
		"M365_TENANT_ID":       "00000000-0000-0000-0000-000000000001",
		"M365_DEBUG_MODE":      "true",
		"M365_TELEMETRY_OPTOUT": "true",
	}

	// Set environment variables
	for key, value := range envVars {
		t.Setenv(key, value)
	}

	config := mocks.LoadTerraformConfigFile("provider_env_precedence.tf")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			// Verify env vars are set
			for key := range envVars {
				if v := os.Getenv(key); v == "" {
					t.Fatalf("%s must be set for this test", key)
				}
			}
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acceptance.ConfiguredM365ProviderBlock(config),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccM365Provider_ClientSecretAuth(t *testing.T) {
	config := mocks.LoadTerraformConfigFile("provider_client_secret.tf")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acceptance.ConfiguredM365ProviderBlock(config),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccM365Provider_ClientCertificateAuth(t *testing.T) {
	config := mocks.LoadTerraformConfigFile("provider_client_certificate.tf")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acceptance.ConfiguredM365ProviderBlock(config),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccM365Provider_ProxyConfiguration(t *testing.T) {
	config := mocks.LoadTerraformConfigFile("provider_proxy.tf")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acceptance.ConfiguredM365ProviderBlock(config),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccM365Provider_InvalidTenantID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "microsoft365" {
  tenant_id = "invalid-tenant-id"
  auth_method = "device_code"
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type = "display_name"
  filter_value = "Test Script"
}
`,
				ExpectError: regexp.MustCompile("must be a valid GUID|Invalid Attribute Value"),
			},
		},
	})
}

func TestAccM365Provider_InvalidClientID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "microsoft365" {
  auth_method = "client_secret"
  entra_id_options = {
    client_id = "invalid-client-id"
    client_secret = "test-secret"
  }
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type = "display_name"
  filter_value = "Test Script"
}
`,
				ExpectError: regexp.MustCompile("must be a valid GUID|Invalid Attribute Value"),
			},
		},
	})
}

func TestAccM365Provider_InvalidProxyURL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "microsoft365" {
  auth_method = "device_code"
  client_options = {
    use_proxy = true
    proxy_url = "invalid-proxy-url"
  }
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type = "display_name"
  filter_value = "Test Script"
}
`,
				ExpectError: regexp.MustCompile("Unable to create credentials|secret can't be empty|not a valid URL|Invalid Attribute Value"),
			},
		},
	})
}

func TestAccM365Provider_InvalidRedirectURL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "microsoft365" {
  auth_method = "interactive_browser"
  entra_id_options = {
    redirect_url = "invalid-redirect-url"
  }
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type = "display_name"
  filter_value = "Test Script"
}
`,
				ExpectError: regexp.MustCompile("Invalid Redirect URL|not a valid URL|Invalid Attribute Value"),
			},
		},
	})
}

func TestAccM365Provider_CompleteConfiguration(t *testing.T) {
	config := mocks.LoadTerraformConfigFile("provider_maximal.tf")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acceptance.ConfiguredM365ProviderBlock(config),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccM365Provider_WithDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "microsoft365" {
  # Configuration from environment variables
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type = "display_name"
  filter_value = "Test Script"
}
`,
				ExpectError: regexp.MustCompile("Unable to create credentials|secret can't be empty"),
			},
		},
	})
}