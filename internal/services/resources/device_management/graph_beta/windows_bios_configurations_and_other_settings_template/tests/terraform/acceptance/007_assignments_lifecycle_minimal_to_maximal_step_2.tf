
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
  description      = "Test group 1 for windows bios configuration assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_2" {
  display_name     = "acc-test-group-007-2-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-2-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for windows bios configuration assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_3" {
  display_name     = "acc-test-group-007-3-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-3-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 3 for windows bios configuration assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_4" {
  display_name     = "acc-test-group-007-4-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-4-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 4 for windows bios configuration assignments"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_5" {
  display_name     = "acc-test-group-007-5-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-5-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 5 for windows bios configuration assignments"
  hard_delete      = true
}

# ==============================================================================
# Assignment Filter Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_device_management_assignment_filter" "acc_test_filter_007_1" {
  display_name                      = "acc-test-filter-007-1-${random_string.test_suffix.result}"
  description                       = "Assignment filter 1 for windows bios configuration assignments"
  platform                          = "windows10AndLater"
  rule                              = "(device.osVersion -startsWith \"10.0\")"
  assignment_filter_management_type = "devices"
  role_scope_tags                   = ["0"]
}

resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "test_007" {
  display_name               = "acc-test-windows-bios-configuration-007-assignments-lifecycle-${random_string.test_suffix.result}"
  file_name                  = "test-bios-007.cctk"
  configuration_file_content = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtaW5pbWFsIHVuaXQgdGVzdCBwYWNrYWdlClNlY3VyZUJvb3Q9RW5hYmxlZApUcG09T24K"

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_007_1.id
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_filter_007_1.id
      filter_type = "include"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_007_2.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_007_3.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_007_4.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_007_5.id
    }
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}
