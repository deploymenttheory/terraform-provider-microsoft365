resource "random_uuid" "suffix" {}

resource "microsoft365_graph_beta_applications_application" "target" {
  display_name = "acc-test-b2b-target-${random_uuid.suffix.result}"
  hard_delete  = true
}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "test" {
  display_name            = "acc-test-b2b-assignment-${random_uuid.suffix.result}"
  is_organization_default = false
  definition = [jsonencode({
    B2BManagementPolicy = {
      InvitationsAllowedAndBlockedDomainsPolicy = { BlockedDomains = ["blocked.example"] }
    }
  })]
}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "application" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.test.id
  directory_object_id      = microsoft365_graph_beta_applications_application.target.id
}
