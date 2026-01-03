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

resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "acc-test-role-assignment-maximal"
  description        = "Comprehensive role assignment for acceptance testing with all features"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be" # Policy and Profile manager

  members = [
    microsoft365_graph_beta_groups_group.acc_test_group_1.id,
    microsoft365_graph_beta_groups_group.acc_test_group_3.id
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