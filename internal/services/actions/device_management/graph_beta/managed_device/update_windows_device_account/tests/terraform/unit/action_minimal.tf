action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "minimal" {
  config {
    managed_devices = [
      {
        device_id                 = "00000000-0000-0000-0000-000000000001"
        device_account_email      = "conference-room-01@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]
  }
}

