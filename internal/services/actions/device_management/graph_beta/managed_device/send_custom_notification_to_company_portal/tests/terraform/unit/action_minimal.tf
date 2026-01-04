action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "minimal" {
  config {
    managed_devices = [
      {
        device_id          = "00000000-0000-0000-0000-000000000001"
        notification_title = "Action Required"
        notification_body  = "Please update your device password to maintain access to company resources."
      }
    ]
  }
}

