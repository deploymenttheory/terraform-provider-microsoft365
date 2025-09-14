resource "random_string" "office_suite_app_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Example 1: Office Suite App with Configuration Designer (individual settings)
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_config_designer" {
  display_name       = "acc-test-office-suite-app-${random_string.office_suite_app_suffix.result}"
  description        = "Microsoft 365 Apps deployed with Configuration Designer settings"
  is_featured        = true
  information_url    = "https://support.microsoft.com/office"
  notes              = "Microsoft 365 Apps configured with specific settings using Configuration Designer approach."
  role_scope_tag_ids = ["0"] # Default role scope tag

  categories = [
    "Business",
    "Productivity",
  ]

  # Configuration Designer block - use this for individual configuration settings
  configuration_designer = {
    auto_accept_eula = true

    excluded_apps = {
      access               = true  # Exclude Microsoft Access
      bing                 = false # Include Microsoft Search in Bing
      excel                = false # Include Excel
      groove               = true  # Exclude OneDrive for Business (Groove)
      info_path            = true  # Exclude InfoPath
      lync                 = false # Include Skype for Business
      one_drive            = false # Include OneDrive
      one_note             = false # Include OneNote
      outlook              = false # Include Outlook
      power_point          = false # Include PowerPoint
      publisher            = true  # Exclude Publisher
      share_point_designer = true  # Exclude SharePoint Designer
      teams                = false # Include Teams
      visio                = true  # Exclude Visio
      word                 = false # Include Word
    }

    locales_to_install = [
      "en-us", # English (United States)
      "fr-fr", # French (France)
      "de-de", # German (Germany)
    ]

    office_platform_architecture         = "x64"
    office_suite_app_default_file_format = "officeOpenXMLFormat"

    product_ids = [
      "o365ProPlusRetail"
    ]

    should_uninstall_older_versions_of_office = true
    target_version                            = "16.0.19029.20244"
    update_channel                            = "current"
    update_version                            = "" // for latest version, use empty string
    use_shared_computer_activation            = false
  }

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}