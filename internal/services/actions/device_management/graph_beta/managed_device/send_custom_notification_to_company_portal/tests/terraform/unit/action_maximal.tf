action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "maximal" {
  config {
    managed_devices = [
      {
        device_id          = "00000000-0000-0000-0000-000000000001"
        notification_title = "Action Required: Update Password"
        notification_body  = "Your device password will expire in 3 days. Please update it to maintain access to company resources."
      },
      {
        device_id          = "00000000-0000-0000-0000-000000000002"
        notification_title = "Compliance Alert"
        notification_body  = "Your device is not compliant with corporate security policies. Please contact IT support."
      }
    ]
    comanaged_devices = [
      {
        device_id          = "00000000-0000-0000-0000-000000000003"
        notification_title = "Maintenance Window"
        notification_body  = "A scheduled maintenance will occur tonight from 10 PM to 2 AM. Please save your work."
      }
    ]
    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

