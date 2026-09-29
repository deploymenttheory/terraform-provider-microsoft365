# Stage an allow-list policy for selected partner domains.
# Set is_organization_default to true to enforce it tenant-wide.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "allow_domains" {
  display_name            = "B2B invitations - allowed partner domains"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          AllowedDomains = ["example.com", "example.net"]
        }
      }
    })
  ]
}
