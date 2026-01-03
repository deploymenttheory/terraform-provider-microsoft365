resource "microsoft365_graph_beta_windows_365_cloud_pc_role_definition" "test" {
  display_name = "unit-test-cloud-pc-role-definition-invalid-prefix"
  description  = "Test description"

  role_permissions = [
    {
      allowed_resource_actions = [
        "InvalidPrefix_Permission"
      ]
    }
  ]
}

