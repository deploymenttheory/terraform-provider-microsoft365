# Example: Create a certificate credential with HEX encoding
# Use this when you have a hex-encoded certificate string

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example_hex" {
  display_name     = "my-agent-blueprint-hex"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for HEX certificate authentication"
  hard_delete      = true
}

# Create a certificate credential using a pre-generated hex-encoded certificate
# To generate a hex-encoded certificate:
#   1. openssl x509 -in cert.pem -outform DER -out cert.der
#   2. xxd -p cert.der | tr -d '\n' > cert.hex
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_certificate_credential" "example_hex" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example_hex.id
  display_name = "my-hex-certificate"

  # Pre-generated hex-encoded certificate string
  # Generate with: openssl x509 -in cert.pem -outform DER | xxd -p | tr -d '\n'
  key      = file("path/to/certificate.hex")
  encoding = "hex"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"
}

output "hex_certificate_key_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_certificate_credential.example_hex.key_id
  description = "The key ID of the hex-encoded certificate credential"
}

