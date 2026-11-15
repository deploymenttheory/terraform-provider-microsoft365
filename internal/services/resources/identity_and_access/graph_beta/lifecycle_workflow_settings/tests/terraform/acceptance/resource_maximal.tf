# Maximal configuration: every configurable setting is managed.
# microsoft.com is the default Microsoft sender domain and is accepted in every tenant.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings" "test" {
  workflow_schedule_interval_in_hours = 7

  email_settings = {
    sender_domain        = "microsoft.com"
    use_company_branding = true
  }
}
