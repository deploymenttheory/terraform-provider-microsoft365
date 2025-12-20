
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group and Filter Dependencies (create all groups/filters needed for lifecycle)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_assignment_filter" "acc_test_filter_007_1" {
  display_name                      = "acc-test-filter-007-1-${random_string.test_suffix.result}"
  description                       = "Test filter 1 for windows remediation script assignment lifecycle"
  platform                          = "windows10AndLater"
  rule                              = "(device.osVersion -startsWith \"10.0\")"
  assignment_filter_management_type = "devices"


  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_1" {
  display_name     = "acc-test-group-007-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for windows remediation script assignment lifecycle"
  hard_delete      = true

}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_2" {
  display_name     = "acc-test-group-007-2-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-2-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 2 for windows remediation script assignment lifecycle"
  hard_delete      = true

}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_007_3" {
  display_name     = "acc-test-group-007-3-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-007-3-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 3 for windows remediation script assignment lifecycle"
  hard_delete      = true

}

# ==============================================================================
# Windows Remediation Script Resource - Assignment Lifecycle Step 1 (Minimal)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_007" {
  display_name               = "acc-test-windows-remediation-script-007-assignments-lifecycle-${random_string.test_suffix.result}"
  description                = "Scenario 7 Step 1: Starting with minimal assignments"
  publisher                  = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Simple detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Simple remediation script\nWrite-Host 'Remediation complete'\nexit 0"

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

