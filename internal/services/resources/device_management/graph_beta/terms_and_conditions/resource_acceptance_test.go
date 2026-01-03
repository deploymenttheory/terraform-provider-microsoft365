package graphBetaTermsAndConditions_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaTermsAndConditions "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/terms_and_conditions"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourceType      = graphBetaTermsAndConditions.ResourceName
	groupResourceType = graphBetaGroup.ResourceName
	testResource      = graphBetaTermsAndConditions.TermsAndConditionsTestResource{}
	groupTestResource = graphBetaGroup.GroupTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccTermsAndConditionsResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "~> 3.6",
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			15*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: resourceType,
				TestResource: testResource,
			},
			destroy.ResourceTypeMapping{
				ResourceType: groupResourceType,
				TestResource: groupTestResource,
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("acc-test-terms-and-conditions-minimal"),
					check.That(resourceType+".test").Key("title").HasValue("Company Terms"),
					check.That(resourceType+".test").Key("body_text").HasValue("These are the basic terms and conditions."),
					check.That(resourceType+".test").Key("acceptance_statement").HasValue("I accept these terms"),
					check.That(resourceType+".test").Key("version").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("acc-test-terms-and-conditions-maximal"),
					check.That(resourceType+".test").Key("description").HasValue("Updated description for acceptance testing"),
					check.That(resourceType+".test").Key("title").HasValue("Complete Company Terms and Conditions"),
					check.That(resourceType+".test").Key("body_text").HasValue("These are the comprehensive terms and conditions that all users must read and accept before accessing company resources."),
					check.That(resourceType+".test").Key("acceptance_statement").HasValue("I have read and agree to abide by all terms and conditions outlined above"),
					check.That(resourceType+".test").Key("version").HasValue("2"),
					check.That(resourceType+".test").Key("assignments.#").HasValue("3"),
				),
			},
		},
	})
}

func TestAccTermsAndConditionsResource_Description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "~> 3.6",
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			0,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_description.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".description").Key("id").Exists(),
					check.That(resourceType+".description").Key("display_name").HasValue("acc-test-terms-and-conditions-description"),
					check.That(resourceType+".description").Key("description").HasValue("This is a test terms and conditions with description"),
				),
			},
		},
	})
}

func TestAccTermsAndConditionsResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "~> 3.6",
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			15*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: resourceType,
				TestResource: testResource,
			},
			destroy.ResourceTypeMapping{
				ResourceType: groupResourceType,
				TestResource: groupTestResource,
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("resource_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".assignments").Key("id").Exists(),
					check.That(resourceType+".assignments").Key("display_name").HasValue("acc-test-terms-and-conditions-assignments"),
					check.That(resourceType+".assignments").Key("description").HasValue("Terms and conditions policy with assignments for acceptance testing"),
					check.That(resourceType+".assignments").Key("assignments.#").HasValue("3"),
				),
			},
		},
	})
}
