
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_003" {
  display_name     = "acc-test-windows-quality-update-policy-003-lifecycle-${random_string.test_suffix.result}"
  hotpatch_enabled = false

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

