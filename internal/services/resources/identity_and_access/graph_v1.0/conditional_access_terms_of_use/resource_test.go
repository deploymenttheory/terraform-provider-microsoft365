package graphConditionalAccessTermsOfUse_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	termsOfUseMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_v1.0/conditional_access_terms_of_use/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment() {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *termsOfUseMocks.ConditionalAccessTermsOfUseMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	termsOfUseMock := &termsOfUseMocks.ConditionalAccessTermsOfUseMock{}
	termsOfUseMock.RegisterMocks()

	return mockClient, termsOfUseMock
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

// TestConditionalAccessTermsOfUseResource_Schema validates the resource schema
func TestConditionalAccessTermsOfUseResource_Schema(t *testing.T) {
	setupUnitTestEnvironment()
	_, termsOfUseMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsOfUseMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal", "display_name", "Minimal Terms of Use"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal", "file.localizations.#", "1"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

// TestConditionalAccessTermsOfUseResource_Minimal tests basic CRUD operations
func TestConditionalAccessTermsOfUseResource_Minimal(t *testing.T) {
	setupUnitTestEnvironment()
	_, termsOfUseMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsOfUseMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal", "display_name", "Minimal Terms of Use"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal", "is_viewing_before_acceptance_required", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal", "is_per_device_acceptance_required", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal"]
					if !ok {
						return "", fmt.Errorf("not found: microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal")
					}
					return rs.Primary.ID, nil
				},
			},
		},
	})
}

// TestConditionalAccessTermsOfUseResource_Maximal tests maximal configuration
func TestConditionalAccessTermsOfUseResource_Maximal(t *testing.T) {
	setupUnitTestEnvironment()
	_, termsOfUseMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsOfUseMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.maximal", "display_name", "Maximal Terms of Use Agreement"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.maximal", "is_viewing_before_acceptance_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.maximal", "is_per_device_acceptance_required", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.maximal", "user_reaccept_required_frequency", "P90D"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.maximal", "file.localizations.#", "2"),
				),
			},
		},
	})
}

// TestConditionalAccessTermsOfUseResource_UpdateInPlace tests in-place updates
func TestConditionalAccessTermsOfUseResource_UpdateInPlace(t *testing.T) {
	setupUnitTestEnvironment()
	_, termsOfUseMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsOfUseMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal", "display_name", "Minimal Terms of Use"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.minimal", "is_viewing_before_acceptance_required", "false"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.maximal", "display_name", "Maximal Terms of Use Agreement"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.maximal", "is_viewing_before_acceptance_required", "true"),
				),
			},
		},
	})
}

// TestConditionalAccessTermsOfUseResource_FileValidation tests file configuration validation
func TestConditionalAccessTermsOfUseResource_FileValidation(t *testing.T) {
	setupUnitTestEnvironment()
	_, termsOfUseMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsOfUseMock.CleanupMockState()

	testCases := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "missing_file_configuration",
			config: `
resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "test" {
  display_name                          = "Test Terms of Use"
  is_viewing_before_acceptance_required = false
  is_per_device_acceptance_required     = false
}
`,
			expectedError: `Missing required argument`,
		},
		{
			name: "empty_localizations",
			config: `
resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "test" {
  display_name                          = "Test Terms of Use"
  is_viewing_before_acceptance_required = false
  is_per_device_acceptance_required     = false
  
  file = {
    localizations = []
  }
}
`,
			expectedError: `set must contain at least 1 elements`,
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

// TestConditionalAccessTermsOfUseResource_TermsExpiration tests terms expiration configuration
func TestConditionalAccessTermsOfUseResource_TermsExpiration(t *testing.T) {
	setupUnitTestEnvironment()
	_, termsOfUseMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsOfUseMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "test" {
  display_name                          = "Terms with Expiration"
  is_viewing_before_acceptance_required = false
  is_per_device_acceptance_required     = false
  
  terms_expiration = {
    start_date_time = "2025-12-31"
    frequency       = "P365D"
  }
  
  file = {
    localizations = [
      {
        file_name        = "terms.pdf"
        display_name     = "Terms of Use"
        language         = "en-US"
        is_default       = true
        is_major_version = false
        file_data = {
          data = "%PDF-1.4\nTest content"
        }
      }
    ]
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "display_name", "Terms with Expiration"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "terms_expiration.start_date_time", "2025-12-31"),
					resource.TestCheckResourceAttr("microsoft365_graph_identity_and_access_conditional_access_terms_of_use.test", "terms_expiration.frequency", "P365D"),
				),
			},
		},
	})
}
