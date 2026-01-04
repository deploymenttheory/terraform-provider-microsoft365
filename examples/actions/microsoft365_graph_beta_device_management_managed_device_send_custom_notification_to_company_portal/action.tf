# Example 1: Send custom notification to a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "send_single" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        notification_title = "Action Required"
        notification_body  = "Please update your device password to maintain access to company resources."
      }
    ]
  }
}

# Example 2: Send custom notifications to multiple devices with different messages
action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "send_multiple" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        notification_title = "Action Required: Update Password"
        notification_body  = "Your device password will expire in 3 days. Please update it to maintain access to company resources."
      },
      {
        device_id          = "87654321-4321-4321-4321-ba9876543210"
        notification_title = "Compliance Alert"
        notification_body  = "Your device is not compliant with corporate security policies. Please contact IT support."
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "send_maximal" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        notification_title = "Action Required: Update Password"
        notification_body  = "Your device password will expire in 3 days. Please update it to maintain access to company resources."
      },
      {
        device_id          = "87654321-4321-4321-4321-ba9876543210"
        notification_title = "Compliance Alert"
        notification_body  = "Your device is not compliant with corporate security policies. Please contact IT support."
      }
    ]

    comanaged_devices = [
      {
        device_id          = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
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

# Example 4: Send notification to non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "noncompliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "notify_noncompliant" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.noncompliant_devices.items : {
        device_id          = device.id
        notification_title = "Compliance Action Required"
        notification_body  = "Your device is not compliant. Please ensure all required policies are applied."
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Send notification for upcoming maintenance
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "maintenance_notification" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.windows_devices.items : {
        device_id          = device.id
        notification_title = "Scheduled Maintenance"
        notification_body  = "System maintenance is scheduled for this weekend. Please ensure your work is saved."
      }
    ]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 6: Send notification to co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "send_comanaged" {
  config {
    comanaged_devices = [
      {
        device_id          = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        notification_title = "Important Update"
        notification_body  = "Please restart your device to apply critical security updates."
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Output examples
output "notifications_sent_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal.send_multiple.config.managed_devices)
  description = "Number of notifications sent to managed devices"
}

