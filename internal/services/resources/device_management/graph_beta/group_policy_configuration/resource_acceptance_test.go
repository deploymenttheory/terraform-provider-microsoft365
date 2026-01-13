package graphBetaGroupPolicyConfiguration_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroupPolicyConfiguration "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_configuration"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testResource = graphBetaGroupPolicyConfiguration.GroupPolicyConfigurationTestResource{}

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccGroupPolicyConfigurationResource_Minimal tests creating a minimal group policy configuration
func TestAccGroupPolicyConfigurationResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal group policy configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-001-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".minimal").Key("policy_configuration_ingestion_type").Exists(),
					check.That(resourceType+".minimal").Key("created_date_time").Exists(),
					check.That(resourceType+".minimal").Key("last_modified_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal group policy configuration")
				},
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccGroupPolicyConfigurationResource_Maximal tests creating a maximal group policy configuration
func TestAccGroupPolicyConfigurationResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal group policy configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-002-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".maximal").Key("description").HasValue("acc-test-002-maximal"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".maximal").Key("policy_configuration_ingestion_type").Exists(),
					check.That(resourceType+".maximal").Key("created_date_time").Exists(),
					check.That(resourceType+".maximal").Key("last_modified_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal group policy configuration")
				},
				ResourceName:            resourceType + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccGroupPolicyConfigurationResource_MinimalAssignment tests creating a configuration with minimal assignment
func TestAccGroupPolicyConfigurationResource_MinimalAssignment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating group policy configuration with minimal assignment")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal_assignment").ExistsInGraph(testResource),
					check.That(resourceType+".minimal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal_assignment").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-003-minimal-assignment-[a-z0-9]{8}$`)),
					check.That(resourceType+".minimal_assignment").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".minimal_assignment").Key("assignments.0.type").HasValue("allDevicesAssignmentTarget"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("group policy configuration assignment", 30*time.Second)
					time.Sleep(30 * time.Second)
					testlog.StepAction(resourceType, "Importing group policy configuration with minimal assignment")
				},
				ResourceName:            resourceType + ".minimal_assignment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccGroupPolicyConfigurationResource_MaximalAssignment tests creating a configuration with maximal assignments
func TestAccGroupPolicyConfigurationResource_MaximalAssignment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating group policy configuration with maximal assignments")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal_assignment").ExistsInGraph(testResource),
					check.That(resourceType+".maximal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal_assignment").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-004-maximal-assignment-[a-z0-9]{8}$`)),
					check.That(resourceType+".maximal_assignment").Key("description").HasValue("acc-test-004-maximal-assignment"),
					check.That(resourceType+".maximal_assignment").Key("assignments.#").HasValue("4"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("group policy configuration assignments", 30*time.Second)
					time.Sleep(30 * time.Second)
					testlog.StepAction(resourceType, "Importing group policy configuration with maximal assignments")
				},
				ResourceName:            resourceType + ".maximal_assignment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccGroupPolicyConfigurationResource_MinimalToMaximal tests transitioning from minimal to maximal configuration
func TestAccGroupPolicyConfigurationResource_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal configuration for transition test")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-001-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("group policy configuration before update", 30*time.Second)
					time.Sleep(30 * time.Second)
					testlog.StepAction(resourceType, "Updating to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".transition").ExistsInGraph(testResource),
					check.That(resourceType+".transition").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".transition").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-005-lifecycle-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".transition").Key("description").HasValue("acc-test-005-lifecycle-maximal"),
					check.That(resourceType+".transition").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".transition").Key("assignments.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("group policy configuration after update", 30*time.Second)
					time.Sleep(30 * time.Second)
					testlog.StepAction(resourceType, "Importing transitioned configuration")
				},
				ResourceName:            resourceType + ".transition",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccGroupPolicyConfigurationResource_MaximalToMinimal tests transitioning from maximal to minimal configuration
func TestAccGroupPolicyConfigurationResource_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal configuration for transition test")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".transition").ExistsInGraph(testResource),
					check.That(resourceType+".transition").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".transition").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-005-lifecycle-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".transition").Key("description").HasValue("acc-test-005-lifecycle-maximal"),
					check.That(resourceType+".transition").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".transition").Key("assignments.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("group policy configuration before downgrade", 30*time.Second)
					time.Sleep(30 * time.Second)
					testlog.StepAction(resourceType, "Downgrading to minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal_to_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".transition").ExistsInGraph(testResource),
					check.That(resourceType+".transition").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".transition").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-006-lifecycle-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".transition").Key("description").HasValue("acc-test-005-lifecycle-maximal"), // Description persists from step 1
					check.That(resourceType+".transition").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.WaitForConsistency("group policy configuration after downgrade", 30*time.Second)
					time.Sleep(30 * time.Second)
					testlog.StepAction(resourceType, "Importing downgraded configuration")
				},
				ResourceName:            resourceType + ".transition",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}
