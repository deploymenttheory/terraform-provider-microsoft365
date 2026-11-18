package graphBetaMacOSPKGApp_test

import (
	"strings"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaMacOSPKGApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/macos_pkg_app"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceType = graphBetaMacOSPKGApp.ResourceName

func TestAccResourceMacOSPkgApp_01_IsFeaturedLifecycle(t *testing.T) {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		t.Fatal(err)
	}
	name := resourceType + ".test"
	testResource := macOSPKGAppTestResource{}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					check.That(name).ExistsInGraph(testResource),
					resource.TestCheckResourceAttr(name, "is_featured", "false"),
					resource.TestCheckResourceAttr(name, "publishing_state", "published"),
					resource.TestCheckResourceAttr(name, "macos_pkg_app.primary_bundle_id", "com.github.macadmins.Nudge"),
					resource.TestCheckResourceAttr(name, "macos_pkg_app.primary_bundle_version", "2.1.4.81863"),
				),
			},
			{Config: config, PlanOnly: true, ExpectNonEmptyPlan: false},
			{
				ResourceName: name, ImportState: true, ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"app_installer", "timeouts"},
			},
			{
				Config: strings.Replace(config, "# FEATURED", "is_featured = true", 1),
				Check:  resource.TestCheckResourceAttr(name, "is_featured", "true"),
			},
			{
				Config: strings.Replace(config, "# FEATURED", "is_featured = false", 1),
				Check:  resource.TestCheckResourceAttr(name, "is_featured", "false"),
			},
			{
				Config: strings.Replace(config, "# FEATURED", "is_featured = true", 1),
				Check:  resource.TestCheckResourceAttr(name, "is_featured", "true"),
			},
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(name, "is_featured", "false"),
			},
			{Config: config, PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}
