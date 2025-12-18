
resource "random_string" "test_012_invalid_mode" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "test_012_invalid_mode" {
  display_name                            = "unit-test-windows-update-ring-012-invalid-mode-${random_string.test_012_invalid_mode.result}"
  description                             = "Test invalid automatic_update_mode"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 0
  feature_updates_deferral_period_in_days = 0
  allow_windows11_upgrade                 = true
  skip_checks_before_restart              = false
  automatic_update_mode                   = "invalidMode"
  feature_updates_rollback_window_in_days = 10

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

