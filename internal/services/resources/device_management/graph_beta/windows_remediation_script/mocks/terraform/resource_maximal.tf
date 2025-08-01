resource "microsoft365_graph_beta_device_management_windows_remediation_script" "maximal" {
  display_name             = "Test Maximal Windows Remediation Script - Unique"
  description              = "Maximal Windows remediation script for testing with all features"
  publisher                = "Terraform Provider Test Suite"
  run_as_account           = "user"
  run_as_32_bit            = true
  enforce_signature_check  = true
  detection_script_content = <<-EOT
    # Comprehensive detection script
    $computerName = $env:COMPUTERNAME
    $osVersion = (Get-WmiObject Win32_OperatingSystem).Version
    Write-Host "Computer: $computerName"
    Write-Host "OS Version: $osVersion"
    
    # Check for specific condition
    if (Test-Path "C:\temp\marker.txt") {
        Write-Host "Marker file found - issue detected"
        exit 1
    } else {
        Write-Host "No issues detected"
        exit 0
    }
  EOT

  remediation_script_content = <<-EOT
    # Comprehensive remediation script
    $logPath = "C:\temp\remediation.log"
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    
    # Create directory if it doesn't exist
    if (!(Test-Path "C:\temp")) {
        New-Item -ItemType Directory -Path "C:\temp" -Force
    }
    
    # Log the remediation action
    Add-Content -Path $logPath -Value "$timestamp - Remediation started"
    
    # Remove the marker file
    if (Test-Path "C:\temp\marker.txt") {
        Remove-Item "C:\temp\marker.txt" -Force
        Add-Content -Path $logPath -Value "$timestamp - Marker file removed"
    }
    
    Add-Content -Path $logPath -Value "$timestamp - Remediation completed"
    Write-Host "Remediation completed successfully"
    exit 0
  EOT

  role_scope_tag_ids = ["0", "1"]

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
      type = "allDevicesAssignmentTarget"
      run_once_schedule = {
        date    = "2024-12-31"
        time    = "23:59:00"
        use_utc = false
      }
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}