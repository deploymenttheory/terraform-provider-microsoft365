resource "microsoft365_graph_beta_device_and_app_management_windows_remediation_script" "basic_example" {
  display_name            = "windows remediation script with no assignments"
  description             = "Simple script applied to no devices"
  publisher               = "Contoso IT"
  run_as_account          = "system"
  run_as_32_bit           = false
  enforce_signature_check = true

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

  assignment {
    all_devices = false
    all_users   = false

    exclude_group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}


resource "microsoft365_graph_beta_device_and_app_management_windows_remediation_script" "windows_remediation_script_with_scoping" {
  display_name            = "windows remediation script with assignment options"
  description             = "Simple script applied to scoped devices"
  publisher               = "Contoso IT"
  run_as_account          = "system"
  run_as_32_bit           = false
  enforce_signature_check = true

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

  assignment {
    all_devices = false
    all_users   = false

    include_groups = [
      {
        group_id                   = "11111111-2222-3333-4444-555555555555"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "11111111-2222-3333-4444-555555555555"
        run_remediation_script     = true
        run_schedule = {
          schedule_type = "once"
          date          = "2025-05-01"
          time          = "14:30"
          use_utc       = true
        }
      },
      {
        group_id                   = "11111111-2222-3333-4444-555555555555"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "11111111-2222-3333-4444-555555555555"
        run_remediation_script     = true
        run_schedule = {
          schedule_type = "daily"
          interval      = "1"
          time          = "14:30"
          use_utc       = true
        }
      },
      {
        group_id                   = "11111111-2222-3333-4444-555555555555"
        include_groups_filter_type = "exclude"
        include_groups_filter_id   = "11111111-2222-3333-4444-555555555555"
        run_remediation_script     = true
        run_schedule = {
          schedule_type = "hourly"
          interval      = "1"
        }
      }
    ]

    exclude_group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}
