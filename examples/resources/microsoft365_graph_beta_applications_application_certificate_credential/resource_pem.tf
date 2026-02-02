# Example: Create a certificate credential with PEM encoding (default)

# First, create or reference an existing application
resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-application"
  description  = "Application for certificate authentication"
  hard_delete  = true
}

# Generate a self-signed certificate using the tls provider
resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "example" {
  private_key_pem = tls_private_key.example.private_key_pem

  subject {
    common_name  = "my-application-certificate"
    organization = "My Organization"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

# Create a certificate credential using PEM encoding
resource "microsoft365_graph_beta_applications_application_certificate_credential" "example" {
  application_id = microsoft365_graph_beta_applications_application.example.id
  display_name   = "my-application-certificate"

  key      = tls_self_signed_cert.example.cert_pem
  encoding = "pem" # Default, can be omitted
  type     = "AsymmetricX509Cert"
  usage    = "Verify"
}

output "certificate_key_id" {
  value       = microsoft365_graph_beta_applications_application_certificate_credential.example.key_id
  description = "The key ID of the certificate credential"
}

output "certificate_thumbprint" {
  value       = microsoft365_graph_beta_applications_application_certificate_credential.example.thumbprint
  description = "The thumbprint (SHA-1 hash) of the certificate"
}
