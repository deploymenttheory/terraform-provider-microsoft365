package graphBetaDeviceEnrollmentNotification_test

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

// Android Platform Tests
func TestAccDeviceEnrollmentNotificationResource_AndroidEmailMinimal(t *testing.T) {
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
				Config: testAccAndroidEnrollmentNotificationsConfig_androidEmailMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_android", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_android", "display_name", "email minimal android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_android", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_android", "notification_templates.*", "email"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_android", "branding_options.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_android", "branding_options.*", "none"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_android", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"branding_options", "platform_type"}},
		},
	})
}

func TestAccDeviceEnrollmentNotificationResource_AndroidEmailMaximal(t *testing.T) {
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
				Config: testAccAndroidEnrollmentNotificationsConfig_androidEmailMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_android", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_android", "display_name", "email maximal android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_android", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_android", "notification_templates.*", "email"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_android", "branding_options.#", "5"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_android", "branding_options.*", "includeCompanyLogo"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_android", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"branding_options", "platform_type"}},
		},
	})
}

func TestAccDeviceEnrollmentNotificationResource_AndroidPushMaximal(t *testing.T) {
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
				Config: testAccAndroidEnrollmentNotificationsConfig_androidPushMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_android", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_android", "display_name", "push maximal android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_android", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_android", "notification_templates.*", "push"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_android", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"branding_options", "platform_type"}},
		},
	})
}

func TestAccDeviceEnrollmentNotificationResource_AndroidAllMaximal(t *testing.T) {
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
				Config: testAccAndroidEnrollmentNotificationsConfig_androidAllMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_enrollment_notification.all_android", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_android", "display_name", "Complete Test - all android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_android", "notification_templates.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_android", "notification_templates.*", "email"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_android", "notification_templates.*", "push"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_enrollment_notification.all_android", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"branding_options", "platform_type"}},
		},
	})
}

// AndroidForWork Platform Tests
func TestAccDeviceEnrollmentNotificationResource_AndroidForWorkEmailMinimal(t *testing.T) {
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
				Config: testAccAndroidEnrollmentNotificationsConfig_androidForWorkEmailMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_androidforwork", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_androidforwork", "display_name", "email minimal androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_androidforwork", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_androidforwork", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_androidforwork", "notification_templates.*", "email"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_enrollment_notification.email_minimal_androidforwork", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"branding_options", "platform_type"}},
		},
	})
}

func TestAccDeviceEnrollmentNotificationResource_AndroidForWorkEmailMaximal(t *testing.T) {
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
				Config: testAccAndroidEnrollmentNotificationsConfig_androidForWorkEmailMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_androidforwork", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_androidforwork", "display_name", "email maximal androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_androidforwork", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_androidforwork", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_androidforwork", "notification_templates.*", "email"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_androidforwork", "branding_options.#", "5"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_enrollment_notification.email_maximal_androidforwork", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"branding_options", "platform_type"}},
		},
	})
}

func TestAccDeviceEnrollmentNotificationResource_AndroidForWorkPushMaximal(t *testing.T) {
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
				Config: testAccAndroidEnrollmentNotificationsConfig_androidForWorkPushMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_androidForWork", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_androidForWork", "display_name", "push maximal android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_androidForWork", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_androidForWork", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_androidForWork", "notification_templates.*", "push"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_enrollment_notification.push_maximal_androidForWork", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"branding_options", "platform_type"}},
		},
	})
}

func TestAccDeviceEnrollmentNotificationResource_AndroidForWorkAllMaximal(t *testing.T) {
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
				Config: testAccAndroidEnrollmentNotificationsConfig_androidForWorkAllMaximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_device_enrollment_notification.all_androidforwork", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_androidforwork", "display_name", "Complete Test - all androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_androidforwork", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_androidforwork", "notification_templates.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_androidforwork", "notification_templates.*", "email"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_enrollment_notification.all_androidforwork", "notification_templates.*", "push"),
				),
			},
			{ResourceName: "microsoft365_graph_beta_device_management_device_enrollment_notification.all_androidforwork", ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"branding_options", "platform_type"}},
		},
	})
}

// Android Platform Configuration Functions
func testAccAndroidEnrollmentNotificationsConfig_androidEmailMinimal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_android_email_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load android email minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccAndroidEnrollmentNotificationsConfig_androidEmailMaximal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_android_email_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load android email maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccAndroidEnrollmentNotificationsConfig_androidPushMaximal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_android_push_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load android push maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccAndroidEnrollmentNotificationsConfig_androidAllMaximal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_android_all_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load android all maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

// AndroidForWork Platform Configuration Functions
func testAccAndroidEnrollmentNotificationsConfig_androidForWorkEmailMinimal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_androidForWork_email_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load androidForWork email minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccAndroidEnrollmentNotificationsConfig_androidForWorkEmailMaximal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_androidForWork_email_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load androidForWork email maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccAndroidEnrollmentNotificationsConfig_androidForWorkPushMaximal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_androidForWork_push_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load androidForWork push maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccAndroidEnrollmentNotificationsConfig_androidForWorkAllMaximal() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_androidForWork_all_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load androidForWork all maximal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccCheckAndroidEnrollmentNotificationsDestroy(s *terraform.State) error {

	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return fmt.Errorf("error creating Graph client for CheckDestroy: %v", err)
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "microsoft365_graph_beta_device_management_device_enrollment_notification" {
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
			return fmt.Errorf("error checking if android enrollment notifications %s was destroyed: %v", rs.Primary.ID, err)
		}
		return fmt.Errorf("android enrollment notifications %s still exists", rs.Primary.ID)
	}
	return nil
}
