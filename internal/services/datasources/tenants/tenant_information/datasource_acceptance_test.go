package graphBetaTenantInformation_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

func TestAccDatasourceTenantInformation_Lookup(t *testing.T) {
	id := os.Getenv("M365_TENANT_ID")
	config := fmt.Sprintf(
		"data %q \"by_id\" { tenant_id = %q }\ndata %q \"by_domain\" { domain_name = data.%s.by_id.items[0].default_domain_name }",
		dataSourceType,
		id,
		dataSourceType,
		dataSourceType,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{Config: "provider \"microsoft365\" {}\n" + config, Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("data."+dataSourceType+".by_id", "items.0.tenant_id", id),
			resource.TestCheckResourceAttr("data."+dataSourceType+".by_domain", "items.0.tenant_id", id),
			resource.TestCheckResourceAttrPair("data."+dataSourceType+".by_id", "id", "data."+dataSourceType+".by_domain", "id"),
		)}},
	})
}
