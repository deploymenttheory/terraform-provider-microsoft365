# Create the application registration used by this example.
resource "microsoft365_graph_beta_applications_application" "both_targets" {
  display_name = "B2B policy target - both targets"
}

# Allow the application registration to replicate before creating its service principal.
resource "time_sleep" "both_application_replication" {
  depends_on      = [microsoft365_graph_beta_applications_application.both_targets]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_service_principal" "both_targets" {
  app_id = microsoft365_graph_beta_applications_application.both_targets.app_id

  depends_on = [time_sleep.both_application_replication]
}

# This non-default policy does not change tenant-wide invitation restrictions.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "both_targets" {
  display_name            = "B2B policy - both targets"
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
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "both_application" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.both_targets.id
  directory_object_id      = microsoft365_graph_beta_applications_application.both_targets.id
}

# Use the service principal object ID, not its app_id (client ID).
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "both_service_principal" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.both_targets.id
  directory_object_id      = microsoft365_graph_beta_applications_service_principal.both_targets.id
}
