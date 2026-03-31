resource "microsoft365_graph_beta_windows_updates_update_policy" "test" {
  audience_id        = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
  compliance_changes = true

  compliance_change_rules = [
    {
      content_filter = {
        filter_type = "driverUpdateFilter"
      }
      duration_before_deployment_start = "P7D"
    }
  ]

  deployment_settings = {
    schedule = {
      gradual_rollout = {
        duration_between_offers = "P2D"
        devices_per_offer       = 2000
      }
    }
  }
}
