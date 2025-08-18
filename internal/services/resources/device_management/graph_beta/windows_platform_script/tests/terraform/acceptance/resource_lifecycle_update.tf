resource "random_uuid" "lifecycle" {
}

resource "microsoft365_graph_beta_device_management_windows_platform_script" "lifecycle" {
  display_name            = "Acceptance - Windows Platform Script - Updated"
  description             = "Updated description for acceptance testing"
  file_name               = "acceptance-test-script-updated.ps1"
  script_content          = "Write-Host 'Updated acceptance test script'"
  run_as_account          = "user"
  enforce_signature_check = true
  run_as_32_bit           = true

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