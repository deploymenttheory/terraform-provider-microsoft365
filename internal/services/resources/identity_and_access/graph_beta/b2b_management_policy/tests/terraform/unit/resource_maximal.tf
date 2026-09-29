resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "maximal" {
  display_name            = "unit-test-b2b-management-policy-max"
  is_organization_default = true
  definition              = ["{\"B2BManagementPolicy\":{\"InvitationsAllowedAndBlockedDomainsPolicy\":{\"AllowedDomains\":[\"example.net\",\"example.org\"]},\"AutoRedeemPolicy\":{\"AdminConsentedForUsersIntoTenantIds\":[],\"NoAADConsentForUsersFromTenantsIds\":[]}}}"]
}
