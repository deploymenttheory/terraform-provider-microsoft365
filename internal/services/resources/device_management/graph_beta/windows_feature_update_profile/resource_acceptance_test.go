package graphBetaWindowsFeatureUpdateProfile_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWindowsFeatureUpdateProfileResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsFeatureUpdateProfileDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsFeatureUpdateProfileConfig_windows10_22h2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", "display_name", "Acceptance - Windows 10 22H2 Feature Update Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", "feature_update_version", "Windows 10, version 22H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", "install_feature_updates_optional", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", "install_latest_windows10_on_windows11_ineligible_device", "false"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", ImportState: true, ImportStateVerify: true},
			{
				Config: testAccWindowsFeatureUpdateProfileConfig_windows11_22h2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "display_name", "Acceptance - Windows 11 22H2 Feature Update Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "feature_update_version", "Windows 11, version 22H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "install_feature_updates_optional", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "install_latest_windows10_on_windows11_ineligible_device", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "rollout_settings.offer_start_date_time_in_utc", "2029-08-01T00:00:00Z"),
				),
			},
			{
				Config: testAccWindowsFeatureUpdateProfileConfig_windows11_23h2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_23h2", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_23h2", "display_name", "Acceptance - Windows 11 23H2 Feature Update Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_23h2", "feature_update_version", "Windows 11, version 23H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_23h2", "rollout_settings.offer_start_date_time_in_utc", "2029-08-01T00:00:00Z"),
				),
			},
			{
				Config: testAccWindowsFeatureUpdateProfileConfig_windows11_24h2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_24h2", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_24h2", "display_name", "Acceptance - Windows 11 24H2 Feature Update Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_24h2", "feature_update_version", "Windows 11, version 24H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_24h2", "rollout_settings.offer_start_date_time_in_utc", "2029-08-01T00:00:00Z"),
				),
			},
		},
	})
}

func TestAccWindowsFeatureUpdateProfileResource_Assignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsFeatureUpdateProfileDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsFeatureUpdateProfileConfig_withAssignments(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_assignments", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_assignments", "display_name", "Acceptance - Windows Feature Update Profile with Assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_assignments", "assignments.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_assignments", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
				),
			},
		},
	})
}

func TestAccWindowsFeatureUpdateProfileResource_AllFeatureVersions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWindowsFeatureUpdateProfileDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWindowsFeatureUpdateProfileConfig_windows11_24h2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_24h2", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_24h2", "display_name", "Acceptance - Windows 11 24H2 Feature Update Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_24h2", "feature_update_version", "Windows 11, version 24H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_24h2", "rollout_settings.offer_start_date_time_in_utc", "2029-08-01T00:00:00Z"),
				),
			},
			{
				Config: testAccWindowsFeatureUpdateProfileConfig_windows11_23h2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_23h2", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_23h2", "display_name", "Acceptance - Windows 11 23H2 Feature Update Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_23h2", "feature_update_version", "Windows 11, version 23H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_23h2", "rollout_settings.offer_start_date_time_in_utc", "2029-08-01T00:00:00Z"),
				),
			},
			{
				Config: testAccWindowsFeatureUpdateProfileConfig_windows11_22h2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "display_name", "Acceptance - Windows 11 22H2 Feature Update Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "feature_update_version", "Windows 11, version 22H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_22h2", "rollout_settings.offer_start_date_time_in_utc", "2029-08-01T00:00:00Z"),
				),
			},
			{
				Config: testAccWindowsFeatureUpdateProfileConfig_windows10_22h2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", "display_name", "Acceptance - Windows 10 22H2 Feature Update Profile"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", "feature_update_version", "Windows 10, version 22H2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_windows_feature_update_profile.test_win10_22h2", "install_feature_updates_optional", "false"),
				),
			},
		},
	})
}


func testAccWindowsFeatureUpdateProfileConfig_withAssignments() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_with_assignments.tf")
	if err != nil {
		log.Fatalf("Failed to load assignments test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccWindowsFeatureUpdateProfileConfig_windows11_24h2() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_windows11_24h2.tf")
	if err != nil {
		log.Fatalf("Failed to load Windows 11 24H2 test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccWindowsFeatureUpdateProfileConfig_windows11_23h2() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_windows11_23h2.tf")
	if err != nil {
		log.Fatalf("Failed to load Windows 11 23H2 test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccWindowsFeatureUpdateProfileConfig_windows11_22h2() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_windows11_22h2.tf")
	if err != nil {
		log.Fatalf("Failed to load Windows 11 22H2 test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccWindowsFeatureUpdateProfileConfig_windows10_22h2() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_windows10_22h2.tf")
	if err != nil {
		log.Fatalf("Failed to load Windows 10 22H2 test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccCheckWindowsFeatureUpdateProfileDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_windows_feature_update_profile" {
			continue
		}
		_, err := graphClient.DeviceManagement().WindowsFeatureUpdateProfiles().ByWindowsFeatureUpdateProfileId(rs.Primary.ID).Get(ctx, nil)
		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", rs.Primary.ID)
				continue
			}
			return fmt.Errorf("error checking if windows feature update profile %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("windows feature update profile %s still exists", rs.Primary.ID)
	}
	return nil
}
