# ==============================================================================
# Update Policy with Scheduled Start Time
# ==============================================================================
# Creates a Windows Update policy with a specific start date/time for the
# deployment schedule. This allows you to control when updates begin rolling
# out to devices. The gradual rollout settings control the pace of deployment.

resource "microsoft365_graph_beta_windows_updates_update_policy" "with_scheduled_start" {
  audience_id        = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
  compliance_changes = true

  deployment_settings = {
    schedule = {
      start_date_time = "2026-04-01T00:00:00Z"
      gradual_rollout = {
        duration_between_offers = "P2D"
        devices_per_offer       = 500
      }
    }
  }
}
