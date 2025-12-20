resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_008" {
  display_name               = "unit-test-windows-remediation-script-008-assignments-downgrade"
  description                = "Scenario 8 Step 1: Starting with maximal assignments"
  publisher                  = "Terraform Provider Test"
  run_as_account             = "system"
  detection_script_content   = "# Simple detection script\nWrite-Host 'Detection complete'\nexit 0"
  remediation_script_content = "# Simple remediation script\nWrite-Host 'Remediation complete'\nexit 0"

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "44444444-4444-4444-4444-444444444444"
      filter_id   = "55555555-5555-5555-5555-555555555555"
      filter_type = "include"
      daily_schedule = {
        interval = 1
        time     = "09:00:00"
        use_utc  = true
      }
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_id   = "66666666-6666-6666-6666-666666666666"
      filter_type = "exclude"
      hourly_schedule = {
        interval = 4
      }
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type = "allDevicesAssignmentTarget"
      run_once_schedule = {
        date    = "2024-12-31"
        time    = "23:59:00"
        use_utc = false
      }
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "77777777-7777-7777-7777-777777777777"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

