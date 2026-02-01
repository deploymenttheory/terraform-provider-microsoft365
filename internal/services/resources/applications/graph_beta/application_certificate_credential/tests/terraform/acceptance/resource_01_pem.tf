# Acceptance test: Application Certificate Credential with PEM encoding
# Full dependency chain: random_string -> application -> time_sleep -> certificate_credential

resource "random_string" "test_id_pem" {
  length  = 8
  special = false
  upper   = false
}

# Minimal application for certificate testing
resource "microsoft365_graph_beta_applications_application" "test_app_pem" {
  display_name = "acc-test-app-cert-pem-${random_string.test_id_pem.result}"
  hard_delete  = true
}

# Wait for eventual consistency after application creation
resource "time_sleep" "wait_for_app_pem" {
  depends_on = [microsoft365_graph_beta_applications_application.test_app_pem]

  create_duration = "30s"
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

resource "microsoft365_graph_beta_applications_application_certificate_credential" "test_pem" {
  application_id = microsoft365_graph_beta_applications_application.test_app_pem.id
  display_name   = "acc-test-certificate-pem-${random_string.test_id_pem.result}"

  key      = tls_self_signed_cert.test_cert_pem.cert_pem
  encoding = "pem"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"

  depends_on = [time_sleep.wait_for_app_pem]
}
