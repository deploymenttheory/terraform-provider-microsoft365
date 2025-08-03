resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "maximal" {
  alert_rule_template = "cloudPcProvisionScenario"
  display_name        = "Test Maximal Cloud PC Alert Rule - Unique"
  description         = "Comprehensive alert rule for testing Cloud PC provisioning failures with all features"
  severity            = "critical"
  enabled             = true
  is_system_rule      = false

  notification_channels = [
    {
      notification_channel_type = "portal"
      notification_receivers = [
        {
          contact_information = "admin@test.com"
          locale              = "en-US"
        },
        {
          contact_information = "manager@test.com"
          locale              = "en-US"
        }
      ]
    },
    {
      notification_channel_type = "email"
      notification_receivers = [
        {
          contact_information = "alerts@test.com"
          locale              = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"
    operator    = "greaterOrEqual"
    target      = 5
  }

  conditions = [
    {
      relationship_type  = "and"
      condition_category = "provisionFailures"
      aggregation        = "count"
      operator           = "greaterOrEqual"
      threshold_value    = "3"
    },
    {
      relationship_type  = "or"
      condition_category = "cloudPcConnectionErrors"
      aggregation        = "percentage"
      operator           = "less"
      threshold_value    = "95"
    }
  ]

  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "2m"
  }
}