
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_001" {
  display_name     = "acc-test-windows-quality-update-policy-001-minimal-${random_string.test_suffix.result}"
  hotpatch_enabled = false

}

