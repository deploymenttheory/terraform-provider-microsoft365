resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "example" {
  display_name        = "Example Cloud PC Alert Rule"
  description         = "This is an example Cloud PC alert rule."
  alert_rule_template = "cloudPcFrontlineInsufficientLicensesScenario" # Use a valid template value from the API
  enabled             = true
  is_system_rule      = false
  severity            = "warning" # Use a valid severity value from the API

  notification_channels = [
    {
      notification_channel_type = "email" # Use a valid channel type from the API
      notification_receivers = [
        {
          contact_information = "admin@example.com"
          locale              = "en-US"
        }
      ]
    }
  ]

  threshold = {
    aggregation = "count"          # Use a valid aggregation type from the API
    operator    = "greaterOrEqual" # Use a valid operator from the API
    target      = 10
  }

  conditions = [
    {
      condition_category = "cloudPcInGracePeriod"
      relationship_type  = "and"
      aggregation        = "count"
      operator           = "greater"
      threshold_value    = "5"
    }
  ]
} 