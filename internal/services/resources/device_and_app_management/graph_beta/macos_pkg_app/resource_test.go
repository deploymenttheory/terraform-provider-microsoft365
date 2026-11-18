package graphBetaMacOSPKGApp_test

import (
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestUnitResourceMacOSPkgApp_01_IsFeaturedPlan(t *testing.T) {
	for _, test := range []struct {
		name   string
		config string
		want   bool
	}{
		{name: "omitted"},
		{name: "null", config: "is_featured = null"},
		{name: "false", config: "is_featured = false"},
		{name: "true", config: "is_featured = true", want: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
			if err != nil {
				t.Fatal(err)
			}
			config = strings.Replace(config, "# FEATURED", test.config, 1)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config:             config,
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPreRefresh: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue(
								"microsoft365_graph_beta_device_and_app_management_macos_pkg_app.test",
								tfjsonpath.New("is_featured"),
								knownvalue.Bool(test.want),
							),
						},
					},
				}},
			})
		})
	}
}
