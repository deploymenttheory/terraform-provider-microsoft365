# Invalid configuration: boolean setting with a non canonical value
resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "non_canonical_value_example" {
  display_name = "unit-test-windows-custom-configuration-non-canonical-value"

  oma_settings = [
    {
      odata_type   = "#microsoft.graph.omaSettingBoolean"
      display_name = "AllowTelemetry"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/System/AllowTelemetry"
      value        = "True"
    }
  ]
}
