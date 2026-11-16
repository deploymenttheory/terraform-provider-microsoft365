package graphBetaM365AdminProfileCardProperty_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	profileCardPropertyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/m365_admin/graph_beta/profile_card_property/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *profileCardPropertyMocks.ProfileCardPropertyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	profileCardPropertyMock := &profileCardPropertyMocks.ProfileCardPropertyMock{}
	profileCardPropertyMock.RegisterMocks()
	return mockClient, profileCardPropertyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *profileCardPropertyMocks.ProfileCardPropertyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	profileCardPropertyMock := &profileCardPropertyMocks.ProfileCardPropertyMock{}
	profileCardPropertyMock.RegisterErrorMocks()
	return mockClient, profileCardPropertyMock
}

func setupForbiddenMockEnvironment() (*mocks.Mocks, *profileCardPropertyMocks.ProfileCardPropertyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	profileCardPropertyMock := &profileCardPropertyMocks.ProfileCardPropertyMock{}
	profileCardPropertyMock.RegisterForbiddenMocks()
	return mockClient, profileCardPropertyMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Minimal: directory property name only; is_visible defaults to true and no annotations are sent
func TestUnitResourceProfileCardProperty_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileCardPropertyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileCardPropertyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").HasValue("CustomAttribute1"),
					check.That(resourceType+".minimal").Key("directory_property_name").HasValue("CustomAttribute1"),
					check.That(resourceType+".minimal").Key("is_visible").HasValue("true"),
					check.That(resourceType+".minimal").Key("annotations.#").DoesNotExist(),
				),
			},
			{
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Maximal: hidden, with a custom label and two localized labels
func TestUnitResourceProfileCardProperty_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileCardPropertyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileCardPropertyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").HasValue("CustomAttribute2"),
					check.That(resourceType+".maximal").Key("directory_property_name").HasValue("CustomAttribute2"),
					check.That(resourceType+".maximal").Key("is_visible").HasValue("false"),
					check.That(resourceType+".maximal").Key("annotations.#").HasValue("1"),
					check.That(resourceType+".maximal").Key("annotations.0.display_name").HasValue("Cost Center"),
					check.That(resourceType+".maximal").Key("annotations.0.localizations.#").HasValue("2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".maximal", "annotations.0.localizations.*", map[string]string{"language_tag": "de-DE", "display_name": "Kostenstelle"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".maximal", "annotations.0.localizations.*", map[string]string{"language_tag": "nb-NO", "display_name": "Kostnadssenter"}),
				),
			},
			{
				ResourceName:            resourceType + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Minimal to maximal and back: in-place updates, including clearing all annotations
func TestUnitResourceProfileCardProperty_03_MinimalToMaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileCardPropertyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileCardPropertyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("is_visible").HasValue("true"),
					check.That(resourceType+".minimal").Key("annotations.#").DoesNotExist(),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").HasValue("CustomAttribute1"),
					check.That(resourceType+".minimal").Key("is_visible").HasValue("false"),
					check.That(resourceType+".minimal").Key("annotations.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("annotations.0.display_name").HasValue("Cost Center"),
					check.That(resourceType+".minimal").Key("annotations.0.localizations.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("is_visible").HasValue("true"),
					check.That(resourceType+".minimal").Key("annotations.#").DoesNotExist(),
				),
			},
			{
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Annotation without localizations: Graph returns an empty localizations collection, which must not drift
func TestUnitResourceProfileCardProperty_04_AnnotationWithoutLocalizations(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileCardPropertyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileCardPropertyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_annotation_without_localizations.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".label_only").Key("annotations.#").HasValue("1"),
					check.That(resourceType+".label_only").Key("annotations.0.display_name").HasValue("Employee Number"),
					check.That(resourceType+".label_only").Key("annotations.0.localizations").DoesNotExist(),
				),
			},
			{
				ResourceName:            resourceType + ".label_only",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Import of a pre-existing property by its exact, case-sensitive key, followed by an empty plan
func TestUnitResourceProfileCardProperty_05_ImportExistingCaseSensitive(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileCardPropertyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileCardPropertyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             loadUnitTestTerraform("resource_existing_import.tf"),
				ResourceName:       resourceType + ".existing",
				ImportState:        true,
				ImportStateId:      "customAttribute9",
				ImportStatePersist: true,
			},
			{
				Config:   loadUnitTestTerraform("resource_existing_import.tf"),
				PlanOnly: true,
			},
		},
	})
}

// Import of a key with the wrong casing: Graph keys are case-sensitive, so the property is not found
func TestUnitResourceProfileCardProperty_06_ImportWrongCaseNotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileCardPropertyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileCardPropertyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        loadUnitTestTerraform("resource_existing_import.tf"),
				ResourceName:  resourceType + ".existing",
				ImportState:   true,
				ImportStateId: "CustomAttribute9",
				ExpectError:   regexp.MustCompile(`Cannot import non-existent remote object`),
			},
		},
	})
}

// Invalid directory property name: rejected by schema validation at plan time
func TestUnitResourceProfileCardProperty_07_ValidationErrorInvalidPropertyName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileCardPropertyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileCardPropertyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_property_name.tf"),
				ExpectError: regexp.MustCompile(`Attribute directory_property_name value must be one of`),
			},
		},
	})
}

// Graph error: a 400 on create surfaces as a Terraform error
func TestUnitResourceProfileCardProperty_08_ErrorHandlingBadRequest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileCardPropertyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileCardPropertyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Invalid request|Bad Request|400`),
			},
		},
	})
}

// Graph error: a 403 on create surfaces as a Terraform error naming the write permission
func TestUnitResourceProfileCardProperty_09_ErrorHandlingForbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, profileCardPropertyMock := setupForbiddenMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer profileCardPropertyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`PeopleSettings\.ReadWrite\.All`),
			},
		},
	})
}
