resource "microsoft365_graph_beta_device_management_windows_enrollment_status_page" "minimal" {
  display_name                                     = "unit-test-windows-enrollment-status-page-minimal"
  description                                      = "Test description for minimal enrollment status page"
  show_installation_progress                       = true
  block_device_setup_retry_by_user                 = false
  allow_device_reset_on_install_failure            = true
  allow_log_collection_on_install_failure          = true
  custom_error_message                             = "Contact IT support for assistance"
  install_progress_timeout_in_minutes              = 120
  allow_device_use_on_install_failure              = false
  track_install_progress_for_autopilot_only        = true
  disable_user_status_tracking_after_first_user    = false
  
  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}