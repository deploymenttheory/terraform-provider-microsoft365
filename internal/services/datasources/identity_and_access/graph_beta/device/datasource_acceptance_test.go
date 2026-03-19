package graphBetaDevice_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const dataSourceType = "data.microsoft365_graph_beta_identity_and_access_device"

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccDatasourceDevice_01_ListAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_02_ByObjectId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_by_object_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").Exists(),
					testCheckItemsCountExists(dataSourceType+".test"),
					check.That(dataSourceType+".test").Key("items.0.id").Exists(),
					check.That(dataSourceType+".test").Key("items.0.display_name").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_03_ByDisplayName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("03_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("display_name").Exists(),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_04_ByDeviceId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("04_by_device_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("device_id").Exists(),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_05_ODataSimpleFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("05_odata_simple_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("operatingSystem eq 'Windows'"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_06_ODataAndFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("06_odata_and_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("operatingSystem eq 'Windows' and accountEnabled eq true"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_07_ODataCompliant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("07_odata_compliant.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("isCompliant eq true"),
					testCheckItemsCountExists(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_08_WithMemberOf(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("08_with_member_of.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").Exists(),
					check.That(dataSourceType+".test").Key("list_member_of").HasValue("true"),
					testCheckItemsCountExists(dataSourceType+".test"),
					// member_of count should exist (may be 0 if device is not a member of any groups)
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "member_of.#"),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_09_WithRegisteredOwners(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("09_with_registered_owners.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").Exists(),
					check.That(dataSourceType+".test").Key("list_registered_owners").HasValue("true"),
					testCheckItemsCountExists(dataSourceType+".test"),
					// registered_owners count should exist (may be 0 if no owners)
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "registered_owners.#"),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_10_WithRegisteredUsers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("10_with_registered_users.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").Exists(),
					check.That(dataSourceType+".test").Key("list_registered_users").HasValue("true"),
					testCheckItemsCountExists(dataSourceType+".test"),
					// registered_users count should exist (may be 0 if no users)
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "registered_users.#"),
				),
			},
		},
	})
}

func TestAccDatasourceDevice_11_Comprehensive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("11_comprehensive.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("object_id").Exists(),
					check.That(dataSourceType+".test").Key("list_member_of").HasValue("true"),
					check.That(dataSourceType+".test").Key("list_registered_owners").HasValue("true"),
					check.That(dataSourceType+".test").Key("list_registered_users").HasValue("true"),
					testCheckItemsCountExists(dataSourceType+".test"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "member_of.#"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "registered_owners.#"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "registered_users.#"),
					check.That(dataSourceType+".test").Key("items.0.id").Exists(),
					check.That(dataSourceType+".test").Key("items.0.display_name").Exists(),
					check.That(dataSourceType+".test").Key("items.0.operating_system").Exists(),
				),
			},
		},
	})
}

// testCheckItemsCountExists verifies that the items list exists and has at least one item
func testCheckItemsCountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", resourceName, s.RootModule().Resources)
		}

		itemsCount := rs.Primary.Attributes["items.#"]
		if itemsCount == "" {
			return fmt.Errorf("items count not found in state")
		}

		count, err := strconv.Atoi(itemsCount)
		if err != nil {
			return fmt.Errorf("failed to parse items count: %v", err)
		}

		if count < 1 {
			return fmt.Errorf("expected at least 1 item, got %d", count)
		}

		return nil
	}
}
