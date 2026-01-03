# Test Group 1 - Policy Managers
resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-role-assignment-policy-managers"
  description      = "Test group for policy managers used in role assignments"
  mail_nickname    = "acc-test-policy-mgrs"
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

# Test Group 2 - Device Administrators
resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "acc-test-role-assignment-device-admins"
  description      = "Test group for device administrators used in role assignments"
  mail_nickname    = "acc-test-device-admins"
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

resource "microsoft365_graph_beta_device_management_role_assignment" "resource_scopes" {
  display_name       = "acc-test-role-assignment-resource-scopes"
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