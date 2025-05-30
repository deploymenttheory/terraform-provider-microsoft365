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
    installer_file_path_source = "/path/to/Docker Desktop 4.36.0.dmg"
  }

  macos_dmg_app = {
    ignore_version_detection = true

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