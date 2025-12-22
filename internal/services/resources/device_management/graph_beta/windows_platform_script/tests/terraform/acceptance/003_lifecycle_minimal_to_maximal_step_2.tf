
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_platform_script" "test_003" {
  display_name            = "acc-test-windows-platform-script-003-${random_string.test_suffix.result}"
  description             = "Maximal lifecycle test configuration"
  file_name               = "test-script-003.ps1"
  script_content          = "Write-Host 'Maximal script'"
  run_as_account          = "user"
  enforce_signature_check = true
  run_as_32_bit           = true
  role_scope_tag_ids      = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

