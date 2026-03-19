# ==============================================================================
# Minimal Update Policy
# ==============================================================================
# Creates a Windows Update policy with only the required fields.
# This minimal configuration creates a policy with an audience reference
# and enables compliance changes, but does not configure any compliance
# change rules or deployment settings.

resource "microsoft365_graph_beta_windows_updates_update_policy" "minimal" {
  audience_id        = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
  compliance_changes = true
}
