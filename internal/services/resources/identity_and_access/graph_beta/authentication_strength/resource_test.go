package graphBetaAuthenticationStrength_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	authStrengthMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authentication_strength/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *authStrengthMocks.AuthenticationStrengthMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	authStrengthMock := &authStrengthMocks.AuthenticationStrengthMock{}
	authStrengthMock.RegisterMocks()
	return mockClient, authStrengthMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *authStrengthMocks.AuthenticationStrengthMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	authStrengthMock := &authStrengthMocks.AuthenticationStrengthMock{}
	authStrengthMock.RegisterErrorMocks()
	return mockClient, authStrengthMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestAuthenticationStrengthResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authStrengthMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authStrengthMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAuthStrengthMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "display_name", "unit-test-authentication-strength-minimal"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "description", "Unit test minimal authentication strength policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "policy_type", "custom"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "requirements_satisfied", "mfa"),

					// Allowed combinations
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "allowed_combinations.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "allowed_combinations.*", "password,sms"),

					// Computed timestamps
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "modified_date_time"),
				),
			},
		},
	})
}

func TestAuthenticationStrengthResource_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authStrengthMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authStrengthMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAuthStrengthMaximal(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					testCheckExists("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "display_name", "unit-test-authentication-strength-maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "description", "Unit test maximal authentication strength policy with all combinations"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "policy_type", "custom"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "requirements_satisfied", "mfa"),

					// Allowed combinations (all of them)
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.#", "22"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "deviceBasedPush"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "federatedMultiFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "fido2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "hardwareOath,federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "microsoftAuthenticatorPush,federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "password"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "password,hardwareOath"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "password,microsoftAuthenticatorPush"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "password,sms"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "password,softwareOath"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "password,voice"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "qrCodePin"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "sms"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "sms,federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "softwareOath,federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "temporaryAccessPassMultiUse"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "temporaryAccessPassOneTime"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "voice,federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "windowsHelloForBusiness"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "x509CertificateMultiFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "allowed_combinations.*", "x509CertificateSingleFactor"),

					// Empty combination configurations
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_maximal", "combination_configurations.#", "0"),
				),
			},
		},
	})
}

func TestAuthenticationStrengthResource_MFAOnly(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authStrengthMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authStrengthMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAuthStrengthMFAOnly(),
				Check: resource.ComposeTestCheckFunc(
					// Basic attributes
					testCheckExists("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "display_name", "unit-test-authentication-strength-mfa-only"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "description", "Unit test MFA-only authentication strength policy"),

					// MFA-only combinations
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.#", "4"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.*", "fido2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.*", "windowsHelloForBusiness"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.*", "microsoftAuthenticatorPush,federatedSingleFactor"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.*", "x509CertificateMultiFactor"),
				),
			},
		},
	})
}

// Configuration helper functions
func testConfigAuthStrengthMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_auth_strength_minimal.tf")
	if err != nil {
		panic("failed to load authentication strength minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAuthStrengthMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_auth_strength_maximal.tf")
	if err != nil {
		panic("failed to load authentication strength maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAuthStrengthMFAOnly() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_auth_strength_mfa_only.tf")
	if err != nil {
		panic("failed to load authentication strength MFA-only config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAuthStrengthWithConfigurations() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_auth_strength_with_configurations.tf")
	if err != nil {
		panic("failed to load authentication strength with configurations config: " + err.Error())
	}
	return unitTestConfig
}

// TestAuthenticationStrengthResource_Update tests updating an authentication strength policy
func TestAuthenticationStrengthResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, authStrengthMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer authStrengthMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create minimal config
				Config: testConfigAuthStrengthMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal", "allowed_combinations.#", "1"),
					testCheckExists("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_minimal"),
				),
			},
			{
				// Update to MFA-only config
				Config: testConfigAuthStrengthMFAOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "allowed_combinations.#", "4"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_authentication_strength.auth_strength_mfa_only", "description", "Unit test MFA-only authentication strength policy"),
				),
			},
		},
	})
}
