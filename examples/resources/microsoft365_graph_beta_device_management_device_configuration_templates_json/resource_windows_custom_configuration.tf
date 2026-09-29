# Windows Custom Configuration.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_custom_configuration" {
  display_name       = "Windows Custom Configuration"
  description        = "Windows custom configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingInteger",
        "displayName" = "Allow Camera",
        "description" = "Wire Study",
        "omaUri"      = "./Device/Vendor/MSFT/Policy/Config/Camera/AllowCamera",
        "value"       = 1
      }
    ]
  })
}
