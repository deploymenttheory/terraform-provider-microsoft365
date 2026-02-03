# Look up an application and output its authentication configuration
data "microsoft365_graph_beta_applications_application" "with_auth_config" {
  app_id = "00000000-0000-0000-0000-000000000000" # Replace with actual app ID
}

# Output web application authentication configuration
output "web_auth_config" {
  value = {
    id                        = data.microsoft365_graph_beta_applications_application.with_auth_config.id
    display_name              = data.microsoft365_graph_beta_applications_application.with_auth_config.display_name
    sign_in_audience          = data.microsoft365_graph_beta_applications_application.with_auth_config.sign_in_audience
    is_fallback_public_client = data.microsoft365_graph_beta_applications_application.with_auth_config.is_fallback_public_client
    web                       = data.microsoft365_graph_beta_applications_application.with_auth_config.web
    spa                       = data.microsoft365_graph_beta_applications_application.with_auth_config.spa
    public_client             = data.microsoft365_graph_beta_applications_application.with_auth_config.public_client
  }
}
