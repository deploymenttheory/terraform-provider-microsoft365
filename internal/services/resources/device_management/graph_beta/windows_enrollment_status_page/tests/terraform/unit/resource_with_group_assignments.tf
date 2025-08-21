resource "microsoft365_graph_beta_device_management_windows_enrollment_status_page" "full_with_group_assignments" {
  display_name                                  = "unit-test-windows-enrollment-status-page-with-group-assignments"
  description                                   = "Test description for full enrollment status page with group assignments"
  show_installation_progress                    = true
  block_device_setup_retry_by_user              = false
  allow_device_reset_on_install_failure         = true
  allow_log_collection_on_install_failure       = true
  custom_error_message                          = "Contact IT support for assistance"
  install_progress_timeout_in_minutes           = 120
  allow_device_use_on_install_failure           = false
  track_install_progress_for_autopilot_only     = true
  disable_user_status_tracking_after_first_user = false

  selected_mobile_app_ids = [
    "12345678-1234-1234-1234-123456789012",
    "87654321-4321-4321-4321-210987654321"
  ]

  role_scope_tag_ids = ["0", "1"]

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "12345678-1234-1234-1234-123456789012"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "12345678-1234-1234-1234-123456789013"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}