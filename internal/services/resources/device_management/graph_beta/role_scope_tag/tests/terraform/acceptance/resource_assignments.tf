resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
  numeric = true
}

resource "random_string" "group_suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
  numeric = true
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-rst-it-${random_string.group_suffix.result}"
  description      = "Test group for IT support staff used in role scope tag assignments"
  mail_nickname    = "rst-it-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "acc-test-rst-dm-${random_string.group_suffix.result}"
  description      = "Test group for device management staff used in role scope tag assignments"
  mail_nickname    = "rst-dm-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "assignments" {
  display_name = "acc-test-role-scope-tag-assignments-${random_string.suffix.result}"
  description  = "acc-test-role-scope-tag-assignments-${random_string.suffix.result}"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    }
  ]

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_1,
    microsoft365_graph_beta_groups_group.acc_test_group_2
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}