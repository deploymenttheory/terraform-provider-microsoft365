
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "example" {
  display_name            = "Windows Security Remediation Script"
  description             = "Detects and remediates common security issues"
  publisher               = "IT Security Team"
  run_as_32_bit           = false
  enforce_signature_check = true
  role_scope_tag_ids      = ["0"]
  run_as_account          = "system" // Optional values:"system", "user"

  detection_script_content = <<-EOT
    # Detection script logic
    if (Test-Path "C:\Temp\issues.txt") {
      Write-Host "Issue detected"
      Exit 1
    } else {
      Write-Host "No issues found"
      Exit 0
    }
  EOT

  remediation_script_content = <<-EOT
    # Remediation script logic
    Remove-Item "C:\Temp\issues.txt" -Force
    Write-Host "Issue remediated"
    Exit 0
  EOT


  # Assignments are defined as a set
  assignments = [
    # Optional: Assignment targeting all devices with a daily schedule
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"

      daily_schedule = {
        interval = 1
        time     = "23:59:59"
        use_utc  = true
      }
    },
    # Optional: Assignment targeting all licensed users with an hourly schedule
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"

      hourly_schedule = {
        interval = 2
      }
    },
    # Optional: Assignment targeting a specific group with a run-once schedule
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"

      run_once_schedule = {
        date    = "2023-12-31"
        time    = "23:59:59"
        use_utc = true
      }
    },
    # Optional: Assignment targeting a specific group with a daily schedule
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"

      daily_schedule = {
        interval = 1
        time     = "23:59:59"
        use_utc  = true
      }
    },
    # Optional: Assignment targeting a specific group with an hourly schedule
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "exclude"

      hourly_schedule = {
        interval = 1
      }
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },

  ]

  timeouts = {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
