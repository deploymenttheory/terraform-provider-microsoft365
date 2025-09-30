# Example 1: Office Suite App with Configuration Designer (individual settings)
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_config_designer" {
  display_name       = "Microsoft 365 Apps for Enterprise - Configuration Designer"
  description        = "Microsoft 365 Apps deployed with Configuration Designer settings"
  is_featured        = true
  information_url    = "https://support.microsoft.com/office"
  notes              = "Microsoft 365 Apps configured with specific settings using Configuration Designer approach."
  role_scope_tag_ids = ["0"] # Default role scope tag

  categories = [
    "Business",
    "Productivity",
  ]

  app_icon = {
    web_url_source = "https://your-website.com/office-icon.png"
    //or
    icon_file_path_source = "/path/to/office365-icon.png"
  }

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

# Example 2: Office Suite App with XML Configuration
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_xml_config" {
  display_name       = "Microsoft 365 Apps for Enterprise - XML Configuration"
  description        = "Microsoft 365 Apps deployed with XML configuration file"
  is_featured        = true
  information_url    = "https://support.microsoft.com/office"
  notes              = "Microsoft 365 Apps configured using Office Deployment Tool (ODT) XML configuration."
  role_scope_tag_ids = ["0"] # Default role scope tag

  categories = [
    "Business",
    "Productivity",
  ]

  app_icon = {
    web_url_source = "https://your-website.com/office-icon.png"
    //or
    icon_file_path_source = "/path/to/office365-icon.png"
  }

  # XML Configuration block - use this for ODT XML-based configuration
  xml_configuration = {
    office_configuration_xml = <<EOF
<Configuration>
  <Add SourcePath="\\Server\Share" 
      OfficeClientEdition="64"
    Channel="MonthlyEnterprise" >
    <Product ID="O365ProPlusRetail">
      <Language ID="en-us" />
      <Language ID="ja-jp" />
    </Product>
    <Product ID="VisioProRetail">
      <Language ID="en-us" />
      <Language ID="ja-jp" />
    </Product>
  </Add>
  <Updates Enabled="TRUE" 
           UpdatePath="\\Server\Share" />
   <Display Level="None" AcceptEULA="TRUE" />  
</Configuration>
EOF
  }

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example 3: Office Suite App with Business Apps (O365BusinessRetail)
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_business" {
  display_name       = "Microsoft 365 Apps for Business"
  description        = "Microsoft 365 Apps for Business with basic configuration"
  is_featured        = false
  information_url    = "https://support.microsoft.com/office"
  notes              = "Basic Microsoft 365 Apps for Business deployment."
  role_scope_tag_ids = ["0"]

  categories = [
    "Business",
    "Productivity",
  ]

  configuration_designer = {
    auto_accept_eula = true

    excluded_apps = {
      access               = true # Exclude Access (not available in Business)
      bing                 = false
      excel                = false
      groove               = false
      info_path            = true # Exclude InfoPath
      lync                 = false
      one_drive            = false
      one_note             = false
      outlook              = false
      power_point          = false
      publisher            = true # Exclude Publisher
      share_point_designer = true
      teams                = false
      visio                = true # Exclude Visio (not available in Business)
      word                 = false
    }

    locales_to_install                   = ["en-us"]
    office_platform_architecture         = "x64"
    office_suite_app_default_file_format = "officeOpenXMLFormat"

    product_ids = [
      "o365BusinessRetail"
    ]

    should_uninstall_older_versions_of_office = false
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

# Example 4: Office Suite App with Project and Visio
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_with_project_visio" {
  display_name       = "Microsoft 365 Apps with Project and Visio"
  description        = "Microsoft 365 Apps including Project and Visio applications"
  is_featured        = true
  information_url    = "https://support.microsoft.com/office"
  notes              = "Complete Microsoft 365 suite including Project and Visio for power users."
  role_scope_tag_ids = ["0"]

  categories = [
    "Business",
    "Productivity",
  ]

  configuration_designer = {
    auto_accept_eula = true

    excluded_apps = {
      access               = true # Exclude Access
      bing                 = false
      excel                = false
      groove               = false
      info_path            = true
      lync                 = false
      one_drive            = false
      one_note             = false
      outlook              = false
      power_point          = false
      publisher            = true
      share_point_designer = true
      teams                = false
      visio                = false # Include Visio
      word                 = false
    }

    locales_to_install = [
      "en-us",
      "es-es", # Spanish (Spain)
      "it-it", # Italian (Italy)
    ]

    office_platform_architecture         = "x64"
    office_suite_app_default_file_format = "officeOpenXMLFormat"

    product_ids = [
      "o365ProPlusRetail",
      "projectProRetail",
      "visioProRetail"
    ]

    should_uninstall_older_versions_of_office = true
    update_channel                            = "monthlyEnterprise"
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

# Example 5: Office Suite App for Shared Computer Activation (VDI/Terminal Services)
resource "microsoft365_graph_beta_device_and_app_management_office_suite_app" "office_365_shared_activation" {
  display_name       = "Microsoft 365 Apps - Shared Computer Activation"
  description        = "Microsoft 365 Apps configured for shared computer environments (VDI/Terminal Services)"
  is_featured        = false
  information_url    = "https://support.microsoft.com/office"
  notes              = "Microsoft 365 Apps optimized for shared computer activation scenarios."
  role_scope_tag_ids = ["0"]

  categories = [
    "Business",
    "Productivity",
  ]

  configuration_designer = {
    auto_accept_eula = true

    excluded_apps = {
      access               = true
      bing                 = false
      excel                = false
      groove               = true # Exclude OneDrive for shared environments
      info_path            = true
      lync                 = false
      one_drive            = true # Exclude OneDrive for shared environments
      one_note             = false
      outlook              = false
      power_point          = false
      publisher            = true
      share_point_designer = true
      teams                = false
      visio                = true
      word                 = false
    }

    locales_to_install                   = ["en-us"]
    office_platform_architecture         = "x64"
    office_suite_app_default_file_format = "officeOpenXMLFormat"

    product_ids = [
      "o365ProPlusRetail"
    ]

    should_uninstall_older_versions_of_office = true
    update_channel                            = "deferred" # Use deferred channel for stability in shared environments
    update_version                            = ""         // for latest version, use empty string
    use_shared_computer_activation            = true       # Enable shared computer activation
  }

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}