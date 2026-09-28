# Example: Block B2B invitations to specific domains
# The definition is a single JSON string stored by Microsoft Graph exactly as supplied.

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "blocked_domains" {
  display_name = "example-b2b-blocked-domains"

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

# Example: Only allow B2B invitations to specific domains
# Only one B2B management policy can be the organization default. Setting
# is_organization_default = true changes the tenant's external collaboration settings.

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "allowed_domains" {
  display_name            = "example-b2b-allowed-domains"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          AllowedDomains = ["contoso.com", "fabrikam.com"]
        }
      }
    })
  ]
}
