resource "microsoft365_graph_beta_device_management_windows_enrollment_status_page" "unknown_app_id" {
  display_name                                                        = "unit-test-windows-enrollment-status-page-unknown-app-id"
  description                                                         = "Test description for unknown app id enrollment status page"
  show_installation_progress                                          = true
  custom_error_message                                                = "Contact IT support for assistance"
  install_quality_updates                                             = true
  install_progress_timeout_in_minutes                                 = 120
  allow_log_collection_on_install_failure                             = true
  only_show_page_to_devices_provisioned_by_out_of_box_experience_oobe = true

  block_device_use_until_all_apps_and_profiles_are_installed = false
  allow_device_reset_on_install_failure                      = true
  allow_device_use_on_install_failure                        = true

  # Not present in the tenant
  selected_mobile_app_ids = [
    "99999999-9999-9999-9999-999999999999",
  ]

  only_fail_selected_blocking_apps_in_technician_phase = true

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
