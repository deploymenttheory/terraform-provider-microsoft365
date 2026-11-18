resource "microsoft365_graph_beta_identity_and_access_admin_consent_request_policy" "test" {
  is_enabled               = false
  notify_reviewers         = false
  reminders_enabled        = false
  request_duration_in_days = 0
  reviewers                = []
}
