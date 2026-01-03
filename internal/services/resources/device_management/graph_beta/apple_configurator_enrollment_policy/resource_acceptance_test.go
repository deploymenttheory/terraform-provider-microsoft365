package graphBetaAppleConfiguratorEnrollmentPolicy_test

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

func TestAccAppleConfiguratorEnrollmentPolicyResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAppleConfiguratorEnrollmentPolicyDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAppleConfiguratorEnrollmentPolicyConfig_enrollWithoutUserAffinity(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.enroll_without_user_affinity", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.enroll_without_user_affinity", "display_name", "acc-test-apple-configurator-enrollment-policy-enroll-without-user-affinity"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.enroll_without_user_affinity", "description", "apple configurator enrollment policy without user affinity"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.enroll_without_user_affinity", "requires_user_authentication", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.enroll_without_user_affinity", "enable_authentication_via_company_portal", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.enroll_without_user_affinity", "require_company_portal_on_setup_assistant_enrolled_devices", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.enroll_without_user_affinity", "configuration_endpoint_url"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.enroll_without_user_affinity", "dep_onboarding_settings_id"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.enroll_without_user_affinity", ImportState: true, ImportStateVerify: true},
			{
				Config: testAccAppleConfiguratorEnrollmentPolicyConfig_userAffinityWithCompanyPortal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_company_portal", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_company_portal", "display_name", "acc-test-apple-configurator-enrollment-policy-user-affinity-with-company-portal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_company_portal", "description", "apple configurator enrollment policy with user affinity via company portal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_company_portal", "requires_user_authentication", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_company_portal", "enable_authentication_via_company_portal", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_company_portal", "require_company_portal_on_setup_assistant_enrolled_devices", "false"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_company_portal", "configuration_endpoint_url"),
				),
			},
			{
				Config: testAccAppleConfiguratorEnrollmentPolicyConfig_userAffinityWithSetupAssistant(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_setup_assistant", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_setup_assistant", "display_name", "acc-test-apple-configurator-enrollment-policy-user-affinity-with-setup-assistant"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_setup_assistant", "description", "apple configurator enrollment policy with user affinity via setup assistant"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_setup_assistant", "requires_user_authentication", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_setup_assistant", "enable_authentication_via_company_portal", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_setup_assistant", "require_company_portal_on_setup_assistant_enrolled_devices", "true"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy.user_affinity_with_setup_assistant", "configuration_endpoint_url"),
				),
			},
		},
	})
}

func testAccAppleConfiguratorEnrollmentPolicyConfig_enrollWithoutUserAffinity() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_enroll_without_user_affinity.tf")
	if err != nil {
		log.Fatalf("Failed to load enroll without user affinity test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAppleConfiguratorEnrollmentPolicyConfig_userAffinityWithCompanyPortal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_user_affinity_with_company_portal.tf")
	if err != nil {
		log.Fatalf("Failed to load user affinity with company portal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAppleConfiguratorEnrollmentPolicyConfig_userAffinityWithSetupAssistant() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_user_affinity_with_setup_assistant.tf")
	if err != nil {
		log.Fatalf("Failed to load user affinity with setup assistant test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccCheckAppleConfiguratorEnrollmentPolicyDestroy(s *terraform.State) error {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy" {
			continue
		}

		// Parse the resource ID to get both depOnboardingSettingsId and enrollmentProfileId
		// The ID format is: depOnboardingSettingsId_enrollmentProfileId
		resourceId := rs.Primary.ID
		depOnboardingSettingsId := rs.Primary.Attributes["dep_onboarding_settings_id"]

		if depOnboardingSettingsId == "" {
			// Try to resolve from device management if not in state
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
			return fmt.Errorf("error checking if apple configurator enrollment policy %s was destroyed: %v", resourceId, err)
		}
		return fmt.Errorf("apple configurator enrollment policy %s still exists", resourceId)
	}
	return nil
}
