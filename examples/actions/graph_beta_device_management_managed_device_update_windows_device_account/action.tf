# Example 1: Update device account for a single Microsoft Teams Room with Exchange Online
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_teams_room" {

  managed_devices {
    device_id                           = "12345678-1234-1234-1234-123456789abc"
    device_account_email                = "conference-room-01@company.com"
    password                            = var.teams_room_password # Use sensitive variable
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "outlook.office365.com"
    session_initiation_protocol_address = "sip:conference-room-01@company.com"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Update multiple Teams Rooms in bulk
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_multiple_teams_rooms" {

  managed_devices {
    device_id                           = "11111111-1111-1111-1111-111111111111"
    device_account_email                = "meeting-room-a@company.com"
    password                            = var.room_a_password
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "outlook.office365.com"
    session_initiation_protocol_address = "sip:meeting-room-a@company.com"
  }

  managed_devices {
    device_id                           = "22222222-2222-2222-2222-222222222222"
    device_account_email                = "meeting-room-b@company.com"
    password                            = var.room_b_password
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "outlook.office365.com"
    session_initiation_protocol_address = "sip:meeting-room-b@company.com"
  }

  managed_devices {
    device_id                           = "33333333-3333-3333-3333-333333333333"
    device_account_email                = "meeting-room-c@company.com"
    password                            = var.room_c_password
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "outlook.office365.com"
    session_initiation_protocol_address = "sip:meeting-room-c@company.com"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Update co-managed device (managed by both Intune and SCCM)
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_comanaged_device" {

  comanaged_devices {
    device_id                           = "55555555-5555-5555-5555-555555555555"
    device_account_email                = "hybrid-room@company.com"
    password                            = var.hybrid_room_password
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "mail.company.local"
    session_initiation_protocol_address = "sip:hybrid-room@company.com"
  }

  timeouts = {
    invoke = "5m"
  }
}