resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "custom_configuration_example" {
  display_name = "unit-test-windows-custom-configuration-missing-group-id"

  oma_settings = [
    {
      odata_type   = "#microsoft.graph.omaSettingInteger"
      display_name = "Example integer setting"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/Example/IntegerSetting"
      value        = "30"
    }
  ]

  assignments = [
    {
      type = "groupAssignmentTarget"
    }
  ]
}
