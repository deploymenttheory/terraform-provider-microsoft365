package graphBetaRoleDefinitions_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleDefinitionsDataSource_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "filter_type", "all"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.#"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.0.id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.0.display_name"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_identity_and_access_role_definitions.all", "items.0.template_id"),
				),
			},
		},
	})
}

func TestAccRoleDefinitionsDataSource_ById(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigById(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "filter_type", "id"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "filter_value", "62e90394-69f5-4237-9190-012177145e10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.id", "62e90394-69f5-4237-9190-012177145e10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.display_name", "Global Administrator"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.template_id", "62e90394-69f5-4237-9190-012177145e10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.is_built_in", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_id", "items.0.is_privileged", "true"),
				),
			},
		},
	})
}

func TestAccRoleDefinitionsDataSource_ByDisplayName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigByDisplayName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "filter_type", "display_name"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "filter_value", "Security Administrator"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "items.#"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "items.0.id", "194ae4cb-b126-40b2-bd5b-6091b380977d"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.by_display_name", "items.0.display_name", "Security Administrator"),
				),
			},
		},
	})
}

func TestAccRoleDefinitionsDataSource_ODataFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigODataFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_filter", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_filter", "odata_filter", "isPrivileged eq true"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_filter", "items.#"),
				),
			},
		},
	})
}

func TestAccRoleDefinitionsDataSource_ODataAdvanced(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigODataAdvanced(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "odata_filter", "isBuiltIn eq true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "odata_orderby", "displayName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "odata_select", "id,displayName,description,isPrivileged"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_advanced", "items.#"),
				),
			},
		},
	})
}

func TestAccRoleDefinitionsDataSource_ODataComprehensive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source:            "hashicorp/azuread",
				VersionConstraint: ">= 2.47.0",
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigODataComprehensive(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "odata_filter", "isBuiltIn eq true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "odata_count", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "odata_orderby", "displayName"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_identity_and_access_role_definitions.odata_comprehensive", "items.#"),
				),
			},
		},
	})
}

// Acceptance test configuration functions
func testAccConfigAll() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_all.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigById() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_by_id.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigByDisplayName() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_by_display_name.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigODataFilter() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/04_odata_filter.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigODataAdvanced() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/05_odata_advanced.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigODataComprehensive() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/06_odata_comprehensive.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
