resource "random_uuid" "lifecycle" {
}

resource "microsoft365_graph_beta_device_management_windows_device_compliance_script" "lifecycle" {
  display_name             = "Acceptance - Windows Device Compliance Script - Updated"
  description              = "Updated description for acceptance testing"
  publisher                = "Updated Test Publisher"
  detection_script_content = "Get-Service | Where-Object {$_.Status -eq 'Running'}"
  run_as_account           = "user"
  enforce_signature_check  = true
  run_as_32_bit            = true

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