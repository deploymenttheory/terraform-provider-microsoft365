resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_004_1" {
  display_name     = "acc-test-group-004-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-004-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for maximal assignment"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_004_2" {
  display_name     = "acc-test-group-004-2-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-004-2-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for maximal assignment exclusion"
  hard_delete      = true
}

# ==============================================================================
# Group Policy Configuration Resource - Maximal Assignment
# ==============================================================================

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "maximal_assignment" {
  display_name       = "acc-test-004-maximal-assignment-${random_string.test_suffix.result}"
  description        = "acc-test-004-maximal-assignment"
  role_scope_tag_ids = ["0"]

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_004_1,
    microsoft365_graph_beta_groups_group.acc_test_group_004_2
  ]

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_004_1.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_004_2.id
    }
  ]
}
