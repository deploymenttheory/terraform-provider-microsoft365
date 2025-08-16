resource "random_uuid" "lifecycle" {
}

resource "microsoft365_graph_beta_device_management_windows_device_compliance_script" "lifecycle" {
  display_name           = "Acceptance - Windows Device Compliance Script"
  description            = "Acceptance test for Windows Device Compliance Script lifecycle"
  publisher              = "Acceptance Test Publisher"
  detection_script_content = "Get-Process | Select-Object -First 10"
  run_as_account         = "system"
  enforce_signature_check = false
  run_as_32_bit          = false
  
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }

  lifecycle {
    ignore_changes = [
      role_scope_tag_ids
    ]
  }
}