resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-windows-encrypted"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingString",
        "displayName" = "Provider test string",
        "description" = "Disposable unassigned test",
        "omaUri"      = "./Device/Vendor/MSFT/Policy/Config/Experience/ConfigureWindowsSpotlightOnLockScreen",
        "value"       = "provider-test-string"
      }
    ]
  })
}
