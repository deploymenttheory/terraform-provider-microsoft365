# ==============================================================================
# Update Policy with Gradual Rollout Settings
# ==============================================================================
# Creates a Windows Update policy with deployment settings that control how
# updates are rolled out to devices. The gradual rollout settings specify that
# updates should be offered to 1000 devices at a time, with 1 day between each
# batch of offers.

resource "microsoft365_graph_beta_windows_updates_update_policy" "with_gradual_rollout" {
  audience_id        = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
  compliance_changes = true

  deployment_settings = {
    schedule = {
      gradual_rollout = {
        duration_between_offers = "P1D"
        devices_per_offer       = 1000
      }
    }
  }
}
