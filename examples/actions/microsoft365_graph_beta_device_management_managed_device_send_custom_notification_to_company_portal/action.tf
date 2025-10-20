# REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-sendcustomnotificationtocompanyportal?view=graph-rest-beta

# Data source to find devices for targeted messaging
data "microsoft365_graph_beta_device_management_managed_device" "all_devices" {}

# Example 1: Send compliance reminder to specific non-compliant devices
# Use this to remind users about compliance issues
action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "compliance_reminder" {
  managed_devices {
    device_id          = "12345678-1234-1234-1234-123456789abc"
    notification_title = "Action Required: Device Compliance"
    notification_body  = "Your device is not compliant with company security policies. Please update your device settings to regain access to company resources."
  }

  managed_devices {
    device_id          = "87654321-4321-4321-4321-ba9876543210"
    notification_title = "Action Required: Device Compliance"
    notification_body  = "Your device requires immediate attention. Please update your antivirus software to maintain compliance."
  }
}

# Example 2: Send password expiration warning with different messages per device
# Use this to send personalized password expiration reminders
action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "password_expiration_warning" {
  managed_devices {
    device_id          = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
    notification_title = "Password Expires in 3 Days"
    notification_body  = "Your device password will expire on October 22, 2025. Please update it now to avoid losing access to company resources."
  }

  managed_devices {
    device_id          = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
    notification_title = "Password Expires Tomorrow"
    notification_body  = "URGENT: Your device password expires tomorrow! Update it immediately to prevent account lockout."
  }
}

# Example 3: Send security alert to all Windows devices
# Use this for critical security communications
data "microsoft365_graph_beta_device_management_managed_device" "windows_devices" {
  filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "security_alert" {
  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.windows_devices.managed_devices
    content {
      device_id          = managed_devices.value.id
      notification_title = "Security Alert: Update Required"
      notification_body  = "A critical security update is available for your Windows device. Please install updates through the Company Portal app within 24 hours."
    }
  }

  timeouts = {
    invoke = "20m"
  }
}

# Example 4: Send maintenance window notification to specific departments
# Use this to communicate scheduled maintenance
data "microsoft365_graph_beta_device_management_managed_device" "finance_devices" {
  filter = "startswith(deviceName, 'FIN-')"
}

action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "maintenance_notice" {
  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.finance_devices.managed_devices
    content {
      device_id          = managed_devices.value.id
      notification_title = "Scheduled Maintenance: October 25"
      notification_body  = "Your device will undergo maintenance on October 25 from 10 PM to 2 AM. Please save your work and leave your device powered on."
    }
  }
}

# Example 5: Send app update notification to iOS devices
# Use this to notify users about required app updates
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices" {
  filter = "operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "app_update_required" {
  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.ios_devices.managed_devices
    content {
      device_id          = managed_devices.value.id
      notification_title = "App Update Available"
      notification_body  = "Important security updates are available for your corporate apps. Open the Company Portal to update your apps now."
    }
  }
}

# Example 6: Send policy change notification to Android devices
# Use this to communicate new policies to Android users
data "microsoft365_graph_beta_device_management_managed_device" "android_devices" {
  filter = "operatingSystem eq 'Android'"
}

action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "policy_update" {
  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.android_devices.managed_devices
    content {
      device_id          = managed_devices.value.id
      notification_title = "New Device Policy Effective October 30"
      notification_body  = "New security policies will be enforced on your Android device starting October 30. Please review the policy changes in the Company Portal."
    }
  }
}

# Example 7: Send custom notifications to both managed and co-managed devices
# Use this for mixed management scenarios
action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "mixed_notification" {
  managed_devices {
    device_id          = "11111111-1111-1111-1111-111111111111"
    notification_title = "IT Support Notification"
    notification_body  = "Your device enrollment is about to expire. Please contact IT support at ext. 5555 to renew."
  }

  comanaged_devices {
    device_id          = "22222222-2222-2222-2222-222222222222"
    notification_title = "IT Support Notification"
    notification_body  = "Your device enrollment is about to expire. Please contact IT support at ext. 5555 to renew."
  }
}

# Example 8: Send BitLocker recovery key reminder
# Use this to inform users about BitLocker key escrow
data "microsoft365_graph_beta_device_management_managed_device" "windows_bitlocker" {
  filter = "operatingSystem eq 'Windows' and encryptionState eq 'encrypted'"
}

action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "bitlocker_reminder" {
  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.windows_bitlocker.managed_devices
    content {
      device_id          = managed_devices.value.id
      notification_title = "BitLocker Recovery Key Information"
      notification_body  = "Your BitLocker recovery key is safely stored. If you need to recover your device, contact IT support with your device ID."
    }
  }
}

# Example 9: Send enrollment completion notification
# Use this to welcome newly enrolled devices
data "microsoft365_graph_beta_device_management_managed_device" "recently_enrolled" {
  filter = "enrolledDateTime gt 2024-10-15T00:00:00Z"
}

action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "welcome_notification" {
  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.recently_enrolled.managed_devices
    content {
      device_id          = managed_devices.value.id
      notification_title = "Welcome to Company Portal"
      notification_body  = "Your device has been successfully enrolled. You now have access to company apps and resources. For help, visit the IT portal at portal.company.com."
    }
  }
}

# Example 10: Send certificate expiration warning
# Use this to notify users about expiring certificates
action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "certificate_expiration" {
  managed_devices {
    device_id          = "cert-device-1"
    notification_title = "Certificate Expires in 7 Days"
    notification_body  = "Your device authentication certificate will expire on November 1. It will be automatically renewed, but please ensure your device stays connected."
  }

  managed_devices {
    device_id          = "cert-device-2"
    notification_title = "Certificate Renewal Required"
    notification_body  = "Your device certificate requires manual renewal. Please sync your device through the Company Portal within 48 hours."
  }
}

# Example 11: Send VPN configuration update notice
# Use this to communicate VPN changes to users
action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "vpn_update" {
  managed_devices {
    device_id          = "vpn-device-1"
    notification_title = "VPN Configuration Update"
    notification_body  = "Your VPN configuration will be updated automatically on October 28. You may need to reconnect to the VPN after the update completes."
  }
}

# Example 12: Send custom notification with extended timeout for large deployments
# Use this for organization-wide notifications
data "microsoft365_graph_beta_device_management_managed_device" "all_managed" {
  filter = "managementAgent eq 'mdm'"
}

action "microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal" "org_wide_notification" {
  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.all_managed.managed_devices
    content {
      device_id          = managed_devices.value.id
      notification_title = "Important Company Announcement"
      notification_body  = "Our company is implementing new security measures. All devices will receive updated policies over the next week. Thank you for your cooperation."
    }
  }

  timeouts = {
    invoke = "45m"
  }
}

