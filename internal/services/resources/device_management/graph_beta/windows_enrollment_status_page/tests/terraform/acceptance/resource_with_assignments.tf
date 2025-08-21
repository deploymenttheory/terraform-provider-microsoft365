resource "microsoft365_graph_beta_device_management_windows_enrollment_status_page" "with_assignments" {
  display_name                        = "acc-test-windows-enrollment-status-page-assignments-${random_string.test_suffix.result}"
  description                        = "Test enrollment status page with group assignments"
  show_installation_progress         = true
  block_device_setup_retry_by_user   = false
  allow_device_reset_on_install_failure = true
  allow_log_collection_on_install_failure = true
  allow_device_use_on_install_failure = false
  track_install_progress_for_autopilot_only = true
  disable_user_status_tracking_after_first_user = false
  custom_error_message              = "Contact IT support for device enrollment assistance"
  install_progress_timeout_in_minutes = 120
  
  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"  
    }
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