resource "microsoft365_graph_beta_windows_365_cloud_pc_role_definition" "test" {
  display_name = "Test Acceptance Role Definition"
  description  = ""

  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.CloudPC/CloudPCs/Read",
        "Microsoft.CloudPC/CloudPCs/Reboot"
      ]
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}