# Stage a block-list policy while allowing invitations to other domains.
# AllowedDomains and BlockedDomains are alternative modes; configure only one.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "block_domains" {
  display_name            = "B2B invitations - blocked domains"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          BlockedDomains = ["example.net", "example.org"]
        }
      }
    })
  ]
}
