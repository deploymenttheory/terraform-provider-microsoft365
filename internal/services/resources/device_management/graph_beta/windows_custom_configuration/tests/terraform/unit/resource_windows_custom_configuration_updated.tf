# Updated configuration used by the update unit test
resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "custom_configuration_example" {
  display_name = "unit-test-windows-custom-configuration-example-updated"
  description  = "Example Windows custom configuration profile using OMA-URI settings (updated)"

  oma_settings = [
    {
      odata_type   = "#microsoft.graph.omaSettingString"
      display_name = "ExtensionsAutoUpdate"
      description  = "Enable VSCode extension auto update"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application~extensionsConfigurationTitle/ExtensionsAutoUpdate"
      value        = "<enabled/>\n<data id=\"ExtensionsAutoUpdate\" value=\"on\"/>"
    }
  ]

  role_scope_tag_ids = ["00000000-0000-0000-0000-000000000001"]

  assignments = [
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000004"
    }
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
