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

resource "microsoft365_graph_beta_groups_group" "acc_test_group_3" {
  display_name     = "acc-test-rst-hd-${random_string.group_suffix.result}"
  description      = "Test group for helpdesk staff used in role scope tag assignments"
  mail_nickname    = "rst-hd-${random_string.group_suffix.result}"
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

resource "microsoft365_graph_beta_groups_group" "acc_test_group_4" {
  display_name     = "acc-test-rst-sec-${random_string.group_suffix.result}"
  description      = "Test group for security staff used in role scope tag assignments"
  mail_nickname    = "rst-sec-${random_string.group_suffix.result}"
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

resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "acc-test-role-scope-tag-maximal-${random_string.suffix.result}"
  description  = "acc-test-role-scope-tag-maximal-${random_string.suffix.result}"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_4.id
    }
  ]

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_3,
    microsoft365_graph_beta_groups_group.acc_test_group_4
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}