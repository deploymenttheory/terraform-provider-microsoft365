resource "microsoft365_graph_beta_device_management_windows_platform_script" "minimal" {
  display_name  = "Test Minimal Windows Platform Script - Unique"
  file_name     = "test-script.ps1"
  script_content = "Write-Host 'Hello World'"
  run_as_account = "system"
  
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}