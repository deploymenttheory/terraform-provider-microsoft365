---
page_title: "Resource: microsoft365_graph_beta_device_and_app_management_m365_apps_installation_options"
description: |-
    
---

# Resource: microsoft365_graph_beta_device_and_app_management_m365_apps_installation_options





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `apps_for_mac` (Attributes) The Microsoft 365 apps installation options container object for a MAC platform. (see [below for nested schema](#nestedatt--apps_for_mac))
- `apps_for_windows` (Attributes) The Microsoft 365 apps installation options container object for a Windows platform. (see [below for nested schema](#nestedatt--apps_for_windows))
- `update_channel` (String) Specifies how often users get feature updates for Microsoft 365 apps installed on devices running Windows. The possible values are: `current`, `monthlyEnterprise`, or `semiAnnual`, with corresponding update frequencies of `As soon as they're ready`, `Once a month`, and `Every six months`.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for the M365AppsInstallationOptions.

<a id="nestedatt--apps_for_mac"></a>
### Nested Schema for `apps_for_mac`

Required:

- `is_microsoft_365_apps_enabled` (Boolean) Specifies whether users can install Microsoft 365 apps on their MAC devices. The default value is `true`.
- `is_skype_for_business_enabled` (Boolean) Specifies whether users can install Skype for Business on their MAC devices running OS X El Capitan 10.11 or later. The default value is `true`.


<a id="nestedatt--apps_for_windows"></a>
### Nested Schema for `apps_for_windows`

Required:

- `is_microsoft_365_apps_enabled` (Boolean) Specifies whether users can install Microsoft 365 apps, including Skype for Business, on their Windows devices. The default value is `true`.
- `is_project_enabled` (Boolean) Specifies whether users can install Microsoft Project on their Windows devices. The default value is `true`.
- `is_skype_for_business_enabled` (Boolean) Specifies whether users can install Skype for Business (standalone) on their Windows devices. The default value is `true`.
- `is_visio_enabled` (Boolean) Specifies whether users can install Visio on their Windows devices. The default value is `true`.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

