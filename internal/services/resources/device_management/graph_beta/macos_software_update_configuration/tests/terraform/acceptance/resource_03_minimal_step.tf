# ==============================================================================
# Test 03: Minimal to Maximal in Steps - Step 1 (Minimal)
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_03_progression" {
  display_name                             = "acc-test-03-progression-${random_string.suffix.result}"
  update_schedule_type                     = "alwaysUpdate"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"
}
