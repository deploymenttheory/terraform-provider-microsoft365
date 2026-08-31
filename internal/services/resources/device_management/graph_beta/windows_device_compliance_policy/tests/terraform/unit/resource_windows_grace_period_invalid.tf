resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "grace_period_invalid" {
  display_name       = "unit-test-wdcp-grace-period-invalid"
  description        = "unit-test-wdcp-grace-period-invalid"
  role_scope_tag_ids = ["0"]

  microsoft_defender_for_endpoint = {
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
  }

  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        # 8761 exceeds the Graph maximum of 8760 hours (365 days)
        {
          action_type        = "block"
          grace_period_hours = 8761
        },
      ]
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
