---
page_title: "microsoft365_graph_beta_m365_admin_sharepoint_settings Resource - terraform-provider-microsoft365"
subcategory: "M365 Admin"
description: |-
  Manages the tenant-wide SharePoint and OneDrive settings using the Microsoft Graph beta /admin/sharepoint/settings endpoint. The settings control external sharing, site creation, OneDrive storage and retention, the OneDrive sync app and idle session sign-out.
  This is a singleton resource: one settings object exists per tenant. create adopts the existing settings and create/update apply the configuration with PATCH. Settings that are not configured keep their current tenant values and are not sent to Graph. On destroy, Terraform only removes the resource from state and leaves the tenant settings as they are. Manage only one instance per tenant across all Terraform states.
---

# microsoft365_graph_beta_m365_admin_sharepoint_settings (Resource)

Manages the tenant-wide SharePoint and OneDrive settings using the Microsoft Graph beta `/admin/sharepoint/settings` endpoint. The settings control external sharing, site creation, OneDrive storage and retention, the OneDrive sync app and idle session sign-out.

This is a **singleton resource**: one settings object exists per tenant. `create` adopts the existing settings and `create`/`update` apply the configuration with PATCH. Settings that are not configured keep their current tenant values and are not sent to Graph. On `destroy`, Terraform only removes the resource from state and leaves the tenant settings as they are. Manage only one instance per tenant across all Terraform states.

## Microsoft Documentation

- [SharePoint admin center settings](https://learn.microsoft.com/en-us/sharepoint/sharepoint-admin-role)
- [sharepointSettings resource type](https://learn.microsoft.com/en-us/graph/api/resources/sharepointsettings?view=graph-rest-beta)
- [Get sharepointSettings](https://learn.microsoft.com/en-us/graph/api/sharepointsettings-get?view=graph-rest-beta)
- [Update sharepointSettings](https://learn.microsoft.com/en-us/graph/api/sharepointsettings-update?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `SharePointTenantSettings.Read.All`
- `SharePointTenantSettings.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v1.7.0 | Experimental | Initial release |

## Example Usage

### Minimal

```terraform
# Keep a deleted user's OneDrive for one year and hide the OneDrive sync button.
# Every other setting keeps its current tenant value.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "example" {
  deleted_user_personal_site_retention_period_in_days = 365
  is_sync_button_hidden_on_personal_site              = true
}
```

### Maximal

```terraform
# Manage every SharePoint and OneDrive tenant setting.
# Review each value before applying: sharing, legacy authentication and sync app
# restrictions affect every user in the tenant.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "example" {
  # OneDrive sync app
  allowed_domain_guids_for_sync_app           = ["bdd1ab9b-3fd0-4def-a761-ec8d7471732c"]
  excluded_file_extensions_for_sync_app       = ["*.pst", "*.mp3"]
  is_mac_sync_app_enabled                     = true
  is_sync_button_hidden_on_personal_site      = false
  is_unmanaged_sync_app_for_tenant_restricted = true

  # OneDrive storage and retention
  deleted_user_personal_site_retention_period_in_days = 365
  personal_site_default_storage_limit_in_mb           = 1048576

  # External sharing
  sharing_capability                                      = "externalUserSharingOnly"
  sharing_domain_restriction_mode                         = "allowList"
  sharing_allowed_domain_list                             = ["contoso.com", "fabrikam.com"]
  sharing_blocked_domain_list                             = []
  is_resharing_by_external_users_enabled                  = false
  is_require_accepting_user_to_match_invited_user_enabled = true

  # Sites and pages
  is_site_creation_enabled                  = true
  is_site_creation_ui_enabled               = true
  is_site_pages_creation_enabled            = true
  is_sites_storage_limit_automatic          = true
  site_creation_default_managed_path        = "/sites/"
  site_creation_default_storage_limit_in_mb = 26214400
  tenant_default_timezone                   = "(UTC) Coordinated Universal Time"
  is_commenting_on_site_pages_enabled       = true
  is_sharepoint_newsfeed_enabled            = false
  image_tagging_option                      = "basic"

  # Notifications, Loop and authentication
  is_file_activity_notification_enabled     = true
  is_sharepoint_mobile_notification_enabled = true
  is_loop_enabled                           = true
  is_legacy_auth_protocols_enabled          = false

  # Sign out users after one hour of inactivity, with a warning five minutes before.
  idle_session_sign_out = {
    is_enabled                = true
    warn_after_in_seconds     = 3300
    sign_out_after_in_seconds = 3600
  }

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `allowed_domain_guids_for_sync_app` (Set of String) The trusted domain GUIDs for the OneDrive sync app, in lower case. Used with `is_unmanaged_sync_app_for_tenant_restricted`. Set to an empty set to clear the list. When not configured, the current tenant value is kept.
- `deleted_user_personal_site_retention_period_in_days` (Number) The number of days a deleted user's OneDrive is kept. Must be between `30` and `3650`. When not configured, the current tenant value is kept.
- `excluded_file_extensions_for_sync_app` (Set of String) The file extensions that the OneDrive sync app does not upload, for example `*.pst`. Set to an empty set to clear the list. When not configured, the current tenant value is kept.
- `idle_session_sign_out` (Attributes) The idle session sign-out policy for SharePoint and OneDrive. When not configured, the current tenant value is kept. (see [below for nested schema](#nestedatt--idle_session_sign_out))
- `image_tagging_option` (String) The image tagging option for the tenant. Possible values are `disabled`, `basic` and `enhanced`. When not configured, the current tenant value is kept.
- `is_commenting_on_site_pages_enabled` (Boolean) Whether comments are allowed on modern site pages in SharePoint. When not configured, the current tenant value is kept.
- `is_file_activity_notification_enabled` (Boolean) Whether push notifications are enabled for OneDrive events. When not configured, the current tenant value is kept.
- `is_legacy_auth_protocols_enabled` (Boolean) Whether legacy authentication protocols are enabled for the tenant. Enabling legacy authentication weakens conditional access enforcement for SharePoint and OneDrive. When not configured, the current tenant value is kept.
- `is_loop_enabled` (Boolean) Whether the Fluid Framework (Loop components) is allowed on SharePoint sites. When not configured, the current tenant value is kept.
- `is_mac_sync_app_enabled` (Boolean) Whether files can be synced using the OneDrive sync app for Mac. When not configured, the current tenant value is kept.
- `is_require_accepting_user_to_match_invited_user_enabled` (Boolean) Whether guests must sign in using the same account to which sharing invitations are sent. Graph may return no value for this setting, in which case it is null in state. When not configured, the current tenant value is kept.
- `is_resharing_by_external_users_enabled` (Boolean) Whether guests are allowed to reshare files, folders and sites they don't own. When not configured, the current tenant value is kept.
- `is_sharepoint_mobile_notification_enabled` (Boolean) Whether mobile push notifications are enabled for SharePoint. When not configured, the current tenant value is kept.
- `is_sharepoint_newsfeed_enabled` (Boolean) Whether the newsfeed is allowed on modern site pages in SharePoint. When not configured, the current tenant value is kept.
- `is_site_creation_enabled` (Boolean) Whether users are allowed to create sites. When not configured, the current tenant value is kept.
- `is_site_creation_ui_enabled` (Boolean) Whether the UI commands for creating sites are shown. When not configured, the current tenant value is kept.
- `is_site_pages_creation_enabled` (Boolean) Whether creating new modern pages is allowed on SharePoint sites. When not configured, the current tenant value is kept.
- `is_sites_storage_limit_automatic` (Boolean) Whether site storage space is managed automatically (`true`) or with specific storage limits per site (`false`). When not configured, the current tenant value is kept.
- `is_sync_button_hidden_on_personal_site` (Boolean) Whether the sync button in OneDrive is hidden. When not configured, the current tenant value is kept.
- `is_unmanaged_sync_app_for_tenant_restricted` (Boolean) Whether users are allowed to sync files only on PCs joined to specific domains (see `allowed_domain_guids_for_sync_app`). When not configured, the current tenant value is kept.
- `personal_site_default_storage_limit_in_mb` (Number) The default OneDrive storage limit in megabytes for all new and existing users who are assigned a qualifying licence. When not configured, the current tenant value is kept.
- `sharing_allowed_domain_list` (Set of String) The email domains that are allowed for sharing outside the organization. Applies when `sharing_domain_restriction_mode` is `allowList`. Set to an empty set to clear the list. When not configured, the current tenant value is kept.
- `sharing_blocked_domain_list` (Set of String) The email domains that are blocked for sharing outside the organization. Applies when `sharing_domain_restriction_mode` is `blockList`. Set to an empty set to clear the list. When not configured, the current tenant value is kept.
- `sharing_capability` (String) The external sharing capability for the tenant. Possible values are `disabled` (only people in the organization), `existingExternalUserSharingOnly` (existing guests), `externalUserSharingOnly` (new and existing guests) and `externalUserAndGuestSharing` (anyone, including links that don't require sign-in). When not configured, the current tenant value is kept.
- `sharing_domain_restriction_mode` (String) The external sharing mode for domains. Possible values are `none`, `allowList` and `blockList`. When not configured, the current tenant value is kept.
- `site_creation_default_managed_path` (String) The managed path under which new team sites are created, for example `/sites/`. It must be one of `available_managed_paths_for_site_creation`. When not configured, the current tenant value is kept.
- `site_creation_default_storage_limit_in_mb` (Number) The default storage quota in megabytes for a new site. When not configured, the current tenant value is kept.
- `tenant_default_timezone` (String) The default time zone for newly created sites, for example `(UTC-08:00) Pacific Time (US and Canada)`. See [SPRegionalSettings.TimeZones](https://learn.microsoft.com/en-us/sharepoint/dev/schema/regional-settings-schema) for the possible values. When not configured, the current tenant value is kept.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `available_managed_paths_for_site_creation` (Set of String) The managed paths available for site creation, for example `/sites/` and `/teams/`. Read-only.
- `id` (String) The fixed singleton identifier `sharepointSettings`.

<a id="nestedatt--idle_session_sign_out"></a>
### Nested Schema for `idle_session_sign_out`

Optional:

- `is_enabled` (Boolean) Whether the idle session sign-out policy is enabled. When not configured, the current tenant value is kept.
- `sign_out_after_in_seconds` (Number) The number of seconds of inactivity after which a user is signed out. When not configured, the current tenant value is kept.
- `warn_after_in_seconds` (Number) The number of seconds of inactivity after which a user is warned that they will be signed out. Must be less than `sign_out_after_in_seconds`. When not configured, the current tenant value is kept.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

### Singleton Resource
- One SharePoint settings object exists per tenant. It cannot be created or deleted.
- `create` adopts the existing settings and applies the configured values with `PATCH`.
- `destroy` only removes the resource from Terraform state; the tenant settings keep their current values. Nothing is reset to a default.
- Manage only one instance per tenant across all Terraform states.

### Unconfigured Values
- Every setting is optional. A setting that is not configured keeps its current tenant value and is not sent to Graph.
- When a configuration stops managing a setting, the value that was last applied stays in place.
- Collections (`allowed_domain_guids_for_sync_app`, `excluded_file_extensions_for_sync_app`, `sharing_allowed_domain_list` and `sharing_blocked_domain_list`) are only sent when configured. Set a collection to `[]` to clear it.
- Inside `idle_session_sign_out`, only the configured values are sent.
- `available_managed_paths_for_site_creation` is read-only.

### High-Impact Settings
- `sharing_capability`, the sharing domain lists, `is_resharing_by_external_users_enabled` and `is_require_accepting_user_to_match_invited_user_enabled` control what can be shared outside the organization.
- `is_legacy_auth_protocols_enabled` allows clients that cannot enforce conditional access.
- `is_unmanaged_sync_app_for_tenant_restricted` with `allowed_domain_guids_for_sync_app` limits which PCs can sync OneDrive files.
- `idle_session_sign_out` signs users out of SharePoint and OneDrive on unmanaged devices after inactivity.
- Review these settings carefully, because they apply to every user in the tenant.

### Value Formats
- `allowed_domain_guids_for_sync_app` must contain lower-case GUIDs, as Graph returns them.
- `deleted_user_personal_site_retention_period_in_days` must be between `30` and `3650`.
- `idle_session_sign_out.warn_after_in_seconds` must be less than `sign_out_after_in_seconds`.
- `is_require_accepting_user_to_match_invited_user_enabled` is not returned by Graph in some tenants; it is then null in state until it is configured.

## Import

Import is supported using the following syntax:

```shell
# Import the tenant-wide SharePoint and OneDrive settings using their fixed singleton ID.
terraform import microsoft365_graph_beta_m365_admin_sharepoint_settings.example sharepointSettings
```
