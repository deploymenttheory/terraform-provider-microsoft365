package graphBetaM365AdminProfileCardProperty_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaM365AdminProfileCardProperty "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/m365_admin/graph_beta/profile_card_property"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaM365AdminProfileCardProperty.ResourceName

	// testResource is the test resource implementation for profile card properties
	testResource = graphBetaM365AdminProfileCardProperty.ProfileCardPropertyTestResource{}

	// externalProviders are the providers the acceptance configs use for unique labels
	externalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
	}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// requireUnusedPropertyName stops the test when the tenant already has a profile card property for the key.
// Profile card properties are keyed by directoryPropertyName, so a test may only create and delete keys it owns.
func requireUnusedPropertyName(t *testing.T, directoryPropertyName string) {
	t.Helper()
	exists, err := testResource.Exists(context.Background(), nil, &terraform.InstanceState{ID: directoryPropertyName})
	if err != nil {
		t.Fatalf("could not check whether profile card property %s exists: %v", directoryPropertyName, err)
	}
	if exists != nil && *exists {
		t.Fatalf("profile card property %s already exists in the tenant; refusing to run because the test would modify and delete it", directoryPropertyName)
	}
}

// Minimal: directory property name only; is_visible defaults to true
func TestAccResourceProfileCardProperty_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			requireUnusedPropertyName(t, "CustomAttribute13")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal profile card property")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").HasValue("CustomAttribute13"),
					check.That(resourceType+".minimal").Key("directory_property_name").HasValue("CustomAttribute13"),
					check.That(resourceType+".minimal").Key("is_visible").HasValue("true"),
					check.That(resourceType+".minimal").Key("annotations.#").DoesNotExist(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal profile card property")
				},
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Maximal: hidden, with a custom label and two localized labels
func TestAccResourceProfileCardProperty_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			requireUnusedPropertyName(t, "CustomAttribute14")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal profile card property")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").HasValue("CustomAttribute14"),
					check.That(resourceType+".maximal").Key("is_visible").HasValue("false"),
					check.That(resourceType+".maximal").Key("annotations.#").HasValue("1"),
					check.That(resourceType+".maximal").Key("annotations.0.display_name").MatchesRegex(regexp.MustCompile(`^acc-test-profile-card-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".maximal").Key("annotations.0.localizations.#").HasValue("2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".maximal", "annotations.0.localizations.*", map[string]string{"language_tag": "de-DE"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".maximal", "annotations.0.localizations.*", map[string]string{"language_tag": "nb-NO"}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal profile card property")
				},
				ResourceName:            resourceType + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Minimal to maximal and back: in-place updates, including clearing all annotations
func TestAccResourceProfileCardProperty_03_MinimalToMaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			requireUnusedPropertyName(t, "CustomAttribute15")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal profile card property")
				},
				Config: loadAcceptanceTestTerraform("resource_update_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update").ExistsInGraph(testResource),
					check.That(resourceType+".update").Key("is_visible").HasValue("true"),
					check.That(resourceType+".update").Key("annotations.#").DoesNotExist(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating profile card property to maximal in place")
				},
				Config: loadAcceptanceTestTerraform("resource_update_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update").Key("id").HasValue("CustomAttribute15"),
					check.That(resourceType+".update").Key("is_visible").HasValue("false"),
					check.That(resourceType+".update").Key("annotations.#").HasValue("1"),
					check.That(resourceType+".update").Key("annotations.0.display_name").MatchesRegex(regexp.MustCompile(`^acc-test-profile-card-updated-[a-z0-9]{8}$`)),
					check.That(resourceType+".update").Key("annotations.0.localizations.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating profile card property back to minimal in place")
				},
				Config: loadAcceptanceTestTerraform("resource_update_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".update").Key("is_visible").HasValue("true"),
					check.That(resourceType+".update").Key("annotations.#").DoesNotExist(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing updated profile card property")
				},
				ResourceName:            resourceType + ".update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Invalid directory property name: rejected by schema validation before any write
func TestAccResourceProfileCardProperty_04_ValidationErrorInvalidPropertyName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating profile card property with an invalid directory property name")
				},
				Config:      loadAcceptanceTestTerraform("resource_invalid_property_name.tf"),
				ExpectError: regexp.MustCompile(`Attribute directory_property_name value must be one of`),
			},
		},
	})
}
