package graphBetaCloudPcUserSetting_test

import (
	"os"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserSettingResource_Create_Minimal(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.minimal", "display_name", "Test Minimal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.minimal", "local_admin_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.minimal", "reset_enabled", "false"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.minimal", "self_service_enabled", "false"),
				),
			},
		},
	})
}

func TestAccUserSettingResource_Create_Maximal(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "display_name", "Test Maximal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "local_admin_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "reset_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "self_service_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "restore_point_setting.0.frequency_in_hours", "24"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "restore_point_setting.0.frequency_type", "twentyFourHours"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.maximal", "restore_point_setting.0.user_restore_enabled", "true"),
				),
			},
		},
	})
}

func TestAccUserSettingResource_Update_MinimalToMaximal(t *testing.T) {
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
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "display_name", "Test Minimal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "local_admin_enabled", "false"),
				),
			},
			// Update to maximal configuration
			{
				Config: testConfigMaximalWithResourceName("test"),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.test"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "display_name", "Test Maximal User Setting"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "local_admin_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "reset_enabled", "true"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_windows_365_user_setting.test", "restore_point_setting.0.frequency_in_hours", "24"),
				),
			},
		},
	})
}

func TestAccUserSettingResource_Delete_Minimal(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.minimal"),
				),
			},
		},
	})
}

func TestAccUserSettingResource_Import(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless TF_ACC=1")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_windows_365_user_setting.minimal"),
				),
			},
			{
				ResourceName:      "microsoft365_graph_beta_windows_365_user_setting.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}