# Windows Binary Base64 OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_base64" {
  display_name       = "Windows Binary Base64 OMA Setting"
  description        = "Windows binary base64 oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingBase64",
        "displayName" = "Provider Base64",
        "description" = "Disposable encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Base64",
        "fileName"    = null,
        "value"       = "U3ludGhldGljIGJpbmFyeQ=="
      }
    ]
  })
}
