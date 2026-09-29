resource "random_uuid" "suffix" {}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "test" {
  display_name = "acc-test-b2b-${random_uuid.suffix.result}"
  definition   = ["{\"B2BManagementPolicy\":{\"InvitationsAllowedAndBlockedDomainsPolicy\":{\"BlockedDomains\":[\"example.com\"]}}}"]
}
