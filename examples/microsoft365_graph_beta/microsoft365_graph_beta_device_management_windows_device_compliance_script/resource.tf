resource "microsoft365_graph_beta_device_management_windows_device_compliance_script" "example" {
  display_name            = "Windows Security Compliance Check"
  description             = "Checks for critical security configurations on Windows devices"
  publisher               = "IT Security Team"
  run_as_32_bit           = false
  enforce_signature_check = true
  run_as_account          = "user" // Optional values: "system", "user"

  detection_script_content = <<-EOT
    # Detection script for Windows security compliance
    # Check if Windows Defender is enabled
    $defenderStatus = Get-MpComputerStatus
    if ($defenderStatus.AntivirusEnabled -eq $false) {
      Write-Host "Windows Defender is disabled - Non-compliant"
      Exit 1
    }
    
    # Check if Windows Updates are configured properly
    $auSettings = Get-ItemProperty -Path "HKLM:\SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate\AU" -ErrorAction SilentlyContinue
    if ($auSettings.NoAutoUpdate -eq 1) {
      Write-Host "Automatic updates are disabled - Non-compliant"
      Exit 1
    }
    
    # Check if firewall is enabled
    $firewallProfiles = Get-NetFirewallProfile
    $disabledProfiles = $firewallProfiles | Where-Object { $_.Enabled -eq $false }
    if ($disabledProfiles.Count -gt 0) {
      Write-Host "One or more firewall profiles are disabled - Non-compliant"
      Exit 1
    }
    
    Write-Host "All security checks passed - Compliant"
    Exit 0
  EOT

  timeouts = {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}