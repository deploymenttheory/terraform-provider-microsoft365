# Windows Boolean OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_boolean" {
  display_name       = "Windows Boolean OMA Setting"
  description        = "Windows boolean oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingBoolean",
        "displayName" = "Provider Boolean",
        "description" = "Disposable value test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Boolean",
        "value"       = true
      }
    ]
  })
}
