package graphBetaWindowsEnrollmentStatusPage_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWindowsEnrollmentStatusPageResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsEnrollmentStatusPageDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsEnrollmentStatusPageConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "display_name", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("display_name should not be empty")
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "show_installation_progress", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "block_device_setup_retry_by_user", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "allow_device_reset_on_install_failure", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "install_progress_timeout_in_minutes", "120"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "custom_error_message", "Contact IT support for assistance"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestAccWindowsEnrollmentStatusPageResource_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsEnrollmentStatusPageDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsEnrollmentStatusPageConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "display_name", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("display_name should not be empty")
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "show_installation_progress", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "install_progress_timeout_in_minutes", "180"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "track_install_progress_for_autopilot_only", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "allow_log_collection_on_install_failure", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestAccWindowsEnrollmentStatusPageResource_WithAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsEnrollmentStatusPageDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsEnrollmentStatusPageConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_enrollment_status_page.with_assignments", "id"),
					resource.TestCheckResourceAttrWith("microsoft365_graph_beta_device_management_windows_enrollment_status_page.with_assignments", "display_name", func(value string) error {
						if len(value) == 0 {
							return fmt.Errorf("display_name should not be empty")
						}
						return nil
					}),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.with_assignments", "assignments.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_enrollment_status_page.with_assignments", "assignments.*", map[string]string{
						"type": "allDevicesAssignmentTarget",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_enrollment_status_page.with_assignments", "assignments.*", map[string]string{
						"type": "allLicensedUsersAssignmentTarget",
					}),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_windows_enrollment_status_page.with_assignments", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestAccWindowsEnrollmentStatusPageResource_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsEnrollmentStatusPageDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsEnrollmentStatusPageConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.minimal", "install_progress_timeout_in_minutes", "120"),
				),
			},
			{
				Config: testAccWindowsEnrollmentStatusPageConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", "install_progress_timeout_in_minutes", "180"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_windows_enrollment_status_page.maximal", ImportState: true, ImportStateVerify: true},
		},
	})
}

// Configuration Functions
func testAccWindowsEnrollmentStatusPageConfig_minimal() string {
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal acceptance test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccWindowsEnrollmentStatusPageConfig_maximal() string {
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load maximal acceptance test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccWindowsEnrollmentStatusPageConfig_withAssignments() string {
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_with_assignments.tf")
	if err != nil {
		log.Fatalf("Failed to load assignments acceptance test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccCheckWindowsEnrollmentStatusPageDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_windows_enrollment_status_page" {
			continue
		}
		_, err := graphClient.
			DeviceManagement().
			DeviceEnrollmentConfigurations().
			ByDeviceEnrollmentConfigurationId(rs.Primary.ID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if windows enrollment status page %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("windows enrollment status page %s still exists", rs.Primary.ID)
	}
	return nil
}
