resource "microsoft365_graph_beta_device_management_windows_device_compliance_script" "system_account" {
  display_name           = "Test System Account Windows Device Compliance Script - Unique"
  description            = "Test description for system account script"
  publisher              = "Test Publisher"
  detection_script_content = "Get-ComputerInfo"
  run_as_account         = "system"
  enforce_signature_check = false
  run_as_32_bit          = false
  
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}