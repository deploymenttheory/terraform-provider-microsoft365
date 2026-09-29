resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "minimal" {
  display_name            = "unit-test-b2b-management-policy-min"
  is_organization_default = true
  definition              = ["{\"B2BManagementPolicy\":{\"InvitationsAllowedAndBlockedDomainsPolicy\":{\"BlockedDomains\":[\"example.com\"]}}}"]
}
