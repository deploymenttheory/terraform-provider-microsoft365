# An empty domain policy specifies no domain restrictions.
# This non-default example does not change the active organization policy.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "any_domain" {
  display_name            = "B2B invitations - any domain"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {}
      }
    })
  ]
}
