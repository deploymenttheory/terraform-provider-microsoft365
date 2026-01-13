package graphBetaTenantInformation_test

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTenantInformationDataSource_ByTenantId(t *testing.T) {
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
				Config: testAccConfigByTenantId(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_tenant_id").Key("filter_type").HasValue("tenant_id"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("filter_value").HasValue("6babcaad-604b-40ac-a9d7-9fd97c0b779f"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("tenant_id").HasValue("6babcaad-604b-40ac-a9d7-9fd97c0b779f"),
					check.That("data."+dataSourceType+".by_tenant_id").Key("display_name").IsSet(),
					check.That("data."+dataSourceType+".by_tenant_id").Key("default_domain_name").IsSet(),
					check.That("data."+dataSourceType+".by_tenant_id").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccTenantInformationDataSource_ByDomainName(t *testing.T) {
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
				Config: testAccConfigByDomainName(),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".by_domain_name").Key("filter_type").HasValue("domain_name"),
					check.That("data."+dataSourceType+".by_domain_name").Key("filter_value").HasValue("deploymenttheory.com"),
					check.That("data."+dataSourceType+".by_domain_name").Key("tenant_id").IsSet(),
					check.That("data."+dataSourceType+".by_domain_name").Key("display_name").IsSet(),
					check.That("data."+dataSourceType+".by_domain_name").Key("default_domain_name").IsSet(),
					check.That("data."+dataSourceType+".by_domain_name").Key("id").IsSet(),
				),
			},
		},
	})
}

// Acceptance test configuration functions
func testAccConfigByTenantId() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_by_tenant_id.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigByDomainName() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_by_domain_name.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance test config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
