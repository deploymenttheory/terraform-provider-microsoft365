resource "microsoft365_graph_beta_device_management_windows_remediation_script" "assignments" {
  display_name         = "Test All Assignment Types Windows Remediation Script"
  description          = "Windows remediation script with comprehensive assignments for acceptance testing"
  publisher            = "Terraform Provider Test"
  run_as_account       = "system"
  run_as_32_bit        = false
  enforce_signature_check = false
  detection_script_content   = "# Comprehensive detection script with all assignment types\nWrite-Host 'Detection complete for all assignment types'\nexit 0"
  remediation_script_content = "# Comprehensive remediation script with all assignment types\nWrite-Host 'Remediation complete for all assignment types'\nexit 0"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
      daily_schedule = {
        interval = 1
        time     = "09:00:00"
        use_utc  = false
      }
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
      daily_schedule = {
        interval = 1
        time     = "15:00:00"
        use_utc  = false
      }
    },
    {
      type = "allLicensedUsersAssignmentTarget"
      daily_schedule = {
        interval = 1
        time     = "12:00:00"
        use_utc  = false
      }
    },
    {
      type = "allDevicesAssignmentTarget"
      daily_schedule = {
        interval = 1
        time     = "02:00:00"
        use_utc  = false
      }
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}