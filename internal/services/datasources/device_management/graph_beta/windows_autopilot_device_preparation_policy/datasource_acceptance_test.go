package graphBetaWindowsAutopilotDevicePreparationPolicy_test

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

const dataSourceType = "data.microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy"

const resourceType = "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy"

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccDatasourceWindowsAutopilotDevicePreparationPolicy_01_ListAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "items.#"),
					testCheckAllItemsAreDevicePreparationPolicies(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsAutopilotDevicePreparationPolicy_02_ByPolicyId(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("02_by_policy_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					resource.TestCheckResourceAttrPair(
						dataSourceType+".test", "items.0.id",
						resourceType+".test", "id",
					),
					resource.TestCheckResourceAttrPair(
						dataSourceType+".test", "items.0.name",
						resourceType+".test", "name",
					),
					check.That(dataSourceType+".test").Key("items.0.template_reference.template_family").HasValue("enrollmentConfiguration"),
					check.That(dataSourceType+".test").Key("items.0.template_reference.deployment_mode").HasValue("automatic"),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsAutopilotDevicePreparationPolicy_03_ByName(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("03_by_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					resource.TestCheckResourceAttrPair(
						dataSourceType+".test", "items.0.id",
						resourceType+".test", "id",
					),
					resource.TestCheckResourceAttrPair(
						dataSourceType+".test", "items.0.name",
						resourceType+".test", "name",
					),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsAutopilotDevicePreparationPolicy_04_ODataQuery(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("04_odata_query.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("technologies has 'enrollment'"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "items.#"),
					testCheckAllItemsAreDevicePreparationPolicies(dataSourceType+".test"),
				),
			},
		},
	})
}

func TestAccDatasourceWindowsAutopilotDevicePreparationPolicy_05_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("05_with_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_assignments").HasValue("true"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "items.#"),
					resource.TestCheckResourceAttrSet(dataSourceType+".test", "assignments.#"),
				),
			},
		},
	})
}

// testCheckAllItemsAreDevicePreparationPolicies asserts that every returned item was created from an
// Autopilot device preparation template, so unrelated settings catalog policies are never returned.
func testCheckAllItemsAreDevicePreparationPolicies(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		itemsCount := rs.Primary.Attributes["items.#"]
		if itemsCount == "" {
			return fmt.Errorf("items count not found in state")
		}

		count, err := strconv.Atoi(itemsCount)
		if err != nil {
			return fmt.Errorf("failed to parse items count: %v", err)
		}

		for i := 0; i < count; i++ {
			family := rs.Primary.Attributes[fmt.Sprintf("items.%d.template_reference.template_family", i)]
			if family != "enrollmentConfiguration" {
				return fmt.Errorf("item %d has template family %q, expected enrollmentConfiguration", i, family)
			}

			mode := rs.Primary.Attributes[fmt.Sprintf("items.%d.template_reference.deployment_mode", i)]
			if mode != "automatic" && mode != "user_driven" {
				return fmt.Errorf("item %d has deployment mode %q, expected automatic or user_driven", i, mode)
			}
		}

		return nil
	}
}
