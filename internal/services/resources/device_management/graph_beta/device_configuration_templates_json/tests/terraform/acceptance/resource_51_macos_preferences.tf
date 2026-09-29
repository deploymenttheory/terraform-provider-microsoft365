resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-macos-preferences-${random_string.suffix.result}"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSCustomAppConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "bundleId"                                    = "com.example.provider-wire-study",
    "fileName"                                    = "preferences.plist",
    "configurationXml"                            = "PD94bWwgdmVyc2lvbj0iMS4wIj8+PHBsaXN0IHZlcnNpb249IjEuMCI+PGRpY3Q+PGtleT5FeGFtcGxlRW5hYmxlZDwva2V5Pjx0cnVlLz48L2RpY3Q+PC9wbGlzdD4="
  })
}
