package graphBetaAssignmentFilter_test

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
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAssignmentFilter "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/assignment_filter"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name constructed from exported constants
	resourceType = constants.PROVIDER_NAME + "_" + graphBetaAssignmentFilter.ResourceName

	// testResource is the test resource implementation for assignment filters
	testResource = graphBetaAssignmentFilter.AssignmentFilterTestResource{}
)

func TestAccAssignmentFilterResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			0,
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
				Config: testAccAssignmentFilterConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("id").Exists(),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("display_name").HasValue("Test Acceptance Assignment Filter"),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("platform").HasValue("windows10AndLater"),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("rule").HasValue("(device.osVersion -startsWith \"10.0\")"),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("assignment_filter_management_type").HasValue("devices"),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("created_date_time").Exists(),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("last_modified_date_time").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      "microsoft365_graph_beta_device_management_assignment_filter.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating")
				},
				Config: testAccAssignmentFilterConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("display_name").HasValue("Test Acceptance Assignment Filter - Updated"),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("description").HasValue("Updated description for acceptance testing"),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("platform").HasValue("windows10AndLater"),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("assignment_filter_management_type").HasValue("devices"),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.test").Key("role_scope_tags.#").HasValue("2"),
				),
			},
		},
	})
}

func TestAccAssignmentFilterResource_MultiPlatform(t *testing.T) {
	platforms := []string{
		"android",
		"androidForWork",
		"iOS",
		"macOS",
		//"windowsPhone81", // causes a 500 error in acc tests
		//"windows81AndLater", // causes a 500 error in acc tests
		"windows10AndLater",
		//"androidWorkProfile", // causes a 500 error in acc tests
		//"unknown", // causes a 500 error in acc tests
		"androidAOSP",
		"androidMobileApplicationManagement",
		"iOSMobileApplicationManagement",
		"windowsMobileApplicationManagement",
	}

	for _, platform := range platforms {
		t.Run(fmt.Sprintf("platform_%s", platform), func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { mocks.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
				CheckDestroy: destroy.CheckDestroyedAllFunc(
					testResource,
					resourceType,
					0,
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
							testlog.StepAction(resourceType, fmt.Sprintf("Creating with platform: %s", platform))
						},
						Config: testAccAssignmentFilterConfig_platform(platform),
						Check: resource.ComposeTestCheckFunc(
							check.That(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", platform)).ExistsInGraph(testResource),
							check.That(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", platform)).Key("id").Exists(),
							check.That(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", platform)).Key("platform").HasValue(platform),
							func(state *terraform.State) error {
								// Check rule and management type based on platform type
								resourceName := fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", platform)
								rs, ok := state.RootModule().Resources[resourceName]
								if !ok {
									return fmt.Errorf("resource not found: %s", resourceName)
								}

								isAppPlatform := platform == "androidMobileApplicationManagement" ||
									platform == "iOSMobileApplicationManagement" ||
									platform == "windowsMobileApplicationManagement"

								if isAppPlatform {
									expectedRule := `(app.osVersion -startsWith "14.0")`
									expectedManagementType := "apps"
									if rs.Primary.Attributes["rule"] != expectedRule {
										return fmt.Errorf("expected rule %s, got %s", expectedRule, rs.Primary.Attributes["rule"])
									}
									if rs.Primary.Attributes["assignment_filter_management_type"] != expectedManagementType {
										return fmt.Errorf("expected management type %s, got %s", expectedManagementType, rs.Primary.Attributes["assignment_filter_management_type"])
									}
								} else {
									expectedRule := `(device.osVersion -startsWith "10.0")`
									expectedManagementType := "devices"
									if rs.Primary.Attributes["rule"] != expectedRule {
										return fmt.Errorf("expected rule %s, got %s", expectedRule, rs.Primary.Attributes["rule"])
									}
									if rs.Primary.Attributes["assignment_filter_management_type"] != expectedManagementType {
										return fmt.Errorf("expected management type %s, got %s", expectedManagementType, rs.Primary.Attributes["assignment_filter_management_type"])
									}
								}
								return nil
							},
						),
					},
				},
			})
		})
	}
}

func TestAccAssignmentFilterResource_ComplexRule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			0,
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
					testlog.StepAction(resourceType, "Creating with complex rule")
				},
				Config: testAccAssignmentFilterConfig_complexRule(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_device_management_assignment_filter.complex").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.complex").Key("id").Exists(),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.complex").Key("display_name").HasValue("Test Complex Rule Assignment Filter"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.complex", "rule", regexp.MustCompile(`device\.osVersion`)),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.complex", "rule", regexp.MustCompile(`device\.manufacturer`)),
				),
			},
		},
	})
}

func TestAccAssignmentFilterResource_RoleScopeTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			0,
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
					testlog.StepAction(resourceType, "Creating with role scope tags")
					testlog.WaitForConsistency("Microsoft Graph", 15*time.Second)
					time.Sleep(15 * time.Second)
				},
				Config: testAccAssignmentFilterConfig_roleScopeTags(),
				Check: resource.ComposeTestCheckFunc(
					check.That("microsoft365_graph_beta_device_management_assignment_filter.role_tags").ExistsInGraph(testResource),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.role_tags").Key("id").Exists(),
					check.That("microsoft365_graph_beta_device_management_assignment_filter.role_tags").Key("role_scope_tags.#").HasValue("3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_assignment_filter.role_tags", "role_scope_tags.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_assignment_filter.role_tags", "role_scope_tags.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_assignment_filter.role_tags", "role_scope_tags.*", "2"),
				),
			},
		},
	})
}

func TestAccAssignmentFilterResource_ManagementTypes(t *testing.T) {
	managementTypes := []string{"devices", "apps"}

	for _, mgmtType := range managementTypes {
		t.Run(fmt.Sprintf("management_type_%s", mgmtType), func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { mocks.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
				CheckDestroy: destroy.CheckDestroyedAllFunc(
					testResource,
					resourceType,
					0,
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
							testlog.StepAction(resourceType, fmt.Sprintf("Creating with management type: %s", mgmtType))
						},
						Config: testAccAssignmentFilterConfig_managementType(mgmtType),
						Check: resource.ComposeTestCheckFunc(
							check.That(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", mgmtType)).ExistsInGraph(testResource),
							check.That(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", mgmtType)).Key("id").Exists(),
							check.That(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", mgmtType)).Key("assignment_filter_management_type").HasValue(mgmtType),
							func(state *terraform.State) error {
								// Check platform and rule based on management type
								resourceName := fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", mgmtType)
								rs, ok := state.RootModule().Resources[resourceName]
								if !ok {
									return fmt.Errorf("resource not found: %s", resourceName)
								}

								if mgmtType == "apps" {
									expectedPlatform := "androidMobileApplicationManagement"
									expectedRule := `(app.osVersion -startsWith "14.0")`
									if rs.Primary.Attributes["platform"] != expectedPlatform {
										return fmt.Errorf("expected platform %s, got %s", expectedPlatform, rs.Primary.Attributes["platform"])
									}
									if rs.Primary.Attributes["rule"] != expectedRule {
										return fmt.Errorf("expected rule %s, got %s", expectedRule, rs.Primary.Attributes["rule"])
									}
								} else {
									expectedPlatform := "windows10AndLater"
									expectedRule := `(device.osVersion -startsWith "10.0")`
									if rs.Primary.Attributes["platform"] != expectedPlatform {
										return fmt.Errorf("expected platform %s, got %s", expectedPlatform, rs.Primary.Attributes["platform"])
									}
									if rs.Primary.Attributes["rule"] != expectedRule {
										return fmt.Errorf("expected rule %s, got %s", expectedRule, rs.Primary.Attributes["rule"])
									}
								}
								return nil
							},
						),
					},
				},
			})
		})
	}
}

// Test configuration functions
func testAccAssignmentFilterConfig_minimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccAssignmentFilterConfig_maximal() string {
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccAssignmentFilterConfig_platform(platform string) string {
	data := struct {
		Platform       string
		Rule           string
		ManagementType string
	}{
		Platform: platform,
	}

	// Check if this is an app platform
	isAppPlatform := platform == "androidMobileApplicationManagement" ||
		platform == "iOSMobileApplicationManagement" ||
		platform == "windowsMobileApplicationManagement"

	if isAppPlatform {
		data.Rule = `(app.osVersion -startsWith \"14.0\")`
		data.ManagementType = "apps"
	} else {
		data.Rule = `(device.osVersion -startsWith \"10.0\")`
		data.ManagementType = "devices"
	}

	config := mocks.LoadTerraformTemplateFile("resource_platform_template.tf", data)
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccAssignmentFilterConfig_complexRule() string {
	config := mocks.LoadTerraformConfigFile("resource_complex_rule.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccAssignmentFilterConfig_roleScopeTags() string {
	config := mocks.LoadTerraformConfigFile("resource_role_scope_tags.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccAssignmentFilterConfig_managementType(managementType string) string {
	data := struct {
		ManagementType string
		Platform       string
		Rule           string
	}{
		ManagementType: managementType,
	}

	if managementType == "apps" {
		data.Platform = "androidMobileApplicationManagement"
		data.Rule = `(app.osVersion -startsWith \"14.0\")`
	} else {
		data.Platform = "windows10AndLater"
		data.Rule = `(device.osVersion -startsWith \"10.0\")`
	}

	config := mocks.LoadTerraformTemplateFile("resource_management_type_template.tf", data)
	return acceptance.ConfiguredM365ProviderBlock(config)
}
