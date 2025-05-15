# Extract metadata from a local PKG file
data "microsoft365_utility_macos_pkg_app_metadata" "local_example" {
  installer_file_path_source = "/path/to/your/application.pkg"
  
  timeouts = {
    read = "4m"  # Extended timeout for large PKG files
  }
}

# Output all extracted metadata fields
output "pkg_metadata" {
  description = "All metadata extracted from the PKG file"
  value = {
    # Basic app identification
    bundle_id = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.cf_bundle_identifier
    version   = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.cf_bundle_short_version_string
    name      = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.name
    
    # Installation details
    install_location = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.install_location
    package_ids      = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.package_ids
    app_paths        = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.app_paths
    min_os_version   = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.min_os_version
    
    # File information
    size_mb        = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.size_mb
    md5_checksum   = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.md5_checksum
    sha256_checksum = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.sha256_checksum
    
    # Bundles information
    included_bundles = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.included_bundles
  }
}

# Individual outputs for easy reference
output "bundle_id" {
  description = "The bundle identifier of the PKG application"
  value       = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.cf_bundle_identifier
}

output "version" {
  description = "The version of the PKG application"
  value       = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.cf_bundle_short_version_string
}

output "name" {
  description = "The name of the PKG application"
  value       = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.name
}

output "size_mb" {
  description = "The size of the PKG file in megabytes"
  value       = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.size_mb
}

output "sha256_checksum" {
  description = "SHA256 hash of the PKG file content"
  value       = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.sha256_checksum
}

output "md5_checksum" {
  description = "MD5 hash of the PKG file content"
  value       = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.md5_checksum
}

# Example of using the extracted metadata with a macOS PKG app resource
resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "firefox" {
  display_name    = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.name
  description     = "Firefox browser deployed via Intune"
  publisher       = "Mozilla"
  
  # Use the extracted metadata directly
  bundle_id       = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.cf_bundle_identifier
  version         = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.cf_bundle_short_version_string
  
  # Use the same file for the actual deployment
  installer_file_path_source = "/Users/dafyddwatkins/Downloads/Firefox 134.0.pkg"
  
  # Optional: Use other extracted metadata
  minimum_os_version = data.microsoft365_utility_macos_pkg_app_metadata.local_example.metadata.min_os_version
}