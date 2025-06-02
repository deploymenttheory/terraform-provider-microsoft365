resource "microsoft365_graph_beta_device_and_app_management_macos_lob_app" "example_app" {
  display_name            = "Example macOS LOB App"
  description             = "Example macOS Line of Business application"
  publisher               = "Example Publisher"
  is_featured             = true
  privacy_information_url = "https://example.com/privacy"
  information_url         = "https://example.com/info"
  owner                   = "Example Owner"
  developer               = "Example Developer"
  notes                   = "This is a macOS LOB application managed through Terraform."
  role_scope_tag_ids      = [microsoft365_graph_beta_device_management_role_scope_tag.example.id, "2"]

  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.example.id, # custom app category
    "Business", # builtin category
    "Productivity",
  ]

  app_icon = {
    icon_file_path_source = "/path/to/app_icon.png"
  }

  app_installer = {
    installer_file_path_source = "/path/to/example_app.pkg"
  }

  macos_lob_app = {
    bundle_id                  = "com.example.app"
    build_number              = "1.0.0"
    version_number            = "1.0.0"
    ignore_version_detection  = false
    install_as_managed        = true

    minimum_supported_operating_system = {
      v14_0 = true
    }

    child_apps = [
      {
        bundle_id      = "com.example.app.helper"
        build_number   = "1.0.0"
        version_number = "1.0.0"
      }
    ]
  }

  # Optional: Add timeouts block
  timeouts = {
    create = "10m"
    read   = "20s"
    update = "10m"
    delete = "20s"
  }
} 