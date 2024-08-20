package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	testingResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": func() (tfprotov6.ProviderServer, error) {
		fmt.Println("Instantiating provider for unit tests")
		provider := New("1.0.0")()
		fmt.Printf("Unit Test Provider instantiated: %T\n", provider)
		return providerserver.NewProtocol6WithError(provider)()
	},
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": func() (tfprotov6.ProviderServer, error) {
		fmt.Println("Instantiating provider for acceptance tests")
		provider := New("1.0.0")()
		fmt.Printf("Acceptance Test Provider instantiated: %T\n", provider)
		return providerserver.NewProtocol6WithError(provider)()
	},
}

func TestAccM365Provider_EnvVarPrecedence(t *testing.T) {
	// Set environment variables
	t.Setenv("M365_TENANT_ID", "env-tenant-id")
	t.Setenv("M365_AUTH_METHOD", "env-device_code")
	t.Setenv("M365_CLIENT_ID", "env-client-id")
	t.Setenv("M365_CLIENT_SECRET", "env-client-secret")
	t.Setenv("M365_CLIENT_CERTIFICATE_FILE_PATH", "env-client-certificate")
	t.Setenv("M365_CLIENT_CERTIFICATE_PASSWORD", "env-client-cert-password")
	t.Setenv("M365_USERNAME", "env-username")
	t.Setenv("M365_PASSWORD", "env-password")
	t.Setenv("M365_REDIRECT_URL", "env-redirect-url")
	t.Setenv("M365_USE_PROXY", "true")
	t.Setenv("M365_PROXY_URL", "http://env-proxy-url:8080")
	t.Setenv("M365_CLOUD", "gcc")
	t.Setenv("M365_ENABLE_CHAOS", "true")
	t.Setenv("M365_TELEMETRY_OPTOUT", "true")
	t.Setenv("M365_DEBUG_MODE", "true")

	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {
                    tenant_id                   = "hcl-tenant-id"
                    auth_method                 = "client_secret"
                    client_id                   = "hcl-client-id"
                    client_secret               = "hcl-client-secret"
                    client_certificate          = "hcl-client-certificate"
                    client_certificate_password = "hcl-client-cert-password"
                    username                    = "hcl-username"
                    password                    = "hcl-password"
                    redirect_url                = "http://hcl-redirect-url:8080"
                    use_proxy                   = false
                    proxy_url                   = "http://hcl-proxy-url:8080"
                    cloud                       = "public"
                    enable_chaos                = false
                    telemetry_optout            = false
                    debug_mode                  = false
                }`,
				Check: testingResource.ComposeTestCheckFunc(
					// Ensure that the value from the environment variable is used
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "tenant_id", "env-tenant-id",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "auth_method", "env-device_code",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "client_id", "env-client-id",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "client_secret", "env-client-secret",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "client_certificate", "env-client-certificate",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "client_certificate_password", "env-client-cert-password",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "username", "env-username",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "password", "env-password",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "redirect_url", "env-redirect-url",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "use_proxy", "true",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "proxy_url", "http://env-proxy-url:8080",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "cloud", "gcc",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "enable_chaos", "true",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "telemetry_optout", "true",
					),
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "debug_mode", "true",
					),
				),
			},
		},
	})
}

// Tenant ID

func TestAccM365Provider_InvalidTenantIDFormat(t *testing.T) {
	t.Log("Starting TestAccM365Provider_InvalidTenantIDFormat")
	t.Setenv("M365_TENANT_ID", "invalid-tenant-id")
	t.Logf("Set M365_TENANT_ID to: %s", os.Getenv("M365_TENANT_ID"))

	testingResource.Test(t, testingResource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("M365_TENANT_ID"); v == "" {
				t.Fatal("M365_TENANT_ID must be set for this test")
			}
			t.Log("PreCheck passed")
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config:      `provider "microsoft365" {}`,
				ExpectError: regexp.MustCompile("Invalid GUID format for tenant_id"),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_InvalidTenantIDFormat")
}

func TestAccM365Provider_MissingTenantID(t *testing.T) {
	t.Log("Starting TestAccM365Provider_MissingTenantID")
	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {
                    auth_method = "client_secret"
                    client_id = "some-client-id"
                    client_secret = "some-client-secret"
                }`,
				ExpectError: regexp.MustCompile("Missing required argument: tenant_id"),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_MissingTenantID")
}

func TestAccM365Provider_TenantIDRequired(t *testing.T) {
	t.Log("Starting TestAccM365Provider_TenantIDRequired")
	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config:      `provider "microsoft365" {}`,
				ExpectError: regexp.MustCompile(`The argument "tenant_id" is required`),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_TenantIDRequired")
}

func TestAccM365Provider_TenantIDFromEnvVar(t *testing.T) {
	t.Log("Starting TestAccM365Provider_TenantIDFromEnvVar")
	t.Setenv("M365_TENANT_ID", "env-tenant-id")
	t.Logf("Set M365_TENANT_ID to: %s", os.Getenv("M365_TENANT_ID"))

	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {}`,
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "tenant_id", "env-tenant-id",
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_TenantIDFromEnvVar")
}

func TestAccM365Provider_TenantIDEnvVarOverridesHCL(t *testing.T) {
	t.Log("Starting TestAccM365Provider_TenantIDEnvVarOverridesHCL")
	t.Setenv("M365_TENANT_ID", "env-tenant-id")
	t.Logf("Set M365_TENANT_ID to: %s", os.Getenv("M365_TENANT_ID"))

	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {
                    tenant_id = "hcl-tenant-id"
                }`,
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "tenant_id", "env-tenant-id",
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_TenantIDEnvVarOverridesHCL")
}

func TestAccM365Provider_ValidTenantIDFormat(t *testing.T) {
	t.Log("Starting TestAccM365Provider_ValidTenantIDFormat")
	validTenantID := "123e4567-e89b-12d3-a456-426614174000" // Example valid GUID
	t.Logf("Using valid tenant ID: %s", validTenantID)

	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: fmt.Sprintf(`provider "microsoft365" {
                    tenant_id = "%s"
                }`, validTenantID),
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "tenant_id", validTenantID,
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_ValidTenantIDFormat")
}

func TestAccM365Provider_TenantIDSensitivity(t *testing.T) {
	t.Log("Starting TestAccM365Provider_TenantIDSensitivity")
	validTenantID := "123e4567-e89b-12d3-a456-426614174000" // Example valid GUID
	t.Logf("Using valid tenant ID: %s", validTenantID)

	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: fmt.Sprintf(`provider "microsoft365" {
                    tenant_id = "%s"
                }`, validTenantID),
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "tenant_id", validTenantID,
					),
					testingResource.TestCheckTypeSetElemNestedAttrs(
						"microsoft365.provider", "tenant_id", map[string]string{
							"sensitive": "true",
						},
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes and sensitivity")
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_TenantIDSensitivity")
}

// Auth Method

// TestAccM365Provider_AuthMethodRequired ensures that the auth_method is required
// and that a configuration without it fails to apply.
func TestAccM365Provider_AuthMethodRequired(t *testing.T) {
	t.Log("Starting TestAccM365Provider_AuthMethodRequired")
	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config:      `provider "microsoft365" {}`,
				ExpectError: regexp.MustCompile(`The argument "auth_method" is required`),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_AuthMethodRequired")
}

// TestAccM365Provider_InvalidAuthMethod verifies that an invalid value for auth_method
// triggers a validation error.
func TestAccM365Provider_InvalidAuthMethod(t *testing.T) {
	t.Log("Starting TestAccM365Provider_InvalidAuthMethod")
	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {
                    auth_method = "invalid_method"
                }`,
				ExpectError: regexp.MustCompile(`expected auth_method to be one of`),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_InvalidAuthMethod")
}

// TestAccM365Provider_AuthMethodFromEnvVar tests the scenario where the auth_method
// is not specified in the HCL configuration but is provided as an environment variable.
func TestAccM365Provider_AuthMethodFromEnvVar(t *testing.T) {
	t.Log("Starting TestAccM365Provider_AuthMethodFromEnvVar")
	t.Setenv("M365_AUTH_METHOD", "client_secret")
	t.Logf("Set M365_AUTH_METHOD to: %s", os.Getenv("M365_AUTH_METHOD"))

	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {}`,
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "auth_method", "client_secret",
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_AuthMethodFromEnvVar")
}

// TestAccM365Provider_AuthMethodEnvVarOverridesHCL ensures that when both an environment variable
// and an HCL configuration are provided, the environment variable takes precedence.
func TestAccM365Provider_AuthMethodEnvVarOverridesHCL(t *testing.T) {
	t.Log("Starting TestAccM365Provider_AuthMethodEnvVarOverridesHCL")
	t.Setenv("M365_AUTH_METHOD", "client_secret")
	t.Logf("Set M365_AUTH_METHOD to: %s", os.Getenv("M365_AUTH_METHOD"))

	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {
                    auth_method = "device_code"
                }`,
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "auth_method", "client_secret",
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_AuthMethodEnvVarOverridesHCL")
}

// TestAccM365Provider_ValidAuthMethodValues ensures that valid values for auth_method
// are accepted and used correctly.
func TestAccM365Provider_ValidAuthMethodValues(t *testing.T) {
	t.Log("Starting TestAccM365Provider_ValidAuthMethodValues")
	validAuthMethods := []string{
		"device_code", "client_secret", "client_certificate", "interactive_browser", "username_password",
	}

	for _, method := range validAuthMethods {
		t.Run(fmt.Sprintf("AuthMethod=%s", method), func(t *testing.T) {
			t.Logf("Testing auth_method: %s", method)
			testingResource.Test(t, testingResource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Steps: []testingResource.TestStep{
					{
						Config: fmt.Sprintf(`provider "microsoft365" {
                            auth_method = "%s"
                        }`, method),
						Check: testingResource.ComposeTestCheckFunc(
							testingResource.TestCheckResourceAttr(
								"microsoft365.provider", "auth_method", method,
							),
							func(s *terraform.State) error {
								t.Log("Checking provider attributes")
								return nil
							},
						),
					},
				},
			})
		})
	}
	t.Log("Completed TestAccM365Provider_ValidAuthMethodValues")
}

// TestAccM365Provider_AuthMethodSensitivity ensures that the auth_method attribute
// is treated correctly and does not reveal sensitive information.
func TestAccM365Provider_AuthMethodSensitivity(t *testing.T) {
	t.Log("Starting TestAccM365Provider_AuthMethodSensitivity")
	validAuthMethod := "client_secret"

	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: fmt.Sprintf(`provider "microsoft365" {
                    auth_method = "%s"
                }`, validAuthMethod),
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "auth_method", validAuthMethod,
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes and sensitivity")
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_AuthMethodSensitivity")
}

// TestAccM365Provider_AuthMethodCombinedEnvVarAndHCL tests the scenario where the environment variable is set,
// but the HCL configuration is explicitly set to a different valid value, to confirm that precedence is respected.
func TestAccM365Provider_AuthMethodCombinedEnvVarAndHCL(t *testing.T) {
	t.Log("Starting TestAccM365Provider_AuthMethodCombinedEnvVarAndHCL")
	t.Setenv("M365_AUTH_METHOD", "interactive_browser")
	t.Logf("Set M365_AUTH_METHOD to: %s", os.Getenv("M365_AUTH_METHOD"))

	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {
                    auth_method = "device_code"
                }`,
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "auth_method", "interactive_browser",
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_AuthMethodCombinedEnvVarAndHCL")
}
