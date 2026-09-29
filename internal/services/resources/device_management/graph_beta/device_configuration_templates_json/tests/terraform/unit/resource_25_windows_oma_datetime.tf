resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-windows-oma-datetime"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingDateTime",
        "displayName" = "Provider DateTime",
        "description" = "Disposable value test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/DateTime",
        "value"       = "2026-01-01T00:00:00Z"
      }
    ]
  })
}
