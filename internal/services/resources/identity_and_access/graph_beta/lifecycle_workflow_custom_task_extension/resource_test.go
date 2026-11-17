package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension_test

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
	customTaskExtensionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/lifecycle_workflow_custom_task_extension/mocks"
)

var (
	errResourceNotInState = errors.New("resource not found in state")
	errIdChanged          = errors.New("custom task extension was replaced instead of updated in place")
)

func setupMockEnvironment() (*mocks.Mocks, *customTaskExtensionMocks.LifecycleWorkflowCustomTaskExtensionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	customTaskExtensionMock := &customTaskExtensionMocks.LifecycleWorkflowCustomTaskExtensionMock{}
	customTaskExtensionMock.RegisterMocks()
	return mockClient, customTaskExtensionMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *customTaskExtensionMocks.LifecycleWorkflowCustomTaskExtensionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	customTaskExtensionMock := &customTaskExtensionMocks.LifecycleWorkflowCustomTaskExtensionMock{}
	customTaskExtensionMock.RegisterErrorMocks()
	return mockClient, customTaskExtensionMock
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

// sameIdAs fails when the resource ID changed, proving an update happened in place.
func sameIdAs(previous *string) func(value string) error {
	return func(value string) error {
		if value != *previous {
			return fmt.Errorf("%w: expected id %s, got %s", errIdChanged, *previous, value)
		}
		return nil
	}
}

// CTX001: Minimal custom task extension with only required attributes
func TestUnitResourceLifecycleWorkflowCustomTaskExtension_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, customTaskExtensionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer customTaskExtensionMock.CleanupMockState()

	resourceName := resourceType + ".ctx001_minimal"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_ctx001-minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceName).Key("display_name").HasValue("unit-test-ctx001-minimal"),
					check.That(resourceName).Key("description").HasValue(""),
					check.That(resourceName).Key("endpoint_configuration.subscription_id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceName).Key("endpoint_configuration.resource_group_name").HasValue("acc-test-rg-placeholder"),
					check.That(resourceName).Key("endpoint_configuration.logic_app_workflow_name").HasValue("acc-test-la-placeholder"),
					check.That(resourceName).Key("endpoint_configuration.url").MatchesRegex(regexp.MustCompile(`^https://prod-00\.westeurope\.logic\.azure\.com`)),
					check.That(resourceName).Key("authentication_configuration.type").HasValue("azureAdPopTokenAuthentication"),
					check.That(resourceName).Key("authentication_configuration.resource_id").DoesNotExist(),
					check.That(resourceName).Key("client_configuration.maximum_retries").HasValue("1"),
					check.That(resourceName).Key("client_configuration.timeout_in_milliseconds").HasValue("1000"),
					check.That(resourceName).Key("callback_configuration.%").DoesNotExist(),
					check.That(resourceName).Key("reply_mode").HasValue("none"),
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

// CTX002: Maximal custom task extension with token authentication, client settings and a callback
func TestUnitResourceLifecycleWorkflowCustomTaskExtension_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, customTaskExtensionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer customTaskExtensionMock.CleanupMockState()

	resourceName := resourceType + ".ctx002_maximal"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_ctx002-maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("display_name").HasValue("unit-test-ctx002-maximal"),
					check.That(resourceName).Key("description").HasValue("Calls a Logic App and waits for it to resume the workflow"),
					check.That(resourceName).Key("endpoint_configuration.subscription_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceName).Key("authentication_configuration.type").HasValue("azureAdTokenAuthentication"),
					check.That(resourceName).Key("authentication_configuration.resource_id").HasValue("00000003-0000-0000-c000-000000000000"),
					check.That(resourceName).Key("client_configuration.maximum_retries").HasValue("0"),
					check.That(resourceName).Key("client_configuration.timeout_in_milliseconds").HasValue("2000"),
					check.That(resourceName).Key("callback_configuration.timeout_duration").HasValue("PT1H"),
					check.That(resourceName).Key("callback_configuration.authorized_app_ids.#").HasValue("1"),
					check.That(resourceName).Key("callback_configuration.authorized_app_ids.*").ContainsTypeSetElement("00000003-0000-0000-c000-000000000000"),
					check.That(resourceName).Key("reply_mode").HasValue("callback"),
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

// CTX003: Every property is updated in place with PATCH, then the callback is removed
func TestUnitResourceLifecycleWorkflowCustomTaskExtension_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, customTaskExtensionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer customTaskExtensionMock.CleanupMockState()

	resourceName := resourceType + ".ctx003_update"
	var firstId string

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_ctx003-update_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					captureAttribute(resourceName, "id", &firstId),
					check.That(resourceName).Key("display_name").HasValue("unit-test-ctx003-update"),
					check.That(resourceName).Key("reply_mode").HasValue("none"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_ctx003-update_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("id").ValidatesWith(sameIdAs(&firstId)),
					check.That(resourceName).Key("display_name").HasValue("unit-test-ctx003-update-renamed"),
					check.That(resourceName).Key("description").HasValue("Calls a Logic App and waits for it to resume the workflow"),
					check.That(resourceName).Key("endpoint_configuration.logic_app_workflow_name").HasValue("acc-test-la-placeholder-2"),
					check.That(resourceName).Key("endpoint_configuration.url").MatchesRegex(regexp.MustCompile(`^https://prod-01\.westeurope\.logic\.azure\.com`)),
					check.That(resourceName).Key("authentication_configuration.type").HasValue("azureAdTokenAuthentication"),
					check.That(resourceName).Key("client_configuration.maximum_retries").HasValue("0"),
					check.That(resourceName).Key("callback_configuration.timeout_duration").HasValue("PT1H"),
					check.That(resourceName).Key("reply_mode").HasValue("callback"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_ctx003-update_step3.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).Key("id").ValidatesWith(sameIdAs(&firstId)),
					check.That(resourceName).Key("description").HasValue("Callback removed"),
					check.That(resourceName).Key("authentication_configuration.type").HasValue("azureAdPopTokenAuthentication"),
					check.That(resourceName).Key("authentication_configuration.resource_id").DoesNotExist(),
					check.That(resourceName).Key("client_configuration.maximum_retries").HasValue("1"),
					check.That(resourceName).Key("client_configuration.timeout_in_milliseconds").HasValue("1000"),
					check.That(resourceName).Key("callback_configuration.%").DoesNotExist(),
					check.That(resourceName).Key("reply_mode").HasValue("none"),
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

// CTX004: A 403 from Graph surfaces through the Kiota error handler
func TestUnitResourceLifecycleWorkflowCustomTaskExtension_04_Forbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, customTaskExtensionMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer customTaskExtensionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_ctx001-minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)forbidden|permission|403`),
			},
		},
	})
}

// CTX005: A 400 from Graph surfaces through the Kiota error handler
func TestUnitResourceLifecycleWorkflowCustomTaskExtension_05_BadRequest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, customTaskExtensionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer customTaskExtensionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_ctx904-missing_url.tf"),
				ExpectError: regexp.MustCompile(`(?i)missing url|400|bad.?request`),
			},
		},
	})
}

// CTX006: A custom task extension deleted outside Terraform (404 on read) is removed from state and planned for creation
func TestUnitResourceLifecycleWorkflowCustomTaskExtension_06_NotFoundOnRead(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, customTaskExtensionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer customTaskExtensionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_ctx001-minimal.tf"),
			},
			{
				PreConfig:          customTaskExtensionMock.CleanupMockState,
				Config:             loadUnitTestTerraform("resource_ctx001-minimal.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// CTX007: Invalid configurations fail at plan time
func TestUnitResourceLifecycleWorkflowCustomTaskExtension_07_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, customTaskExtensionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer customTaskExtensionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_ctx900-token_without_resource_id.tf"),
				ExpectError: regexp.MustCompile(`resource_id is required when type is\s+azureAdTokenAuthentication`),
			},
			{
				Config:      loadUnitTestTerraform("resource_ctx901-pop_with_resource_id.tf"),
				ExpectError: regexp.MustCompile(`resource_id must not be set when type is\s+azureAdPopTokenAuthentication`),
			},
			{
				Config:      loadUnitTestTerraform("resource_ctx902-invalid_duration.tf"),
				ExpectError: regexp.MustCompile(`must be an ISO 8601\s+duration`),
			},
			{
				Config:      loadUnitTestTerraform("resource_ctx903-invalid_retries.tf"),
				ExpectError: regexp.MustCompile(`must be between 0 and 1`),
			},
		},
	})
}
