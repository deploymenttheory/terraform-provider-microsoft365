resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "test" {
  display_name = "acc-test-macos-pkg-featured-${random_string.suffix.result}"
  description  = "Nudge macOS PKG featured lifecycle regression"
  publisher    = "MacAdmins"
  # FEATURED

  app_installer = {
    installer_url_source = "https://github.com/macadmins/nudge/releases/download/v2.1.4.81863/Nudge-2.1.4.81863.pkg"
  }

  macos_pkg_app = {
    minimum_supported_operating_system = {
      v12_0 = true
    }
  }

  timeouts = {
    create = "10m"
    update = "10m"
  }
}
