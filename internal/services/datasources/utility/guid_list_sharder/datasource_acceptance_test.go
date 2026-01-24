package utilityGuidListSharder_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityGuidListSharder "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/guid_list_sharder"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Helper function to load acceptance test Terraform configurations
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccGuidListSharderDataSource_01_UsersRoundRobinNoSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_users_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource ID exists (deterministic hash)
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("id").Exists(),

					// Verify exactly 3 shards exist
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.%").HasValue("3"),
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.shard_0.#").Exists(),
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.shard_1.#").Exists(),
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.shard_2.#").Exists(),

					// Verify total distributed = 9 test users
					resource.TestCheckOutput("total_users_distributed", "9"),
					
					// Verify each shard has at least 1 member
					resource.TestMatchOutput("shard_0_count", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchOutput("shard_1_count", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchOutput("shard_2_count", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify first GUID in each shard is valid GUID format
					resource.TestMatchOutput("shard_0_first_guid", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
					resource.TestMatchOutput("shard_1_first_guid", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
					resource.TestMatchOutput("shard_2_first_guid", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),

					// Verify ALL GUIDs in ALL shards are valid GUID format (comprehensive validation)
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}

func TestAccGuidListSharderDataSource_02_UsersRoundRobinWithSeed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_users_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource ID exists (deterministic hash)
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("id").Exists(),

					// Verify exactly 2 shards exist
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.%").HasValue("2"),
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.shard_0.#").Exists(),
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.shard_1.#").Exists(),

					// Verify total distributed = 6 test users
					resource.TestCheckOutput("total_users_distributed", "6"),

					// Verify balanced distribution (difference should be 0 or 1)
					resource.TestMatchOutput("is_balanced", regexp.MustCompile(`^[01]$`)),
					
					// Verify each group has exactly 3 users (balanced distribution)
					resource.TestCheckOutput("group_a_count", "3"),
					resource.TestCheckOutput("group_b_count", "3"),

					// Verify first GUID in each shard is valid GUID format
					resource.TestMatchOutput("shard_0_first_guid", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
					resource.TestMatchOutput("shard_1_first_guid", regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),

					// Verify ALL GUIDs in ALL shards are valid GUID format (comprehensive validation)
					resource.TestCheckOutput("all_guids_valid", "true"),
				),
			},
		},
	})
}
