resource "microsoft365_graph_beta_device_management_role_definition" "description" {
  display_name                = "Test Description Role Definition"
  description                 = "This is a test role definition with description"
  is_built_in_role_definition = false
  is_built_in                 = false

  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.Intune_ManagedDevices_Read"
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