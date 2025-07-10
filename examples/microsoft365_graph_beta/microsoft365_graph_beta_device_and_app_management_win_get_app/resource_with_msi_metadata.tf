# Example showing how to use both MSI metadata and WinGet metadata together
# First, extract metadata from a local MSI file
data "microsoft365_utility_windows_msi_app_metadata" "firefox_msi" {
  # This would be a path to a Firefox MSI installer
  installer_file_path_source = "C:/path/to/Firefox_Setup_115.0.msi"

  timeouts = {
    read = "4m"
  }
}

# Then, get the WinGet package metadata for the same application
data "microsoft365_utility_microsoft_store_package_manifest_metadata" "firefox_winget" {
  package_identifier = "Mozilla.Firefox" # Firefox package ID in WinGet

  timeouts = {
    read = "2m"
  }
}

# Create a WinGet app resource using both metadata sources
resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "firefox_with_metadata" {
  # Use the WinGet package identifier
  package_identifier = data.microsoft365_utility_microsoft_store_package_manifest_metadata.firefox_winget.manifests[0].package_identifier

  # Disable automatic metadata generation since we're providing our own
  automatically_generate_metadata = false

  # Use metadata from both sources
  display_name = data.microsoft365_utility_windows_msi_app_metadata.firefox_msi.metadata.product_name
  description  = data.microsoft365_utility_microsoft_store_package_manifest_metadata.firefox_winget.manifests[0].default_locale.short_description
  publisher    = data.microsoft365_utility_windows_msi_app_metadata.firefox_msi.metadata.publisher

  # Use WinGet metadata for URLs
  information_url         = data.microsoft365_utility_microsoft_store_package_manifest_metadata.firefox_winget.manifests[0].default_locale.public_website_url
  privacy_information_url = "https://www.mozilla.org/privacy/firefox/"

  # Additional metadata
  developer = "Mozilla Corporation"
  owner     = "IT Department"
  notes     = "Deployed using combined metadata from MSI and WinGet"

  # Required install experience settings
  install_experience = {
    run_as_account = "user" # Allowed values: "system" or "user"
  }

  # Add to the Web Browsers category
  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.web_browser.id,
    "Business",
  ]

  # Optional timeouts
  timeouts = {
    create = "30s"
    update = "30s"
    read   = "30s"
    delete = "30s"
  }
}

# Application category resource for reference
resource "microsoft365_graph_beta_device_and_app_management_application_category" "web_browser" {
  display_name = "Web Browsers"
} 