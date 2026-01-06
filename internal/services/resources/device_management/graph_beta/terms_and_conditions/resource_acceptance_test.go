package graphBetaTermsAndConditions_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaTermsAndConditions "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/terms_and_conditions"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

// TestAccTermsAndConditionsResource_Minimal tests creating a minimal terms and conditions
func TestAccTermsAndConditionsResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 15*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "~> 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal terms and conditions")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("display_name").HasValue("acc-test-terms-and-conditions-minimal"),
					check.That(resourceType+".test").Key("title").HasValue("Company Terms"),
					check.That(resourceType+".test").Key("body_text").HasValue("These are the basic terms and conditions."),
					check.That(resourceType+".test").Key("acceptance_statement").HasValue("I accept these terms"),
					check.That(resourceType+".test").Key("version").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal terms and conditions")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTermsAndConditionsResource_Maximal tests creating a maximal terms and conditions
func TestAccTermsAndConditionsResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal terms and conditions")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("display_name").HasValue("acc-test-terms-and-conditions-maximal"),
					check.That(resourceType+".test").Key("description").HasValue("Updated description for acceptance testing"),
					check.That(resourceType+".test").Key("title").HasValue("Complete Company Terms and Conditions"),
					check.That(resourceType+".test").Key("body_text").HasValue("These are the comprehensive terms and conditions that all users must read and accept before accessing company resources."),
					check.That(resourceType+".test").Key("acceptance_statement").HasValue("I have read and agree to abide by all terms and conditions outlined above"),
					check.That(resourceType+".test").Key("version").HasValue("1"),
					check.That(resourceType+".test").Key("assignments.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("terms and conditions assignments", 10*time.Second)
					time.Sleep(10 * time.Second)
					testlog.StepAction(resourceType, "Importing maximal terms and conditions")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTermsAndConditionsResource_MinimalAssignment tests creating terms and conditions with minimal assignment
func TestAccTermsAndConditionsResource_MinimalAssignment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 15*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating terms and conditions with minimal assignment")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal_assignment").ExistsInGraph(testResource),
					check.That(resourceType+".minimal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal_assignment").Key("display_name").HasValue("acc-test-terms-and-conditions-minimal-assignment"),
					check.That(resourceType+".minimal_assignment").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".minimal_assignment").Key("assignments.0.type").HasValue("allLicensedUsersAssignmentTarget"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("terms and conditions assignment", 10*time.Second)
					time.Sleep(10 * time.Second)
					testlog.StepAction(resourceType, "Importing terms and conditions with minimal assignment")
				},
				ResourceName:      resourceType + ".minimal_assignment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTermsAndConditionsResource_MaximalAssignment tests creating terms and conditions with maximal assignments
func TestAccTermsAndConditionsResource_MaximalAssignment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			20*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: resourceType,
				TestResource: testResource,
			},
			destroy.ResourceTypeMapping{
				ResourceType: groupResourceType,
				TestResource: groupTestResource,
			},
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating terms and conditions with maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("terms and conditions with groups", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".maximal_assignment").ExistsInGraph(testResource),
					check.That(resourceType+".maximal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal_assignment").Key("display_name").HasValue("acc-test-terms-and-conditions-maximal-assignment"),
					check.That(resourceType+".maximal_assignment").Key("description").HasValue("Terms and conditions with comprehensive assignments for acceptance testing"),
					check.That(resourceType+".maximal_assignment").Key("assignments.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("terms and conditions assignments", 15*time.Second)
					time.Sleep(15 * time.Second)
					testlog.StepAction(resourceType, "Importing terms and conditions with maximal assignments")
				},
				ResourceName:      resourceType + ".maximal_assignment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTermsAndConditionsResource_MinimalToMaximal tests transitioning from minimal to maximal configuration
func TestAccTermsAndConditionsResource_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal configuration for transition test")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("display_name").HasValue("acc-test-terms-and-conditions-minimal"),
					check.That(resourceType+".test").Key("version").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("terms and conditions before update", 10*time.Second)
					time.Sleep(10 * time.Second)
					testlog.StepAction(resourceType, "Updating to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".transition").ExistsInGraph(testResource),
					check.That(resourceType+".transition").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".transition").Key("display_name").HasValue("acc-test-terms-and-conditions-transition"),
					check.That(resourceType+".transition").Key("description").HasValue("Configuration that transitions from minimal to maximal for acceptance testing"),
					check.That(resourceType+".transition").Key("version").HasValue("1"),
					check.That(resourceType+".transition").Key("assignments.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("terms and conditions after update", 10*time.Second)
					time.Sleep(10 * time.Second)
					testlog.StepAction(resourceType, "Importing transitioned configuration")
				},
				ResourceName:      resourceType + ".transition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTermsAndConditionsResource_MaximalToMinimal tests transitioning from maximal to minimal configuration
func TestAccTermsAndConditionsResource_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal configuration for transition test")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("terms and conditions with groups", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".transition").ExistsInGraph(testResource),
					check.That(resourceType+".transition").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".transition").Key("display_name").HasValue("acc-test-terms-and-conditions-transition"),
					check.That(resourceType+".transition").Key("description").HasValue("Configuration that transitions from minimal to maximal for acceptance testing"),
					check.That(resourceType+".transition").Key("version").HasValue("1"),
					check.That(resourceType+".transition").Key("assignments.#").HasValue("3"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("terms and conditions before downgrade", 10*time.Second)
					time.Sleep(10 * time.Second)
					testlog.StepAction(resourceType, "Downgrading to minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal_to_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".transition").ExistsInGraph(testResource),
					check.That(resourceType+".transition").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".transition").Key("display_name").HasValue("acc-test-terms-and-conditions-transition"),
					check.That(resourceType+".transition").Key("description").IsEmpty(),
					check.That(resourceType+".transition").Key("version").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("terms and conditions after downgrade", 10*time.Second)
					time.Sleep(10 * time.Second)
					testlog.StepAction(resourceType, "Importing downgraded configuration")
				},
				ResourceName:      resourceType + ".transition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTermsAndConditionsResource_Description tests creating terms and conditions with description attribute
func TestAccTermsAndConditionsResource_Description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 0),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating terms and conditions with description")
				},
				Config: loadAcceptanceTestTerraform("resource_description.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".description").ExistsInGraph(testResource),
					check.That(resourceType+".description").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".description").Key("display_name").HasValue("acc-test-terms-and-conditions-description"),
					check.That(resourceType+".description").Key("description").HasValue("This is a test terms and conditions with description"),
				),
			},
		},
	})
}
