package graphBetaAssignmentFilter_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAssignmentFilterResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		CheckDestroy: testAccCheckAssignmentFilterDestroy,
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccAssignmentFilterConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_assignment_filter.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "display_name", "Test Acceptance Assignment Filter"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "platform", "windows10AndLater"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "rule", "(device.osVersion -startsWith \"10.0\")"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "assignment_filter_management_type", "devices"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_assignment_filter.test", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_assignment_filter.test", "last_modified_date_time"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_assignment_filter.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update to maximal configuration
			{
				Config: testAccAssignmentFilterConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_assignment_filter.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "display_name", "Test Acceptance Assignment Filter - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "platform", "windows10AndLater"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "assignment_filter_management_type", "devices"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "role_scope_tags.#", "2"),
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
		//"windowsPhone81", , causes a 500 error in acc tests
		//"windows81AndLater", , causes a 500 error in acc tests
		"windows10AndLater",
		// "androidWorkProfile", causes a 500 error in acc tests
		//"unknown", causes a 500 error in acc tests
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
				ExternalProviders: map[string]resource.ExternalProvider{
					"random": {
						Source:            "hashicorp/random",
						VersionConstraint: ">= 3.7.2",
					},
				},
				CheckDestroy: testAccCheckAssignmentFilterDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccAssignmentFilterConfig_platform(platform),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", platform), "id"),
							resource.TestCheckResourceAttr(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", platform), "platform", platform),
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAssignmentFilterConfig_complexRule(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_assignment_filter.complex", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.complex", "display_name", "Test Complex Rule Assignment Filter"),
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
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAssignmentFilterConfig_roleScopeTags(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_assignment_filter.role_tags", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.role_tags", "role_scope_tags.#", "3"),
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
				ExternalProviders: map[string]resource.ExternalProvider{
					"random": {
						Source:            "hashicorp/random",
						VersionConstraint: ">= 3.7.2",
					},
				},
				Steps: []resource.TestStep{
					{
						Config: testAccAssignmentFilterConfig_managementType(mgmtType),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", mgmtType), "id"),
							resource.TestCheckResourceAttr(fmt.Sprintf("microsoft365_graph_beta_device_management_assignment_filter.%s", mgmtType), "assignment_filter_management_type", mgmtType),
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

// testAccCheckAssignmentFilterDestroy verifies that assignment filters have been destroyed
func testAccCheckAssignmentFilterDestroy(s *terraform.State) error {
	// Get a Graph client using the same configuration as acceptance tests
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_assignment_filter" {
			continue
		}

		// Attempt to get the assignment filter by ID
		_, err := graphClient.
			DeviceManagement().
			AssignmentFilters().
			ByDeviceAndAppManagementAssignmentFilterId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 ||
				errorInfo.ErrorCode == "ResourceNotFound" ||
				errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue // Resource successfully destroyed
			}
			return fmt.Errorf("error checking if assignment filter %s was destroyed: %v", rs.Primary.ID, err)
		}

		// If we can still get the resource, it wasn't destroyed
		return fmt.Errorf("assignment filter %s still exists", rs.Primary.ID)
	}

	return nil
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
