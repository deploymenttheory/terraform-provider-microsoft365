# Example: Create a password credential for an Application

# First, create or reference an existing application
resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-application"
  description  = "Application for automated workflows"
  hard_delete  = true
}

# Create a password credential for the application
resource "microsoft365_graph_beta_applications_application_password_credential" "example" {
  application_id = microsoft365_graph_beta_applications_application.example.id
  display_name   = "api-access-credential"
}

# IMPORTANT: Store the secret securely - it is only available at creation time
# You can use a secrets manager or output to a secure location
output "client_secret" {
  value       = microsoft365_graph_beta_applications_application_password_credential.example.secret_text
  description = "The generated client secret - store this securely!"
  sensitive   = true
}

output "key_id" {
  value       = microsoft365_graph_beta_applications_application_password_credential.example.key_id
  description = "The key ID of the password credential"
}
