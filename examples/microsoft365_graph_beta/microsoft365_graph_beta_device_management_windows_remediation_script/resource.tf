resource "microsoft365_graph_beta_device_management_windows_remediation_script" "basic_example" {
  display_name            = "windows remediation script with no assignments"
  description             = "Simple script applied to no devices"
  publisher               = "Contoso IT"
  run_as_account          = "system"
  run_as_32_bit           = false
  enforce_signature_check = true
  role_scope_tag_ids      = [8, 9] # Optional

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
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "microsoft365_graph_beta_device_management_windows_remediation_script" "example_with_filters" {
  display_name            = "windows remediation script with no assignments"
  description             = "Simple script applied to no devices"
  publisher               = "Contoso IT"
  run_as_account          = "system"
  run_as_32_bit           = false
  enforce_signature_check = true
  role_scope_tag_ids      = [8, 9] # Optional

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
    all_devices             = true
    all_devices_filter_type = "include"
    all_devices_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"

    all_users             = true
    all_users_filter_type = "exclude"
    all_users_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}


resource "microsoft365_graph_beta_device_management_windows_remediation_script" "windows_remediation_script_with_scoping" {
  display_name            = "windows remediation script with assignment options"
  description             = "Simple script applied to scoped devices"
  publisher               = "Contoso IT"
  run_as_account          = "system"
  run_as_32_bit           = false
  enforce_signature_check = true
  role_scope_tag_ids      = [8, 9] # Optional

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
        group_id                   = "bae7a85a-8284-4f58-9873-a84bd4d22585"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        run_remediation_script     = true
        run_schedule = {
          schedule_type = "once"
          date          = "2025-05-01"
          time          = "14:30"
          use_utc       = true
        }
      },
      {
        group_id                   = "6117fcd2-2812-44b2-a0d7-3c57ca81c015"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        run_remediation_script     = true
        run_schedule = {
          schedule_type = "daily"
          interval      = "1"
          time          = "14:30"
          use_utc       = true
        }
      },
      {
        group_id                   = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
        include_groups_filter_type = "exclude"
        include_groups_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        run_remediation_script     = true
        run_schedule = {
          schedule_type = "hourly"
          interval      = "1"
        }
      }
    ]

    exclude_group_ids = [
      "b8c661c2-fa9a-4351-af86-adc1729c343f",
      "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f"
    ]
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}