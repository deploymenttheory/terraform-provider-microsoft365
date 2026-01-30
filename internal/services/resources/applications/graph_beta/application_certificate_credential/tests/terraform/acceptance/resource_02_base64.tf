# Acceptance test: Application Certificate Credential with base64 encoding
# Full dependency chain: random_string -> application -> time_sleep -> certificate_credential

resource "random_string" "test_id_base64" {
  length  = 8
  special = false
  upper   = false
}

# Minimal application for certificate testing
resource "microsoft365_graph_beta_applications_application" "test_app_base64" {
  display_name = "acc-test-app-cert-b64-${random_string.test_id_base64.result}"
  hard_delete  = true

  lifecycle {
    ignore_changes = [key_credentials]
  }
}

# Wait for eventual consistency after application creation
resource "time_sleep" "wait_for_app_base64" {
  depends_on = [microsoft365_graph_beta_applications_application.test_app_base64]

  create_duration = "30s"
}

# Generate a self-signed certificate for testing
resource "tls_private_key" "test_key_base64" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test_cert_base64" {
  private_key_pem = tls_private_key.test_key_base64.private_key_pem

  subject {
    common_name  = "acc-test-certificate-base64-${random_string.test_id_base64.result}"
    organization = "Terraform Provider Test"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

resource "microsoft365_graph_beta_applications_application_certificate_credential" "test_base64" {
  application_id = microsoft365_graph_beta_applications_application.test_app_base64.id
  display_name   = "acc-test-certificate-base64-${random_string.test_id_base64.result}"

  # Convert PEM to base64-encoded DER format
  key      = base64encode(tls_self_signed_cert.test_cert_base64.cert_pem)
  encoding = "base64"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"

  depends_on = [time_sleep.wait_for_app_base64]
}
