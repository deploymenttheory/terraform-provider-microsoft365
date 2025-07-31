package graphBetaCloudPcOrganizationSettings_test

import (
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationSettingsResource_Create_Minimal(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal", "enable_mem_auto_enroll", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal", "enable_single_sign_on", "false"),
				),
			},
		},
	})
}

func TestAccOrganizationSettingsResource_Create_Maximal(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "enable_mem_auto_enroll", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "enable_single_sign_on", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "os_version", "windows11"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "user_account_type", "administrator"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.maximal", "windows_settings.0.language", "en-GB"),
				),
			},
		},
	})
}

func TestAccOrganizationSettingsResource_Update_MinimalToMaximal(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Start with minimal configuration
			{
				Config: testConfigMinimalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_mem_auto_enroll", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_single_sign_on", "false"),
				),
			},
			// Update to maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_mem_auto_enroll", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "enable_single_sign_on", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "os_version", "windows11"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.test", "user_account_type", "administrator"),
				),
			},
		},
	})
}

func TestAccOrganizationSettingsResource_Read(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal"),
				),
			},
		},
	})
}

func TestAccOrganizationSettingsResource_Import(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}