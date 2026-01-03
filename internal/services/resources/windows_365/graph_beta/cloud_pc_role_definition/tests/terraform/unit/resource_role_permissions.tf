resource "microsoft365_graph_beta_windows_365_cloud_pc_role_definition" "test" {
  display_name = "unit-test-cloud-pc-role-definition-role-permissions"
  description  = "Test description"

  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.CloudPC/CloudPCs/Read",
        "Microsoft.CloudPC/CloudPCs/Reboot",
        "Microsoft.CloudPC/DeviceImages/Read"
      ]
    }
  ]
}

