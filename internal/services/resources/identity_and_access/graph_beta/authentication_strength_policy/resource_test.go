package graphBetaAuthenticationStrengthPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	authStrengthMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authentication_strength_policy/mocks"
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

func TestUnitResourceAuthenticationStrengthPolicy_01_Minimal(t *testing.T) {
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
					check.That(resourceType+".auth_strength_minimal").Key("display_name").HasValue("unit-test-auth-strength-min"),
					check.That(resourceType+".auth_strength_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".auth_strength_minimal").Key("description").HasValue("Unit test minimal authentication strength policy"),
					check.That(resourceType+".auth_strength_minimal").Key("policy_type").HasValue("custom"),
					check.That(resourceType+".auth_strength_minimal").Key("requirements_satisfied").HasValue("mfa"),

					// Allowed combinations
					check.That(resourceType+".auth_strength_minimal").Key("allowed_combinations.#").HasValue("1"),
					check.That(resourceType+".auth_strength_minimal").Key("allowed_combinations.*").ContainsTypeSetElement("password,sms"),

					// Computed timestamps
					check.That(resourceType+".auth_strength_minimal").Key("created_date_time").Exists(),
					check.That(resourceType+".auth_strength_minimal").Key("modified_date_time").Exists(),
				),
			},
			{
				ResourceName:      resourceType + ".auth_strength_minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitResourceAuthenticationStrengthPolicy_02_Maximal(t *testing.T) {
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
					check.That(resourceType+".auth_strength_maximal").Key("id").Exists(),
					check.That(resourceType+".auth_strength_maximal").Key("display_name").HasValue("unit-test-auth-strength-max"),
					check.That(resourceType+".auth_strength_maximal").Key("description").HasValue("Unit test maximal authentication strength policy with all combinations and configurations"),
					check.That(resourceType+".auth_strength_maximal").Key("policy_type").HasValue("custom"),
					check.That(resourceType+".auth_strength_maximal").Key("requirements_satisfied").HasValue("mfa"),

					// Allowed combinations (all of them)
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.#").HasValue("22"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("deviceBasedPush"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("fido2"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("password,sms"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("windowsHelloForBusiness"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("x509CertificateMultiFactor"),
					check.That(resourceType+".auth_strength_maximal").Key("allowed_combinations.*").ContainsTypeSetElement("x509CertificateSingleFactor"),

					// Combination configurations
					check.That(resourceType+".auth_strength_maximal").Key("combination_configurations.#").HasValue("3"),
				),
			},
			{
				ResourceName:      resourceType + ".auth_strength_maximal",
				ImportState:       true,
				ImportStateVerify: true,
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
