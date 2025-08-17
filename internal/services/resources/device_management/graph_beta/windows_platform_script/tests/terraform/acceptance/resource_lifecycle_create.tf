resource "random_uuid" "lifecycle" {
}

resource "microsoft365_graph_beta_device_management_windows_platform_script" "lifecycle" {
  display_name            = "Acceptance - Windows Platform Script"
  description             = "Acceptance test for Windows Platform Script lifecycle"
  file_name               = "acceptance-test-script.ps1"
  script_content          = "Write-Host 'Acceptance test script'"
  run_as_account          = "system"
  enforce_signature_check = false
  run_as_32_bit           = false

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