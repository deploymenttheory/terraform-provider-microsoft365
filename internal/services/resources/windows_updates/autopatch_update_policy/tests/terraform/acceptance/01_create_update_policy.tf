resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "test" {
}

resource "microsoft365_graph_beta_windows_updates_update_policy" "test" {
  audience_id         = microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test.id
  compliance_changes  = true

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
        duration_between_offers = "P1D"
        devices_per_offer       = 1000
      }
    }
  }
}
