resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "test" {
  display_name = "acc-test-macos-pkg-metadata-${random_string.suffix.result}"
  description  = "Original app metadata"
  publisher    = "MacAdmins"

  app_installer = {
    installer_file_path_source = "INSTALLER_SOURCE"
  }

  macos_pkg_app = {
    ignore_version_detection = false
    minimum_supported_operating_system = {
      v12_0 = true
    }
  }

  timeouts = {
    create = "10m"
    update = "3m"
  }
}
