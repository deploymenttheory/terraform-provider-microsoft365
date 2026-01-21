---
page_title: "Microsoft 365_microsoft365_graph_beta_device_management_managed_device_send_custom_notification_to_company_portal Action"
subcategory: "Device Management"
description: |-
  Sends custom notifications to the Company Portal app on managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal and /deviceManagement/comanagedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal endpoints. This action is used to enable IT administrators to send targeted messages to end users through the Company Portal app.
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

Sends custom notifications to the Company Portal app on managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal` and `/deviceManagement/comanagedDevices/{managedDeviceId}/sendCustomNotificationToCompanyPortal` endpoints. This action is used to enable IT administrators to send targeted messages to end users through the Company Portal app.

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

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.42.0-alpha | Experimental | Added missing version history |

## Related Documentation

- [Custom notifications - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/custom-notifications?pivots=ios)
- [Custom notifications - Android](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/custom-notifications?pivots=android)

## Example Usage

```terraform
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
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed devices to send custom notifications to. These are devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and the custom notification title and body.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) When set to `true`, the action will complete successfully even if some devices fail to receive notifications. When `false` (default), the action will fail if any device notification delivery fails. Use this flag when sending notifications to multiple devices and you want the action to succeed even if some deliveries fail.
- `managed_devices` (Attributes List) List of managed devices to send custom notifications to. These are devices fully managed by Intune only. Each entry specifies a device ID and the custom notification title and body for that device.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. You can provide both to send notifications to different types of devices in one action. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) When set to `true` (default), the action will validate that all specified devices exist and support custom notifications (iOS, iPadOS, Android only) before sending notifications. When `false`, device validation is skipped and the action will attempt to send notifications directly. Disabling validation can improve performance but may result in errors if devices don't exist or are unsupported.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to send the notification to. Example: `"12345678-1234-1234-1234-123456789abc"`
- `notification_body` (String) The body/content of the custom notification to display in the Company Portal app. Should provide clear instructions or information to the user. Maximum recommended length: 200-300 characters. Example: `"Your device is not compliant with corporate security policies. Please contact IT support."`
- `notification_title` (String) The title of the custom notification to display in the Company Portal app. Should be concise and descriptive. Maximum recommended length: 50-60 characters. Example: `"Compliance Alert"`


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed device to send the notification to. Example: `"12345678-1234-1234-1234-123456789abc"`
- `notification_body` (String) The body/content of the custom notification to display in the Company Portal app. Should provide clear instructions or information to the user. Maximum recommended length: 200-300 characters. Example: `"Your device password will expire in 3 days. Please update it to maintain access to company resources."`
- `notification_title` (String) The title of the custom notification to display in the Company Portal app. Should be concise and descriptive. Maximum recommended length: 50-60 characters. Example: `"Action Required: Update Your Password"`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

