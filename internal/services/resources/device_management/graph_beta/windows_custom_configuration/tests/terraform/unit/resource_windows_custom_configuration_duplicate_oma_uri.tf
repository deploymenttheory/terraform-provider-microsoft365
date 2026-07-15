# Invalid configuration: two settings targeting the same OMA-URI
resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "duplicate_oma_uri_example" {
  display_name = "unit-test-windows-custom-configuration-duplicate-oma-uri"

  oma_settings = [
    {
      odata_type   = "#microsoft.graph.omaSettingInteger"
      display_name = "ExtensionsAutoUpdateDelay"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application~extensionsConfigurationTitle/ExtensionsAutoUpdateDelay"
      value        = "30"
    },
    {
      odata_type   = "#microsoft.graph.omaSettingInteger"
      display_name = "ExtensionsAutoUpdateDelay (duplicate)"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application~extensionsConfigurationTitle/ExtensionsAutoUpdateDelay"
      value        = "60"
    }
  ]
}
