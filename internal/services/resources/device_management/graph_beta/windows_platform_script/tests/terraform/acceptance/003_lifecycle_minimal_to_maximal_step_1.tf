
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_platform_script" "test_003" {
  display_name   = "acc-test-windows-platform-script-003-${random_string.test_suffix.result}"
  file_name      = "test-script-003.ps1"
  script_content = "Write-Host 'Hello World'"
  run_as_account = "system"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

