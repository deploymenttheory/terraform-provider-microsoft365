package graphBetaMacOSDepEnrollmentProfile_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceMacOSDepEnrollmentProfile_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMacOSDepEnrollmentProfileDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccMacOSDepEnrollmentProfileConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "display_name", "acc-test-macos-dep-enrollment-profile-minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "description", "macOS DEP enrollment profile minimal acceptance test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "requires_user_authentication", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "configuration_endpoint_url"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", "dep_onboarding_settings_id"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.minimal", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"admin_account_password"}},
			{
				Config: testAccMacOSDepEnrollmentProfileConfig_skipSetup(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "display_name", "acc-test-macos-dep-enrollment-profile-skip-setup"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "await_device_configured", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "supervised_mode_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "admin_account_user_name", "localadmin"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.skip_setup", "enabled_skip_keys.#", "6"),
				),
			},
		},
	})
}

func testAccMacOSDepEnrollmentProfileConfig_minimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccMacOSDepEnrollmentProfileConfig_skipSetup() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_skip_setup.tf")
	if err != nil {
		log.Fatalf("Failed to load skip setup test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccCheckMacOSDepEnrollmentProfileDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" {
			continue
		}

		// The ID format is: depOnboardingSettingsId_enrollmentProfileId
		resourceId := rs.Primary.ID
		depOnboardingSettingsId := rs.Primary.Attributes["dep_onboarding_settings_id"]

		if depOnboardingSettingsId == "" {
			dm, err := graphClient.DeviceManagement().Get(ctx, nil)
			if err != nil {
				return fmt.Errorf("error resolving dep onboarding settings id: %v", err)
			}
			depOnboardingSettingsId = dm.GetIntuneAccountId().String()
		}

		_, err := graphClient.
			DeviceManagement().
			DepOnboardingSettings().
			ByDepOnboardingSettingId(depOnboardingSettingsId).
			EnrollmentProfiles().
			ByEnrollmentProfileId(resourceId).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "ItemNotFound" {
				fmt.Printf("DEBUG: Resource %s successfully destroyed (404/NotFound)\n", resourceId)
				continue
			}
			return fmt.Errorf("error checking if macos dep enrollment profile %s was destroyed: %v", resourceId, err)
		}
		return fmt.Errorf("macos dep enrollment profile %s still exists", resourceId)
	}
	return nil
}
