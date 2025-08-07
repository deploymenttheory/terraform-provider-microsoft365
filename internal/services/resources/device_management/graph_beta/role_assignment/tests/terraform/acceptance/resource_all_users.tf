resource "microsoft365_graph_beta_device_management_role_assignment" "all_users" {
  display_name       = "Test All Users Role Assignment - ACC"
  description        = "Role assignment for all licensed users scope for acceptance testing"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be" # Policy and Profile manager

  members = [
    microsoft365_graph_beta_groups_group.acc_test_group_3.id,
    microsoft365_graph_beta_groups_group.acc_test_group_4.id
  ]

  scope_configuration {
    type = "AllLicensedUsers"
  }

  timeouts = {
    create = "300s"
    read   = "300s"
    update = "300s"
    delete = "300s"
  }
}