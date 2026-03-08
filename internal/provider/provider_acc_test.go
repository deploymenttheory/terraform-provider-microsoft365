package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// loadAccTestTerraform loads a Terraform configuration file from the acceptance tests directory.
func loadAccTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccM365Provider_AuthMethods(t *testing.T) {
	validAuthMethods := []string{
		"azure_developer_cli",
		"azure_cli",
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
						PreConfig: func() {
							testlog.StepAction("provider", "Testing auth method: "+method)
						},
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
				PreConfig: func() {
					testlog.StepAction("provider", "Testing invalid auth method validation error")
				},
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
						PreConfig: func() {
							testlog.StepAction("provider", "Testing cloud environment: "+cloud)
						},
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
				PreConfig: func() {
					testlog.StepAction("provider", "Testing invalid cloud validation error")
				},
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
		"M365_CLOUD":            "gcc",
		"M365_AUTH_METHOD":      "device_code",
		"M365_CLIENT_ID":        "env-client-id",
		"M365_TENANT_ID":        "00000000-0000-0000-0000-000000000001",
		"M365_DEBUG_MODE":       "true",
		"M365_TELEMETRY_OPTOUT": "true",
	}

	// Set environment variables
	for key, value := range envVars {
		t.Setenv(key, value)
	}

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
				PreConfig: func() {
					testlog.StepAction("provider", "Testing environment variable precedence over HCL config")
				},
				Config: acceptance.ConfiguredM365ProviderBlock(loadAccTestTerraform("provider_env_precedence.tf")),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccM365Provider_ClientSecretAuth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing client_secret authentication")
				},
				Config: acceptance.ConfiguredM365ProviderBlock(loadAccTestTerraform("provider_client_secret.tf")),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccM365Provider_ClientCertificateAuth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing client_certificate authentication")
				},
				Config: acceptance.ConfiguredM365ProviderBlock(loadAccTestTerraform("provider_client_certificate.tf")),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccM365Provider_ProxyConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing proxy configuration")
				},
				Config: acceptance.ConfiguredM365ProviderBlock(loadAccTestTerraform("provider_proxy.tf")),
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
				PreConfig: func() {
					testlog.StepAction("provider", "Testing invalid tenant_id validation error")
				},
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
				PreConfig: func() {
					testlog.StepAction("provider", "Testing invalid client_id validation error")
				},
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
				PreConfig: func() {
					testlog.StepAction("provider", "Testing invalid proxy URL error handling")
				},
				Config: `
provider "microsoft365" {
  cloud = "public"
  auth_method = "device_code"
  client_options = {
    use_proxy = true
    proxy_url = "://invalid-proxy-url"
  }
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type = "display_name"
  filter_value = "Test Script"
}
`,
				ExpectError: regexp.MustCompile("failed to parse proxy URL|missing protocol scheme|invalid proxy URL"),
			},
		},
	})
}

func TestAccM365Provider_InvalidRedirectURL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing invalid redirect_url validation error")
				},
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
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing complete maximal provider configuration")
				},
				Config: acceptance.ConfiguredM365ProviderBlock(loadAccTestTerraform("provider_maximal.tf")),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccM365Provider_WithDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing empty credentials validation")
				},
				Config: `
provider "microsoft365" {
  auth_method = "client_secret"
  entra_id_options = {
    client_id = ""
    client_secret = ""
  }
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type = "display_name"
  filter_value = "Test Script"
}
`,
				ExpectError: regexp.MustCompile("client_id.*cannot be empty|client_secret.*cannot be empty|Invalid Attribute Value|Attribute client_id string length must be at least"),
			},
		},
	})
}

func TestAccM365Provider_AuthWithoutProxy(t *testing.T) {
	dataSourceName := "data.microsoft365_graph_beta_device_management_windows_remediation_script.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing authentication without proxy configuration")
				},
				Config: loadAccTestTerraform("provider_auth_without_proxy.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Data source configuration
					check.That(dataSourceName).Key("filter_type").HasValue("display_name"),
					check.That(dataSourceName).Key("filter_value").HasValue("NonExistentScript"),
				),
			},
		},
	})
}

func TestAccM365Provider_CompressionDisabledForAuth(t *testing.T) {
	dataSourceName := "data.microsoft365_graph_beta_device_management_windows_remediation_script.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing auth succeeds with compression enabled (validates gzip fix)")
				},
				Config: loadAccTestTerraform("provider_compression_enabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Validates that auth succeeds even with compression enabled
					// (compression is excluded from auth middleware but enabled for Graph API calls)
					check.That(dataSourceName).Key("filter_type").HasValue("display_name"),
					check.That(dataSourceName).Key("filter_value").HasValue("NonExistentScript"),
				),
			},
		},
	})
}

func TestAccM365Provider_AllAuthMethodsWithoutProxy(t *testing.T) {
	dataSourceName := "data.microsoft365_graph_beta_device_management_windows_remediation_script.test"

	authMethods := []struct {
		name       string
		configFile string
	}{
		{
			name:       "client_secret",
			configFile: "provider_auth_client_secret_no_proxy.tf",
		},
		{
			name:       "device_code",
			configFile: "provider_auth_device_code_no_proxy.tf",
		},
		{
			name:       "azure_cli",
			configFile: "provider_auth_azure_cli_no_proxy.tf",
		},
		{
			name:       "azure_developer_cli",
			configFile: "provider_auth_azure_developer_cli_no_proxy.tf",
		},
	}

	for _, tt := range authMethods {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { mocks.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						PreConfig: func() {
							testlog.StepAction("provider", "Testing "+tt.name+" authentication without proxy")
						},
						Config: loadAccTestTerraform(tt.configFile),
						Check: resource.ComposeTestCheckFunc(
							// Data source configuration
							check.That(dataSourceName).Key("filter_type").HasValue("display_name"),
							check.That(dataSourceName).Key("filter_value").HasValue("NonExistentScript"),
						),
					},
				},
			})
		})
	}
}

func TestAccM365Provider_RetryConfiguration(t *testing.T) {
	dataSourceName := "data.microsoft365_graph_beta_device_management_windows_remediation_script.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing retry configuration with max_retries=5")
				},
				Config: loadAccTestTerraform("provider_retry_config.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Data source configuration
					check.That(dataSourceName).Key("filter_type").HasValue("display_name"),
					check.That(dataSourceName).Key("filter_value").HasValue("NonExistentScript"),
				),
			},
		},
	})
}

func TestAccM365Provider_CustomUserAgent(t *testing.T) {
	dataSourceName := "data.microsoft365_graph_beta_device_management_windows_remediation_script.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing custom user agent configuration")
				},
				Config: loadAccTestTerraform("provider_custom_user_agent.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Data source configuration
					check.That(dataSourceName).Key("filter_type").HasValue("display_name"),
					check.That(dataSourceName).Key("filter_value").HasValue("NonExistentScript"),
				),
			},
		},
	})
}

func TestAccM365Provider_ClientOptions_AllMiddleware(t *testing.T) {
	dataSourceName := "data.microsoft365_graph_beta_device_management_windows_remediation_script.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction("provider", "Testing all middleware options enabled together")
				},
				Config: loadAccTestTerraform("provider_all_middleware.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Validates that all middleware options work together
					// Data source configuration
					check.That(dataSourceName).Key("filter_type").HasValue("display_name"),
					check.That(dataSourceName).Key("filter_value").HasValue("NonExistentScript"),
				),
			},
		},
	})
}
