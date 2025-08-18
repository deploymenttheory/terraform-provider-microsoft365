resource "microsoft365_graph_beta_device_management_windows_device_compliance_script" "minimal" {
  display_name             = "Test Minimal Windows Device Compliance Script - Unique"
  publisher                = "Test Publisher"
  detection_script_content = "Get-Process"
  run_as_account           = "system"

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}