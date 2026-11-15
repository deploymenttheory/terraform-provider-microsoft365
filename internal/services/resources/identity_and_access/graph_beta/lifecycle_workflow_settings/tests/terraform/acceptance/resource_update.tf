# Update: the interval and company branding change in place.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings" "test" {
  workflow_schedule_interval_in_hours = 9

  email_settings = {
    sender_domain        = "microsoft.com"
    use_company_branding = false
  }
}
