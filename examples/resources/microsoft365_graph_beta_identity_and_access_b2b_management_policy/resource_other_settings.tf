# Keep an allow list and explicitly supply empty preview and auto-redemption lists.
# The resource accepts these additional settings in the same JSON definition.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "other_settings" {
  display_name            = "B2B invitations - explicit empty additional settings"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          AllowedDomains = ["example.com", "example.net"]
        }
        PreviewPolicy = {
          Features = []
        }
        AutoRedeemPolicy = {
          AdminConsentedForUsersIntoTenantIds = []
          NoAADConsentForUsersFromTenantsIds  = []
        }
      }
    })
  ]
}
