---
page_title: "microsoft365_graph_beta_device_management_notification_message_template Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages an Intune notification message template for compliance notifications
---

# microsoft365_graph_beta_device_management_notification_message_template (Resource)

Manages an Intune notification message template for compliance notifications


## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All` , `DeviceManagementConfiguration.Read.All`

## Example Usage

```terraform
# Example: Basic Notification Message Template
resource "microsoft365_graph_beta_device_management_notification_message_template" "basic" {
  display_name     = "Basic Compliance Notification"
  description      = "Basic notification template for device compliance violations"
  default_locale   = "en-US"
  branding_options = "includeCompanyLogo"

  role_scope_tag_ids = ["0"]

  localized_notification_messages = [
    {
      locale           = "en-US"
      subject          = "Device Compliance Issue Detected"
      message_template = "Hello {UserName},\n\nYour device '{DeviceName}' has been found to be non-compliant with company policies. Please take action to resolve the following issues:\n\n{ComplianceReasons}\n\nFor assistance, please contact IT support.\n\nThank you,\nIT Security Team"
      is_default       = true
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example: Multi-language Notification Message Template
resource "microsoft365_graph_beta_device_management_notification_message_template" "multilingual" {
  display_name     = "Multi-language Compliance Notification"
  description      = "Notification template with multiple language support"
  default_locale   = "en-US"
  branding_options = "includeCompanyLogo"

  role_scope_tag_ids = ["0"]

  localized_notification_messages = [
    {
      locale           = "en-US"
      subject          = "Device Compliance Issue"
      message_template = "Hello {UserName},\n\nYour device '{DeviceName}' is not compliant. Please resolve: {ComplianceReasons}\n\nContact IT for help.\n\nIT Security Team"
      is_default       = true
    },
    {
      locale           = "es-ES"
      subject          = "Problema de Cumplimiento del Dispositivo"
      message_template = "Hola {UserName},\n\nTu dispositivo '{DeviceName}' no cumple las normas. Por favor resuelve: {ComplianceReasons}\n\nContacta con IT para ayuda.\n\nEquipo de Seguridad IT"
      is_default       = false
    },
    {
      locale           = "fr-FR"
      subject          = "Problème de Conformité de l'Appareil"
      message_template = "Bonjour {UserName},\n\nVotre appareil '{DeviceName}' n'est pas conforme. Veuillez résoudre: {ComplianceReasons}\n\nContactez l'IT pour aide.\n\nÉquipe de Sécurité IT"
      is_default       = false
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example: Advanced Notification Template with Full Branding
resource "microsoft365_graph_beta_device_management_notification_message_template" "advanced" {
  display_name     = "Advanced Compliance Notification"
  description      = "Advanced notification template with comprehensive branding and device details"
  default_locale   = "en-US"
  branding_options = "includeCompanyLogo"

  role_scope_tag_ids = ["0", "1"]

  localized_notification_messages = [
    {
      locale           = "en-US"
      subject          = "Immediate Action Required: Device Compliance"
      message_template = <<-EOT
        Dear {UserName},

        Your device '{DeviceName}' (Serial: {DeviceSerialNumber}) has been identified as non-compliant with our security policies.

        Compliance Issues:
        {ComplianceReasons}

        Required Actions:
        1. Update your device to the latest security patches
        2. Enable BitLocker encryption if not already enabled
        3. Ensure Windows Defender is active and up-to-date
        4. Contact IT support if you need assistance

        Device Details:
        - Device Name: {DeviceName}
        - Operating System: {DeviceOSVersion}
        - Last Check-in: {LastCheckInTime}

        Failure to address these issues within 24 hours may result in restricted access to company resources.

        For immediate assistance, contact:
        - IT Helpdesk: +1-555-0123
        - Email: it-support@company.com
        - Portal: https://company.com/it-support

        Best regards,
        IT Security Team
        Company Name
      EOT
      is_default       = true
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "5m"
  }
}

# Output examples
output "basic_template_id" {
  description = "ID of the basic notification message template"
  value       = microsoft365_graph_beta_device_management_notification_message_template.basic.id
}

output "multilingual_template_id" {
  description = "ID of the multi-language notification message template"
  value       = microsoft365_graph_beta_device_management_notification_message_template.multilingual.id
}

output "advanced_template_id" {
  description = "ID of the advanced notification message template"
  value       = microsoft365_graph_beta_device_management_notification_message_template.advanced.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `default_locale` (String) The default locale to fallback onto when the requested locale is not available
- `display_name` (String) Display name for the notification message template

### Optional

- `branding_options` (String) The branding options for the message template. Possible values are: none, includeCompanyLogo, includeCompanyName, includeContactInformation, includeCompanyPortalLink, includeDeviceDetails
- `description` (String) Description of the notification message template
- `localized_notification_messages` (Attributes Set) The list of localized notification messages for this template (see [below for nested schema](#nestedatt--localized_notification_messages))
- `role_scope_tag_ids` (Set of String) List of scope tag IDs for this notification message template
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) Manages notification message templates in Microsoft Intune using the `/deviceManagement/notificationMessageTemplates` endpoint. Notification message templates define the content and branding of compliance notifications sent to users.
- `last_modified_date_time` (String) DateTime the notification message template was last modified

<a id="nestedatt--localized_notification_messages"></a>
### Nested Schema for `localized_notification_messages`

Required:

- `locale` (String) The locale for the notification message (e.g., en-US, es-ES)
- `message_template` (String) The message template text that can include tokens like {DeviceName}, {UserName}, etc.
- `subject` (String) The subject of the notification message

Optional:

- `is_default` (Boolean) Indicates if this is the default message for the template

Read-Only:

- `id` (String) Unique identifier for the localized notification message
- `last_modified_date_time` (String) DateTime the localized message was last modified


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
| v0.25.0-alpha | Experimental | Initial release |

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import examples for Microsoft 365 Graph Beta Device Management Notification Message Template

# Basic import using template ID
# Replace "12345678-1234-1234-1234-123456789012" with the actual notification message template ID from Microsoft Graph
terraform import microsoft365_graph_beta_device_management_notification_message_template.basic "12345678-1234-1234-1234-123456789012"

# Import multilingual template
terraform import microsoft365_graph_beta_device_management_notification_message_template.multilingual "87654321-4321-4321-4321-210987654321"

# Import advanced template
terraform import microsoft365_graph_beta_device_management_notification_message_template.advanced "11111111-2222-3333-4444-555555555555"

# To find existing notification message template IDs, you can use Microsoft Graph Explorer:
# https://developer.microsoft.com/en-us/graph/graph-explorer
#
# Use this query to list all notification message templates:
# GET https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates
#
# Example response will include template IDs like:
# {
#   "value": [
#     {
#       "id": "12345678-1234-1234-1234-123456789012",
#       "displayName": "Basic Compliance Notification",
#       "description": "Basic notification template for device compliance violations",
#       "defaultLocale": "en-US",
#       "brandingOptions": "includeCompanyLogo"
#     }
#   ]
# }

# Note: Ensure you have appropriate permissions to read notification message templates:
# - DeviceManagementServiceConfig.Read.All (for reading)
# - DeviceManagementServiceConfig.ReadWrite.All (for importing and managing)
```