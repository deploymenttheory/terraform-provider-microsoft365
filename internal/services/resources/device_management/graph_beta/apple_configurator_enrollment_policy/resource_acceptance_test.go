package graphBetaAppleConfiguratorEnrollmentPolicy_test

import (
	"log"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAppleConfiguratorEnrollmentPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/apple_configurator_enrollment_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceType = graphBetaAppleConfiguratorEnrollmentPolicy.ResourceName

var testResource = graphBetaAppleConfiguratorEnrollmentPolicy.AppleConfiguratorEnrollmentPolicyTestResource{}

func TestAccResourceAppleConfiguratorEnrollmentPolicy_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
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
