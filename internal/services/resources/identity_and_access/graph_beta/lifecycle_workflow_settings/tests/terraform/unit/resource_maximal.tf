# Maximal configuration: every configurable setting is managed.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings" "test" {
  workflow_schedule_interval_in_hours = 12

  email_settings = {
    sender_domain        = "contoso.com"
    use_company_branding = true
  }
}
