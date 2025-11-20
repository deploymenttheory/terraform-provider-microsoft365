package graphBetaUserLicenseAssignment_test

import (
	"fmt"
	"os"
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
	testUserID := os.Getenv("TEST_USER_ID_1")
	if testUserID == "" {
		t.Skip("TEST_USER_ID_1 environment variable must be set for acceptance tests")
	}

	testLicenseSkuID := os.Getenv("TEST_LICENSE_SKU_ID_1")
	if testLicenseSkuID == "" {
		t.Skip("TEST_LICENSE_SKU_ID_1 environment variable must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			5*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccLicenseAssignmentConfig_minimal(testUserID, testLicenseSkuID),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("user_id").HasValue(testUserID),
					check.That(resourceType+".minimal").Key("add_licenses.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("add_licenses.0.sku_id").HasValue(testLicenseSkuID),
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
	testUserID := os.Getenv("TEST_USER_ID_2")
	if testUserID == "" {
		t.Skip("TEST_USER_ID_2 environment variable must be set for acceptance tests")
	}

	testLicenseSkuID1 := os.Getenv("TEST_LICENSE_SKU_ID_1")
	if testLicenseSkuID1 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_1 environment variable must be set for acceptance tests")
	}

	testLicenseSkuID2 := os.Getenv("TEST_LICENSE_SKU_ID_2")
	if testLicenseSkuID2 == "" {
		t.Skip("TEST_LICENSE_SKU_ID_2 environment variable must be set for acceptance tests")
	}

	testServicePlanID := os.Getenv("TEST_SERVICE_PLAN_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			5*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating")
				},
				Config: testAccLicenseAssignmentConfig_maximal(testUserID, testLicenseSkuID1, testLicenseSkuID2, testServicePlanID),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("user_id").HasValue(testUserID),
					check.That(resourceType+".maximal").Key("add_licenses.#").HasValue("2"),
					check.That(resourceType+".maximal").Key("add_licenses.0.sku_id").HasValue(testLicenseSkuID1),
					check.That(resourceType+".maximal").Key("add_licenses.1.sku_id").HasValue(testLicenseSkuID2),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLicenseAssignmentConfig_minimal(userID, licenseSkuID string) string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		panic("failed to load acceptance resource_minimal.tf: " + err.Error())
	}

	return fmt.Sprintf(`
%s

provider "microsoft365" {}

`, acceptanceTestConfig) + fmt.Sprintf(`
variable "test_user_id" {
  default = "%s"
}

variable "test_license_sku_id" {
  default = "%s"
}
`, userID, licenseSkuID)
}

func testAccLicenseAssignmentConfig_maximal(userID, licenseSkuID1, licenseSkuID2, servicePlanID string) string {
	acceptanceTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		panic("failed to load acceptance resource_maximal.tf: " + err.Error())
	}

	return fmt.Sprintf(`
%s

provider "microsoft365" {}

`, acceptanceTestConfig) + fmt.Sprintf(`
variable "test_user_id" {
  default = "%s"
}

variable "test_license_sku_id_1" {
  default = "%s"
}

variable "test_license_sku_id_2" {
  default = "%s"
}

variable "test_service_plan_id" {
  default = "%s"
}
`, userID, licenseSkuID1, licenseSkuID2, servicePlanID)
}
