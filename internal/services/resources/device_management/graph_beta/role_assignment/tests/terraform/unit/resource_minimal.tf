resource "microsoft365_graph_beta_device_management_role_assignment" "minimal" {
  display_name       = "unit-test-role-assignment-minimal"
  description        = "Minimal role assignment for unit testing"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be" # Policy and Profile manager

  members = [
    "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
  ]

  scope_configuration {
    type = "AllLicensedUsers"
  }

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}