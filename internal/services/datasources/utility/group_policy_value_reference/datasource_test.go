package utilityGroupPolicyValueReference_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityGroupPolicyValueReference "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/group_policy_value_reference"
	groupPolicyValueReferenceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/group_policy_value_reference/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *groupPolicyValueReferenceMocks.GroupPolicyValueReferenceMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	groupPolicyValueReferenceMock := &groupPolicyValueReferenceMocks.GroupPolicyValueReferenceMock{}
	groupPolicyValueReferenceMock.RegisterMocks()

	return mockClient, groupPolicyValueReferenceMock
}

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 01: Single Definition - RDP Policy
func TestGroupPolicyValueReferenceDataSource_01_SingleDefinition(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_single_definition.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("policy_name").HasValue("Allow users to connect remotely by using Remote Desktop Services"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.#").HasValue("1"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.id").HasValue("bb67ec37-f275-484c-942c-36a07e80add8"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.display_name").HasValue("Allow users to connect remotely by using Remote Desktop Services"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.class_type").HasValue("machine"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.category_path").HasValue("\\Windows Components\\Remote Desktop Services\\Remote Desktop Session Host\\Connections"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.explain_text").HasValue("Allows you to configure user access to Remote Desktop Services."),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.supported_on").HasValue("At least Windows Server 2008 or Windows Vista"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.policy_type").HasValue("admxBacked"),
				),
			},
		},
	})
}

// Test 02: Multiple Definitions - Show Home button (8 variants across browsers)
func TestGroupPolicyValueReferenceDataSource_02_MultipleDefinitions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_multiple_definitions.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("policy_name").HasValue("Show Home button on toolbar"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.#").HasValue("8"),

					// Verify first definition (Edge machine)
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.id").HasValue("eaca0db8-9673-4487-8055-d6dc037a3ef9"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.class_type").HasValue("machine"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.category_path").HasValue("\\Microsoft Edge\\Startup, home page and new tab page"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.policy_type").HasValue("admxIngested"),

					// Verify one user variant exists (Chrome)
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.3.class_type").HasValue("user"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.3.category_path").HasValue("\\Google\\Google Chrome\\Startup, Home page and New Tab page"),
				),
			},
		},
	})
}

// Test 03: No Results (with warning)
func TestGroupPolicyValueReferenceDataSource_03_NoResults(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_no_results.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("id").Exists(),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("policy_name").HasValue("Nonexistent Policy"),
					check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.#").HasValue("0"),
				),
			},
		},
	})
}

// Test 04: Fuzzy Match with Suggestions - Should return error with ranked suggestions
func TestGroupPolicyValueReferenceDataSource_04_FuzzyMatchSuggestions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("04_fuzzy_match_casing.tf"),
				ExpectError: regexp.MustCompile(`No exact match found for policy name 'Show Home button'`),
			},
		},
	})
}

// Test 05: Exact Match with Case Insensitivity (normalized matching)
func TestGroupPolicyValueReferenceDataSource_05_CaseInsensitiveMatch(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupMockEnvironment()
	defer httpmock.DeactivateAndReset()

	// Test that case variations match after normalization
	testCases := []struct {
		name           string
		policyName     string
		expectedResult string
	}{
		{
			name:           "lowercase",
			policyName:     "allow users to connect remotely by using remote desktop services",
			expectedResult: "Allow users to connect remotely by using Remote Desktop Services",
		},
		{
			name:           "UPPERCASE",
			policyName:     "ALLOW USERS TO CONNECT REMOTELY BY USING REMOTE DESKTOP SERVICES",
			expectedResult: "Allow users to connect remotely by using Remote Desktop Services",
		},
		{
			name:           "Mixed Case",
			policyName:     "AlLoW UsErS tO CoNnEcT ReMoTeLy By UsInG ReMoTe DeSkToP SeRvIcEs",
			expectedResult: "Allow users to connect remotely by using Remote Desktop Services",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
data "microsoft365_utility_group_policy_value_reference" "test" {
  policy_name = "%s"
}
`, tc.policyName),
						Check: resource.ComposeTestCheckFunc(
							check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.#").HasValue("1"),
							check.That("data."+utilityGroupPolicyValueReference.DataSourceName+".test").Key("definitions.0.display_name").HasValue(tc.expectedResult),
						),
					},
				},
			})
		})
	}
}
