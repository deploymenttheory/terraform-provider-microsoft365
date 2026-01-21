---
page_title: "microsoft365_graph_beta_users_user_mailbox_settings Resource - terraform-provider-microsoft365"
subcategory: "Users"
description: |-
  Manages Microsoft 365 user mailbox settings using the /users/{id}/mailboxSettings endpoint. This resource is used to allows you to configure automatic replies, date/time formats, locale, time zone, working hours, and other mailbox preferences for a user. Note: This resource manages settings that may also be modified by users through Outlook clients.
---

# microsoft365_graph_beta_users_user_mailbox_settings (Resource)

Manages Microsoft 365 user mailbox settings using the `/users/{id}/mailboxSettings` endpoint. This resource is used to allows you to configure automatic replies, date/time formats, locale, time zone, working hours, and other mailbox preferences for a user. Note: This resource manages settings that may also be modified by users through Outlook clients.

## Microsoft Documentation

- [mailboxSettings resource type](https://learn.microsoft.com/en-us/graph/api/resources/mailboxsettings?view=graph-rest-beta)
- [Get mailboxSettings](https://learn.microsoft.com/en-us/graph/api/user-get-mailboxsettings?view=graph-rest-beta)
- [Update mailboxSettings](https://learn.microsoft.com/en-us/graph/api/user-update-mailboxsettings?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `MailboxSettings.ReadWrite`

**Optional:**
- `None` `[N/A]`, `User.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.40.0-alpha | Experimental | Initial release |

## Example Usage

### Minimal Example

```terraform
# Example 1: Minimal mailbox settings configuration
# This example shows the minimum required configuration for user mailbox settings
resource "microsoft365_graph_beta_users_user_mailbox_settings" "minimal" {
  user_id   = "john.doe@example.com" # Can be user ID (UUID) or UPN
  time_zone = "UTC"
}
```

### Maximal Example

```terraform
# Example 2: Maximal mailbox settings configuration
# This example shows all available configuration options for user mailbox settings
resource "microsoft365_graph_beta_users_user_mailbox_settings" "maximal" {
  user_id                                   = "jane.smith@example.com"
  time_zone                                 = "Pacific Standard Time"
  date_format                               = "dd/MM/yyyy"
  time_format                               = "HH:mm"
  delegate_meeting_message_delivery_options = "sendToDelegateOnly"

  # Configure automatic replies (Out of Office)
  automatic_replies_setting = {
    status            = "scheduled"
    external_audience = "all"

    scheduled_start_date_time = {
      date_time = "2024-12-20T00:00:00"
      time_zone = "Pacific Standard Time"
    }

    scheduled_end_date_time = {
      date_time = "2024-12-30T00:00:00"
      time_zone = "Pacific Standard Time"
    }

    internal_reply_message = "<html><body><p>I'm out of office and will respond when I return.</p></body></html>"
    external_reply_message = "<html><body><p>I'm currently out of office. For urgent matters, please contact support@example.com.</p></body></html>"
  }

  # Configure language/locale settings
  language = {
    locale = "en-GB"
  }

  # Configure working hours
  working_hours = {
    days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday"]
    start_time   = "09:00:00"
    end_time     = "17:00:00"

    time_zone = {
      name = "Pacific Standard Time"
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
```

### Complete Workflow with Dependencies

```terraform
# Example 3: Complete workflow with user creation and license assignment
# This example demonstrates creating a user, assigning an Exchange Online license,
# and then configuring mailbox settings

# Step 1: Look up the license SKU using the licensing service plan reference datasource
# This ensures you always have the correct GUID without hardcoding it
data "microsoft365_utility_licensing_service_plan_reference" "m365_e3" {
  string_id = "ENTERPRISEPACK" # Microsoft 365 E3

  # Alternative search options:
  # product_name = "Microsoft 365 E3"
  # guid = "6fd2c87f-b296-42f0-b197-1e91e994b900"
}

# Step 2: Create the user
resource "microsoft365_graph_beta_users_user" "example_user" {
  display_name        = "Example User"
  user_principal_name = "example.user@yourdomain.com"
  mail_nickname       = "example.user"
  account_enabled     = true
  usage_location      = "US" # Field is required for license assignment

  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Step 3: Assign a license that includes Exchange Online
resource "microsoft365_graph_beta_users_user_license_assignment" "example_user_license" {
  user_id = microsoft365_graph_beta_users_user.example_user.id

  # Use the dynamically looked-up SKU ID from the datasource
  # This is more maintainable than hardcoding GUIDs and ensures accuracy
  sku_id = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].guid

  # Optional: Disable specific service plans
  disabled_plans = []
}

# Step 4: Wait for mailbox provisioning
# Exchange Online mailboxes can take 1-2 minutes to provision after license assignment
resource "time_sleep" "wait_for_mailbox_provisioning" {
  depends_on = [microsoft365_graph_beta_users_user_license_assignment.example_user_license]

  create_duration = "2m"
}

# Step 5: Configure mailbox settings
resource "microsoft365_graph_beta_users_user_mailbox_settings" "example_user_settings" {
  depends_on = [time_sleep.wait_for_mailbox_provisioning]

  user_id                                   = microsoft365_graph_beta_users_user.maximal_dependency_user.id
  time_zone                                 = "Greenwich Standard Time"
  delegate_meeting_message_delivery_options = "sendToDelegateOnly"

  automatic_replies_setting = {
    status            = "scheduled"
    external_audience = "all"

    scheduled_start_date_time = {
      date_time = "2030-03-14T07:00:00"
      time_zone = "UTC"
    }

    scheduled_end_date_time = {
      date_time = "2030-03-28T07:00:00"
      time_zone = "UTC"
    }

    internal_reply_message = "<html>\n<body>\n<p>I'm at our company's worldwide reunion and will respond to your message as soon as I return.<br>\n</p></body>\n</html>\n"
    external_reply_message = "<html>\n<body>\n<p>I'm at the Deployment Theory worldwide reunion and will respond to your message as soon as I return.<br>\n</p></body>\n</html>\n"
  }

  language = {
    locale = "en-US"
  }

  working_hours = {
    days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday"]
    start_time   = "08:00:00"
    end_time     = "17:00:00"

    time_zone = {
      name = "Greenwich Standard Time"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `user_id` (String) The ID of the user whose mailbox settings are being managed. This can be the user's object ID or userPrincipalName.

### Optional

- `automatic_replies_setting` (Attributes) Configuration for automatic replies (also known as Out of Office or OOF) for the user's mailbox. (see [below for nested schema](#nestedatt--automatic_replies_setting))
- `date_format` (String) The date format for the user's mailbox. This uses [.NET standard date format patterns](https://learn.microsoft.com/en-us/dotnet/standard/base-types/standard-date-and-time-format-strings#ShortDate) that are culture-specific. Common values include: `M/d/yyyy` (US), `dd/MM/yyyy` (UK/EU), `yyyy-MM-dd` (ISO), `dd.MM.yyyy` (German). The format determines how dates are displayed in the user's mailbox.
- `delegate_meeting_message_delivery_options` (String) Specifies how meeting messages and responses are delivered to delegates. Possible values: `sendToDelegateAndInformationToPrincipal`, `sendToDelegateAndPrincipal`, `sendToDelegateOnly`.
- `language` (Attributes) The locale (language and country/region) information for the user. (see [below for nested schema](#nestedatt--language))
- `time_format` (String) The time format for the user's mailbox. This uses [.NET standard time format patterns](https://learn.microsoft.com/en-us/dotnet/standard/base-types/standard-date-and-time-format-strings#ShortTime) that are culture-specific. Common examples include: `h:mm tt` (1:45 PM - US 12-hour), `HH:mm` (13:45 - European 24-hour). The format determines how times are displayed in the user's mailbox.
- `time_zone` (String) The default time zone for the user's mailbox. Must be one of the Windows time zone names supported by Microsoft Graph API. Common values include `Pacific Standard Time`, `Eastern Standard Time`, `UTC`, etc. See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/outlookuser-supportedtimezones) for the full list of supported time zones.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `working_hours` (Attributes) The working hours configured for the user's calendar. (see [below for nested schema](#nestedatt--working_hours))

### Read-Only

- `id` (String) Computed identifier for this resource (format: users/{user_id}/mailboxSettings). Read-only.
- `user_purpose` (String) The purpose of the mailbox. Differentiates a mailbox for a single user from a shared mailbox and equipment mailbox in Exchange Online. Possible values are: user, linked, shared, room, equipment, others, unknownFutureValue. Read-only.

<a id="nestedatt--automatic_replies_setting"></a>
### Nested Schema for `automatic_replies_setting`

Optional:

- `external_audience` (String) The audience that will receive external automatic reply messages. Possible values: `none`, `contactsOnly`, `all`.
- `external_reply_message` (String) The automatic reply message to send to external recipients. Supports HTML formatting.
- `internal_reply_message` (String) The automatic reply message to send to internal recipients. Supports HTML formatting.
- `scheduled_end_date_time` (Attributes) The end date and time when automatic replies are scheduled to stop being sent. Required when status is `scheduled`. (see [below for nested schema](#nestedatt--automatic_replies_setting--scheduled_end_date_time))
- `scheduled_start_date_time` (Attributes) The start date and time when automatic replies are scheduled to be sent. Required when status is `scheduled`. (see [below for nested schema](#nestedatt--automatic_replies_setting--scheduled_start_date_time))
- `status` (String) The status of automatic replies. Possible values: `disabled`, `alwaysEnabled`, `scheduled`.

<a id="nestedatt--automatic_replies_setting--scheduled_end_date_time"></a>
### Nested Schema for `automatic_replies_setting.scheduled_end_date_time`

Optional:

- `date_time` (String) The date and time value in ISO 8601 format (e.g., `2026-03-20T02:00:00`). The timezone is specified separately in the `time_zone` field.
- `time_zone` (String) The time zone for the date time value. Defaults to `UTC` if not specified.


<a id="nestedatt--automatic_replies_setting--scheduled_start_date_time"></a>
### Nested Schema for `automatic_replies_setting.scheduled_start_date_time`

Optional:

- `date_time` (String) The date and time value in ISO 8601 format (e.g., `2026-03-19T02:00:00`). The timezone is specified separately in the `time_zone` field.
- `time_zone` (String) The time zone for the date time value. Defaults to `UTC` if not specified.



<a id="nestedatt--language"></a>
### Nested Schema for `language`

Required:

- `locale` (String) A locale representation for the user, which includes the user's preferred language and country/region. For example, `en-US`. The language component follows 2-letter codes as defined in [ISO 639-1](https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes), and the country component follows 2-letter codes as defined in [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2).

Read-Only:

- `display_name` (String) The display name of the locale. Read-only.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--working_hours"></a>
### Nested Schema for `working_hours`

Optional:

- `days_of_week` (Set of String) The days of the week on which the user works. Possible values: `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`.
- `end_time` (String) The time the user stops working each day, in HH:mm:ss format (e.g., `17:00:00`).
- `start_time` (String) The time the user starts working each day, in HH:mm:ss format (e.g., `09:00:00`).
- `time_zone` (Attributes) The time zone for the working hours. Must be one of the Windows time zone names supported by Microsoft Graph API. (see [below for nested schema](#nestedatt--working_hours--time_zone))

<a id="nestedatt--working_hours--time_zone"></a>
### Nested Schema for `working_hours.time_zone`

Optional:

- `name` (String) The name of the time zone. Must be one of the supported Windows time zone names (e.g., `Pacific Standard Time`, `Eastern Standard Time`, `UTC`). See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/outlookuser-supportedtimezones) for the full list.

## Important Notes

- **Exchange Online Requirement**: Mailbox settings can only be configured for users with an active Exchange Online mailbox. Users must have a license that includes Exchange Online (e.g., Microsoft 365 E3/E5, Exchange Online Plan 1/2).
- **Mailbox Provisioning**: After assigning an Exchange Online license to a user, allow 1-2 minutes for the mailbox to provision before configuring settings.
- **Resource Behavior**: This resource uses PATCH operations for both create and update. The mailbox settings always exist for a user with an Exchange Online mailbox, so constants.TfOperationCreate effectively updates the existing settings.
- **Delete Behavior**: Deleting this resource only removes it from Terraform state without affecting the actual mailbox settings. Mailbox settings cannot be deleted, only reset to defaults manually.
- **Time Zones**: Use IANA time zone identifiers or Windows time zone names. Common examples: `UTC`, `Pacific Standard Time`, `Eastern Standard Time`.
- **Date/Time Formats**: Use standard format strings compatible with Outlook. Refer to [.NET standard date and time format strings](https://learn.microsoft.com/en-us/dotnet/standard/base-types/standard-date-and-time-format-strings) for valid formats.
- **Automatic Replies**: When configuring scheduled automatic replies, both `scheduled_start_date_time` and `scheduled_end_date_time` must be provided.
- **Working Hours**: The `working_hours` configuration is used by scheduling assistants and affects meeting time suggestions.
- **Locale Format**: Language locale should follow RFC 4646 format (e.g., `en-US`, `fr-FR`), combining ISO 639 language code and ISO 3166 country code.
- **User Purpose**: This field is read-only and managed by Microsoft 365 based on license assignment and user type.
- **Computed Fields**: Fields like `language.display_name` and `user_purpose` are computed and populated by the API.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import format: {user_id}
# The user_id can be either the user's object ID (UUID) or User Principal Name (UPN)

# Example 1: Import using user object ID (UUID)
terraform import microsoft365_graph_beta_users_user_mailbox_settings.example "12345678-1234-1234-1234-123456789012"

# Example 2: Import using User Principal Name (UPN)
terraform import microsoft365_graph_beta_users_user_mailbox_settings.example "john.doe@example.com"
```

