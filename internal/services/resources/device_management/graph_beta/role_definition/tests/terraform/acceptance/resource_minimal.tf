resource "random_string" "role_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_role_definition" "test" {
  display_name = "acc-test-role-definition-minimal-${random_string.role_suffix.result}"
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
