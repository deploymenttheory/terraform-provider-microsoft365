resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  display_name             = "Test Acceptance Windows Remediation Script - Updated"
  description              = "Updated description for acceptance testing"
  publisher                = "Terraform Provider Test Suite"
  run_as_account           = "user"
  run_as_32_bit            = true
  enforce_signature_check  = true
  detection_script_content = <<-EOT
    # Comprehensive detection script for acceptance testing
    $computerName = $env:COMPUTERNAME
    Write-Host "Computer: $computerName"
    
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
    # Comprehensive remediation script for acceptance testing
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
}