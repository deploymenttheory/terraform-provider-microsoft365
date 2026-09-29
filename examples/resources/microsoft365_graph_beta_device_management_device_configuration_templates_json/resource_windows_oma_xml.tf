# Windows Cleartext XML OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.
# XML values are cleartext; the provider Base64-encodes XML when sending it to Graph.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_xml" {
  display_name       = "Windows Cleartext XML OMA Setting"
  description        = "Windows cleartext xml oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingStringXml",
        "displayName" = "Provider StringXml",
        "description" = "Disposable encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/StringXml",
        "fileName"    = null,
        "value"       = "<test>\n  <name>é 日本語</name>\n</test>\n"
      }
    ]
  })
}
