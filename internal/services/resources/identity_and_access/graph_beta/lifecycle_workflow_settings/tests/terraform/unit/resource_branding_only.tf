# Only use_company_branding is managed inside email_settings; sender_domain keeps its tenant value.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings" "test" {
  email_settings = {
    use_company_branding = true
  }
}
