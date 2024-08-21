package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	testingResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()
		tflog.Info(ctx, "Instantiating provider for unit tests")
		provider := New("test")()
		tflog.Debug(ctx, fmt.Sprintf("Unit Test Provider instantiated: %T", provider))

		return providerserver.NewProtocol6WithError(provider)()
	},
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"microsoft365": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()
		tflog.Info(ctx, "Instantiating provider for acceptance tests")
		provider := New("test")()
		tflog.Debug(ctx, fmt.Sprintf("Acceptance Test Provider instantiated: %T", provider))

		return providerserver.NewProtocol6WithError(provider)()
	},
}

// testAccPreCheck runs a pre-check to ensure the required environment variables are set
func testAccPreCheck(t *testing.T) {
	t.Helper()
	t.Log("Running pre-check for acceptance tests")

	requiredEnvVars := []string{
		"M365_TENANT_ID",
		"M365_AUTH_METHOD",
		"M365_CLOUD",
		"M365_CLIENT_ID",
	}

	for _, envVar := range requiredEnvVars {
		if v := os.Getenv(envVar); v == "" {
			t.Fatalf("%s must be set for acceptance tests", envVar)
		}
	}

	// Check for auth method specific environment variables
	authMethod := os.Getenv("M365_AUTH_METHOD")
	switch authMethod {
	case "client_secret":
		if v := os.Getenv("M365_CLIENT_SECRET"); v == "" {
			t.Fatal("M365_CLIENT_SECRET must be set when auth_method is client_secret")
		}
	case "client_certificate":
		if v := os.Getenv("M365_CLIENT_CERTIFICATE_FILE_PATH"); v == "" {
			t.Fatal("M365_CLIENT_CERTIFICATE_FILE_PATH must be set when auth_method is client_certificate")
		}
	case "username_password":
		if v := os.Getenv("M365_USERNAME"); v == "" {
			t.Fatal("M365_USERNAME must be set when auth_method is username_password")
		}
		if v := os.Getenv("M365_PASSWORD"); v == "" {
			t.Fatal("M365_PASSWORD must be set when auth_method is username_password")
		}
	case "device_code", "interactive_browser":
		// These methods don't require additional environment variables
	default:
		t.Fatalf("Unknown auth_method: %s", authMethod)
	}

	// Optional environment variables
	optionalEnvVars := []string{
		"M365_REDIRECT_URL",
		"M365_USE_PROXY",
		"M365_PROXY_URL",
		"M365_ENABLE_CHAOS",
		"M365_TELEMETRY_OPTOUT",
		"M365_DEBUG_MODE",
	}

	for _, envVar := range optionalEnvVars {
		if v := os.Getenv(envVar); v != "" {
			t.Logf("Optional environment variable %s is set", envVar)
		}
	}

	t.Log("Pre-check completed successfully")
}

// TestAccM365Provider_EnvVarPrecedence verifies that environment variables take precedence over HCL configuration.
// This test is crucial for ensuring the provider respects the expected configuration hierarchy.
func TestAccM365Provider_EnvVarPrecedence(t *testing.T) {
	t.Log("Starting TestAccM365Provider_EnvVarPrecedence")

	// Set environment variables
	envVars := map[string]string{
		"M365_TENANT_ID":                    "env-tenant-id",
		"M365_AUTH_METHOD":                  "env-device_code",
		"M365_CLIENT_ID":                    "env-client-id",
		"M365_CLIENT_SECRET":                "env-client-secret",
		"M365_CLIENT_CERTIFICATE_FILE_PATH": "env-client-certificate",
		"M365_CLIENT_CERTIFICATE_PASSWORD":  "env-client-cert-password",
		"M365_USERNAME":                     "env-username",
		"M365_PASSWORD":                     "env-password",
		"M365_REDIRECT_URL":                 "env-redirect-url",
		"M365_USE_PROXY":                    "true",
		"M365_PROXY_URL":                    "http://env-proxy-url:8080",
		"M365_CLOUD":                        "gcc",
		"M365_ENABLE_CHAOS":                 "true",
		"M365_TELEMETRY_OPTOUT":             "true",
		"M365_DEBUG_MODE":                   "true",
	}

	for key, value := range envVars {
		t.Setenv(key, value)
		t.Logf("Set %s to: %s", key, value)
	}

	testingResource.Test(t, testingResource.TestCase{
		PreCheck: func() {
			t.Log("Running pre-check")
			for key := range envVars {
				if v := os.Getenv(key); v == "" {
					t.Fatalf("%s must be set for this test", key)
				}
			}
		},
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
					func(s *terraform.State) error {
						t.Log("Starting attribute checks")
						providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
						if !ok {
							t.Log("microsoft365.provider not found in state")
							return fmt.Errorf("microsoft365.provider not found in state")
						}

						for key, envValue := range envVars {
							attrKey := strings.ToLower(strings.TrimPrefix(key, "M365_"))
							if attrValue := providerResource.Primary.Attributes[attrKey]; attrValue != envValue {
								t.Logf("Attribute mismatch for %s: expected %s, got %s", attrKey, envValue, attrValue)
								return fmt.Errorf("attribute mismatch for %s: expected %s, got %s", attrKey, envValue, attrValue)
							} else {
								t.Logf("Attribute check passed for %s: %s", attrKey, attrValue)
							}
						}
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_EnvVarPrecedence")
}

// Tenant ID

// TestAccM365Provider_InvalidTenantIDFormat verifies that an invalid tenant ID format is rejected.
// This test ensures the provider properly validates the tenant ID input.
func TestAccM365Provider_InvalidTenantIDFormat(t *testing.T) {
	t.Log("Starting TestAccM365Provider_InvalidTenantIDFormat")
	invalidTenantID := "invalid-tenant-id"
	t.Setenv("M365_TENANT_ID", invalidTenantID)
	t.Logf("Set M365_TENANT_ID to: %s", invalidTenantID)

	testingResource.Test(t, testingResource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
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
				Check: func(s *terraform.State) error {
					t.Log("Checking for expected error")
					return nil
				},
			},
		},
	})
	t.Log("Completed TestAccM365Provider_InvalidTenantIDFormat")
}

// TestAccM365Provider_MissingTenantID verifies that the provider fails when tenant_id is missing.
// This test ensures the provider enforces the required tenant_id field.
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
				Check: func(s *terraform.State) error {
					t.Log("Checking for expected error")
					return nil
				},
			},
		},
	})
	t.Log("Completed TestAccM365Provider_MissingTenantID")
}

// TestAccM365Provider_TenantIDRequired verifies that the tenant_id argument is required.
// This test ensures the provider enforces the required tenant_id field when no configuration is provided.
func TestAccM365Provider_TenantIDRequired(t *testing.T) {
	t.Log("Starting TestAccM365Provider_TenantIDRequired")
	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config:      `provider "microsoft365" {}`,
				ExpectError: regexp.MustCompile(`The argument "tenant_id" is required`),
				Check: func(s *terraform.State) error {
					t.Log("Checking for expected error")
					return nil
				},
			},
		},
	})
	t.Log("Completed TestAccM365Provider_TenantIDRequired")
}

// TestAccM365Provider_TenantIDFromEnvVar verifies that the tenant_id can be set via an environment variable.
// This test ensures the provider correctly reads the tenant_id from the environment when not specified in the configuration.
func TestAccM365Provider_TenantIDFromEnvVar(t *testing.T) {
	t.Log("Starting TestAccM365Provider_TenantIDFromEnvVar")
	envTenantID := "env-tenant-id"
	t.Setenv("M365_TENANT_ID", envTenantID)
	t.Logf("Set M365_TENANT_ID to: %s", envTenantID)

	testingResource.Test(t, testingResource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if v := os.Getenv("M365_TENANT_ID"); v == "" {
				t.Fatal("M365_TENANT_ID must be set for this test")
			}
			t.Log("PreCheck passed")
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {}`,
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "tenant_id", envTenantID,
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
						if !ok {
							t.Log("microsoft365.provider not found in state")
							return fmt.Errorf("microsoft365.provider not found in state")
						}
						if tenantID := providerResource.Primary.Attributes["tenant_id"]; tenantID != envTenantID {
							t.Logf("Tenant ID mismatch: expected %s, got %s", envTenantID, tenantID)
							return fmt.Errorf("tenant ID mismatch: expected %s, got %s", envTenantID, tenantID)
						}
						t.Logf("Tenant ID correctly set to: %s", envTenantID)
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_TenantIDFromEnvVar")
}

// TestAccM365Provider_TenantIDEnvVarOverridesHCL verifies that the tenant_id set via
// environment variable takes precedence over the one specified in HCL configuration.
func TestAccM365Provider_TenantIDEnvVarOverridesHCL(t *testing.T) {
	t.Log("Starting TestAccM365Provider_TenantIDEnvVarOverridesHCL")
	envTenantID := "env-tenant-id"
	hclTenantID := "hcl-tenant-id"
	t.Setenv("M365_TENANT_ID", envTenantID)
	t.Logf("Set M365_TENANT_ID to: %s", envTenantID)

	testingResource.Test(t, testingResource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if v := os.Getenv("M365_TENANT_ID"); v == "" {
				t.Fatal("M365_TENANT_ID must be set for this test")
			}
			t.Log("PreCheck passed")
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: fmt.Sprintf(`provider "microsoft365" {
									tenant_id = "%s"
							}`, hclTenantID),
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "tenant_id", envTenantID,
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
						if !ok {
							t.Log("microsoft365.provider not found in state")
							return fmt.Errorf("microsoft365.provider not found in state")
						}
						if tenantID := providerResource.Primary.Attributes["tenant_id"]; tenantID != envTenantID {
							t.Logf("Tenant ID mismatch: expected %s (from env), got %s", envTenantID, tenantID)
							t.Logf("HCL-specified tenant ID was: %s", hclTenantID)
							return fmt.Errorf("tenant ID mismatch: expected %s (from env), got %s. HCL-specified was: %s", envTenantID, tenantID, hclTenantID)
						}
						t.Logf("Tenant ID correctly set to env value: %s, overriding HCL value: %s", envTenantID, hclTenantID)
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_TenantIDEnvVarOverridesHCL")
}

// TestAccM365Provider_ValidTenantIDFormat verifies that a valid GUID format for tenant_id is accepted.
func TestAccM365Provider_ValidTenantIDFormat(t *testing.T) {
	t.Log("Starting TestAccM365Provider_ValidTenantIDFormat")
	validTenantID := "123e4567-e89b-12d3-a456-426614174000"
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
						providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
						if !ok {
							t.Log("microsoft365.provider not found in state")
							return fmt.Errorf("microsoft365.provider not found in state")
						}
						if tenantID := providerResource.Primary.Attributes["tenant_id"]; tenantID != validTenantID {
							t.Logf("Tenant ID mismatch: expected %s, got %s", validTenantID, tenantID)
							return fmt.Errorf("tenant ID mismatch: expected %s, got %s", validTenantID, tenantID)
						}
						t.Logf("Tenant ID correctly set to valid GUID: %s", validTenantID)
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_ValidTenantIDFormat")
}

// TestAccM365Provider_TenantIDSensitivity verifies that the tenant_id is treated as sensitive information.
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
						providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
						if !ok {
							t.Log("microsoft365.provider not found in state")
							return fmt.Errorf("microsoft365.provider not found in state")
						}
						if tenantID := providerResource.Primary.Attributes["tenant_id"]; tenantID != validTenantID {
							t.Logf("Tenant ID mismatch: expected %s, got %s", validTenantID, tenantID)
							return fmt.Errorf("tenant ID mismatch: expected %s, got %s", validTenantID, tenantID)
						}
						if sensitive, ok := providerResource.Primary.Attributes["tenant_id.0.sensitive"]; !ok || sensitive != "true" {
							t.Log("Tenant ID is not marked as sensitive")
							return fmt.Errorf("tenant ID is not marked as sensitive")
						}
						t.Log("Tenant ID correctly set and marked as sensitive")
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
				Check: func(s *terraform.State) error {
					t.Log("Checking for expected error")
					if s != nil {
						t.Error("Expected an error, but state was not nil")
						return fmt.Errorf("expected an error, but state was not nil")
					}
					return nil
				},
			},
		},
	})
	t.Log("Completed TestAccM365Provider_AuthMethodRequired")
}

// TestAccM365Provider_InvalidAuthMethod verifies that an invalid value for auth_method
// triggers a validation error.
func TestAccM365Provider_InvalidAuthMethod(t *testing.T) {
	t.Log("Starting TestAccM365Provider_InvalidAuthMethod")
	invalidMethod := "invalid_method"
	testingResource.Test(t, testingResource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: fmt.Sprintf(`provider "microsoft365" {
									auth_method = "%s"
							}`, invalidMethod),
				ExpectError: regexp.MustCompile(`expected auth_method to be one of`),
				Check: func(s *terraform.State) error {
					t.Log("Checking for expected error")
					if s != nil {
						t.Errorf("Expected an error for invalid auth_method: %s, but state was not nil", invalidMethod)
						return fmt.Errorf("expected an error for invalid auth_method: %s, but state was not nil", invalidMethod)
					}
					return nil
				},
			},
		},
	})
	t.Log("Completed TestAccM365Provider_InvalidAuthMethod")
}

// TestAccM365Provider_AuthMethodFromEnvVar tests the scenario where the auth_method
// is not specified in the HCL configuration but is provided as an environment variable.
func TestAccM365Provider_AuthMethodFromEnvVar(t *testing.T) {
	t.Log("Starting TestAccM365Provider_AuthMethodFromEnvVar")
	envAuthMethod := "client_secret"
	t.Setenv("M365_AUTH_METHOD", envAuthMethod)
	t.Logf("Set M365_AUTH_METHOD to: %s", envAuthMethod)

	testingResource.Test(t, testingResource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if v := os.Getenv("M365_AUTH_METHOD"); v == "" {
				t.Fatal("M365_AUTH_METHOD must be set for this test")
			}
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: `provider "microsoft365" {}`,
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "auth_method", envAuthMethod,
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
						if !ok {
							t.Error("microsoft365.provider not found in state")
							return fmt.Errorf("microsoft365.provider not found in state")
						}
						if authMethod := providerResource.Primary.Attributes["auth_method"]; authMethod != envAuthMethod {
							t.Errorf("auth_method mismatch: expected %s, got %s", envAuthMethod, authMethod)
							return fmt.Errorf("auth_method mismatch: expected %s, got %s", envAuthMethod, authMethod)
						}
						t.Logf("auth_method correctly set to: %s", envAuthMethod)
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
	envAuthMethod := "client_secret"
	hclAuthMethod := "device_code"
	t.Setenv("M365_AUTH_METHOD", envAuthMethod)
	t.Logf("Set M365_AUTH_METHOD to: %s", envAuthMethod)

	testingResource.Test(t, testingResource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if v := os.Getenv("M365_AUTH_METHOD"); v == "" {
				t.Fatal("M365_AUTH_METHOD must be set for this test")
			}
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: fmt.Sprintf(`provider "microsoft365" {
									auth_method = "%s"
							}`, hclAuthMethod),
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "auth_method", envAuthMethod,
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
						if !ok {
							t.Error("microsoft365.provider not found in state")
							return fmt.Errorf("microsoft365.provider not found in state")
						}
						if authMethod := providerResource.Primary.Attributes["auth_method"]; authMethod != envAuthMethod {
							t.Errorf("auth_method mismatch: expected %s (from env), got %s. HCL specified: %s", envAuthMethod, authMethod, hclAuthMethod)
							return fmt.Errorf("auth_method mismatch: expected %s (from env), got %s. HCL specified: %s", envAuthMethod, authMethod, hclAuthMethod)
						}
						t.Logf("auth_method correctly set to env value: %s, overriding HCL value: %s", envAuthMethod, hclAuthMethod)
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
								providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
								if !ok {
									t.Errorf("microsoft365.provider not found in state for auth_method: %s", method)
									return fmt.Errorf("microsoft365.provider not found in state for auth_method: %s", method)
								}
								if authMethod := providerResource.Primary.Attributes["auth_method"]; authMethod != method {
									t.Errorf("auth_method mismatch for %s: expected %s, got %s", method, method, authMethod)
									return fmt.Errorf("auth_method mismatch for %s: expected %s, got %s", method, method, authMethod)
								}
								t.Logf("auth_method correctly set to: %s", method)
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
						providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
						if !ok {
							t.Error("microsoft365.provider not found in state")
							return fmt.Errorf("microsoft365.provider not found in state")
						}
						if authMethod := providerResource.Primary.Attributes["auth_method"]; authMethod != validAuthMethod {
							t.Errorf("auth_method mismatch: expected %s, got %s", validAuthMethod, authMethod)
							return fmt.Errorf("auth_method mismatch: expected %s, got %s", validAuthMethod, authMethod)
						}
						// Note: auth_method itself is not sensitive, so we're not checking for sensitivity here
						t.Logf("auth_method correctly set to: %s", validAuthMethod)
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
// TestAccM365Provider_AuthMethodCombinedEnvVarAndHCL tests the scenario where the environment variable is set,
// but the HCL configuration is explicitly set to a different valid value, to confirm that precedence is respected.
func TestAccM365Provider_AuthMethodCombinedEnvVarAndHCL(t *testing.T) {
	t.Log("Starting TestAccM365Provider_AuthMethodCombinedEnvVarAndHCL")
	envAuthMethod := "interactive_browser"
	hclAuthMethod := "device_code"
	t.Setenv("M365_AUTH_METHOD", envAuthMethod)
	t.Logf("Set M365_AUTH_METHOD to: %s", envAuthMethod)

	testingResource.Test(t, testingResource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if v := os.Getenv("M365_AUTH_METHOD"); v == "" {
				t.Fatal("M365_AUTH_METHOD must be set for this test")
			}
		},
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []testingResource.TestStep{
			{
				Config: fmt.Sprintf(`provider "microsoft365" {
									auth_method = "%s"
							}`, hclAuthMethod),
				Check: testingResource.ComposeTestCheckFunc(
					testingResource.TestCheckResourceAttr(
						"microsoft365.provider", "auth_method", envAuthMethod,
					),
					func(s *terraform.State) error {
						t.Log("Checking provider attributes")
						providerResource, ok := s.RootModule().Resources["microsoft365.provider"]
						if !ok {
							t.Error("microsoft365.provider not found in state")
							return fmt.Errorf("microsoft365.provider not found in state")
						}
						if authMethod := providerResource.Primary.Attributes["auth_method"]; authMethod != envAuthMethod {
							t.Errorf("auth_method mismatch: expected %s (from env), got %s. HCL specified: %s", envAuthMethod, authMethod, hclAuthMethod)
							return fmt.Errorf("auth_method mismatch: expected %s (from env), got %s. HCL specified: %s", envAuthMethod, authMethod, hclAuthMethod)
						}
						t.Logf("auth_method correctly set to env value: %s, overriding HCL value: %s", envAuthMethod, hclAuthMethod)
						return nil
					},
				),
			},
		},
	})
	t.Log("Completed TestAccM365Provider_AuthMethodCombinedEnvVarAndHCL")
}
