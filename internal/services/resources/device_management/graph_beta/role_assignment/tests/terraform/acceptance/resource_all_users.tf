# Test Group 3 - Application Managers
resource "microsoft365_graph_beta_groups_group" "acc_test_group_3" {
  display_name     = "acc-test-role-assignment-app-managers"
  description      = "Test group for application managers used in role assignments"
  mail_nickname    = "acc-test-app-mgrs"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Test Group 4 - Security Operators
resource "microsoft365_graph_beta_groups_group" "acc_test_group_4" {
  display_name     = "acc-test-role-assignment-sec-ops"
  description      = "Test group for security operators used in role assignments"
  mail_nickname    = "acc-test-sec-ops"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_role_assignment" "all_users" {
  display_name       = "acc-test-role-assignment-all-users"
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