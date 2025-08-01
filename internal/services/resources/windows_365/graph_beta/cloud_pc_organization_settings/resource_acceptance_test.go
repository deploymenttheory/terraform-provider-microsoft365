package graphBetaCloudPcOrganizationSettings_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudPcOrganizationSettingsResource_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccCloudPcOrganizationSettingsConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_mem_auto_enroll", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_single_sign_on", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "os_version", "windows11"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "user_account_type", "standardUser"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "windows_settings.language", "en-US"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudPcOrganizationSettingsResource_RequiredFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudPcOrganizationSettingsConfig_invalidOsVersion(),
				ExpectError: regexp.MustCompile("Attribute os_version value must be one of"),
			},
			{
				Config:      testAccCloudPcOrganizationSettingsConfig_invalidUserAccountType(),
				ExpectError: regexp.MustCompile("Attribute user_account_type value must be one of"),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("M365_TENANT_ID") == "" {
		t.Skip("M365_TENANT_ID must be set for acceptance tests")
	}
	if os.Getenv("M365_CLIENT_ID") == "" {
		t.Skip("M365_CLIENT_ID must be set for acceptance tests")
	}
	if os.Getenv("M365_CLIENT_SECRET") == "" {
		t.Skip("M365_CLIENT_SECRET must be set for acceptance tests")
	}
}

func testAccCloudPcOrganizationSettingsConfig_maximal() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings" "test" {
  enable_mem_auto_enroll = true
  enable_single_sign_on  = true
  os_version             = "windows11"
  user_account_type      = "standardUser"
  windows_settings = {
    language = "en-US"
  }
}
`
}

func testAccCloudPcOrganizationSettingsConfig_invalidOsVersion() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings" "test" {
  enable_mem_auto_enroll = true
  enable_single_sign_on  = true
  os_version             = "invalid"
  user_account_type      = "standardUser"
  windows_settings = {
    language = "en-US"
  }
}
`
}

func testAccCloudPcOrganizationSettingsConfig_invalidUserAccountType() string {
	return `
resource "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings" "test" {
  enable_mem_auto_enroll = true
  enable_single_sign_on  = true
  os_version             = "windows11"
  user_account_type      = "invalid"
  windows_settings = {
    language = "en-US"
  }
}
`
}