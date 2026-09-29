resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-ios-custom"
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

  assignments = []
}
