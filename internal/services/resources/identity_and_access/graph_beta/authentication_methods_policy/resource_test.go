package graphBetaIdentityAndAccessAuthenticationMethodsPolicy_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	authenticationMethodsPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authentication_methods_policy/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *authenticationMethodsPolicyMocks.AuthenticationMethodsPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	authenticationMethodsPolicyMock := &authenticationMethodsPolicyMocks.AuthenticationMethodsPolicyMock{}
	authenticationMethodsPolicyMock.RegisterMocks()
	return mockClient, authenticationMethodsPolicyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *authenticationMethodsPolicyMocks.AuthenticationMethodsPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	authenticationMethodsPolicyMock := &authenticationMethodsPolicyMocks.AuthenticationMethodsPolicyMock{}
	authenticationMethodsPolicyMock.RegisterErrorMocks()
	return mockClient, authenticationMethodsPolicyMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// importStep verifies import by the fixed singleton ID.
func importStep() resource.TestStep {
	return resource.TestStep{
		ResourceName:            resourceType + ".test",
		ImportState:             true,
		ImportStateId:           "authenticationMethodsPolicy",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// checkPatchedMethods asserts which endpoints received PATCH requests, in order. "" is the policy itself.
func checkPatchedMethods(mock *authenticationMethodsPolicyMocks.AuthenticationMethodsPolicyMock, expected ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		patches := mock.PatchRequests()
		if len(patches) != len(expected) {
			return fmt.Errorf("expected %d PATCH requests %v, got %d", len(expected), expected, len(patches))
		}
		for i, patch := range patches {
			if patch.MethodID != expected[i] {
				return fmt.Errorf("expected PATCH %d to target %q, got %q", i, expected[i], patch.MethodID)
			}
		}
		return nil
	}
}

// checkLastPatchBody asserts a value in the most recent PATCH request body sent to the mock.
func checkLastPatchBody(mock *authenticationMethodsPolicyMocks.AuthenticationMethodsPolicyMock, check func(body map[string]any) error) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		patches := mock.PatchRequests()
		if len(patches) == 0 {
			return fmt.Errorf("expected at least one PATCH request")
		}
		return check(patches[len(patches)-1].Body)
	}
}

// Minimal: nothing is configured, so create adopts the policy without sending any PATCH.
func TestUnitResourceAuthenticationMethodsPolicy_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchedMethods(authenticationMethodsPolicyMock),
					check.That(resourceType+".test").Key("id").HasValue("authenticationMethodsPolicy"),
					check.That(resourceType+".test").Key("display_name").HasValue("Authentication Methods Policy"),
					check.That(resourceType+".test").Key("policy_migration_state").HasValue("preMigration"),
					check.That(resourceType+".test").Key("registration_enforcement.authentication_methods_registration_campaign.state").HasValue("default"),
					check.That(resourceType+".test").Key("registration_enforcement.authentication_methods_registration_campaign.snooze_duration_in_days").HasValue("1"),
					check.That(resourceType+".test").Key("registration_enforcement.authentication_methods_registration_campaign.include_targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("report_suspicious_activity_settings.state").HasValue("default"),
					check.That(resourceType+".test").Key("report_suspicious_activity_settings.include_target.id").HasValue("all_users"),
					check.That(resourceType+".test").Key("system_credential_preferences.include_targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("email.state").HasValue("enabled"),
					check.That(resourceType+".test").Key("email.allow_external_id_to_use_email_otp").HasValue("default"),
					check.That(resourceType+".test").Key("email.include_targets.#").HasValue("0"),
					check.That(resourceType+".test").Key("temporary_access_pass.default_length").HasValue("8"),
					check.That(resourceType+".test").Key("temporary_access_pass.maximum_lifetime_in_minutes").HasValue("480"),
					check.That(resourceType+".test").Key("sms.include_targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("sms.include_targets.0.is_usable_for_sign_in").HasValue("false"),
					check.That(resourceType+".test").Key("software_oath.state").HasValue("enabled"),
					check.That(resourceType+".test").Key("hardware_oath.state").HasValue("enabled"),
				),
			},
			importStep(),
		},
	})
}

// Email: only the Email method configuration is PATCHed, with its @odata.type.
func TestUnitResourceAuthenticationMethodsPolicy_02_Email(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_email.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchedMethods(authenticationMethodsPolicyMock, "Email"),
					checkLastPatchBody(authenticationMethodsPolicyMock, func(body map[string]any) error {
						if body["@odata.type"] != "#microsoft.graph.emailAuthenticationMethodConfiguration" {
							return fmt.Errorf("unexpected @odata.type %v", body["@odata.type"])
						}
						excludeTargets, ok := body["excludeTargets"].([]any)
						if !ok || len(excludeTargets) != 0 {
							return fmt.Errorf("expected an empty excludeTargets list, got %v", body["excludeTargets"])
						}
						return nil
					}),
					check.That(resourceType+".test").Key("email.state").HasValue("enabled"),
					check.That(resourceType+".test").Key("email.allow_external_id_to_use_email_otp").HasValue("disabled"),
					check.That(resourceType+".test").Key("email.include_targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("email.include_targets.0.id").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".test").Key("email.include_targets.0.is_registration_required").HasValue("false"),
					check.That(resourceType+".test").Key("email.exclude_targets.#").HasValue("0"),
					check.That(resourceType+".test").Key("temporary_access_pass.default_length").HasValue("8"),
				),
			},
			importStep(),
		},
	})
}

// EmailUpdate: changing one method configuration PATCHes only that method; the others are untouched.
func TestUnitResourceAuthenticationMethodsPolicy_03_EmailUpdate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check:  checkPatchedMethods(authenticationMethodsPolicyMock),
			},
			{
				Config: loadUnitTestTerraform("resource_email.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchedMethods(authenticationMethodsPolicyMock, "Email"),
					check.That(resourceType+".test").Key("email.include_targets.0.id").HasValue("33333333-3333-3333-3333-333333333333"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_email_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchedMethods(authenticationMethodsPolicyMock, "Email", "Email"),
					check.That(resourceType+".test").Key("email.include_targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("email.include_targets.0.id").HasValue("44444444-4444-4444-4444-444444444444"),
					check.That(resourceType+".test").Key("email.exclude_targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("email.exclude_targets.0.id").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".test").Key("sms.state").HasValue("enabled"),
				),
			},
			{
				// Removing the block stops managing the method; nothing is written and the tenant value is kept.
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchedMethods(authenticationMethodsPolicyMock, "Email", "Email"),
					check.That(resourceType+".test").Key("email.include_targets.0.id").HasValue("44444444-4444-4444-4444-444444444444"),
				),
			},
			importStep(),
		},
	})
}

// TemporaryAccessPass: every Temporary Access Pass setting is sent and read back.
func TestUnitResourceAuthenticationMethodsPolicy_04_TemporaryAccessPass(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_temporary_access_pass.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchedMethods(authenticationMethodsPolicyMock, "TemporaryAccessPass"),
					check.That(resourceType+".test").Key("temporary_access_pass.state").HasValue("enabled"),
					check.That(resourceType+".test").Key("temporary_access_pass.default_lifetime_in_minutes").HasValue("120"),
					check.That(resourceType+".test").Key("temporary_access_pass.default_length").HasValue("12"),
					check.That(resourceType+".test").Key("temporary_access_pass.minimum_lifetime_in_minutes").HasValue("60"),
					check.That(resourceType+".test").Key("temporary_access_pass.maximum_lifetime_in_minutes").HasValue("480"),
					check.That(resourceType+".test").Key("temporary_access_pass.is_usable_once").HasValue("true"),
					check.That(resourceType+".test").Key("temporary_access_pass.include_targets.0.id").HasValue("33333333-3333-3333-3333-333333333333"),
				),
			},
			importStep(),
		},
	})
}

// SmsAndOath: the SMS and both OATH method configurations are PATCHed individually.
func TestUnitResourceAuthenticationMethodsPolicy_05_SmsAndOath(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_sms_oath.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchedMethods(authenticationMethodsPolicyMock, "Sms", "SoftwareOath", "HardwareOath"),
					check.That(resourceType+".test").Key("sms.state").HasValue("disabled"),
					check.That(resourceType+".test").Key("sms.include_targets.0.id").HasValue("all_users"),
					check.That(resourceType+".test").Key("sms.include_targets.0.is_usable_for_sign_in").HasValue("true"),
					check.That(resourceType+".test").Key("software_oath.state").HasValue("disabled"),
					check.That(resourceType+".test").Key("software_oath.include_targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("hardware_oath.include_targets.0.id").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".test").Key("hardware_oath.exclude_targets.0.id").HasValue("44444444-4444-4444-4444-444444444444"),
				),
			},
			importStep(),
		},
	})
}

// PolicySettings: policy-level settings are sent in a single PATCH to the policy.
func TestUnitResourceAuthenticationMethodsPolicy_06_PolicySettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_policy_settings.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkPatchedMethods(authenticationMethodsPolicyMock, ""),
					checkLastPatchBody(authenticationMethodsPolicyMock, func(body map[string]any) error {
						if _, ok := body["authenticationMethodConfigurations"]; ok {
							return fmt.Errorf("policy PATCH must not include authenticationMethodConfigurations")
						}
						if body["policyMigrationState"] != "migrationComplete" {
							return fmt.Errorf("unexpected policyMigrationState %v", body["policyMigrationState"])
						}
						return nil
					}),
					check.That(resourceType+".test").Key("policy_migration_state").HasValue("migrationComplete"),
					check.That(resourceType+".test").Key("registration_enforcement.authentication_methods_registration_campaign.state").HasValue("enabled"),
					check.That(resourceType+".test").Key("registration_enforcement.authentication_methods_registration_campaign.snooze_duration_in_days").HasValue("3"),
					check.That(resourceType+".test").Key("registration_enforcement.authentication_methods_registration_campaign.enforce_registration_after_allowed_snoozes").HasValue("true"),
					check.That(resourceType+".test").Key("registration_enforcement.authentication_methods_registration_campaign.include_targets.0.targeted_authentication_method").HasValue("microsoftAuthenticator"),
					check.That(resourceType+".test").Key("registration_enforcement.authentication_methods_registration_campaign.exclude_targets.0.id").HasValue("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".test").Key("report_suspicious_activity_settings.state").HasValue("enabled"),
					check.That(resourceType+".test").Key("report_suspicious_activity_settings.voice_reporting_code").HasValue("8"),
					check.That(resourceType+".test").Key("system_credential_preferences.state").HasValue("enabled"),
					check.That(resourceType+".test").Key("system_credential_preferences.exclude_targets.#").HasValue("0"),
				),
			},
			importStep(),
		},
	})
}

// ErrorForbidden: a 403 on read surfaces as a diagnostic.
func TestUnitResourceAuthenticationMethodsPolicy_07_ErrorForbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)forbidden|insufficient privileges|403`),
			},
		},
	})
}

// ErrorBadRequest: a 400 on a method configuration PATCH surfaces as a diagnostic.
func TestUnitResourceAuthenticationMethodsPolicy_08_ErrorBadRequest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_email.tf"),
				ExpectError: regexp.MustCompile(`(?i)bad request|400|authentication method configuration is invalid`),
			},
		},
	})
}

// ErrorNotFound: a 404 on read surfaces as a diagnostic during create.
func TestUnitResourceAuthenticationMethodsPolicy_09_ErrorNotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupMockEnvironment()
	authenticationMethodsPolicyMock.RegisterNotFoundMocks()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)not found|404|read resource state`),
			},
		},
	})
}

// ValidationErrors: invalid configuration is rejected before any request is sent.
func TestUnitResourceAuthenticationMethodsPolicy_10_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_state.tf"),
				ExpectError: regexp.MustCompile(`(?s)email\.state.*value\s+must\s+be\s+one\s+of`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_length.tf"),
				ExpectError: regexp.MustCompile(`(?s)default_length.*between\s+8\s+and\s+48`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_target.tf"),
				ExpectError: regexp.MustCompile(`(?s)include_targets.*valid\s+GUID`),
			},
		},
	})
}

// InvalidImportID: import only accepts the literal singleton ID.
func TestUnitResourceAuthenticationMethodsPolicy_11_InvalidImportID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authenticationMethodsPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authenticationMethodsPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
			},
			{
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: "somethingElse",
				ExpectError:   regexp.MustCompile(`Invalid import ID`),
			},
		},
	})
}
