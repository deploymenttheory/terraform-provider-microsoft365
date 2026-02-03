# Look up an application and output its credentials information
# Note: Actual secret values are not returned by the API for security reasons
data "microsoft365_graph_beta_applications_application" "with_credentials" {
  display_name = "My Application with Credentials"
}

# Output credentials metadata (not the actual secrets)
output "credentials_info" {
  value = {
    id                   = data.microsoft365_graph_beta_applications_application.with_credentials.id
    display_name         = data.microsoft365_graph_beta_applications_application.with_credentials.display_name
    key_credentials      = data.microsoft365_graph_beta_applications_application.with_credentials.key_credentials
    password_credentials = data.microsoft365_graph_beta_applications_application.with_credentials.password_credentials
  }
  sensitive = true # Credentials information should be marked sensitive
}

# Note: To retrieve the public key value in key_credentials, you must use $select=keyCredentials
# in an OData query. The key value is only returned when explicitly requested.
