resource "microsoft365_graph_beta_device_management_role_assignment" "all_devices" {
  display_name       = "Test All Devices Role Assignment - ACC"
  description        = "Role assignment for all devices scope for acceptance testing"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be" # Policy and Profile manager

  members = [
    microsoft365_graph_beta_groups_group.acc_test_group_2.id,
    microsoft365_graph_beta_groups_group.acc_test_group_4.id
  ]

  scope_configuration {
    type = "AllDevices"
  }

  timeouts = {
    create = "300s"
    read   = "300s"
    update = "300s"
    delete = "300s"
  }
}