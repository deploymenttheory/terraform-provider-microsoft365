package graphBetaGroupLifecycleExpirationPolicy_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroupLifecycleExpirationPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group_lifecycle_expiration_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	resourceType = graphBetaGroupLifecycleExpirationPolicy.ResourceName
)

var (
	testResource = graphBetaGroupLifecycleExpirationPolicy.GroupLifecycleExpirationPolicyTestResource{}
)

func testAccGroupLifecyclePolicyConfig_All() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/acceptance/resource_all.tf")
	return config
}

func testAccGroupLifecyclePolicyConfig_Selected() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/acceptance/resource_selected.tf")
	return config
}

func testAccGroupLifecyclePolicyConfig_None() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/acceptance/resource_none.tf")
	return config
}

// TestAccGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_All tests managed_group_types = "All"
func TestAccGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Creating policy with managed_group_types=All")
				},
				Config: testAccGroupLifecyclePolicyConfig_All(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group lifecycle policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".all").ExistsInGraph(testResource),
					check.That(resourceType+".all").Key("id").Exists(),
					check.That(resourceType+".all").Key("group_lifetime_in_days").HasValue("180"),
					check.That(resourceType+".all").Key("managed_group_types").HasValue("All"),
					check.That(resourceType+".all").Key("alternate_notification_emails").HasValue("admin@deploymenttheory.com"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing policy with managed_group_types=All")
				},
				ResourceName:            resourceType + ".all",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overwrite_existing_policy"},
			},
		},
	})
}

// TestAccGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_Selected tests managed_group_types = "Selected"
func TestAccGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_Selected(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Creating policy with managed_group_types=Selected")
				},
				Config: testAccGroupLifecyclePolicyConfig_Selected(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group lifecycle policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".selected").ExistsInGraph(testResource),
					check.That(resourceType+".selected").Key("id").Exists(),
					check.That(resourceType+".selected").Key("group_lifetime_in_days").HasValue("365"),
					check.That(resourceType+".selected").Key("managed_group_types").HasValue("Selected"),
					check.That(resourceType+".selected").Key("alternate_notification_emails").HasValue("admin@deploymenttheory.com;notifications@deploymenttheory.com"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing policy with managed_group_types=Selected")
				},
				ResourceName:            resourceType + ".selected",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overwrite_existing_policy"},
			},
		},
	})
}

// TestAccGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_None tests managed_group_types = "None"
func TestAccGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_None(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Creating policy with managed_group_types=None")
				},
				Config: testAccGroupLifecyclePolicyConfig_None(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group lifecycle policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".none").ExistsInGraph(testResource),
					check.That(resourceType+".none").Key("id").Exists(),
					check.That(resourceType+".none").Key("group_lifetime_in_days").HasValue("365"),
					check.That(resourceType+".none").Key("managed_group_types").HasValue("None"),
					check.That(resourceType+".none").Key("alternate_notification_emails").HasValue("admin@deploymenttheory.com"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing policy with managed_group_types=None")
				},
				ResourceName:            resourceType + ".none",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overwrite_existing_policy"},
			},
		},
	})
}
