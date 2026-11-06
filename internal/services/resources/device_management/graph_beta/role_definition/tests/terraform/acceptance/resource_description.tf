resource "random_string" "role_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_role_definition" "description" {
  display_name = "acc-test-role-definition-description-${random_string.role_suffix.result}"
  description  = "This is a test role definition with description"

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
