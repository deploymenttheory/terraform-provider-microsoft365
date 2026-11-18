resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "test" {
  display_name = "APP_NAME"
  description  = "Optional URL import regression"
  publisher    = "MacAdmins"
  # PRIVACY_URL

  macos_pkg_app = {
    minimum_supported_operating_system = {
      v12_0 = true
    }
  }
}
