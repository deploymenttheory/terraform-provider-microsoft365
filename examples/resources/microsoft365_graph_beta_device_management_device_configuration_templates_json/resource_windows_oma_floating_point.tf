# Windows Floating-point OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_floating_point" {
  display_name       = "Windows Floating-point OMA Setting"
  description        = "Windows floating-point oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingFloatingPoint",
        "displayName" = "Provider FloatingPoint",
        "description" = "Disposable value test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/FloatingPoint",
        "value"       = 1.5
      }
    ]
  })
}
