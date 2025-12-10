# Acceptance test: Agent Identity Blueprint Certificate Credential with PEM encoding
# Full dependency chain: random_string -> users -> agent_identity_blueprint -> certificate_credential

resource "random_string" "test_id_pem" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "dependency_user_pem_1" {
  display_name        = "acc-test-cert-pem-user1-${random_string.test_id_pem.result}"
  user_principal_name = "acc-test-cert-pem-user1-${random_string.test_id_pem.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-cert-pem-user1-${random_string.test_id_pem.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_users_user" "dependency_user_pem_2" {
  display_name        = "acc-test-cert-pem-user2-${random_string.test_id_pem.result}"
  user_principal_name = "acc-test-cert-pem-user2-${random_string.test_id_pem.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-cert-pem-user2-${random_string.test_id_pem.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_blueprint_pem" {
  display_name = "acc-test-blueprint-cert-pem-${random_string.test_id_pem.result}"
  description  = "Agent identity blueprint for PEM certificate credential acceptance test"
  sponsor_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_pem_1.id,
    microsoft365_graph_beta_users_user.dependency_user_pem_2.id,
  ]
  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_pem_1.id,
    microsoft365_graph_beta_users_user.dependency_user_pem_2.id,
  ]
  hard_delete = true
}

# Generate a self-signed certificate for testing
resource "tls_private_key" "test_key_pem" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test_cert_pem" {
  private_key_pem = tls_private_key.test_key_pem.private_key_pem

  subject {
    common_name  = "acc-test-certificate-pem-${random_string.test_id_pem.result}"
    organization = "Terraform Provider Test"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_certificate_credential" "test_pem" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_blueprint_pem.id
  display_name = "acc-test-certificate-pem-${random_string.test_id_pem.result}"

  key      = tls_self_signed_cert.test_cert_pem.cert_pem
  encoding = "pem"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"
}
