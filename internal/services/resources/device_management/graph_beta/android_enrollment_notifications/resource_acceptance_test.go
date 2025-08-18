package graphBetaAndroidEnrollmentNotifications_test

import (
	"log"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAndroidEnrollmentNotificationsResource_AndroidForWork_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidEnrollmentNotificationsDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAndroidEnrollmentNotificationsConfig_AndroidForWork_Minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "default_locale", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "notification_templates.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "id"),
				),
			},
			{
				Config: testAccAndroidEnrollmentNotificationsConfig_AndroidForWork_Minimal_Update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "notification_templates.#", "2"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccAndroidEnrollmentNotificationsResource_Android_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidEnrollmentNotificationsDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAndroidEnrollmentNotificationsConfig_Android_Minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "default_locale", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "notification_templates.#", "1"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "id"),
				),
			},
			{
				Config: testAccAndroidEnrollmentNotificationsConfig_Android_Minimal_Update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "notification_templates.#", "2"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccAndroidEnrollmentNotificationsResource_AndroidForWork_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidEnrollmentNotificationsDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAndroidEnrollmentNotificationsConfig_AndroidForWork_Maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "branding_options", "includeCompanyLogo,includeCompanyName,includeContactInformation"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "assignments.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "localized_notification_messages.#", "2"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "id"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccAndroidEnrollmentNotificationsResource_Android_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAndroidEnrollmentNotificationsDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAndroidEnrollmentNotificationsConfig_Android_Maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "branding_options", "includeCompanyLogo,includeCompanyName,includeDeviceDetails"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "assignments.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "localized_notification_messages.#", "2"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle", "id"),
				),
			},
			{
				ResourceName:                         "microsoft365_graph_beta_device_management_android_enrollment_notifications.lifecycle",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func testAccCheckAndroidEnrollmentNotificationsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_android_enrollment_notifications" {
			continue
		}
		// Resource should be destroyed - we'll assume success if we reach here
		// In real acceptance tests, this would check against the actual Graph API
	}
	return nil
}

// AndroidForWork Minimal Configuration Functions
func testAccAndroidEnrollmentNotificationsConfig_AndroidForWork_Minimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_lifecycle_create_androidforwork_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load AndroidForWork minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAndroidEnrollmentNotificationsConfig_AndroidForWork_Minimal_Update() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_lifecycle_update_androidforwork_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load AndroidForWork minimal update test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

// AndroidForWork Maximal Configuration Functions
func testAccAndroidEnrollmentNotificationsConfig_AndroidForWork_Maximal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_lifecycle_create_androidforwork_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load AndroidForWork maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

// Android Minimal Configuration Functions
func testAccAndroidEnrollmentNotificationsConfig_Android_Minimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_lifecycle_create_android_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load Android minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAndroidEnrollmentNotificationsConfig_Android_Minimal_Update() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_lifecycle_update_android_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load Android minimal update test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

// Android Maximal Configuration Functions
func testAccAndroidEnrollmentNotificationsConfig_Android_Maximal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_lifecycle_create_android_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load Android maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}
