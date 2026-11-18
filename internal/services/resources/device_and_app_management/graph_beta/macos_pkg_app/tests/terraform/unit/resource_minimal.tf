resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "test" {
  display_name = "unit-test-macos-pkg-featured"
  description  = "macOS PKG featured plan regression"
  publisher    = "MacAdmins"
  # FEATURED

  app_installer = {
    installer_file_path_source = "/test/Nudge.pkg"
  }

  macos_pkg_app = {
    minimum_supported_operating_system = {
      v12_0 = true
    }
  }
}
