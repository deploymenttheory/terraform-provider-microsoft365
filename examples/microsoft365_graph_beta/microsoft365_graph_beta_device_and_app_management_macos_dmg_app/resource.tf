resource "microsoft365_graph_beta_device_and_app_management_macos_dmg_app" "docker_desktop" {
  display_name            = "Docker Desktop 4.36.0.dmg"
  description             = "Docker Desktop for macOS - containerization platform"
  publisher               = "Docker Inc."
  is_featured             = true
  privacy_information_url = "https://www.docker.com/legal/privacy-policy"
  information_url         = "https://www.docker.com/products/docker-desktop"
  owner                   = "Example Owner"
  developer               = "Docker Inc."
  notes                   = "This is a macOS DMG application managed through Terraform."
  role_scope_tag_ids      = [microsoft365_graph_beta_device_management_role_scope_tag.example.id, "2"]

  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.example.id, # custom app category
    "Business", # builtin category
    "Developer Tools",
  ]

  app_icon = {
    icon_file_path_source = "/path/to/Docker_logo.png"
  }

  app_installer = {
    // either installer_file_path_source or installer_url_source must be provided
    installer_file_path_source = "/path/to/Docker Desktop 4.36.0.dmg"
    installer_url_source = "https://example.com/Docker.dmg"
  }
  macos_dmg_app = {
    ignore_version_detection = true

    included_apps = [
      {
        bundle_id = "com.docker.docker"
        bundle_version = "4.36.0"
      }
    ]

    minimum_supported_operating_system = {
      v14_0 = true
    }
  }

  # Optional: Add timeouts block
  timeouts = {
    create = "5m"
    read   = "20s"
    update = "5m"
    delete = "20s"
  }
} 

resource "microsoft365_graph_beta_device_and_app_management_macos_dmg_app" "jamf_connect" {
  display_name            = "Jamf Connect 2.42.0.dmg"
  description             = "Jamf Connect - Identity and network access solution for macOS"
  publisher               = "Jamf Software, LLC"
  is_featured             = true
  privacy_information_url = "https://www.jamf.com/privacy-policy/"
  information_url         = "https://www.jamf.com/products/jamf-connect/"
  owner                   = "Example Owner"
  developer               = "Jamf Software, LLC"
  notes                   = "This is a macOS DMG application for Jamf Connect managed through Terraform."
  //role_scope_tag_ids      = [microsoft365_graph_beta_device_management_role_scope_tag.example.id, "2"]

  categories = [
    //microsoft365_graph_beta_device_and_app_management_application_category.example.id, # custom app category
    "Business", # builtin category
    //"Security",
  ]

  app_icon = {
    icon_file_path_source = "/Users/dafyddwatkins/Downloads/jamf-connect-icon.png"
  }

  app_installer = {
    installer_file_path_source = "/Users/dafyddwatkins/Downloads/JamfConnect-2.42.0.dmg"
  }

  macos_dmg_app = {
    ignore_version_detection = true

    included_apps = [
      {
        bundle_id = "com.jamf.connect"
        bundle_version = "2.42.0"
      }
    ]

    minimum_supported_operating_system = {
      v14_0 = true
    }
  }

  # Optional: Add timeouts block
  timeouts = {
    create = "5m"
    read   = "20s"
    update = "5m"
    delete = "20s"
  }
} 