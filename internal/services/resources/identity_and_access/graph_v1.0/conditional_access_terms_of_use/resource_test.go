package graphConditionalAccessTermsOfUse_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	termsOfUseMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_v1.0/conditional_access_terms_of_use/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// resourceType is declared in resource_acceptance_test.go and shared across the package

func setupMockEnvironment() (*mocks.Mocks, *termsOfUseMocks.ConditionalAccessTermsOfUseMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	termsOfUseMock := &termsOfUseMocks.ConditionalAccessTermsOfUseMock{}
	termsOfUseMock.RegisterMocks()
	return mockClient, termsOfUseMock
}

func TestConditionalAccessTermsOfUseResource_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsOfUseMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsOfUseMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("id").Exists(),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("display_name").HasValue("unit_test_conditional_access_terms_of_use_minimal"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("file.localizations.#").HasValue("1"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("is_viewing_before_acceptance_required").HasValue("true"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("is_per_device_acceptance_required").HasValue("false"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("user_reaccept_required_frequency").HasValue("P10D"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("terms_expiration.start_date_time").HasValue("2025-11-06"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("terms_expiration.frequency").HasValue("P180D"),
				),
			},
			{
				ResourceName:      resourceType + ".unit_test_conditional_access_terms_of_use_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"file",
					"file.%",
					"file.localizations",
					"file.localizations.#",
					"file.localizations.0",
					"file.localizations.0.%",
					"file.localizations.0.file_data",
					"file.localizations.0.file_data.%",
					"file.localizations.0.file_data.data",
					"file.localizations.0.display_name",
					"file.localizations.0.file_name",
					"file.localizations.0.is_default",
					"file.localizations.0.is_major_version",
					"file.localizations.0.language",
				},
			},
		},
	})
}

func TestConditionalAccessTermsOfUseResource_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsOfUseMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsOfUseMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("id").Exists(),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("display_name").HasValue("unit_test_conditional_access_terms_of_use_maximal"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("is_viewing_before_acceptance_required").HasValue("true"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("is_per_device_acceptance_required").HasValue("false"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("user_reaccept_required_frequency").HasValue("P10D"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("file.localizations.#").HasValue("30"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("terms_expiration.start_date_time").HasValue("2025-11-06"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("terms_expiration.frequency").HasValue("P180D"),
				),
			},
			{
				ResourceName:      resourceType + ".unit_test_conditional_access_terms_of_use_maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"file",
					"file.%",
					"file.localizations",
				},
			},
		},
	})
}

func TestConditionalAccessTermsOfUseResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, termsOfUseMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer termsOfUseMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("id").Exists(),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("display_name").HasValue("unit_test_conditional_access_terms_of_use_minimal"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("is_viewing_before_acceptance_required").HasValue("true"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_minimal").Key("file.localizations.#").HasValue("1"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("id").Exists(),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("display_name").HasValue("unit_test_conditional_access_terms_of_use_maximal"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("is_viewing_before_acceptance_required").HasValue("true"),
					check.That(resourceType+".unit_test_conditional_access_terms_of_use_maximal").Key("file.localizations.#").HasValue("30"),
				),
			},
			{
				ResourceName:      resourceType + ".unit_test_conditional_access_terms_of_use_maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"file",
					"file.%",
					"file.localizations",
				},
			},
		},
	})
}

func TestConditionalAccessTermsOfUseResource_FileValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
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

func TestConditionalAccessTermsOfUseResource_TermsExpiration(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
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
          data = "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"
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

func testConfigMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
	if err != nil {
		panic("failed to load resource_minimal.tf: " + err.Error())
	}
	return unitTestConfig
}

func testConfigMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
	if err != nil {
		panic("failed to load resource_maximal.tf: " + err.Error())
	}
	return unitTestConfig
}
