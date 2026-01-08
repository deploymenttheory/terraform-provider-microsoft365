---
page_title: "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy Resource - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Manages iOS managed device mobile app configurations in Microsoft Intune using the /deviceAppManagement/mobileAppConfigurations endpoint.Use app configuration policies in Microsoft Intune to provide custom configuration settings for an iOS/iPadOS app. These configuration settings allow an app  to be customized based on the app suppliers direction. You must get these configuration settings (keys and values) from the supplier of the app. To configure the app, you specify the settings as keys and values, or as XML containing the keys and values. Learn more here: https://learn.microsoft.com/en-us/intune/intune-service/apps/app-configuration-policies-use-ios
---

# microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy (Resource)

Manages iOS managed device mobile app configurations in Microsoft Intune using the `/deviceAppManagement/mobileAppConfigurations` endpoint.Use app configuration policies in Microsoft Intune to provide custom configuration settings for an iOS/iPadOS app. These configuration settings allow an app  to be customized based on the app suppliers direction. You must get these configuration settings (keys and values) from the supplier of the app. To configure the app, you specify the settings as keys and values, or as XML containing the keys and values. Learn more here: https://learn.microsoft.com/en-us/intune/intune-service/apps/app-configuration-policies-use-ios

## Microsoft Documentation

- [iosMobileAppConfiguration resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosmobileappconfiguration?view=graph-rest-beta)
- [Create iosMobileAppConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-apps-iosmobileappconfiguration-create?view=graph-rest-beta)
- [Update iosMobileAppConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-apps-iosmobileappconfiguration-update?view=graph-rest-beta)
- [Delete iosMobileAppConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-apps-iosmobileappconfiguration-delete?view=graph-rest-beta)
- [App Configuration Policies for Managed iOS/iPadOS Devices](https://learn.microsoft.com/en-us/intune/intune-service/apps/app-configuration-policies-use-ios)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`, `DeviceManagementApps.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.36.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example resource with XML encoded settings
resource "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy" "test_xml_encoded" {
  display_name         = "Test Acceptance iOS Mobile App Configuration - xml encoded"
  description          = "Updated description for acceptance testing"
  targeted_mobile_apps = [data.microsoft365_graph_beta_device_and_app_management_mobile_app.company_portal.items[0].id]
  role_scope_tag_ids   = ["0"]

  encoded_setting_xml = <<-XML
    <dict>
      <key>metadata</key>
      <dict>
          <key>version</key>
          <string>1.0</string>
          <key>created</key>
          <string>2025-10-14</string>
          <key>author</key>
          <string>System</string>
      </dict>

      <key>items</key>
      <array>
          <dict>
              <key>id</key>
              <string>001</string>
              <key>category</key>
              <string>electronics</string>
              <key>name</key>
              <string>Wireless Mouse</string>
              <key>description</key>
              <string>Ergonomic wireless mouse with USB receiver</string>
              <key>price</key>
              <real>29.99</real>
              <key>stock</key>
              <integer>150</integer>
              <key>specifications</key>
              <dict>
                  <key>battery</key>
                  <string>AA batteries</string>
                  <key>range</key>
                  <string>10 meters</string>
                  <key>color</key>
                  <string>Black</string>
              </dict>
          </dict>

          <dict>
              <key>id</key>
              <string>002</string>
              <key>category</key>
              <string>books</string>
              <key>name</key>
              <string>The Art of Programming</string>
              <key>description</key>
              <string>A comprehensive guide to software development</string>
              <key>price</key>
              <real>49.99</real>
              <key>stock</key>
              <integer>75</integer>
              <key>specifications</key>
              <dict>
                  <key>pages</key>
                  <integer>500</integer>
                  <key>isbn</key>
                  <string>978-1234567890</string>
                  <key>format</key>
                  <string>Hardcover</string>
              </dict>
          </dict>
      </array>
    </dict>
  XML

}

# Example resource with customsettings
resource "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy" "test_settings" {
  display_name         = "Test Acceptance iOS Mobile App Configuration - settings"
  description          = "Updated description for acceptance testing"
  targeted_mobile_apps = [data.microsoft365_graph_beta_device_and_app_management_mobile_app.company_portal.items[0].id]
  role_scope_tag_ids   = ["0"]

  settings = [
    {
      app_config_key       = "testKey1"
      app_config_key_type  = "stringType"
      app_config_key_value = "testValue1"
    },
    {
      app_config_key       = "testKey2"
      app_config_key_type  = "integerType"
      app_config_key_value = "123"
    }
  ]
}

# Data source to find specific iOS app by display name (e.g., "Company Portal")
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "company_portal" {
  filter_type     = "display_name"
  filter_value    = "Microsoft Intune Company Portal"
  app_type_filter = "iosStoreApp" # Only search iOS store apps
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name of the iOS mobile app configuration

### Optional

- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `encoded_setting_xml` (String, Sensitive) MDM app configuration in Base64 encoded format. **Note:** This field is mutually exclusive with `settings` - only one can be specified.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this iOS mobile app configuration.
- `settings` (Attributes Set) Configuration setting items for the iOS mobile app. **Note:** This field is mutually exclusive with `encoded_setting_xml` - only one can be specified. (see [below for nested schema](#nestedatt--settings))
- `targeted_mobile_apps` (Set of String) Set of mobile app IDs that this configuration targets.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for this iOS mobile app configuration
- `version` (Number) Version of the iOS mobile app configuration.

<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Required:

- `app_config_key` (String) The configuration key name
- `app_config_key_type` (String) The configuration key type (e.g., "stringType", "integerType", "realType", "booleanType", "tokenType")
- `app_config_key_value` (String) The configuration key value


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using ID format: {id}
terraform import microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy.example 00000000-0000-0000-0000-000000000000
```