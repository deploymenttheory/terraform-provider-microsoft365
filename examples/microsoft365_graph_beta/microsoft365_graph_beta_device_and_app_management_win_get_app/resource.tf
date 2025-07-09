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

