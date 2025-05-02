---
page_title: "microsoft365_graph_beta_device_and_app_management_macos_pkg_app Data Source - terraform-provider-microsoft365"
subcategory: "Intune"
description: |-
  Manages an Intune macOS app (PKG), using the mobileapps graph beta API. Apps are deployed using the Microsoft Intune management agent for macOS.
---

# microsoft365_graph_beta_device_and_app_management_macos_pkg_app (Data Source)

The Microsoft 365 Intune macos pkg app data source provides information about a specific macos pkg app.

## Example Usage

```terraform
data "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "mozilla_firefox" {
  # You can use either id or display_name or id to fetch the app
  display_name = "Firefox 136.0.pkg"
  //id = "824024fd-b7d0-4e8a-b53d-980633235765"

  # Optional: Add timeouts
  timeouts = {
    read = "10s"
  }
}

output "firefox_app_details" {
  value = {
    // Basic app information
    id                = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.id
    display_name      = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.display_name
    description       = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.description
    publisher         = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.publisher
    is_featured       = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.is_featured
    owner             = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.owner
    developer         = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.developer
    notes             = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.notes
    created_date_time = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.created_date_time

    // URLs
    privacy_information_url = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.privacy_information_url
    information_url         = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.information_url

    // Status information
    upload_state     = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.upload_state
    publishing_state = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.publishing_state
    is_assigned      = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.is_assigned

    // Related counts
    dependent_app_count   = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.dependent_app_count
    superseding_app_count = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.superseding_app_count
    superseded_app_count  = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.superseded_app_count

    // Tags and categories
    role_scope_tag_ids = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.role_scope_tag_ids
    categories         = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.categories

    // macOS specific configuration
    macos_pkg_app = {
      // Version detection
      ignore_version_detection = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.ignore_version_detection

      // Bundle identifiers
      primary_bundle_id      = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.primary_bundle_id
      primary_bundle_version = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.primary_bundle_version

      // Included apps
      included_apps = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.included_apps

      // Minimum OS requirements
      minimum_supported_operating_system = {
        v10_7  = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v10_7
        v10_8  = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v10_8
        v10_9  = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v10_9
        v10_10 = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v10_10
        v10_11 = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v10_11
        v10_12 = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v10_12
        v10_13 = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v10_13
        v10_14 = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v10_14
        v10_15 = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v10_15
        v11_0  = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v11_0
        v12_0  = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v12_0
        v13_0  = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v13_0
        v14_0  = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.minimum_supported_operating_system.v14_0
      }

      // Installation scripts
      pre_install_script = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.pre_install_script != null ? {
        script_content = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.pre_install_script.script_content
      } : null

      post_install_script = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.post_install_script != null ? {
        script_content = data.microsoft365_graph_beta_device_and_app_management_macos_pkg_app.mozilla_firefox.macos_pkg_app.post_install_script.script_content
      } : null
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) The title of the Intune macOS pkg application.
- `id` (String) The unique identifier of the macOS PKG.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date and time the app was created. This property is read-only.
- `description` (String) A detailed description of the Intune macOS pkg application.
- `role_scope_tag_ids` (Set of String) Set of scope tag ids for this mobile app.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).