# Example: Create a certificate credential with PEM encoding (default)

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "my-agent-blueprint"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for certificate authentication"
  hard_delete      = true
}

# Generate a self-signed certificate using the tls provider
resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "example" {
  private_key_pem = tls_private_key.example.private_key_pem

  subject {
    common_name  = "my-agent-certificate"
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
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_certificate_credential" "example" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  display_name = "my-agent-certificate"

  key      = tls_self_signed_cert.example.cert_pem
  encoding = "pem" # Default, can be omitted
  type     = "AsymmetricX509Cert"
  usage    = "Verify"
}

output "certificate_key_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_certificate_credential.example.key_id
  description = "The key ID of the certificate credential"
}

output "certificate_thumbprint" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_certificate_credential.example.thumbprint
  description = "The thumbprint (SHA-1 hash) of the certificate"
}

