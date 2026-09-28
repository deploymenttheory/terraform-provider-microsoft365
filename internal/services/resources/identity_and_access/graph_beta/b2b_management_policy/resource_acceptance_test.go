package graphBetaIdentityAndAccessB2bManagementPolicy_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaB2bManagementPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	testResource = graphBetaB2bManagementPolicy.B2bManagementPolicyTestResource{}
)

// is_organization_default stays false so the test does not change tenant-wide B2B settings.
func TestAccResourceB2bManagementPolicy_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			10*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating B2B management policy")
				},
				Config: testAccConfig("tests/terraform/acceptance/resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-b2b-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("is_organization_default").HasValue("false"),
					check.That(resourceType+".test").Key("definition.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing B2B management policy")
				},
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating B2B management policy")
				},
				Config: testAccConfig("tests/terraform/acceptance/resource_02_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-b2b-updated-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test").Key("definition.0").MatchesRegex(regexp.MustCompile(`AllowedDomains`)),
				),
			},
		},
	})
}

func testAccConfig(path string) string {
	accTestConfig, err := helpers.ParseHCLFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to load B2B management policy acceptance config %s: %s", path, err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}
