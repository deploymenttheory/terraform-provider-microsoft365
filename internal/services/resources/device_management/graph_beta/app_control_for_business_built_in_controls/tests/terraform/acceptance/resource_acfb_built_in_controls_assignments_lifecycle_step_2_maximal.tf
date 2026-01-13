resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_1" {
  display_name     = "acc-test-group-007-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for app control policy assignments lifecycle"
  hard_delete      = true
}

# ==============================================================================
# App Control Policy Resource - Step 2: Maximal Assignments
# ==============================================================================

resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "assignments_lifecycle" {
  name        = "acc-test-app-control-assignments-lifecycle-${random_string.test_suffix.result}"
  description = "Assignments lifecycle test - Step 2: Maximal assignments"

  enable_app_control = "audit"
  role_scope_tag_ids = ["0"]

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_007_1
  ]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_007_1.id
    },
    {
      type = "allDevicesAssignmentTarget"
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
