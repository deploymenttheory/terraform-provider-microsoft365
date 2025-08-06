resource "microsoft365_graph_beta_device_management_role_assignment" "resource_scopes" {
  display_name       = "Test Resource Scopes Role Assignment - ACC"
  description        = "Role assignment with specific resource scopes for acceptance testing"
  role_definition_id = "c1d9fcbb-cba5-40b0-bf6b-527006585f4b" # Application Manager
  
  members = [
    microsoft365_graph_beta_groups_group.acc_test_group_1.id,
    microsoft365_graph_beta_groups_group.acc_test_group_2.id
  ]
  
  scope_configuration {
    type = "ResourceScopes"
    resource_scopes = [
      microsoft365_graph_beta_groups_group.acc_test_group_3.id,
    microsoft365_graph_beta_groups_group.acc_test_group_4.id
    ]
  }

  timeouts = {
    create = "300s"
    read   = "300s"
    update = "300s"
    delete = "300s"
  }
}