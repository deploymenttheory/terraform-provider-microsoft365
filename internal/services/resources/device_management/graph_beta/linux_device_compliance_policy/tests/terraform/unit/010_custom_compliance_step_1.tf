resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "test_010" {
  name                               = "unit-test-linux-compliance-010"
  device_encryption_required         = false
  custom_compliance_required         = true
  custom_compliance_discovery_script = "55555555-5555-5555-5555-555555555555"
  custom_compliance_rules            = "eyJSdWxlcyI6W3siU2V0dGluZ05hbWUiOiJUZXN0VmVyc2lvbiIsIk9wZXJhdG9yIjoiR3JlYXRlckVxdWFscyIsIkRhdGFUeXBlIjoiSW50NjQiLCJPcGVyYW5kIjoxfV19"
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
