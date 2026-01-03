resource "microsoft365_graph_beta_windows_365_cloud_pc_role_definition" "test" {
  display_name = "unit-test-cloud-pc-role-definition-error-handling"
  description  = "Test description"

  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.CloudPC/Invalid_Permission_Name"
      ]
    }
  ]
}

