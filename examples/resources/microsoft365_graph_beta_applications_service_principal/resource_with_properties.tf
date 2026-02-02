resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-saml-application"
  description  = "SAML-based enterprise application"
}

# Create service principal with full configuration
resource "microsoft365_graph_beta_applications_service_principal" "example" {
  app_id                        = microsoft365_graph_beta_applications_application.example.app_id
  account_enabled               = true
  app_role_assignment_required  = true
  description                   = "Enterprise application for SSO access"
  login_url                     = "https://login.mycompany.com"
  notes                         = "Managed by Terraform - Production environment"
  notification_email_addresses  = ["admin@mycompany.com", "security@mycompany.com"]
  preferred_single_sign_on_mode = "saml"

  tags = [
    "HideApp",
    "WindowsAzureActiveDirectoryIntegratedApp"
  ]
}
