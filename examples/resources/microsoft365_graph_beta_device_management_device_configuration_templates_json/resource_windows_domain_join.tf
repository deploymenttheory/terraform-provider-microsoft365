# Windows Domain Join.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_domain_join" {
  display_name       = "Windows Domain Join"
  description        = "Template configuration example"
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
