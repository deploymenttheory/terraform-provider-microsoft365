---
page_title: "microsoft365_graph_beta_device_management_device_enrollment_notification Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages device enrollment notification configurations using the /deviceManagement/deviceEnrollmentConfigurations endpoint. This resource is used to creates device enrollment notification configurations for Android for Work platform.
---

# microsoft365_graph_beta_device_management_device_enrollment_notification (Resource)

Manages device enrollment notification configurations using the `/deviceManagement/deviceEnrollmentConfigurations` endpoint. This resource is used to creates device enrollment notification configurations for Android for Work platform.

## Microsoft Documentation

- [deviceEnrollmentNotificationConfiguration resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentnotificationconfiguration?view=graph-rest-beta)
- [Create deviceEnrollmentNotificationConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-onboarding-deviceenrollmentnotificationconfiguration-create?view=graph-rest-beta&tabs=http)
- [Update deviceEnrollmentNotificationConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-onboarding-deviceenrollmentnotificationconfiguration-update?view=graph-rest-beta&tabs=http)
- [Delete deviceEnrollmentNotificationConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-onboarding-deviceenrollmentnotificationconfiguration-delete?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `DeviceManagementConfiguration.ReadWrite.All`
- `DeviceManagementConfiguration.Read.All`

**Optional:**
- `None` `[N/A]`

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_management_device_enrollment_notification" "email_minimal" {
  display_name     = "email minimal"
  description      = "minimal configuration for email"
  platform_type    = "mac" // "ios", "windows", "android", "androidForWork", "mac", "linux"
  branding_options = ["none"]

  notification_templates = ["email"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_management_device_enrollment_notification" "email_maximal" {
  display_name  = "email maximal"
  description   = "Complete configuration withall features"
  platform_type = "androidForWork" // "ios", "windows", "android", "androidForWork", "mac", "linux"
  branding_options = ["includeCompanyLogo",
    "includeCompanyName",
    "includeCompanyPortalLink",
    "includeContactInformation",
    "includeDeviceDetails"
  ]

  notification_templates = ["email"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_management_device_enrollment_notification" "push_maximal" {
  display_name           = "push maximal"
  description            = "Complete push configuration"
  platform_type          = "linux"  // "ios", "windows", "android", "androidForWork", "mac", "linux"
  branding_options       = ["none"] // no branding options for push
  notification_templates = ["push"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "push"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_management_device_enrollment_notification" "all" {
  display_name   = "configuration with all enrollment notification features"
  description    = "Complete configuration with all features"
  platform_type  = "androidForWork" // "androidForWork" , "android"
  default_locale = "en-US"
  branding_options = ["includeCompanyLogo",
    "includeCompanyName",
    "includeCompanyPortalLink",
    "includeContactInformation",
    "includeDeviceDetails"
  ]

  notification_templates = ["email", "push"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "push"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `branding_options` (Set of String) The branding options for the message template. Possible values are: none, includeCompanyLogo, includeCompanyName, includeContactInformation, includeCompanyPortalLink, includeDeviceDetails. Defaults to ['none'].
- `display_name` (String) The display name for the device enrollment notification configuration.
- `localized_notification_messages` (Attributes Set) The localized notification messages for the configuration. (see [below for nested schema](#nestedatt--localized_notification_messages))
- `platform_type` (String) The platform type for the notification configuration. Must be either 'ios', 'windows', 'android', 'androidForWork', 'mac', 'linux'.

### Optional

- `assignments` (Attributes Set) Assignments for the compliance policy. Each assignment specifies the target group and schedule for script execution. (see [below for nested schema](#nestedatt--assignments))
- `default_locale` (String) The default locale for the notification configuration (e.g., 'en-US').
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `notification_templates` (Set of String) The notification template types for this configuration. Can be 'email', 'push', or both. Defaults to ['email', 'push'].
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Android Enterprise Notification configuration.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `device_enrollment_configuration_type` (String) The type of device enrollment configuration.
- `id` (String) The unique identifier for the device enrollment notification configuration.
- `priority` (Number) The priority of the notification configuration.

<a id="nestedatt--localized_notification_messages"></a>
### Nested Schema for `localized_notification_messages`

Required:

- `locale` (String) The locale for the notification message (e.g., 'en-us'). Must be in lowercase format.
- `message_template` (String) The template content of the notification message.
- `subject` (String) The subject of the notification message.
- `template_type` (String) The type of template (email or push).

Optional:

- `is_default` (Boolean) Whether this is the default notification message.


<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget'.

Optional:

- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget'.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.27.0-alpha | Experimental | Initial release |
| v0.28-alpha | Experimental | Added support for Windows platform type and updated the resource name to `device_enrollment_notification` |

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_management_device_enrollment_notification.example 00000000-0000-0000-0000-000000000000_EnrollmentNotificationsConfiguration
```