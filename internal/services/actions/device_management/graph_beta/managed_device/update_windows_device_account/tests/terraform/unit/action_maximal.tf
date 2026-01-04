action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "maximal" {
  config {
    managed_devices = [
      {
        device_id                           = "00000000-0000-0000-0000-000000000001"
        device_account_email                = "conference-room-01@company.com"
        password                            = "SecurePassword123!"
        password_rotation_enabled           = true
        calendar_sync_enabled               = true
        exchange_server                     = "outlook.office365.com"
        session_initiation_protocol_address = "sip:conference-room-01@company.com"
      },
      {
        device_id                           = "00000000-0000-0000-0000-000000000002"
        device_account_email                = "surfacehub-lobby@company.com"
        password                            = "AnotherSecurePass456!"
        password_rotation_enabled           = true
        calendar_sync_enabled               = true
        exchange_server                     = "outlook.office365.com"
        session_initiation_protocol_address = "sip:surfacehub-lobby@company.com"
      }
    ]
    comanaged_devices = [
      {
        device_id                           = "00000000-0000-0000-0000-000000000003"
        device_account_email                = "teams-room-03@company.com"
        password                            = "YetAnotherPass789!"
        password_rotation_enabled           = true
        calendar_sync_enabled               = true
        exchange_server                     = "mail.company.com"
        session_initiation_protocol_address = "sip:teams-room-03@company.com"
      }
    ]
    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

