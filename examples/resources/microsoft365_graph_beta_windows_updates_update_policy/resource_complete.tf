# ==============================================================================
# Complete Update Policy Configuration
# ==============================================================================
# Creates a comprehensive Windows Update policy with all available options:
# - Compliance change rules for automatic driver update approval
# - Deployment settings with scheduled start time
# - Gradual rollout configuration for controlled deployment
#
# Note: Compliance change rules cannot be modified after creation. Changes to
# filter_type or duration_before_deployment_start will require resource replacement.

data "microsoft365_graph_beta_windows_updates_deployment_audience" "production" {
  id = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
}

resource "microsoft365_graph_beta_windows_updates_update_policy" "complete" {
  audience_id        = data.microsoft365_graph_beta_windows_updates_deployment_audience.production.id
  compliance_changes = true

  compliance_change_rules = [
    {
      content_filter = {
        filter_type = "driverUpdateFilter"
      }
      duration_before_deployment_start = "P14D"
    }
  ]

  deployment_settings = {
    schedule = {
      start_date_time = "2026-04-15T02:00:00Z"
      gradual_rollout = {
        duration_between_offers = "P7D"
        devices_per_offer       = 2000
      }
    }
  }
}
