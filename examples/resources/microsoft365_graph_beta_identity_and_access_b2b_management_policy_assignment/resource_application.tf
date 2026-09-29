# Create the application registration used by this example.
resource "microsoft365_graph_beta_applications_application" "application" {
  display_name = "B2B policy target - application"
}

# This non-default policy does not change tenant-wide invitation restrictions.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "application" {
  display_name            = "B2B policy - application"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          BlockedDomains = ["example.net"]
        }
      }
    })
  ]
}

# Use the application object ID, not its app_id (client ID).
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "application" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.application.id
  directory_object_id      = microsoft365_graph_beta_applications_application.application.id
}
