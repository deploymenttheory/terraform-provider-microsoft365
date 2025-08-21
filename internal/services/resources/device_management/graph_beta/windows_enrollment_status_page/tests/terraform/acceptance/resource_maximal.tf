resource "microsoft365_graph_beta_device_management_windows_enrollment_status_page" "maximal" {
  display_name                                  = "acc-test-windows-enrollment-status-page-maximal-${random_string.test_suffix.result}"
  description                                   = "Test description for maximal enrollment status page"
  show_installation_progress                    = true
  block_device_setup_retry_by_user              = false
  allow_device_reset_on_install_failure         = true
  allow_log_collection_on_install_failure       = true
  allow_device_use_on_install_failure           = false
  track_install_progress_for_autopilot_only     = true
  disable_user_status_tracking_after_first_user = false
  custom_error_message                          = "Contact IT support for assistance with device enrollment"
  install_progress_timeout_in_minutes           = 180

  role_scope_tag_ids = [
    microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id,
    microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "random_string" "test_suffix" {
  length  = 8
  upper   = false
  special = false
}