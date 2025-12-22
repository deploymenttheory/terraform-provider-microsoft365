resource "microsoft365_graph_beta_device_management_windows_platform_script" "test_001" {
  display_name   = "unit-test-windows-platform-script-001-minimal"
  file_name      = "test-script-001.ps1"
  script_content = "Write-Host 'Hello World'"
  run_as_account = "system"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

