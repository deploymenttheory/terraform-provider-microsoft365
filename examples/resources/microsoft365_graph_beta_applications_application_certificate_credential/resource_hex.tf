# Example: Create a certificate credential with HEX encoding
# Use this when you have a hex-encoded certificate string

# First, create or reference an existing application
resource "microsoft365_graph_beta_applications_application" "example_hex" {
  display_name = "my-application-hex"
  description  = "Application for HEX certificate authentication"
  hard_delete  = true
}

# Create a certificate credential using a pre-generated hex-encoded certificate
# To generate a hex-encoded certificate:
#   1. openssl x509 -in cert.pem -outform DER -out cert.der
#   2. xxd -p cert.der | tr -d '\n' > cert.hex
resource "microsoft365_graph_beta_applications_application_certificate_credential" "example_hex" {
  application_id = microsoft365_graph_beta_applications_application.example_hex.id
  display_name   = "my-hex-certificate"

  # Pre-generated hex-encoded certificate string
  # Generate with: openssl x509 -in cert.pem -outform DER | xxd -p | tr -d '\n'
  key      = file("path/to/certificate.hex")
  encoding = "hex"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"
}

output "hex_certificate_key_id" {
  value       = microsoft365_graph_beta_applications_application_certificate_credential.example_hex.key_id
  description = "The key ID of the hex-encoded certificate credential"
}
