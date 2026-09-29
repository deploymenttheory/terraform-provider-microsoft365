resource "random_uuid" "suffix" {}

resource "microsoft365_graph_beta_applications_application" "target" {
  display_name = "acc-test-b2b-target-${random_uuid.suffix.result}"
  hard_delete  = true
}

resource "time_sleep" "wait_for_application" {
  depends_on      = [microsoft365_graph_beta_applications_application.target]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_service_principal" "target" {
  hard_delete = true
  app_id      = microsoft365_graph_beta_applications_application.target.app_id
  depends_on  = [time_sleep.wait_for_application]
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

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "service_principal" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.test.id
  directory_object_id      = microsoft365_graph_beta_applications_service_principal.target.id
}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "application" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.test.id
  directory_object_id      = microsoft365_graph_beta_applications_application.target.id
}
