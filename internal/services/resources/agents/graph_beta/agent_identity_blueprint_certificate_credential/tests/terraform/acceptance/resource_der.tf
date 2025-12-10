# Acceptance test: Agent Identity Blueprint Certificate Credential with DER encoding (base64)
# DER is binary format, passed as base64-encoded string (e.g., using filebase64())
# Full dependency chain: random_string -> users -> agent_identity_blueprint -> certificate_credential

resource "random_string" "test_id_der" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "dependency_user_der_1" {
  display_name        = "acc-test-cert-der-user1-${random_string.test_id_der.result}"
  user_principal_name = "acc-test-cert-der-user1-${random_string.test_id_der.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-cert-der-user1-${random_string.test_id_der.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_users_user" "dependency_user_der_2" {
  display_name        = "acc-test-cert-der-user2-${random_string.test_id_der.result}"
  user_principal_name = "acc-test-cert-der-user2-${random_string.test_id_der.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-cert-der-user2-${random_string.test_id_der.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_blueprint_der" {
  display_name = "acc-test-blueprint-cert-der-${random_string.test_id_der.result}"
  description  = "Agent identity blueprint for DER certificate credential acceptance test"
  sponsor_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_der_1.id,
    microsoft365_graph_beta_users_user.dependency_user_der_2.id,
  ]
  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_der_1.id,
    microsoft365_graph_beta_users_user.dependency_user_der_2.id,
  ]
  hard_delete = true
}

# Generate a self-signed certificate for testing
resource "tls_private_key" "test_key_der" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test_cert_der" {
  private_key_pem = tls_private_key.test_key_der.private_key_pem

  subject {
    common_name  = "acc-test-certificate-der-${random_string.test_id_der.result}"
    organization = "Terraform Provider Test"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

# Extract the base64 content from PEM (strip headers) - simulates DER format
locals {
  cert_der_base64 = replace(
    replace(
      replace(tls_self_signed_cert.test_cert_der.cert_pem, "-----BEGIN CERTIFICATE-----", ""),
      "-----END CERTIFICATE-----", ""
    ),
    "\n", ""
  )
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_certificate_credential" "test_der" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_blueprint_der.id
  display_name = "acc-test-certificate-der-${random_string.test_id_der.result}"

  key      = local.cert_der_base64
  encoding = "base64"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"
}
