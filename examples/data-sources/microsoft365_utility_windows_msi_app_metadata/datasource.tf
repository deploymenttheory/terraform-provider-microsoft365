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
