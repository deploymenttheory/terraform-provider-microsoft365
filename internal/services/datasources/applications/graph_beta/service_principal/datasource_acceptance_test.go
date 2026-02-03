package graphBetaServicePrincipal_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaServicePrincipal "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/applications/graph_beta/service_principal"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// DataSource type name from the datasource package
	dataSourceType = graphBetaServicePrincipal.DataSourceName
)

func TestAccDatasourceServicePrincipal_01_ByObjectId(t *testing.T) {
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
				Config: loadAccTestTerraform("01_by_object_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_object_id").Key("app_id").HasValue("00000003-0000-0000-c000-000000000000"),
					check.That("data."+dataSourceType+".by_object_id").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("display_name").HasValue("Microsoft Graph"),
					check.That("data."+dataSourceType+".by_object_id").Key("app_display_name").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("account_enabled").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("service_principal_type").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("app_role_assignment_required").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("sign_in_audience").Exists(),
					check.That("data."+dataSourceType+".by_object_id").Key("service_principal_names.#").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceServicePrincipal_02_ByDisplayName(t *testing.T) {
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
				Config: loadAccTestTerraform("02_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_display_name").Key("display_name").HasValue("Microsoft Graph"),
					check.That("data."+dataSourceType+".by_display_name").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("app_id").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("app_display_name").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("service_principal_type").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("account_enabled").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("sign_in_audience").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("app_role_assignment_required").Exists(),
					check.That("data."+dataSourceType+".by_display_name").Key("service_principal_names.#").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceServicePrincipal_03_ByAppId(t *testing.T) {
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
				Config: loadAccTestTerraform("03_by_app_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_app_id").Key("app_id").HasValue("00000003-0000-0000-c000-000000000000"),
					check.That("data."+dataSourceType+".by_app_id").Key("display_name").HasValue("Microsoft Graph"),
					check.That("data."+dataSourceType+".by_app_id").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("app_display_name").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("service_principal_type").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("account_enabled").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("sign_in_audience").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("app_role_assignment_required").Exists(),
					check.That("data."+dataSourceType+".by_app_id").Key("service_principal_names.#").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceServicePrincipal_04_ODataFilter(t *testing.T) {
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
				Config: loadAccTestTerraform("04_odata_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_filter").Key("odata_query").HasValue("appId eq '00000003-0000-0000-c000-000000000000' and servicePrincipalType eq 'Application'"),
					check.That("data."+dataSourceType+".odata_filter").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("display_name").HasValue("Microsoft Graph"),
					check.That("data."+dataSourceType+".odata_filter").Key("app_id").HasValue("00000003-0000-0000-c000-000000000000"),
					check.That("data."+dataSourceType+".odata_filter").Key("service_principal_type").HasValue("Application"),
					check.That("data."+dataSourceType+".odata_filter").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("app_display_name").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("account_enabled").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("sign_in_audience").Exists(),
					check.That("data."+dataSourceType+".odata_filter").Key("service_principal_names.#").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceServicePrincipal_05_ODataAdvanced(t *testing.T) {
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
				Config: loadAccTestTerraform("05_odata_advanced.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_advanced").Key("odata_query").HasValue("displayName eq 'Microsoft Graph' and accountEnabled eq true"),
					check.That("data."+dataSourceType+".odata_advanced").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("display_name").HasValue("Microsoft Graph"),
					check.That("data."+dataSourceType+".odata_advanced").Key("app_id").HasValue("00000003-0000-0000-c000-000000000000"),
					check.That("data."+dataSourceType+".odata_advanced").Key("account_enabled").HasValue("true"),
					check.That("data."+dataSourceType+".odata_advanced").Key("service_principal_type").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("app_display_name").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("sign_in_audience").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("app_role_assignment_required").Exists(),
					check.That("data."+dataSourceType+".odata_advanced").Key("service_principal_names.#").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceServicePrincipal_06_ODataComprehensive(t *testing.T) {
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
				Config: loadAccTestTerraform("06_odata_comprehensive.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".odata_comprehensive").Key("odata_query").HasValue("appId eq '00000003-0000-0000-c000-000000000000'"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("id").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("display_name").HasValue("Microsoft Graph"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("app_id").HasValue("00000003-0000-0000-c000-000000000000"),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("app_display_name").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("service_principal_type").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("account_enabled").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("sign_in_audience").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("app_role_assignment_required").Exists(),
					check.That("data."+dataSourceType+".odata_comprehensive").Key("service_principal_names.#").Exists(),
				),
			},
		},
	})
}

func TestAccDatasourceServicePrincipal_07_ByODataTags(t *testing.T) {
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
				Config: loadAccTestTerraform("07_by_odata_tags.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_odata_tags").Key("odata_query").HasValue("appId eq '00000003-0000-0000-c000-000000000000'"),
					check.That("data."+dataSourceType+".by_odata_tags").Key("id").Exists(),
					check.That("data."+dataSourceType+".by_odata_tags").Key("display_name").HasValue("Microsoft Graph"),
					check.That("data."+dataSourceType+".by_odata_tags").Key("app_id").HasValue("00000003-0000-0000-c000-000000000000"),
					check.That("data."+dataSourceType+".by_odata_tags").Key("service_principal_type").Exists(),
					check.That("data."+dataSourceType+".by_odata_tags").Key("app_display_name").Exists(),
					check.That("data."+dataSourceType+".by_odata_tags").Key("publisher_name").Exists(),
					check.That("data."+dataSourceType+".by_odata_tags").Key("account_enabled").Exists(),
					check.That("data."+dataSourceType+".by_odata_tags").Key("sign_in_audience").Exists(),
				),
			},
		},
	})
}

// Helper function to load acceptance test Terraform configs
func loadAccTestTerraform(filename string) string {
	accTestConfig, err := helpers.ParseHCLFile(fmt.Sprintf("tests/terraform/acceptance/%s", filename))
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
