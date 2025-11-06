---
page_title: "Microsoft 365_microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal Action"
subcategory: "Device Management"
description: |-
  Sends custom notifications to the Company Portal app on managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal and /deviceManagement/comanagedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal endpoints. This action enables IT administrators to send targeted messages to end users through the Company Portal app.
  What This Action Does:
  Sends push notification to Company Portal appDisplays custom message title and bodyTargets specific devices or usersSupports customized messages per deviceProvides in-app notification visibilityEnables two-way communication channel
  When to Use:
  Compliance reminders and deadlinesSecurity alert communicationsPolicy update notificationsMaintenance window announcementsAction required messagesUser guidance and instructionsIncident response communications
  Platform Support:
  Windows: Company Portal app requirediOS/iPadOS: Company Portal app requiredAndroid: Company Portal app requiredmacOS: Company Portal app required
  Important Considerations:
  Company Portal app must be installedDevice must be enrolled in IntuneUser must be signed into Company PortalDevice must have network connectivityNotifications appear in Company Portal appConsider user time zones for timing
  Reference: Microsoft Graph API - Send Custom Notification To Company Portal https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-sendcustomnotificationtocompanyportal?view=graph-rest-beta
---

# Microsoft 365_microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal (Action)

Sends custom notifications to the Company Portal app on managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal` and `/deviceManagement/comanagedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal` endpoints. This action enables IT administrators to send targeted messages to end users through the Company Portal app.

**What This Action Does:**
- Sends push notification to Company Portal app
- Displays custom message title and body
- Targets specific devices or users
- Supports customized messages per device
- Provides in-app notification visibility
- Enables two-way communication channel

**When to Use:**
- Compliance reminders and deadlines
- Security alert communications
- Policy update notifications
- Maintenance window announcements
- Action required messages
- User guidance and instructions
- Incident response communications

**Platform Support:**
- **Windows**: Company Portal app required
- **iOS/iPadOS**: Company Portal app required
- **Android**: Company Portal app required
- **macOS**: Company Portal app required

**Important Considerations:**
- Company Portal app must be installed
- Device must be enrolled in Intune
- User must be signed into Company Portal
- Device must have network connectivity
- Notifications appear in Company Portal app
- Consider user time zones for timing

**Reference:** [Microsoft Graph API - Send Custom Notification To Company Portal](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-sendcustomnotificationtocompanyportal?view=graph-rest-beta)

## Use Cases

This action enables targeted communication with end users through the Company Portal app across all managed devices:

### Compliance & Policy Management
- **Compliance Reminders**: Send targeted notifications to non-compliant devices reminding users to address compliance issues
- **Policy Updates**: Communicate new or updated device policies before enforcement begins
- **Deadline Notifications**: Alert users about upcoming compliance deadlines or grace period expirations
- **Action Required**: Prompt users to take specific actions to maintain device compliance status
- **Conditional Access Changes**: Notify users about changes to conditional access policies affecting their devices
- **Grace Period Warnings**: Alert users when devices are in grace period before enforcement actions
- **Policy Violation Alerts**: Inform users when their devices violate specific security policies

### Security Communications
- **Security Alerts**: Broadcast critical security announcements to all or specific groups of devices
- **Patch Notifications**: Alert users about required security patches or updates
- **Threat Detection**: Communicate when security threats are detected on devices
- **Certificate Expiration**: Warn users about expiring authentication certificates
- **Password Expiration**: Send personalized password expiration reminders with specific dates
- **BitLocker Key Changes**: Inform users when BitLocker keys are rotated or updated
- **VPN Configuration**: Notify users about VPN configuration changes or required reconnections
- **Incident Response**: Communicate security incident information and required user actions

### Device Management
- **Maintenance Windows**: Notify users about scheduled maintenance or service windows
- **App Updates**: Alert users about required app updates available in Company Portal
- **Enrollment Status**: Send welcome messages to newly enrolled devices or users
- **Configuration Changes**: Inform users about pending configuration profile deployments
- **Device Retirement**: Notify users when devices are scheduled for retirement or decommissioning
- **Sync Reminders**: Prompt users to sync their devices to receive latest policies
- **Backup Reminders**: Remind users to backup data before maintenance or updates

### User Support & Guidance
- **Onboarding Messages**: Welcome new users and provide getting-started information
- **Help Resources**: Direct users to support documentation, portals, or contact information
- **Training Announcements**: Notify users about available training sessions or materials
- **Feature Announcements**: Inform users about new Company Portal features or capabilities
- **Self-Service Actions**: Guide users to perform self-service actions through Company Portal
- **FAQ Updates**: Alert users when important FAQ or help documentation is updated

### Operational Communications
- **Scheduled Outages**: Communicate planned system or service outages to affected users
- **Service Restoration**: Notify users when services are restored after outages
- **Migration Notices**: Inform users about device or system migrations
- **Pilot Program**: Communicate with pilot users about new features or configurations
- **Survey Requests**: Request user feedback or participation in surveys
- **Audit Communications**: Notify users about compliance audits or device checks
- **Organization Announcements**: Broadcast company-wide IT announcements through Company Portal

## API Documentation

- [Microsoft Graph API - Send Custom Notification To Company Portal](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-sendcustomnotificationtocompanyportal?view=graph-rest-beta)

## Permissions

The following Microsoft Graph API permissions are required to use this action:

| Permission Type | Permissions (Least Privileged) |
|:----------------|:------------------------------|
| Delegated (work or school account) | DeviceManagementConfiguration.ReadWrite.All, DeviceManagementManagedDevices.ReadWrite.All |
| Delegated (personal Microsoft account) | Not supported |
| Application | DeviceManagementConfiguration.ReadWrite.All, DeviceManagementManagedDevices.ReadWrite.All |

~> **Note:** This action requires both device configuration and device management write permissions as it sends notifications to managed devices.

## Related Documentation

- [Custom notifications - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/custom-notifications?pivots=ios)
- [Custom notifications - Android](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/custom-notifications?pivots=android)

## Notes

### Platform Compatibility

Custom notifications to Company Portal are only supported on mobile platforms:

| Platform | Supported | Company Portal App Required |
|:---------|:----------|:---------------------------|
| **iOS** | ✅ | Company Portal from App Store |
| **iPadOS** | ✅ | Company Portal from App Store |
| **Android** | ✅ | Company Portal from Google Play |
| **Android Enterprise** | ✅ | Company Portal from managed Google Play |
| **Windows 10/11** | ❌ | Not supported for custom notifications |
| **macOS** | ❌ | Not supported for custom notifications |
| **ChromeOS** | ❌ | Not supported for custom notifications |

### Important Considerations

**Company Portal Requirements:**
- **App Installation**: Company Portal app must be installed on the target device
- **User Sign-In**: User must be signed into the Company Portal app to receive notifications
- **App Version**: Company Portal app should be updated to latest version for best compatibility
- **Notification Permissions**: Company Portal must have notification permissions enabled on the device
- **Background App**: Company Portal may need to be running in background on some platforms
- **Device Enrollment**: Device must be actively enrolled and managed by Intune

**Notification Delivery:**
- **Delivery Method**: Notifications are delivered through the Company Portal app, not as OS-level push notifications
- **Timing**: Notifications appear when users open or check the Company Portal app
- **Persistence**: Notifications persist in the Company Portal notifications section until dismissed by user
- **No Guaranteed Delivery**: There is no delivery confirmation; notification may fail silently if requirements not met
- **Network Dependency**: Device requires network connectivity to receive notifications
- **Offline Devices**: Offline devices will receive notifications when they reconnect and sync with Intune

**Per-Device Customization:**
- **Individual Messages**: Each device can have a completely different notification title and body
- **Personalization**: Enables personalized messaging based on device state, user, or other attributes
- **Flexible Targeting**: Mix different message types in single action (e.g., different expiration dates)
- **Dynamic Content**: Use Terraform's dynamic blocks to generate device-specific messages from data sources

**Notification Content Guidelines:**
- **Title Length**: Keep titles concise (50-60 characters recommended, 250 character maximum)
- **Body Length**: Keep body messages clear and actionable (200-300 characters recommended, 1000 character maximum)
- **Clear Language**: Use simple, direct language that users can understand
- **Call to Action**: Include specific next steps or actions for users to take
- **Contact Information**: Provide support contact details when appropriate
- **Urgency Indication**: Clearly indicate if action is urgent or time-sensitive
- **Professional Tone**: Maintain professional, respectful communication tone

**User Experience:**
- **In-App Display**: Notifications appear in the Company Portal app's notification section
- **No OS Notification**: Typically does not trigger operating system push notifications
- **User Dismissal**: Users can dismiss notifications; they won't reappear after dismissal
- **No Response Tracking**: No way to track if users read or acted on notifications
- **Multiple Notifications**: Sending multiple notifications to same device accumulates in notification list

**Operational Considerations:**
- **Bulk Notifications**: This action supports sending different messages to many devices simultaneously
- **Concurrent Processing**: Notifications are sent concurrently for better performance
- **Rate Limiting**: Be mindful of API rate limits when sending to thousands of devices
- **Retry Logic**: Failed notification sends are not automatically retried; use Terraform retry mechanisms if needed
- **Idempotency**: Re-running the action sends duplicate notifications; avoid frequent re-runs
- **Testing**: Test notifications with small groups before large-scale deployments

**Time Zone Considerations:**
- **No Scheduling**: This action sends notifications immediately; no built-in scheduling capability
- **User Time Zones**: Consider user time zones when sending time-sensitive notifications
- **Business Hours**: Best practice to send during business hours in target user's region
- **After-Hours**: Avoid sending non-urgent notifications during nights or weekends
- **Global Deployments**: For global organizations, stage notifications by region/time zone

**Best Practices:**
- **Targeted Messaging**: Use filters and data sources to target specific device groups
- **Message Clarity**: Ensure messages are clear, actionable, and provide necessary context
- **Frequency Management**: Avoid notification fatigue by limiting frequency of messages
- **User-Centric**: Write from user's perspective; focus on what they need to do
- **Test First**: Always test with small group before organization-wide rollout
- **Fallback Communication**: Have backup communication channels (email, Teams) for critical messages
- **Accessibility**: Use clear language and avoid jargon for accessibility
- **Localization**: Consider language localization for multinational organizations
- **Documentation**: Maintain records of notifications sent for audit and reference

**Technical Limitations:**
- **No Read Receipts**: No way to confirm users have read the notification
- **No Response Collection**: Cannot collect user responses or acknowledgments through this action
- **No Notification History**: Limited ability to retrieve sent notification history through API
- **Character Limits**: Titles (250 chars) and bodies (1000 chars) have strict limits enforced by API
- **No Rich Media**: Cannot include images, links, or rich formatting in notifications
- **No Actionable Buttons**: Cannot add buttons or interactive elements to notifications

**Security Considerations:**
- **Information Disclosure**: Avoid including sensitive information in notification content
- **Authentication**: Notifications visible to anyone with access to Company Portal on device
- **Phishing Prevention**: Use consistent, recognizable language to prevent phishing confusion
- **Contact Verification**: Users should verify unexpected notifications through official channels
- **Audit Logging**: All notification actions are logged in Intune audit logs

**Troubleshooting:**
- **Not Receiving Notifications**: Verify Company Portal is installed, user is signed in, and device is enrolled
- **Delayed Notifications**: Check device network connectivity and Company Portal app sync status
- **Permission Issues**: Ensure Company Portal has necessary notification permissions on device
- **Silent Failures**: API may return success even if notification isn't delivered; verify delivery separately
- **App Not Installed**: If Company Portal isn't installed, notification fails silently without error

**Monitoring & Analytics:**
- **Delivery Confirmation**: No built-in delivery confirmation or read receipts available
- **Audit Logs**: Check Intune audit logs to see notification send attempts and any errors
- **User Feedback**: Rely on user feedback or support tickets to gauge notification effectiveness
- **Action Logs**: Check Terraform/provider logs for any errors during notification sending

## Example Usage

```terraform
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
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Block List) List of co-managed devices to send custom notifications to. These are devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and the custom notification title and body.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--comanaged_devices))
- `managed_devices` (Block List) List of managed devices to send custom notifications to. These are devices fully managed by Intune only. Each entry specifies a device ID and the custom notification title and body for that device.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. You can provide both to send notifications to different types of devices in one action. (see [below for nested schema](#nestedblock--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedblock--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to send the notification to. Example: `"12345678-1234-1234-1234-123456789abc"`
- `notification_body` (String) The body/content of the custom notification to display in the Company Portal app. Should provide clear instructions or information to the user. Maximum recommended length: 200-300 characters. Example: `"Your device is not compliant with corporate security policies. Please contact IT support."`
- `notification_title` (String) The title of the custom notification to display in the Company Portal app. Should be concise and descriptive. Maximum recommended length: 50-60 characters. Example: `"Compliance Alert"`


<a id="nestedblock--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed device to send the notification to. Example: `"12345678-1234-1234-1234-123456789abc"`
- `notification_body` (String) The body/content of the custom notification to display in the Company Portal app. Should provide clear instructions or information to the user. Maximum recommended length: 200-300 characters. Example: `"Your device password will expire in 3 days. Please update it to maintain access to company resources."`
- `notification_title` (String) The title of the custom notification to display in the Company Portal app. Should be concise and descriptive. Maximum recommended length: 50-60 characters. Example: `"Action Required: Update Your Password"`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

