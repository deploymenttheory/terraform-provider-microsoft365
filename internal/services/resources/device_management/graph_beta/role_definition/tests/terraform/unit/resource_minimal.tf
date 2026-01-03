resource "microsoft365_graph_beta_device_management_role_definition" "minimal" {
  display_name = "unit-test-role-definition-minimal"
  description  = ""

  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.Intune_ManagedDevices_Read",
        "Microsoft.Intune_ManagedDevices_Update"
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