---
page_title: "microsoft365_graph_beta_device_and_app_management_win_get_app Resource - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Manages WinGet applications from the Microsoft Store using the /deviceAppManagement/mobileApps endpoint. WinGet apps enable deployment of Microsoft Store applications with automatic metadata population, streamlined package management, and integration with the Windows Package Manager for efficient app distribution to managed devices.
---

# microsoft365_graph_beta_device_and_app_management_win_get_app (Resource)

Manages WinGet applications from the Microsoft Store using the `/deviceAppManagement/mobileApps` endpoint. WinGet apps enable deployment of Microsoft Store applications with automatic metadata population, streamlined package management, and integration with the Windows Package Manager for efficient app distribution to managed devices.

## Microsoft Documentation

- [winGetApp resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-wingetapp?view=graph-rest-beta)
- [Create winGetApp](https://learn.microsoft.com/en-us/graph/api/intune-apps-wingetapp-create?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementApps.ReadWrite.All`

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "example_firefox" {
  package_identifier              = "9NZVDKPMR9RD" # Obtained from https://apps.microsoft.com/detail/9nzvdkpmr9rd?hl=en-US&gl=US
  automatically_generate_metadata = true

  # Optional metadata fields (will be auto-populated if automatically_generate_metadata = true)
  # display_name                  = "Mozilla Firefox"
  # description                   = "Choose the browser that prioritizes you, not their bottom line. Don't settle for the default browser. When you choose Firefox, you protect your data while supporting the non-profit"
  # publisher                     = "Mozilla"

  # Optional app information
  is_featured             = true
  privacy_information_url = "https://www.mozilla.org/en-US/privacy/firefox/"
  information_url         = "https://support.mozilla.org/en-US/"
  owner                   = "Workplace Services"
  developer               = "Mozilla Foundation"
  notes                   = "Default browser for all corporate devices"

  # Required install experience settings
  install_experience = {
    run_as_account = "user" # Allowed values: "system" or "user"
  }

  # Optional role scope tag IDs
  role_scope_tag_ids = ["8"]

  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.web_browser.id,
    "Business",
    "Productivity",
  ]

  # Optional timeouts
  timeouts = {
    create = "30s"
    update = "30s"
    read   = "30s"
    delete = "30s"
  }
}

# Example using the Microsoft Store Package Manifest Metadata datasource
data "microsoft365_utility_microsoft_store_package_manifest_metadata" "firefox_metadata" {
  package_identifier = "9NZVDKPMR9RD" # Firefox package ID

  timeouts = {
    read = "2m"
  }
}

# Application category resource for reference
resource "microsoft365_graph_beta_device_and_app_management_application_category" "web_browser" {
  display_name = "Web Browsers"
}

# Then use the metadata to create a winget app with manually specified metadata
resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "firefox_with_datasource" {
  package_identifier              = "9NZVDKPMR9RD"
  automatically_generate_metadata = false # We'll provide the metadata manually

  # Metadata fields populated from the datasource
  display_name = data.microsoft365_utility_microsoft_store_package_manifest_metadata.firefox_metadata.manifests[0].versions[0].default_locale.package_name
  description  = data.microsoft365_utility_microsoft_store_package_manifest_metadata.firefox_metadata.manifests[0].versions[0].default_locale.description
  publisher    = data.microsoft365_utility_microsoft_store_package_manifest_metadata.firefox_metadata.manifests[0].versions[0].default_locale.publisher

  # Additional information from the datasource
  privacy_information_url = data.microsoft365_utility_microsoft_store_package_manifest_metadata.firefox_metadata.manifests[0].versions[0].default_locale.privacy_url
  information_url         = data.microsoft365_utility_microsoft_store_package_manifest_metadata.firefox_metadata.manifests[0].versions[0].default_locale.publisher_support_url

  # Custom fields
  owner     = "IT Department"
  developer = "Mozilla Foundation"
  notes     = "Secondary browser for testing web applications"

  # Required install experience settings
  install_experience = {
    run_as_account = "system" # Different from the first example
  }

  # Optional role scope tag IDs
  role_scope_tag_ids = ["9"]

  # Categories based on the agreements in the metadata
  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.web_browser.id,
    "Business",
    "Productivity",
  ]

  # Optional timeouts
  timeouts = {
    create = "30s"
    update = "30s"
    read   = "30s"
    delete = "30s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `automatically_generate_metadata` (Boolean) When set to `true`, the provider will automatically fetch metadata from the Microsoft Store for Business using the package identifier. This will populate the `display_name`, `description`, `publisher`, and 'icon' fields.
- `install_experience` (Attributes) The install experience settings associated with this application.the value is idempotent and any changes to this field will trigger a recreation of the application. (see [below for nested schema](#nestedatt--install_experience))
- `package_identifier` (String) The **unique package identifier** for the WinGet/Microsoft Store app from the storefront.

For example, for the app Microsoft Edge Browser URL [https://apps.microsoft.com/detail/xpfftq037jwmhs?hl=en-us&gl=US](https://apps.microsoft.com/detail/xpfftq037jwmhs?hl=en-us&gl=US), the package identifier is `xpfftq037jwmhs`.

**Important notes:**
- This identifier is **required** at creation time.
- It **cannot be modified** for existing Terraform-deployed WinGet/Microsoft Store apps.

attempting to modify this value will result in a failed request.

### Optional

- `categories` (Set of String) Set of category names to associate with this application. You can use either thebpredefined Intune category names like 'Business', 'Productivity', etc., or provide specific category UUIDs. Predefined values include: 'Other apps', 'Books & Reference', 'Data management', 'Productivity', 'Business', 'Development & Design', 'Photos & Media', 'Collaboration & Social', 'Computer management'.
- `description` (String) A detailed description of the WinGet/ Microsoft Store for Business app.This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.
- `developer` (String) The developer of the app.
- `display_name` (String) The title of the WinGet app imported from the Microsoft Store for Business.This field value must match the expected title of the app in the Microsoft Store for Business associated with the `package_identifier`.This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.
- `information_url` (String) The more information Url.
- `is_featured` (Boolean) The value indicating whether the app is marked as featured by the admin. Default is false.
- `large_icon` (Attributes) The large icon for the WinGet app, automatically downloaded and set from the Microsoft Store for Business. (see [below for nested schema](#nestedatt--large_icon))
- `notes` (String) Notes for the app.
- `owner` (String) The owner of the app.
- `privacy_information_url` (String) The privacy statement Url.
- `publisher` (String) The publisher of the WinGet/ Microsoft Store for Business app.This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date and time the app was created. This property is read-only.
- `dependent_app_count` (Number) The total number of dependencies the child app has. This property is read-only.
- `id` (String) The unique identifier for this Intune Microsoft Store app
- `is_assigned` (Boolean) The value indicating whether the app is assigned to at least one group. This property is read-only.
- `last_modified_date_time` (String) The date and time the app was last modified. This property is read-only.
- `manifest_hash` (String) Hash of package metadata properties used to validate that the application matches the metadata in the source repository.
- `publishing_state` (String) The publishing state for the app. The app cannot be assigned unless the app is published. Possible values are: notPublished, processing, published.
- `superseded_app_count` (Number) The total number of apps this app is directly or indirectly superseded by. This property is read-only.
- `superseding_app_count` (Number) The total number of apps this app directly or indirectly supersedes. This property is read-only.
- `upload_state` (Number) The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.

<a id="nestedatt--install_experience"></a>
### Nested Schema for `install_experience`

Required:

- `run_as_account` (String) The account type (System or User) that actions should be run as on target devices.  Required at creation time.


<a id="nestedatt--large_icon"></a>
### Nested Schema for `large_icon`

Optional:

- `type` (String) The MIME type of the app's large icon, automatically populated based on the `package_identifier` when `automatically_generate_metadata` is true. Example: `image/png`
- `value` (String, Sensitive) The icon image to use for the winget app. This field is automatically populated based on the `package_identifier` when `automatically_generate_metadata` is set to true.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Windows Specific**: This resource is specifically for managing WinGet applications on Windows devices.
- **Windows Package Manager**: WinGet is Microsoft's official package manager for Windows, providing access to thousands of applications.
- **Package Source**: Apps are sourced from the Windows Package Manager Community Repository or Microsoft Store.
- **Assignment Required**: Apps must be assigned to user or device groups to be deployed through Intune.
- **Package Identifier**: Uses package identifiers from the WinGet repository (e.g., `Microsoft.PowerToys`).
- **Automatic Updates**: WinGet apps can be configured for automatic updates through the Windows Package Manager.
- **Installation Context**: Apps can be installed in user or system context depending on the package configuration.
- **Version Management**: Specific versions can be targeted, or the latest version can be automatically selected.

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_management_win_get_app.example win-get-app-id
```

