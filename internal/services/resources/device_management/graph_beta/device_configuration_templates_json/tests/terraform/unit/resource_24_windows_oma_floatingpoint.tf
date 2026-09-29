resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-windows-oma-floatingpoint"
  description        = "Device configuration template test"
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
