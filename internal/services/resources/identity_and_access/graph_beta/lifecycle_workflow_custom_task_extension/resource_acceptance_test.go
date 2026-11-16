package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension_test

import (
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
	graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/lifecycle_workflow_custom_task_extension"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension.ResourceName

	// testResource is the test resource implementation for custom task extensions
	testResource = graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension.LifecycleWorkflowCustomTaskExtensionTestResource{}

	// externalProviders are the providers acceptance fixtures depend on
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

// CTX001: Minimal custom task extension with only required attributes
func TestAccResourceLifecycleWorkflowCustomTaskExtension_01_Minimal(t *testing.T) {
	resourceName := resourceType + ".ctx001_minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CTX001 minimal custom task extension")
				},
				Config: loadAcceptanceTestTerraform("resource_ctx001-minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceName).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ctx001-minimal-[a-z0-9]{8}$`)),
					check.That(resourceName).Key("description").HasValue(""),
					check.That(resourceName).Key("endpoint_configuration.subscription_id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceName).Key("endpoint_configuration.logic_app_workflow_name").HasValue("acc-test-la-placeholder"),
					check.That(resourceName).Key("authentication_configuration.type").HasValue("azureAdPopTokenAuthentication"),
					check.That(resourceName).Key("client_configuration.maximum_retries").HasValue("1"),
					check.That(resourceName).Key("client_configuration.timeout_in_milliseconds").HasValue("1000"),
					check.That(resourceName).Key("callback_configuration.%").DoesNotExist(),
					check.That(resourceName).Key("reply_mode").HasValue("none"),
				),
			},
			importStep(resourceName, "Importing CTX001 minimal custom task extension"),
		},
	})
}

// CTX002: Maximal custom task extension with token authentication, client settings and a callback
func TestAccResourceLifecycleWorkflowCustomTaskExtension_02_Maximal(t *testing.T) {
	resourceName := resourceType + ".ctx002_maximal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CTX002 maximal custom task extension")
				},
				Config: loadAcceptanceTestTerraform("resource_ctx002-maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("description").HasValue("Calls a Logic App and waits for it to resume the workflow"),
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
			importStep(resourceName, "Importing CTX002 maximal custom task extension"),
		},
	})
}

// CTX003: Every property is updated in place with PATCH, then the callback is removed
func TestAccResourceLifecycleWorkflowCustomTaskExtension_03_Update(t *testing.T) {
	resourceName := resourceType + ".ctx003_update"
	var firstId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CTX003 custom task extension")
				},
				Config: loadAcceptanceTestTerraform("resource_ctx003-update_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					captureAttribute(resourceName, "id", &firstId),
					check.That(resourceName).Key("reply_mode").HasValue("none"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating every CTX003 property and adding a callback")
				},
				Config: loadAcceptanceTestTerraform("resource_ctx003-update_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("id").ValidatesWith(sameIdAs(&firstId)),
					check.That(resourceName).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ctx003-update-renamed-[a-z0-9]{8}$`)),
					check.That(resourceName).Key("endpoint_configuration.resource_group_name").HasValue("acc-test-rg-placeholder-2"),
					check.That(resourceName).Key("authentication_configuration.type").HasValue("azureAdTokenAuthentication"),
					check.That(resourceName).Key("client_configuration.timeout_in_milliseconds").HasValue("2000"),
					check.That(resourceName).Key("callback_configuration.timeout_duration").HasValue("PT1H"),
					check.That(resourceName).Key("callback_configuration.authorized_app_ids.#").HasValue("1"),
					check.That(resourceName).Key("reply_mode").HasValue("callback"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Removing the CTX003 callback")
				},
				Config: loadAcceptanceTestTerraform("resource_ctx003-update_step3.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceName).ExistsInGraph(testResource),
					check.That(resourceName).Key("id").ValidatesWith(sameIdAs(&firstId)),
					check.That(resourceName).Key("description").HasValue("Callback removed"),
					check.That(resourceName).Key("authentication_configuration.type").HasValue("azureAdPopTokenAuthentication"),
					check.That(resourceName).Key("client_configuration.maximum_retries").HasValue("1"),
					check.That(resourceName).Key("client_configuration.timeout_in_milliseconds").HasValue("1000"),
					check.That(resourceName).Key("callback_configuration.%").DoesNotExist(),
					check.That(resourceName).Key("reply_mode").HasValue("none"),
				),
			},
			importStep(resourceName, "Importing CTX003 custom task extension"),
		},
	})
}
