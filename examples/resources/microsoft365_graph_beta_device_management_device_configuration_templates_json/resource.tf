# Metadata and assignments are standard Terraform attributes.
# Settings contain the complete writable JSON object for the selected template.
# These examples are unassigned.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_camera" {
  display_name       = "Windows camera configuration"
  description        = "Configure the Windows camera CSP"
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
        "description" = "Allow use of the camera",
        "omaUri"      = "./Device/Vendor/MSFT/Policy/Config/Camera/AllowCamera",
        "value"       = 1
      }
    ]
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_preferences" {
  display_name       = "macOS application preferences"
  description        = "Configure application preferences"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSCustomAppConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "bundleId"                                    = "com.example.application",
    "fileName"                                    = "preferences.plist",
    "configurationXml"                            = "PD94bWwgdmVyc2lvbj0iMS4wIj8+PHBsaXN0IHZlcnNpb249IjEuMCI+PGRpY3Q+PGtleT5FeGFtcGxlRW5hYmxlZDwva2V5Pjx0cnVlLz48L2RpY3Q+PC9wbGlzdD4="
  })
}
