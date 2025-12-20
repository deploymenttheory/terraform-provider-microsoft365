resource "microsoft365_graph_beta_device_management_windows_remediation_script" "test_004" {
  display_name             = "unit-test-windows-remediation-script-004-downgrade"
  description              = "Downgrade Step 1: Starting with maximal configuration"
  publisher                = "Terraform Provider Test Suite"
  run_as_account           = "user"
  run_as_32_bit            = true
  enforce_signature_check  = true
  detection_script_content = <<-EOT
    # Comprehensive detection script
    $computerName = $env:COMPUTERNAME
    Write-Host "Computer: $computerName"
    exit 0
  EOT

  remediation_script_content = <<-EOT
    # Comprehensive remediation script
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    Write-Host "$timestamp - Remediation completed"
    exit 0
  EOT

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

