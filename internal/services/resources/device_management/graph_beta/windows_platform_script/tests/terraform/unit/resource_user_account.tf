resource "microsoft365_graph_beta_device_management_windows_platform_script" "user_account" {
  display_name           = "Test User Account Windows Platform Script - Unique"
  description            = "Test description for user account script"
  file_name              = "test-user-script.ps1"
  script_content         = "Write-Host 'Running as User'"
  run_as_account         = "user"
  enforce_signature_check = true
  run_as_32_bit          = true
  
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}