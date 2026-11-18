# This configures the existing tenant singleton. Destroy retains the policy.
resource "microsoft365_graph_beta_identity_and_access_admin_consent_request_policy" "example" {
  is_enabled               = true
  notify_reviewers         = true
  reminders_enabled        = true
  request_duration_in_days = 30
  reviewers = [{
    query      = "/users/11111111-1111-4111-8111-111111111111"
    query_type = "MicrosoftGraph"
  }]
}
