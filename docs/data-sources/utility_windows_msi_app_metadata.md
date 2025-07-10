---
page_title: "microsoft365_utility_windows_msi_app_metadata Data Source - terraform-provider-microsoft365"
subcategory: "Utilities"
description: |-
  Use this data source to query the iTunes App Store API for app metadata.
---

# microsoft365_utility_windows_msi_app_metadata

Use this data source to extract metadata from a Windows MSI installer file.
This data source allows you to extract metadata from a Windows MSI installer file,
returning details like product name, version, and other metadata.

## Example Usage

```terraform
```terraform
# # Example 1: Extract metadata from a local MSI file
# data "microsoft365_utility_windows_msi_app_metadata" "example_file_path_msi" {
#   installer_file_path_source = "C:/path/to/your/application.msi"

#   timeouts = {
#     read = "4m"
#   }
# }

# Example 2: Extract metadata from an MSI file at a URL
data "microsoft365_utility_windows_msi_app_metadata" "example_url_msi" {
  installer_url_source = "https://download.mozilla.org/?product=firefox-msi-latest-ssl&os=win64&lang=en-US"

  timeouts = {
    read = "5m" # Longer timeout for download and processing
  }
}

# Output examples showing how to access the extracted metadata
output "product_name" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.product_name
}

output "product_version" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.product_version
}

output "product_code" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.product_code
}

output "publisher" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.publisher
}

output "upgrade_code" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.upgrade_code
}

output "language" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.language
}

output "package_type" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.package_type
}

output "install_location" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.install_location
}

output "install_command" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.install_command
}
output "uninstall_command" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.uninstall_command
}

output "transform_paths" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.transform_paths
}

output "size_mb" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.size_mb
}

output "sha256_checksum" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.sha256_checksum
}

output "md5_checksum" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.md5_checksum
}

output "properties" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.properties
}

output "required_features" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.required_features
}

output "files" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.files
}

output "min_os_version" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.min_os_version
}

output "architecture" {
  value = data.microsoft365_utility_windows_msi_app_metadata.example_url_msi.metadata.architecture
}
```
```

## Argument Reference

* `installer_file_path_source` - (Required) The path to the Windows MSI installer file.
* `installer_url_source` - (Required) The URL to the Windows MSI installer file.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The ID of this resource.
* `metadata` - A list of app results returned from the iTunes App Store API. Each result contains:
  * `product_name` - The name of the app.
  * `product_version` - The version of the app.
  * `product_code` - The product code of the app.
  * `publisher` - The publisher of the app.
  * `upgrade_code` - The upgrade code of the app.

  * `description` - The description of the app.
  * `version` - The version of the app.
  * `price` - The price of the app in the local currency.
  * `formatted_price` - The formatted price of the app (e.g., 'Free', '$0.99').
  * `release_date` - The release date of the app.
  * `average_rating` - The average user rating of the app.
  * `artist_name` - The name of the artist/developer.
  * `minimum_os_version` - The minimum OS version required to run the app.
  * `content_advisory_rating` - The content advisory rating (e.g., '4+', '12+', '17+').
  * `is_vpp_device_based_licensed` - Whether the app supports VPP device-based licensing.
  * `release_notes` - Notes about the latest release of the app.
  * `currency` - The currency code for the price (e.g., 'USD', 'GBP', 'EUR').
  * `user_rating_count` - The number of user ratings for the app. 