package graphBetaServicePrincipal_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServicePrincipalDataSource_All(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.all", "filter_type", "all"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.all", "items.#"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.all", "items.0.id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.all", "items.0.display_name"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.all", "items.0.app_id"),
				),
			},
		},
	})
}

func TestAccServicePrincipalDataSource_ByDisplayName(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "filter_type", "display_name"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "filter_value", "Microsoft Graph"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "items.#"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "items.0.display_name", "Microsoft Graph PowerShell"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "items.0.id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.by_display_name", "items.0.app_id"),
				),
			},
		},
	})
}

func TestAccServicePrincipalDataSource_ByAppId(t *testing.T) {
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
				Config: testAccConfigByAppId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "filter_type", "app_id"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "filter_value", "00000003-0000-0000-c000-000000000000"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "items.#", "1"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "items.0.app_id", "00000003-0000-0000-c000-000000000000"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "items.0.display_name", "Microsoft Graph"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.by_app_id", "items.0.id"),
				),
			},
		},
	})
}

func TestAccServicePrincipalDataSource_ODataFilter(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "odata_filter", "startsWith(displayName,'Microsoft')"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "odata_count", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "odata_orderby", "displayName"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "items.#"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "items.0.id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.odata_filter", "items.0.display_name"),
				),
			},
		},
	})
}

func TestAccServicePrincipalDataSource_ODataAdvanced(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "odata_filter", "startsWith(displayName,'Microsoft')"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "odata_select", "id,appId,displayName,publisherName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "odata_top", "10"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_advanced", "odata_skip", "0"),
				),
			},
		},
	})
}

func TestAccServicePrincipalDataSource_ODataComprehensive(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_filter", "startsWith(displayName,'Microsoft')"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_count", "true"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_orderby", "displayName"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_search", "\"displayName:Graph\""),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_select", "id,appId,displayName,publisherName,servicePrincipalType"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_top", "5"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_comprehensive", "odata_skip", "0"),
				),
			},
		},
	})
}

func TestAccServicePrincipalDataSource_ODataSearchOnly(t *testing.T) {
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
				Config: testAccConfigODataSearchOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_search_only", "filter_type", "odata"),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_search_only", "odata_search", "\"displayName:Intune\""),
					resource.TestCheckResourceAttr("data.microsoft365_graph_beta_applications_service_principal.odata_search_only", "odata_count", "true"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.odata_search_only", "items.#"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.odata_search_only", "items.0.id"),
					resource.TestCheckResourceAttrSet("data.microsoft365_graph_beta_applications_service_principal.odata_search_only", "items.0.display_name"),
				),
			},
		},
	})
}

// Configuration functions
func testAccConfigAll() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_all.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigByDisplayName() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_by_display_name.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigByAppId() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_by_app_id.tf")
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

func testAccConfigODataSearchOnly() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/07_odata_search_only.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
