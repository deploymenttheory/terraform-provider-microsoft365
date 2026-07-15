# Invalid configuration: integer setting with a non numeric value
resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "invalid_integer_example" {
  display_name = "unit-test-windows-custom-configuration-invalid-integer"

  oma_settings = [
    {
      odata_type   = "#microsoft.graph.omaSettingInteger"
      display_name = "ExtensionsAutoUpdateDelay"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application~extensionsConfigurationTitle/ExtensionsAutoUpdateDelay"
      value        = "not-a-number"
    }
  ]
}
