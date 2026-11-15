# Update: the interval and both email settings change in place.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings" "test" {
  workflow_schedule_interval_in_hours = 24

  email_settings = {
    sender_domain        = "microsoft.com"
    use_company_branding = false
  }
}
