resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-windows-domain-join"
  description        = "Updated template description"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windowsDomainJoinConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "computerNameStaticPrefix"                    = "LAB"
    "computerNameSuffixRandomCharCount"           = 8
    "activeDirectoryDomainName"                   = "example.invalid"
    "organizationalUnit"                          = "OU=Testing,DC=example,DC=invalid"
  })
}
