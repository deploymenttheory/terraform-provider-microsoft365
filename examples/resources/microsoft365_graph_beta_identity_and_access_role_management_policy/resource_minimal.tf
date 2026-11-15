# Limit permanent active assignments of the Usage Summary Reports Reader role to 90 days.
# The policy is found through the role definition at the directory scope ("/").
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "usage_summary_reports_reader" {
  role_definition_id = "75934031-6c7e-415a-99d7-48dbd49e875e"

  expiration_rules = {
    Expiration_Admin_Assignment = {
      is_expiration_required = true
      maximum_duration       = "P90D"
    }
  }
}
