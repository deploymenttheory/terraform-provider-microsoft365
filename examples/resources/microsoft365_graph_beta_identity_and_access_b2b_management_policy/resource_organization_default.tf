# Only one organization-default B2B management policy can exist.
# If one already exists, import it instead of creating a second default.
# Changing is_organization_default replaces the policy; do not use create_before_destroy.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "organization_default" {
  display_name            = "B2B invitations - organization default"
  is_organization_default = true

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
