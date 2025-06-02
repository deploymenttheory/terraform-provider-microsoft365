resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "mozilla_firefox" {
  display_name            = "Firefox 136.0.pkg"
  description             = "test"
  publisher               = "Example Publisher"
  is_featured             = true
  privacy_information_url = "https://example.com/privacy"
  information_url         = "https://example.com/info"
  owner                   = "Example Owner"
  developer               = "Example Developer"
  notes                   = "This is a macOS PKG application managed through Terraform."
  role_scope_tag_ids      = [microsoft365_graph_beta_device_management_role_scope_tag.example.id, "2"]

  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.example.id, # custom app category
    "Business", # builtin category
    "Productivity",
  ]

  app_icon = {
    icon_file_path_source = "/path/to/Firefox_logo.png"
  }

  app_installer = {
    // either installer_file_path_source or installer_url_source must be provided
    installer_file_path_source = "/path/to/Firefox_136.0.pkg"
    installer_url_source = "https://example.com/Firefox_136.0.pkg"
  }


  macos_pkg_app = {
    ignore_version_detection = true

    minimum_supported_operating_system = {
      v14_0 = true
    }

    pre_install_script = {
      script_content = base64encode("#!/bin/bash\necho macOS PKG Pre-install script example")
    }

    post_install_script = {
      script_content = base64encode("#!/bin/bash\necho macOS PKG Post-install script example")
    }
  }

  # Optional: Add timeouts block
  timeouts = {
    create = "3m"
    read   = "20s"
    update = "3m"
    delete = "20s"
  }
}