package graphBetaIdentityAndAccessLifecycleWorkflow_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaIdentityAndAccessLifecycleWorkflow "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/lifecycle_workflow"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIdentityAndAccessLifecycleWorkflow.ResourceName

	// testResource is the test resource implementation for lifecycle workflows
	testResource = graphBetaIdentityAndAccessLifecycleWorkflow.LifecycleWorkflowTestResource{}

	// externalProviders are the providers acceptance fixtures depend on
	externalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
		"time": {
			Source:            "hashicorp/time",
			VersionConstraint: constants.ExternalProviderTimeVersion,
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

// sameIdAs fails when the resource ID changed, proving an update happened in place.
func sameIdAs(previous *string) func(value string) error {
	return func(value string) error {
		if value != *previous {
			return fmt.Errorf("%w: expected id %s, got %s", errIdChanged, *previous, value)
		}
		return nil
	}
}

func importStep(resourceName, description string) resource.TestStep {
	return resource.TestStep{
		PreConfig: func() {
			testlog.StepAction(resourceType, description)
		},
		ResourceName:            resourceName,
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// LW001: Joiner workflow with a time based trigger and rule scope
func TestAccResourceLifecycleWorkflow_01_Joiner(t *testing.T) {
	resourceName := resourceType + ".lw001_joiner"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating LW001 joiner workflow")
				},
				Config: loadAcceptanceTestTerraform("resource_lw001-joiner.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceName).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-lw001: joiner pre-hire onboarding [a-z0-9]{8}$`)),
					check.That(resourceName).Key("category").HasValue("joiner"),
					check.That(resourceName).Key("is_enabled").HasValue("false"),
					check.That(resourceName).Key("is_scheduling_enabled").HasValue("false"),
					check.That(resourceName).Key("version").HasValue("1"),
					check.That(resourceName).Key("execution_conditions.scope.rule").HasValue("(department eq 'acc-test-none')"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.attribute").HasValue("employeeHireDate"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.offset_in_days").HasValue("-7"),
					check.That(resourceName).Key("tasks.#").HasValue("2"),
					check.That(resourceName).Key("tasks.0.arguments.#").HasValue("1"),
					check.That(resourceName).Key("tasks.1.continue_on_error").HasValue("true"),
					check.That(resourceName).Key("tasks.1.execution_sequence").HasValue("2"),
				),
			},
			importStep(resourceName, "Importing LW001 joiner workflow"),
		},
	})
}

// LW002: Mover workflow with an attribute change trigger
func TestAccResourceLifecycleWorkflow_02_Mover(t *testing.T) {
	resourceName := resourceType + ".lw002_mover"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating LW002 mover workflow")
				},
				Config: loadAcceptanceTestTerraform("resource_lw002-mover.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("category").HasValue("mover"),
					check.That(resourceName).Key("execution_conditions.trigger.attribute_change.trigger_attribute").HasValue("department"),
					check.That(resourceName).Key("tasks.#").HasValue("1"),
				),
			},
			importStep(resourceName, "Importing LW002 mover workflow"),
		},
	})
}

// LW003: Leaver workflow on the employee leave date
func TestAccResourceLifecycleWorkflow_03_Leaver(t *testing.T) {
	resourceName := resourceType + ".lw003_leaver"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating LW003 leaver workflow")
				},
				Config: loadAcceptanceTestTerraform("resource_lw003-leaver.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("category").HasValue("leaver"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.attribute").HasValue("employeeLeaveDateTime"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.offset_in_days").HasValue("0"),
					check.That(resourceName).Key("tasks.#").HasValue("3"),
					check.That(resourceName).Key("tasks.2.execution_sequence").HasValue("3"),
				),
			},
			importStep(resourceName, "Importing LW003 leaver workflow"),
		},
	})
}

// LW004: Minimal on-demand workflow with defaults
func TestAccResourceLifecycleWorkflow_04_Minimal(t *testing.T) {
	resourceName := resourceType + ".lw004_minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating LW004 minimal workflow")
				},
				Config: loadAcceptanceTestTerraform("resource_lw004-minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("description").HasValue(""),
					check.That(resourceName).Key("execution_conditions.on_demand_execution_only").HasValue("true"),
					check.That(resourceName).Key("tasks.0.is_enabled").HasValue("true"),
					check.That(resourceName).Key("tasks.0.description").HasValue(""),
				),
			},
			importStep(resourceName, "Importing LW004 minimal workflow"),
		},
	})
}

// LW005: Maximal workflow with a group scope, membership change trigger and every task attribute
func TestAccResourceLifecycleWorkflow_05_Maximal(t *testing.T) {
	resourceName := resourceType + ".lw005_maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating LW005 maximal workflow")
				},
				Config: loadAcceptanceTestTerraform("resource_lw005-maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("category").HasValue("mover"),
					check.That(resourceName).Key("is_enabled").HasValue("true"),
					check.That(resourceName).Key("execution_conditions.scope.group_ids.#").HasValue("1"),
					check.That(resourceName).Key("execution_conditions.trigger.membership_change.change_type").HasValue("add"),
					check.That(resourceName).Key("tasks.#").HasValue("2"),
					check.That(resourceName).Key("tasks.0.description").HasValue("Add the user to the team group"),
					check.That(resourceName).Key("tasks.1.is_enabled").HasValue("false"),
					check.That(resourceName).Key("tasks.1.arguments.#").HasValue("2"),
				),
			},
			importStep(resourceName, "Importing LW005 maximal workflow"),
		},
	})
}

// LW006: Basic property changes are applied with PATCH and keep the version
func TestAccResourceLifecycleWorkflow_06_BasicPropertyUpdate(t *testing.T) {
	resourceName := resourceType + ".lw006_basic_update"
	var firstId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating LW006 workflow")
				},
				Config: loadAcceptanceTestTerraform("resource_lw006-basic_update_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					captureAttribute(resourceName, "id", &firstId),
					check.That(resourceName).Key("description").HasValue("Initial description"),
					check.That(resourceName).Key("version").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating LW006 basic properties")
				},
				Config: loadAcceptanceTestTerraform("resource_lw006-basic_update_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("id").ValidatesWith(sameIdAs(&firstId)),
					check.That(resourceName).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-lw006: basic property update renamed [a-z0-9]{8}$`)),
					check.That(resourceName).Key("description").HasValue("Updated description"),
					check.That(resourceName).Key("is_enabled").HasValue("true"),
					check.That(resourceName).Key("version").HasValue("1"),
				),
			},
			importStep(resourceName, "Importing LW006 workflow"),
		},
	})
}

// LW007: Task and trigger changes create a new version of the same workflow
func TestAccResourceLifecycleWorkflow_07_TaskChangeCreatesNewVersion(t *testing.T) {
	resourceName := resourceType + ".lw007_task_change"
	var firstId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating LW007 workflow")
				},
				Config: loadAcceptanceTestTerraform("resource_lw007-task_change_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					captureAttribute(resourceName, "id", &firstId),
					check.That(resourceName).Key("version").HasValue("1"),
					check.That(resourceName).Key("tasks.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Changing LW007 tasks and trigger (new version)")
				},
				Config: loadAcceptanceTestTerraform("resource_lw007-task_change_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("id").ValidatesWith(sameIdAs(&firstId)),
					check.That(resourceName).Key("version").HasValue("2"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.offset_in_days").HasValue("3"),
					check.That(resourceName).Key("tasks.#").HasValue("3"),
					check.That(resourceName).Key("tasks.0.task_definition_id").HasValue("70b29d51-b59a-4773-9280-8841dfd3f2ea"),
					check.That(resourceName).Key("tasks.1.is_enabled").HasValue("false"),
					check.That(resourceName).Key("tasks.2.task_definition_id").HasValue("3C860712-2D37-42A4-928F-5C93935D26A1"),
				),
			},
			importStep(resourceName, "Importing LW007 workflow"),
		},
	})
}
