package graphBetaUser_test

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

const (
	dataSourceType = "data.microsoft365_graph_beta_users_user"
	resourceType   = "microsoft365_graph_beta_users_user"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// externalProviders returns the external providers required by the acceptance tests.
// "random" generates unique resource names and "time" provides the propagation wait
// (time_sleep) between creating the user and querying it via the data source.
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

func TestAccDatasourceUser_01_ListAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
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

func TestAccDatasourceUser_02_ByObjectId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_by_object_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckItemsCountExists(dataSourceType+".test"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.id", resourceType+".test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.display_name", resourceType+".test", "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.user_principal_name", resourceType+".test", "user_principal_name"),
				),
			},
		},
	})
}

func TestAccDatasourceUser_03_ByDisplayName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("03_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckItemsCountExists(dataSourceType+".test"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.id", resourceType+".test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.display_name", resourceType+".test", "display_name"),
				),
			},
		},
	})
}

func TestAccDatasourceUser_04_ByEmployeeId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("04_by_employee_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckItemsCountExists(dataSourceType+".test"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.id", resourceType+".test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.employee_id", resourceType+".test", "employee_id"),
				),
			},
		},
	})
}

func TestAccDatasourceUser_05_ByGivenName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("05_by_given_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckItemsCountExists(dataSourceType+".test"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.id", resourceType+".test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.given_name", resourceType+".test", "given_name"),
				),
			},
		},
	})
}

func TestAccDatasourceUser_06_ByUserPrincipalName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("06_by_user_principal_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckItemsCountExists(dataSourceType+".test"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.id", resourceType+".test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.user_principal_name", resourceType+".test", "user_principal_name"),
				),
			},
		},
	})
}

// Note: lookup by on_premises_immutable_id and on_premises_distinguished_name cannot be
// covered by a deploy-then-query acceptance test, because those properties are only
// populated by Entra Connect synchronization from an on-premises Active Directory and
// cannot be set on a cloud-only user created by the resource. They remain covered by the
// unit tests.

func TestAccDatasourceUser_09_ByODataQuery(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("09_odata_query.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckItemsCountExists(dataSourceType+".test"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.id", resourceType+".test", "id"),
					resource.TestCheckResourceAttrPair(dataSourceType+".test", "items.0.user_principal_name", resourceType+".test", "user_principal_name"),
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
