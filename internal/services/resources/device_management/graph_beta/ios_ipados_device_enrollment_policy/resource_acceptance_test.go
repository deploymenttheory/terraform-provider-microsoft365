package graphBetaIOSiPadOSDeviceEnrollmentPolicy_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaIOSiPadOSDeviceEnrollmentPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/ios_ipados_device_enrollment_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIOSiPadOSDeviceEnrollmentPolicy.ResourceName

	// testResource is the test resource implementation for iOS/iPadOS ADE enrollment policies
	testResource = graphBetaIOSiPadOSDeviceEnrollmentPolicy.IOSiPadOSDeviceEnrollmentPolicyTestResource{}
)

// loadAcceptanceTestTerraform loads a Terraform config from the acceptance test directory.
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Scenario 01: Minimal iOS/iPadOS ADE enrollment policy
func TestAccResourceIOSiPadOSDeviceEnrollmentPolicy_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal iOS/iPadOS ADE enrollment policy")
				},
				Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("ios/ipados device enrollment policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("name").HasValue("acc-test-ios-ade-minimal"),
					check.That(resourceType+".minimal").Key("requires_user_authentication").HasValue("false"),
					check.That(resourceType+".minimal").Key("platforms").HasValue("iOS"),
					check.That(resourceType+".minimal").Key("technologies").HasValue("enrollment"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal iOS/iPadOS ADE enrollment policy")
				},
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Scenario 02: Maximal iOS/iPadOS ADE enrollment policy
func TestAccResourceIOSiPadOSDeviceEnrollmentPolicy_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal iOS/iPadOS ADE enrollment policy with the full settings tree")
				},
				Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("ios/ipados device enrollment policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal").Key("name").HasValue("acc-test-ios-ade-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("iOS/iPadOS ADE enrollment policy exercising the full settings tree"),
					check.That(resourceType+".maximal").Key("require_setup_assistant_with_modern_authentication").HasValue("true"),
					check.That(resourceType+".maximal").Key("await_final_configuration").HasValue("true"),
					check.That(resourceType+".maximal").Key("device_name_template").HasValue("{{DEVICETYPE}}-{{SERIAL}}"),
					check.That(resourceType+".maximal").Key("cellular_data_activation_url").HasValue("http://activation.carrier.net"),
					check.That(resourceType+".maximal").Key("locked_enrollment_enabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("support_department").HasValue("IT Support"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal iOS/iPadOS ADE enrollment policy")
				},
				ResourceName:            resourceType + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Scenario 03: Minimal to Maximal Update
func TestAccResourceIOSiPadOSDeviceEnrollmentPolicy_03_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal iOS/iPadOS ADE enrollment policy for update test")
				},
				Config: loadAcceptanceTestTerraform("003_scenario_minimal_to_maximal_step_01.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("ios/ipados device enrollment policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".update_test").ExistsInGraph(testResource),
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("name").HasValue("acc-test-ios-ade-update"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("004_scenario_minimal_to_maximal_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("ios/ipados device enrollment policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".update_test").ExistsInGraph(testResource),
					check.That(resourceType+".update_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".update_test").Key("name").HasValue("acc-test-ios-ade-update-updated"),
					check.That(resourceType+".update_test").Key("require_setup_assistant_with_modern_authentication").HasValue("true"),
					check.That(resourceType+".update_test").Key("device_name_template").HasValue("{{DEVICETYPE}}-{{SERIAL}}"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing updated iOS/iPadOS ADE enrollment policy")
				},
				ResourceName:            resourceType + ".update_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Scenario 04: Maximal to Minimal Update
func TestAccResourceIOSiPadOSDeviceEnrollmentPolicy_04_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal iOS/iPadOS ADE enrollment policy for downgrade test")
				},
				Config: loadAcceptanceTestTerraform("005_scenario_maximal_to_minimal_step_01.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("ios/ipados device enrollment policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".downgrade_test").ExistsInGraph(testResource),
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("name").HasValue("acc-test-ios-ade-downgrade"),
					check.That(resourceType+".downgrade_test").Key("enable_authentication_via_company_portal").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Downgrading to minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("006_scenario_maximal_to_minimal_step_02.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("ios/ipados device enrollment policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".downgrade_test").ExistsInGraph(testResource),
					check.That(resourceType+".downgrade_test").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".downgrade_test").Key("name").HasValue("acc-test-ios-ade-downgrade-minimal"),
					check.That(resourceType+".downgrade_test").Key("requires_user_authentication").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing downgraded iOS/iPadOS ADE enrollment policy")
				},
				ResourceName:            resourceType + ".downgrade_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Scenario 05: Default Policy Assignment (setDefaultProfile action)
func TestAccResourceIOSiPadOSDeviceEnrollmentPolicy_05_DefaultPolicyAssignment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating iOS/iPadOS ADE enrollment policy as the DEP token default")
				},
				Config: loadAcceptanceTestTerraform("007_scenario_default_policy_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("ios/ipados device enrollment policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".default_assignment").ExistsInGraph(testResource),
					check.That(resourceType+".default_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".default_assignment").Key("is_default_policy_assignment").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing default iOS/iPadOS ADE enrollment policy")
				},
				ResourceName:            resourceType + ".default_assignment",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}
