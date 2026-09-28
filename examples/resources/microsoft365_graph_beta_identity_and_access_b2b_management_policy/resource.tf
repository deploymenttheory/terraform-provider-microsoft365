# Block B2B invitations to specific domains tenant-wide.
# Only the organization default policy takes effect, and only one can exist.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "example" {
  display_name            = "b2b-management-policy"
  is_organization_default = true

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          BlockedDomains = ["example.com", "example.net"]
        }
      }
    })
  ]
}
