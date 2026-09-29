resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-plaintext-stringxml-${random_string.suffix.result}"
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
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/FloatingPoint",
        "value"       = 1.5
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingBoolean",
        "displayName" = "Provider Boolean",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Boolean",
        "value"       = false
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingInteger",
        "displayName" = "Provider Integer",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Integer",
        "value"       = 7
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingBase64",
        "displayName" = "Provider Base64",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Base64",
        "value"       = "AAECA//+/Q==",
        "fileName"    = null
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingString",
        "displayName" = "Provider String",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/String",
        "value"       = "test"
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingStringXml",
        "displayName" = "Provider StringXml",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/StringXml",
        "value"       = "<test>\n  <name>é 日本語</name>\n</test>\n",
        "fileName"    = null
      }
    ]
  })
}
