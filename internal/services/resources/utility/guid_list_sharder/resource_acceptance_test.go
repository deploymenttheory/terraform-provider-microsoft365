package utilityGuidListSharder_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityGuidListSharder "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/utility/guid_list_sharder"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

var guidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
		"time": {
			Source:            "hashicorp/time",
			VersionConstraint: constants.ExternalProviderTimeVersion,
		},
	}
}

// =============================================================================
// Users — Round-Robin
// =============================================================================

func TestAccResourceGuidListSharder_01_UsersRoundRobinNoSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_users_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("3"),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_0.#").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_1.#").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_2.#").Exists(),
					resource.TestCheckOutput("total_users_distributed", "9"),
					resource.TestMatchOutput("shard_0_count", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchOutput("shard_1_count", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchOutput("shard_2_count", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchOutput("shard_0_first_guid", guidRegex),
					resource.TestMatchOutput("shard_1_first_guid", guidRegex),
					resource.TestMatchOutput("shard_2_first_guid", guidRegex),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}

func TestAccResourceGuidListSharder_02_UsersRoundRobinWithSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_users_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("2"),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_0.#").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_1.#").Exists(),
					resource.TestCheckOutput("total_users_distributed", "6"),
					resource.TestMatchOutput("is_balanced", regexp.MustCompile(`^[01]$`)),
					resource.TestCheckOutput("group_a_count", "3"),
					resource.TestCheckOutput("group_b_count", "3"),
					resource.TestMatchOutput("shard_0_first_guid", guidRegex),
					resource.TestMatchOutput("shard_1_first_guid", guidRegex),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}

// =============================================================================
// Users — Percentage
// =============================================================================

// TestAccResourceGuidListSharder_03_UsersPercentageNoSeed verifies percentage-based
// distribution with [10, 30, 60] against 10 test users (no seed).
// Expected: shard_0=1, shard_1=3, shard_2=6 — sizes are deterministic by arithmetic
// regardless of seed; only membership varies.
func TestAccResourceGuidListSharder_03_UsersPercentageNoSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("03_users_percentage_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("total_distributed", "10"),
					resource.TestCheckOutput("shard_0_count", "1"),
					resource.TestCheckOutput("shard_1_count", "3"),
					resource.TestCheckOutput("shard_2_count", "6"),
					resource.TestMatchOutput("shard_0_first_guid", guidRegex),
					resource.TestMatchOutput("shard_1_first_guid", guidRegex),
					resource.TestMatchOutput("shard_2_first_guid", guidRegex),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}

// TestAccResourceGuidListSharder_04_UsersPercentageWithSeed verifies percentage-based
// distribution with [10, 30, 60] and seed="mfa-phased-2024" against 10 test users.
// Shard sizes must be identical to the no-seed variant — seed only varies membership.
func TestAccResourceGuidListSharder_04_UsersPercentageWithSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("04_users_percentage_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("total_distributed", "10"),
					resource.TestCheckOutput("shard_0_count", "1"),
					resource.TestCheckOutput("shard_1_count", "3"),
					resource.TestCheckOutput("shard_2_count", "6"),
					resource.TestMatchOutput("shard_0_first_guid", guidRegex),
					resource.TestMatchOutput("shard_1_first_guid", guidRegex),
					resource.TestMatchOutput("shard_2_first_guid", guidRegex),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}

// =============================================================================
// Users — Size
// =============================================================================

// TestAccResourceGuidListSharder_05_UsersSizeNoSeed verifies absolute size-based
// distribution with sizes [3, 4, -1] against 9 test users (no seed).
// Expected: shard_0=3, shard_1=4, shard_2=2 (all remaining via -1 sentinel).
func TestAccResourceGuidListSharder_05_UsersSizeNoSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("05_users_size_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("total_distributed", "9"),
					resource.TestCheckOutput("shard_0_count", "3"),
					resource.TestCheckOutput("shard_1_count", "4"),
					resource.TestCheckOutput("shard_2_count", "2"),
					resource.TestMatchOutput("shard_0_first_guid", guidRegex),
					resource.TestMatchOutput("shard_1_first_guid", guidRegex),
					resource.TestMatchOutput("shard_2_first_guid", guidRegex),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}

// TestAccResourceGuidListSharder_06_UsersSizeWithSeed verifies absolute size-based
// distribution with sizes [3, 4, -1] and seed="mfa-rollout-2024" against 9 test users.
// Shard sizes must be identical to the no-seed variant — seed only varies membership.
func TestAccResourceGuidListSharder_06_UsersSizeWithSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("06_users_size_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("total_distributed", "9"),
					resource.TestCheckOutput("shard_0_count", "3"),
					resource.TestCheckOutput("shard_1_count", "4"),
					resource.TestCheckOutput("shard_2_count", "2"),
					resource.TestMatchOutput("shard_0_first_guid", guidRegex),
					resource.TestMatchOutput("shard_1_first_guid", guidRegex),
					resource.TestMatchOutput("shard_2_first_guid", guidRegex),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}

// =============================================================================
// Users — Rendezvous (HRW)
// =============================================================================

// TestAccResourceGuidListSharder_07_UsersRendezvousNoSeed verifies that the HRW
// algorithm distributes all 12 users across 4 shards without any explicit seed.
// Per-shard counts are NOT asserted because rendezvous balance is probabilistic;
// only the total count and GUID format are verified.
func TestAccResourceGuidListSharder_07_UsersRendezvousNoSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("07_users_rendezvous_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("4"),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_0.#").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_1.#").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_2.#").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_3.#").Exists(),
					resource.TestCheckOutput("total_distributed", "12"),
					resource.TestMatchOutput("shard_0_count", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchOutput("shard_1_count", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchOutput("shard_2_count", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchOutput("shard_3_count", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}

// TestAccResourceGuidListSharder_08_UsersRendezvousWithSeed verifies that the HRW
// algorithm distributes all 12 users across 4 shards with seed="deployment-ring-2024".
// Per-shard counts are NOT asserted because rendezvous balance is probabilistic;
// only the total count and GUID format are verified.
func TestAccResourceGuidListSharder_08_UsersRendezvousWithSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("08_users_rendezvous_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("4"),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_0.#").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_1.#").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_2.#").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.shard_3.#").Exists(),
					resource.TestCheckOutput("total_distributed", "12"),
					resource.TestMatchOutput("shard_0_count", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchOutput("shard_1_count", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchOutput("shard_2_count", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchOutput("shard_3_count", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}

// =============================================================================
// Users — recalculate_on_next_run lifecycle
// =============================================================================

// TestAccResourceGuidListSharder_09_RecalculateBehaviour proves the lock/unlock
// semantics of recalculate_on_next_run using a 3-step lifecycle:
//
//   Step 1 — Create with 6 users, recalculate_on_next_run = false.
//            Create ALWAYS computes regardless of the flag. total_distributed = 6.
//
//   Step 2 — Add 3 more users (9 total in tenant), flag stays false.
//            Read must NOT re-query the API. Assignments locked. total_distributed = 6.
//
//   Step 3 — Switch flag to true (9 users unchanged in tenant).
//            Update triggers a reshard. total_distributed = 9.
func TestAccResourceGuidListSharder_09_RecalculateBehaviour(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			// Step 1: Create with 6 users — flag=false, assignments always computed on first apply.
			{
				Config: loadAcceptanceTestTerraform("09_recalculate_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("2"),
					resource.TestCheckOutput("total_distributed", "6"),
					resource.TestCheckOutput("recalculate_flag", "false"),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
			// Step 2: 9 users now exist in the tenant; flag stays false — lock must hold.
			{
				Config: loadAcceptanceTestTerraform("09_recalculate_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("2"),
					resource.TestCheckOutput("total_distributed", "6"),
					resource.TestCheckOutput("recalculate_flag", "false"),
				),
			},
			// Step 3: Switch flag to true — Update reshards from current tenant membership.
			{
				Config: loadAcceptanceTestTerraform("09_recalculate_step3.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(utilityGuidListSharder.ResourceName+".test").Key("id").Exists(),
					check.That(utilityGuidListSharder.ResourceName+".test").Key("shards.%").HasValue("2"),
					resource.TestCheckOutput("total_distributed", "9"),
					resource.TestCheckOutput("recalculate_flag", "true"),
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}
