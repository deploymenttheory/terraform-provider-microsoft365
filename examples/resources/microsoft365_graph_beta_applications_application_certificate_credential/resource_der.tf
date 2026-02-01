# Example: Create a certificate credential with DER/Base64 encoding
# Use this when you have a DER (binary) certificate file

# First, create or reference an existing application
resource "microsoft365_graph_beta_applications_application" "example_der" {
  display_name = "my-application-der"
  description  = "Application for DER certificate authentication"
  hard_delete  = true
}

# Option 1: Load a DER certificate from file using filebase64()
resource "microsoft365_graph_beta_applications_application_certificate_credential" "from_file" {
  application_id = microsoft365_graph_beta_applications_application.example_der.id
  display_name   = "my-der-certificate"

  key      = filebase64("path/to/certificate.der")
  encoding = "base64"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"
}

# Option 2: Generate a certificate and convert PEM to base64 (DER equivalent)
resource "tls_private_key" "example_der" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "example_der" {
  private_key_pem = tls_private_key.example_der.private_key_pem

  subject {
    common_name  = "my-application-certificate-der"
    organization = "My Organization"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

# Extract the base64 content from PEM (strip headers) - equivalent to base64-encoded DER
locals {
  cert_base64 = replace(
    replace(
      replace(tls_self_signed_cert.example_der.cert_pem, "-----BEGIN CERTIFICATE-----", ""),
      "-----END CERTIFICATE-----", ""
    ),
    "\n", ""
  )
}

resource "microsoft365_graph_beta_applications_application_certificate_credential" "from_generated" {
  application_id = microsoft365_graph_beta_applications_application.example_der.id
  display_name   = "my-generated-der-certificate"

  key      = local.cert_base64
  encoding = "base64"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"
}
