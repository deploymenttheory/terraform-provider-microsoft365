# Graph rejects a sender domain that is not a verified domain of the tenant.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings" "test" {
  workflow_schedule_interval_in_hours = 6

  email_settings = {
    sender_domain = "acc-test-unverified.example.com"
  }
}
