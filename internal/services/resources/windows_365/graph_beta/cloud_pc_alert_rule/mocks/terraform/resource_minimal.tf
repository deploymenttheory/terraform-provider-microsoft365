resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "minimal" {
  alert_rule_template = "cloudPcProvisionScenario"
  display_name        = "Test Minimal Cloud PC Alert Rule - Unique"
  severity            = "warning"
  enabled             = true
  is_system_rule      = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale              = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 1
  }

  conditions = [
    {
      relationship_type  = "and"
      condition_category = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "1"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}