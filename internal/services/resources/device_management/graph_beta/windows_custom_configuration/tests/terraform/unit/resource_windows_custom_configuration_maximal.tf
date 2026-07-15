# Example: Windows Custom Configuration (OMA-URI) profile
resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "custom_configuration_example" {
  display_name = "unit-test-windows-custom-configuration-example"
  description  = "Example Windows custom configuration profile using OMA-URI settings"

  oma_settings = [
    {
      odata_type   = "#microsoft.graph.omaSettingString"
      display_name = "ADMX Ingest"
      description  = "https://code.visualstudio.com/docs/enterprise/policies#_windows-group-policies"
      oma_uri      = "./Device/Vendor/MSFT/Policy/ConfigOperations/ADMXInstall/VSCode/Policy/VSCodeADMX"
      value        = "<?xml version=\"1.0\" encoding=\"utf-8\"?><policyDefinitions revision=\"1.1\" schemaVersion=\"1.0\"><policies /></policyDefinitions>"
    },
    {
      odata_type   = "#microsoft.graph.omaSettingString"
      display_name = "ExtensionsAutoUpdate"
      description  = "Disable VSCode extension auto update"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application~extensionsConfigurationTitle/ExtensionsAutoUpdate"
      value        = "<enabled/>\n<data id=\"ExtensionsAutoUpdate\" value=\"off\"/>"
    },
    {
      odata_type   = "#microsoft.graph.omaSettingInteger"
      display_name = "ExtensionsAutoUpdateDelay"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application~extensionsConfigurationTitle/ExtensionsAutoUpdateDelay"
      value        = "30"
    },
    {
      odata_type   = "#microsoft.graph.omaSettingBoolean"
      display_name = "AllowTelemetry"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/System/AllowTelemetry"
      value        = "false"
    }
  ]

  role_scope_tag_ids = ["00000000-0000-0000-0000-000000000001"]

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000002"
      filter_id   = "00000000-0000-0000-0000-000000000003"
      filter_type = "include"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000002"
      filter_id   = "00000000-0000-0000-0000-000000000003"
      filter_type = "exclude"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
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
