resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  depends_on         = [time_sleep.directory_consistency]
  display_name       = "acc-test-ios-custom-${random_string.suffix.result}"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosCustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "payloadName"                                 = "Wire Study",
    "payloadFileName"                             = "wire-study.mobileconfig",
    "payload"                                     = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheS8+Cgk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJPHN0cmluZz5XaXJlIFN0dWR5PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPmNvbS5leGFtcGxlLmNvZGV4LndpcmUtc3R1ZHk8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCTxzdHJpbmc+Q29uZmlndXJhdGlvbjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJPHN0cmluZz45ZDBhYjBmOS04YzBjLTQwMzUtOGFlOC03OTA3ZDI2MTc3M2I8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFZlcnNpb248L2tleT4KCTxpbnRlZ2VyPjE8L2ludGVnZXI+CjwvZGljdD4KPC9wbGlzdD4K"
  })

  assignments = [
    { type = "allDevicesAssignmentTarget", filter_type = "none" },
    { type = "allLicensedUsersAssignmentTarget", filter_type = "none" },
    { type = "groupAssignmentTarget", group_id = microsoft365_graph_beta_groups_group.include.id, filter_type = "none" },
    { type = "exclusionGroupAssignmentTarget", group_id = microsoft365_graph_beta_groups_group.exclude.id, filter_type = "none" }
  ]
}

resource "microsoft365_graph_beta_groups_group" "include" {
  display_name     = "acc-test-json-include-${random_string.suffix.result}"
  mail_nickname    = "acc-test-json-include-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}
resource "microsoft365_graph_beta_groups_group" "exclude" {
  display_name     = "acc-test-json-exclude-${random_string.suffix.result}"
  mail_nickname    = "acc-test-json-exclude-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}
resource "time_sleep" "directory_consistency" {
  depends_on      = [microsoft365_graph_beta_groups_group.include, microsoft365_graph_beta_groups_group.exclude]
  create_duration = "30s"
}
