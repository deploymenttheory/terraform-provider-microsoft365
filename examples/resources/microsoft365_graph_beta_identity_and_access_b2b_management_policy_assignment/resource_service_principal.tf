# Create the application registration used by this example.
resource "microsoft365_graph_beta_applications_application" "service_principal" {
  display_name = "B2B policy target - service principal"
}

# Allow the application registration to replicate before creating its service principal.
resource "time_sleep" "service_principal_application_replication" {
  depends_on      = [microsoft365_graph_beta_applications_application.service_principal]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_service_principal" "service_principal" {
  app_id = microsoft365_graph_beta_applications_application.service_principal.app_id

  depends_on = [time_sleep.service_principal_application_replication]
}

# This non-default policy does not change tenant-wide invitation restrictions.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "service_principal" {
  display_name            = "B2B policy - service principal"
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

# Use the service principal object ID, not its app_id (client ID).
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "service_principal" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.service_principal.id
  directory_object_id      = microsoft365_graph_beta_applications_service_principal.service_principal.id
}
