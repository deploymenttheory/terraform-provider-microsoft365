package graphBetaUserLicenseAssignment_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaUserLicenseAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/license_assignment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("user license assignment", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("sku_id").HasValue("f30db892-07e9-47e9-837c-80727f46fd3d"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
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
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("user license assignment", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".dependancy").ExistsInGraph(testResource),
					check.That(resourceType+".dependancy").Key("id").Exists(),
					check.That(resourceType+".dependancy").Key("sku_id").HasValue("f30db892-07e9-47e9-837c-80727f46fd3d"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".dependancy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLicenseAssignmentConfig_minimal() string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		panic("failed to load acceptance resource_minimal.tf: " + err.Error())
	}

	return acceptanceTestConfig
}

func testAccLicenseAssignmentConfig_maximal() string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		panic("failed to load acceptance resource_maximal.tf: " + err.Error())
	}

	return acceptanceTestConfig
}
