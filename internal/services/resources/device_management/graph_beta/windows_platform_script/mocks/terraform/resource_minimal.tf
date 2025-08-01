resource "microsoft365_graph_beta_device_management_windows_platform_script" "minimal" {
  display_name   = "Test Minimal Windows Platform Script - Unique"
  file_name      = "test_minimal.ps1"
  script_content = "# PowerShell Script\nWrite-Host 'Hello World'\nExit 0"
  run_as_account = "system"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}