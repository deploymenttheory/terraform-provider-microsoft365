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

func TestAccGuidListSharderDataSource_01_UsersHashNoSeed(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("01_users_hash_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("id").MatchesRegex(regexp.MustCompile(`^users-3-hash$`)),
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.%").HasValue("3"),
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.shard_0.#").Exists(),
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.shard_1.#").Exists(),
					check.That("data."+utilityGuidListSharder.DataSourceName+".test").Key("shards.shard_2.#").Exists(),
					resource.TestCheckOutput("total_users", "100"),
				),
			},
		},
	})
}
