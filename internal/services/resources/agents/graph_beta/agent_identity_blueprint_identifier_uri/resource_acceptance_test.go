package graphBetaAgentIdentityBlueprintIdentifierUri_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAgentIdentityBlueprintIdentifierUriResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating identifier URI on agent identity blueprint")
				},
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("identifier URI", 10*time.Second)
						time.Sleep(10 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("identifier_uri").MatchesRegex(regexp.MustCompile(`^api://[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("scope.value").HasValue("access_agent"),
					check.That(resourceType+".test_minimal").Key("scope.type").HasValue("User"),
					check.That(resourceType+".test_minimal").Key("scope.is_enabled").HasValue("true"),
				),
			},
			{
				ResourceName: resourceType + ".test_minimal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test_minimal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test_minimal")
					}
					blueprintID := rs.Primary.Attributes["blueprint_id"]
					identifierUri := rs.Primary.Attributes["identifier_uri"]
					return fmt.Sprintf("%s/%s", blueprintID, identifierUri), nil
				},
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "blueprint_id",
			},
		},
	})
}

func testAccConfigMinimal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		panic(err)
	}
	return content
}
