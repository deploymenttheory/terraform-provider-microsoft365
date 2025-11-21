package graphBetaUserLicenseAssignment_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaUserLicenseAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/license_assignment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// resourceType is declared here and shared across both unit and acceptance tests
	// resourceType = graphBetaUserLicenseAssignment.ResourceName // Already declared in resource_test.go

	// testResource is the test resource implementation for user license assignments
	testResource = graphBetaUserLicenseAssignment.UserLicenseAssignmentTestResource{}
)

func TestAccUserLicenseAssignmentResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccLicenseAssignmentConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("add_licenses.#").HasValue("2"),
					check.That(resourceType+".minimal").Key("add_licenses.0.sku_id").HasValue("a403ebcc-fae0-4ca2-8c8c-7a907fd6c235"),
					check.That(resourceType+".minimal").Key("add_licenses.1.sku_id").HasValue("f30db892-07e9-47e9-837c-80727f46fd3d"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"add_licenses",
					"remove_licenses",
				},
			},
		},
	})
}

func TestAccUserLicenseAssignmentResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccLicenseAssignmentConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".dependancy").ExistsInGraph(testResource),
					check.That(resourceType+".dependancy").Key("id").Exists(),
					check.That(resourceType+".dependancy").Key("add_licenses.#").HasValue("2"),
					check.That(resourceType+".dependancy").Key("add_licenses.0.sku_id").HasValue("a403ebcc-fae0-4ca2-8c8c-7a907fd6c235"),
					check.That(resourceType+".dependancy").Key("add_licenses.1.sku_id").HasValue("f30db892-07e9-47e9-837c-80727f46fd3d"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".dependancy",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"add_licenses",
					"remove_licenses",
				},
			},
		},
	})
}

func testAccLicenseAssignmentConfig_minimal() string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		panic("failed to load acceptance resource_minimal.tf: " + err.Error())
	}

	return fmt.Sprintf(`
terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "microsoft365" {}

%s
`, acceptanceTestConfig)
}

func testAccLicenseAssignmentConfig_maximal() string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		panic("failed to load acceptance resource_maximal.tf: " + err.Error())
	}

	return fmt.Sprintf(`
provider "microsoft365" {}

%s
`, acceptanceTestConfig)
}
