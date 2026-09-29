# Windows Encrypted OMA Setting.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_encrypted_oma_setting" {
  display_name       = "Windows Encrypted OMA Setting"
  description        = "Windows encrypted oma setting example"
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
