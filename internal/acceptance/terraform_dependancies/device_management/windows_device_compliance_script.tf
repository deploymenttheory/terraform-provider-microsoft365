resource "random_string" "windows_device_compliance_script_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_device_compliance_script" "acc_test_windows_device_compliance_script" {
  display_name             = "acc-test-windows-device-compliance-script-${random_string.windows_device_compliance_script_suffix.result}"
  description              = "acc-test-windows-device-compliance-script"
  publisher                = "Acceptance Test Publisher"
  detection_script_content = "Get-Process | Select-Object -First 10"
  run_as_account           = "system"
  enforce_signature_check  = false
  run_as_32_bit            = false

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}