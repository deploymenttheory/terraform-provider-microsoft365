# ==============================================================================
# Update Policy with Compliance Change Rules
# ==============================================================================
# Creates a Windows Update policy with compliance change rules that automatically
# approve driver updates after a 7-day waiting period. The compliance change rules
# cannot be modified after creation - any changes will require resource replacement.

data "microsoft365_graph_beta_windows_updates_deployment_audience" "example" {
  id = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
}

resource "microsoft365_graph_beta_windows_updates_update_policy" "with_compliance_rules" {
  audience_id        = data.microsoft365_graph_beta_windows_updates_deployment_audience.example.id
  compliance_changes = true

  compliance_change_rules = [
    {
      content_filter = {
        filter_type = "driverUpdateFilter"
      }
      duration_before_deployment_start = "P7D"
    }
  ]
}
