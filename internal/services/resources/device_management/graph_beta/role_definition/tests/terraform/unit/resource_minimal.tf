resource "microsoft365_graph_beta_device_management_role_definition" "minimal" {
  display_name                = "Test Minimal Role Definition - Unique"
  description                 = ""
  is_built_in_role_definition = false
  is_built_in                 = false

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.management/managedDevices/read",
        "microsoft.management/managedDevices/write"
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