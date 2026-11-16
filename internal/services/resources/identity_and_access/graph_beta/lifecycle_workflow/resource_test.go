package graphBetaIdentityAndAccessLifecycleWorkflow_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	lifecycleWorkflowMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/lifecycle_workflow/mocks"
)

var (
	errResourceNotInState = errors.New("resource not found in state")
	errIdChanged          = errors.New("workflow was replaced instead of updated in place")
)

func setupMockEnvironment() (*mocks.Mocks, *lifecycleWorkflowMocks.LifecycleWorkflowMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	lifecycleWorkflowMock := &lifecycleWorkflowMocks.LifecycleWorkflowMock{}
	lifecycleWorkflowMock.RegisterMocks()
	return mockClient, lifecycleWorkflowMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *lifecycleWorkflowMocks.LifecycleWorkflowMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	lifecycleWorkflowMock := &lifecycleWorkflowMocks.LifecycleWorkflowMock{}
	lifecycleWorkflowMock.RegisterErrorMocks()
	return mockClient, lifecycleWorkflowMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// captureAttribute stores the value of a resource attribute so a later step can compare against it.
func captureAttribute(resourceName, key string, target *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%w: %s", errResourceNotInState, resourceName)
		}
		*target = rs.Primary.Attributes[key]
		return nil
	}
}

// LW001: Joiner workflow with a time based trigger and rule scope
func TestUnitResourceLifecycleWorkflow_01_Joiner(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowMock.CleanupMockState()

	resourceName := resourceType + ".lw001_joiner"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_lw001-joiner.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceName).Key("category").HasValue("joiner"),
					check.That(resourceName).Key("display_name").HasValue("unit-test-lw001: joiner pre-hire onboarding"),
					check.That(resourceName).Key("description").HasValue("Enable the account and send a welcome email before day one"),
					check.That(resourceName).Key("is_enabled").HasValue("false"),
					check.That(resourceName).Key("is_scheduling_enabled").HasValue("false"),
					check.That(resourceName).Key("version").HasValue("1"),
					check.That(resourceName).Key("execution_conditions.on_demand_execution_only").HasValue("false"),
					check.That(resourceName).Key("execution_conditions.scope.rule").HasValue("(department eq 'acc-test-none')"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.attribute").HasValue("employeeHireDate"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.offset_in_days").HasValue("-7"),
					check.That(resourceName).Key("tasks.#").HasValue("2"),
					check.That(resourceName).Key("tasks.0.task_definition_id").HasValue("6fc52c9d-398b-4305-9763-15f42c1676fc"),
					check.That(resourceName).Key("tasks.0.execution_sequence").HasValue("1"),
					check.That(resourceName).Key("tasks.0.is_enabled").HasValue("true"),
					check.That(resourceName).Key("tasks.0.continue_on_error").HasValue("false"),
					check.That(resourceName).Key("tasks.0.arguments.#").HasValue("1"),
					check.That(resourceName).Key("tasks.0.arguments.0.name").HasValue("enableOnPremisesAccount"),
					check.That(resourceName).Key("tasks.0.arguments.0.value").HasValue("false"),
					check.That(resourceName).Key("tasks.1.task_definition_id").HasValue("70b29d51-b59a-4773-9280-8841dfd3f2ea"),
					check.That(resourceName).Key("tasks.1.execution_sequence").HasValue("2"),
					check.That(resourceName).Key("tasks.1.continue_on_error").HasValue("true"),
					check.That(resourceName).Key("tasks.1.description").HasValue(""),
					check.That(resourceName).Key("tasks.1.arguments.#").DoesNotExist(),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// LW002: Mover workflow with an attribute change trigger
func TestUnitResourceLifecycleWorkflow_02_Mover(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowMock.CleanupMockState()

	resourceName := resourceType + ".lw002_mover"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_lw002-mover.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("category").HasValue("mover"),
					check.That(resourceName).Key("execution_conditions.scope.rule").HasValue("(department eq 'acc-test-none')"),
					check.That(resourceName).Key("execution_conditions.trigger.attribute_change.trigger_attribute").HasValue("department"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.%").DoesNotExist(),
					check.That(resourceName).Key("tasks.#").HasValue("1"),
					check.That(resourceName).Key("tasks.0.task_definition_id").HasValue("aab41899-9972-422a-9d97-f626014578b7"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// LW003: Leaver workflow on the employee leave date
func TestUnitResourceLifecycleWorkflow_03_Leaver(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowMock.CleanupMockState()

	resourceName := resourceType + ".lw003_leaver"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_lw003-leaver.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("category").HasValue("leaver"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.attribute").HasValue("employeeLeaveDateTime"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.offset_in_days").HasValue("0"),
					check.That(resourceName).Key("tasks.#").HasValue("3"),
					check.That(resourceName).Key("tasks.0.task_definition_id").HasValue("1dfdfcc7-52fa-4c2e-bf3a-e3919cc12950"),
					check.That(resourceName).Key("tasks.1.task_definition_id").HasValue("8fa97d28-3e52-4985-b3a9-a1126f9b8b4e"),
					check.That(resourceName).Key("tasks.2.task_definition_id").HasValue("509589a4-0466-4471-829e-49c5e502bdee"),
					check.That(resourceName).Key("tasks.2.execution_sequence").HasValue("3"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// LW004: Minimal on-demand workflow with defaults
func TestUnitResourceLifecycleWorkflow_04_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowMock.CleanupMockState()

	resourceName := resourceType + ".lw004_minimal"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_lw004-minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("category").HasValue("leaver"),
					check.That(resourceName).Key("description").HasValue(""),
					check.That(resourceName).Key("is_enabled").HasValue("false"),
					check.That(resourceName).Key("is_scheduling_enabled").HasValue("false"),
					check.That(resourceName).Key("execution_conditions.on_demand_execution_only").HasValue("true"),
					check.That(resourceName).Key("execution_conditions.scope.%").DoesNotExist(),
					check.That(resourceName).Key("execution_conditions.trigger.%").DoesNotExist(),
					check.That(resourceName).Key("tasks.#").HasValue("1"),
					check.That(resourceName).Key("tasks.0.is_enabled").HasValue("true"),
					check.That(resourceName).Key("tasks.0.continue_on_error").HasValue("false"),
					check.That(resourceName).Key("tasks.0.description").HasValue(""),
					check.That(resourceName).Key("created_date_time").IsNotEmpty(),
					check.That(resourceName).Key("last_modified_date_time").IsNotEmpty(),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// LW005: Maximal workflow with a group scope, membership change trigger and every task attribute
func TestUnitResourceLifecycleWorkflow_05_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowMock.CleanupMockState()

	resourceName := resourceType + ".lw005_maximal"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_lw005-maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("category").HasValue("mover"),
					check.That(resourceName).Key("is_enabled").HasValue("true"),
					check.That(resourceName).Key("execution_conditions.scope.group_ids.#").HasValue("1"),
					check.That(resourceName).Key("execution_conditions.scope.group_ids.*").ContainsTypeSetElement("11111111-1111-1111-1111-111111111111"),
					check.That(resourceName).Key("execution_conditions.scope.rule").DoesNotExist(),
					check.That(resourceName).Key("execution_conditions.trigger.membership_change.change_type").HasValue("add"),
					check.That(resourceName).Key("tasks.#").HasValue("2"),
					check.That(resourceName).Key("tasks.0.description").HasValue("Add the user to the team group"),
					check.That(resourceName).Key("tasks.0.arguments.#").HasValue("1"),
					check.That(resourceName).Key("tasks.1.is_enabled").HasValue("false"),
					check.That(resourceName).Key("tasks.1.continue_on_error").HasValue("true"),
					check.That(resourceName).Key("tasks.1.arguments.#").HasValue("2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// LW006: Basic property changes are applied with PATCH and keep the version
func TestUnitResourceLifecycleWorkflow_06_BasicPropertyUpdate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowMock.CleanupMockState()

	resourceName := resourceType + ".lw006_basic_update"
	var firstId string

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_lw006-basic_update_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					captureAttribute(resourceName, "id", &firstId),
					check.That(resourceName).Key("display_name").HasValue("unit-test-lw006: basic property update"),
					check.That(resourceName).Key("description").HasValue("Initial description"),
					check.That(resourceName).Key("is_enabled").HasValue("false"),
					check.That(resourceName).Key("version").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_lw006-basic_update_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("id").ValidatesWith(func(value string) error {
						if value != firstId {
							return fmt.Errorf("%w: expected id %s, got %s", errIdChanged, firstId, value)
						}
						return nil
					}),
					check.That(resourceName).Key("display_name").HasValue("unit-test-lw006: basic property update renamed"),
					check.That(resourceName).Key("description").HasValue("Updated description"),
					check.That(resourceName).Key("is_enabled").HasValue("true"),
					check.That(resourceName).Key("version").HasValue("1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// LW007: Task and trigger changes create a new version of the same workflow
func TestUnitResourceLifecycleWorkflow_07_TaskChangeCreatesNewVersion(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowMock.CleanupMockState()

	resourceName := resourceType + ".lw007_task_change"
	var firstId string

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_lw007-task_change_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					captureAttribute(resourceName, "id", &firstId),
					check.That(resourceName).Key("version").HasValue("1"),
					check.That(resourceName).Key("tasks.#").HasValue("2"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.offset_in_days").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_lw007-task_change_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("id").ValidatesWith(func(value string) error {
						if value != firstId {
							return fmt.Errorf("%w: expected id %s, got %s", errIdChanged, firstId, value)
						}
						return nil
					}),
					check.That(resourceName).Key("version").HasValue("2"),
					check.That(resourceName).Key("execution_conditions.trigger.time_based_attribute.offset_in_days").HasValue("3"),
					check.That(resourceName).Key("tasks.#").HasValue("3"),
					check.That(resourceName).Key("tasks.0.task_definition_id").HasValue("70b29d51-b59a-4773-9280-8841dfd3f2ea"),
					check.That(resourceName).Key("tasks.0.execution_sequence").HasValue("1"),
					check.That(resourceName).Key("tasks.0.continue_on_error").HasValue("true"),
					check.That(resourceName).Key("tasks.1.task_definition_id").HasValue("6fc52c9d-398b-4305-9763-15f42c1676fc"),
					check.That(resourceName).Key("tasks.1.is_enabled").HasValue("false"),
					check.That(resourceName).Key("tasks.2.task_definition_id").HasValue("3C860712-2D37-42A4-928F-5C93935D26A1"),
					check.That(resourceName).Key("tasks.2.execution_sequence").HasValue("3"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// LW008: Graph errors surface through the Kiota error handler
func TestUnitResourceLifecycleWorkflow_08_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_lw001-joiner.tf"),
				ExpectError: regexp.MustCompile(`(?i)forbidden|permission|403`),
			},
		},
	})
}

// LW009: Invalid configurations fail at plan time
func TestUnitResourceLifecycleWorkflow_09_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_lw900-invalid_category.tf"),
				ExpectError: regexp.MustCompile(`value must be one of`),
			},
			{
				Config:      loadUnitTestTerraform("resource_lw901-on_demand_with_scope.tf"),
				ExpectError: regexp.MustCompile(`scope and trigger must not be set when\s+on_demand_execution_only is true`),
			},
			{
				Config:      loadUnitTestTerraform("resource_lw902-scope_without_trigger.tf"),
				ExpectError: regexp.MustCompile(`scope and trigger are both required unless\s+on_demand_execution_only is true`),
			},
		},
	})
}
