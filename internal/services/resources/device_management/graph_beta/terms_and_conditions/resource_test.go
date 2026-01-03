package graphBetaTermsAndConditions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	termsAndConditionsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/terms_and_conditions/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *termsAndConditionsMocks.TermsAndConditionsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	termsAndConditionsMock := &termsAndConditionsMocks.TermsAndConditionsMock{}
	termsAndConditionsMock.RegisterMocks()
	return mockClient, termsAndConditionsMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *termsAndConditionsMocks.TermsAndConditionsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	termsAndConditionsMock := &termsAndConditionsMocks.TermsAndConditionsMock{}
	termsAndConditionsMock.RegisterErrorMocks()
	return mockClient, termsAndConditionsMock
}

func testConfigHelper(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// TestTermsAndConditionsResource_Schema validates the resource schema
func TestTermsAndConditionsResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-terms-and-conditions-minimal"),
					check.That(resourceType+".minimal").Key("title").HasValue("Company Terms"),
					check.That(resourceType+".minimal").Key("body_text").HasValue("These are the basic terms and conditions."),
					check.That(resourceType+".minimal").Key("acceptance_statement").HasValue("I accept these terms"),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("description").HasValue(""),
					check.That(resourceType+".minimal").Key("version").HasValue("1"),
				),
			},
		},
	})
}

// TestTermsAndConditionsResource_Minimal tests basic CRUD operations
func TestTermsAndConditionsResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-terms-and-conditions-minimal"),
					check.That(resourceType+".minimal").Key("title").HasValue("Company Terms"),
					check.That(resourceType+".minimal").Key("body_text").HasValue("These are the basic terms and conditions."),
					check.That(resourceType+".minimal").Key("acceptance_statement").HasValue("I accept these terms"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigHelper("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("display_name").HasValue("unit-test-terms-and-conditions-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("Comprehensive terms and conditions for testing with all features"),
					check.That(resourceType+".maximal").Key("version").HasValue("2"),
				),
			},
		},
	})
}

// TestTermsAndConditionsResource_UpdateInPlace tests in-place updates
func TestTermsAndConditionsResource_UpdateInPlace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-terms-and-conditions-minimal"),
				),
			},
			{
				Config: testConfigHelper("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("display_name").HasValue("unit-test-terms-and-conditions-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("Comprehensive terms and conditions for testing with all features"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.#").HasValue("2"),
				),
			},
		},
	})
}

// TestTermsAndConditionsResource_RequiredFields tests required field validation
func TestTermsAndConditionsResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	testCases := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "missing display_name",
			config: `
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  title               = "Test Terms"
  body_text           = "Test body"
  acceptance_statement = "I accept"
}
`,
			expectedError: `The argument "display_name" is required`,
		},
		{
			name: "missing title",
			config: `
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name        = "Test Terms and Conditions"
  body_text           = "Test body"
  acceptance_statement = "I accept"
}
`,
			expectedError: `The argument "title" is required`,
		},
		{
			name: "missing body_text",
			config: `
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name        = "Test Terms and Conditions"
  title               = "Test Terms"
  acceptance_statement = "I accept"
}
`,
			expectedError: `The argument "body_text" is required`,
		},
		{
			name: "missing acceptance_statement",
			config: `
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name = "Test Terms and Conditions"
  title        = "Test Terms"
  body_text    = "Test body"
}
`,
			expectedError: `The argument "acceptance_statement" is required`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      tc.config,
						ExpectError: regexp.MustCompile(tc.expectedError),
					},
				},
			})
		})
	}
}

// TestTermsAndConditionsResource_ErrorHandling tests error scenarios
func TestTermsAndConditionsResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name         = "Test Terms and Conditions"
  title               = "Test Terms"
  body_text           = "Test body text"
  acceptance_statement = "I accept these terms"
}
`,
				ExpectError: regexp.MustCompile(`Invalid terms and conditions data|BadRequest`),
			},
		},
	})
}

// TestTermsAndConditionsResource_DescriptionValidation tests description length validation
func TestTermsAndConditionsResource_DescriptionValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	// Create a description longer than 1500 characters
	longDescription := ""
	for i := 0; i < 151; i++ {
		longDescription += "0123456789"
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name         = "Test Terms and Conditions"
  description          = "%s"
  title               = "Test Terms"
  body_text           = "Test body text"
  acceptance_statement = "I accept these terms"
}
`, longDescription),
				ExpectError: regexp.MustCompile(`Attribute description string length must be at most 1500`),
			},
		},
	})
}

// TestTermsAndConditionsResource_BodyTextValidation tests body text length validation
func TestTermsAndConditionsResource_BodyTextValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	// Create body text longer than 60000 characters
	longBodyText := ""
	for i := 0; i < 6001; i++ {
		longBodyText += "0123456789"
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name         = "Test Terms and Conditions"
  title               = "Test Terms"
  body_text           = "%s"
  acceptance_statement = "I accept these terms"
}
`, longBodyText),
				ExpectError: regexp.MustCompile(`Attribute body_text string length must be at most 60000`),
			},
		},
	})
}

// TestTermsAndConditionsResource_AcceptanceStatementValidation tests acceptance statement length validation
func TestTermsAndConditionsResource_AcceptanceStatementValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	// Create acceptance statement longer than 500 characters
	longAcceptanceStatement := ""
	for i := 0; i < 51; i++ {
		longAcceptanceStatement += "0123456789"
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name         = "Test Terms and Conditions"
  title               = "Test Terms"
  body_text           = "Test body text"
  acceptance_statement = "%s"
}
`, longAcceptanceStatement),
				ExpectError: regexp.MustCompile(`Attribute acceptance_statement string length must be at most 500`),
			},
		},
	})
}

// TestTermsAndConditionsResource_RoleScopeTagIds tests role scope tag IDs handling
func TestTermsAndConditionsResource_RoleScopeTagIds(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name         = "Test Terms and Conditions with Role Scope Tags"
  title               = "Test Terms"
  body_text           = "Test body text"
  acceptance_statement = "I accept these terms"
  role_scope_tag_ids  = ["0", "1", "2"]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "role_scope_tag_ids.*", "2"),
				),
			},
		},
	})
}

// TestTermsAndConditionsResource_VersionHandling tests version handling
func TestTermsAndConditionsResource_VersionHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name         = "Test Terms and Conditions with Version"
  title               = "Test Terms"
  body_text           = "Test body text"
  acceptance_statement = "I accept these terms"
  version             = 3
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "display_name", "Test Terms and Conditions with Version"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.test", "version", "3"),
				),
			},
		},
	})
}
