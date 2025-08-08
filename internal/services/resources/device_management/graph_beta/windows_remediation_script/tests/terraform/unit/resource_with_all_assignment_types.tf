resource "microsoft365_graph_beta_device_management_windows_remediation_script" "all_assignment_types" {
  display_name         = "Test All Assignment Types Windows Remediation Script - Unique"
  description          = "Windows remediation script with all assignment types for unit testing"
  publisher            = "Terraform Provider Test"
  run_as_account       = "system"
  run_as_32_bit        = false
  enforce_signature_check = false
  detection_script_content   = "# Comprehensive detection script with all assignment types\nWrite-Host 'Detection complete for all assignment types'\nexit 0"
  remediation_script_content = "# Comprehensive remediation script with all assignment types\nWrite-Host 'Remediation complete for all assignment types'\nexit 0"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
      daily_schedule = {
        interval = 1
        time     = "09:00:00"
        use_utc  = false
      }
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
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
      group_id = "33333333-3333-3333-3333-333333333333"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}