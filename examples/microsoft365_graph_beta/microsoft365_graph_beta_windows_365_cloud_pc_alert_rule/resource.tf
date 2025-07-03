resource "microsoft365_graph_beta_windows_365_cloud_pc_alert_rule" "example" {
  display_name        = "Example Cloud PC Alert Rule"
  description         = "This is an example Cloud PC alert rule."
  alert_rule_template = "exampleTemplate" # Use a valid template value from the API
  enabled             = true
  is_system_rule      = false
  severity            = "high" # Use a valid severity value from the API

  notification_channels = [
    {
      odata_type                = "#microsoft.graph.notificationChannel"
      notification_channel_type = "email" # Use a valid channel type from the API
      notification_receivers = [
        {
          odata_type          = "#microsoft.graph.notificationReceiver"
          contact_information = "admin@example.com"
          locale              = "en-US"
        }
      ]
    }
  ]

  threshold = {
    odata_type   = "#microsoft.graph.ruleThreshold"
    aggregation  = "count" # Use a valid aggregation type from the API
    operator     = "greaterThan" # Use a valid operator from the API
    target       = 10
  }

  conditions = [
    {
      odata_type         = "#microsoft.graph.ruleCondition"
      relationship_type  = "device"
      condition_category = "security"
      aggregation        = "count"
      operator           = "greaterThan"
      threshold_value    = "5"
    }
  ]
} 