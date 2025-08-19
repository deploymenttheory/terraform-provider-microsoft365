package graphBetaTermsAndConditions_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	termsAndConditionsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/terms_and_conditions/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment() {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *termsAndConditionsMocks.TermsAndConditionsMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	termsAndConditionsMock := &termsAndConditionsMocks.TermsAndConditionsMock{}
	termsAndConditionsMock.RegisterMocks()

	return mockClient, termsAndConditionsMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *termsAndConditionsMocks.TermsAndConditionsMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	termsAndConditionsMock := &termsAndConditionsMocks.TermsAndConditionsMock{}
	termsAndConditionsMock.RegisterErrorMocks()

	return mockClient, termsAndConditionsMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestTermsAndConditionsResource_Schema validates the resource schema
func TestTermsAndConditionsResource_Schema(t *testing.T) {
	setupUnitTestEnvironment()
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "display_name", "Test Minimal Terms and Conditions - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "title", "Company Terms"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "body_text", "These are the basic terms and conditions."),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "acceptance_statement", "I accept these terms"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "description", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "version", "1"),
				),
			},
		},
	})
}

// TestTermsAndConditionsResource_Minimal tests basic CRUD operations
func TestTermsAndConditionsResource_Minimal(t *testing.T) {
	setupUnitTestEnvironment()
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_terms_and_conditions.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "display_name", "Test Minimal Terms and Conditions - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "title", "Company Terms"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "body_text", "These are the basic terms and conditions."),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "acceptance_statement", "I accept these terms"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_terms_and_conditions.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_terms_and_conditions.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.maximal", "display_name", "Test Maximal Terms and Conditions - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.maximal", "description", "Comprehensive terms and conditions for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.maximal", "version", "2"),
				),
			},
		},
	})
}

// TestTermsAndConditionsResource_UpdateInPlace tests in-place updates
func TestTermsAndConditionsResource_UpdateInPlace(t *testing.T) {
	setupUnitTestEnvironment()
	_, termsAndConditionsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsAndConditionsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_terms_and_conditions.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.minimal", "display_name", "Test Minimal Terms and Conditions - Unique"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_terms_and_conditions.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.maximal", "display_name", "Test Maximal Terms and Conditions - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.maximal", "description", "Comprehensive terms and conditions for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_terms_and_conditions.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestTermsAndConditionsResource_RequiredFields tests required field validation
func TestTermsAndConditionsResource_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment()
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
	setupUnitTestEnvironment()
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
	setupUnitTestEnvironment()
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
	setupUnitTestEnvironment()
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
	setupUnitTestEnvironment()
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
	setupUnitTestEnvironment()
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
	setupUnitTestEnvironment()
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
