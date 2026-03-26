package graphBetaDirectoryRole_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccDatasourceDirectoryRole_01_ListAll(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("01_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("items.#").IsNotEmpty(),
				),
			},
		},
	})
}

func TestAccDatasourceDirectoryRole_02_ByDisplayName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("02_by_display_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.display_name").HasValue("User Administrator"),
					check.That(dataSourceType+".test").Key("items.0.role_template_id").HasValue("fe930be7-5e62-47db-91af-98c3a49a38b1"),
					check.That(dataSourceType+".test").Key("items.0.id").IsNotEmpty(),
				),
			},
		},
	})
}
