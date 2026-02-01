# Acceptance test: Application Certificate Credential - Maximal configuration
# Tests all optional attributes with explicit values
# Full dependency chain: random_string -> application -> time_sleep -> certificate_credential

resource "random_string" "test_id_maximal" {
  length  = 8
  special = false
  upper   = false
}

# Minimal application for certificate testing
resource "microsoft365_graph_beta_applications_application" "test_app_maximal" {
  display_name = "acc-test-app-cert-max-${random_string.test_id_maximal.result}"
  hard_delete  = true
}

# Wait for eventual consistency after application creation
resource "time_sleep" "wait_for_app_maximal" {
  depends_on = [microsoft365_graph_beta_applications_application.test_app_maximal]

  create_duration = "30s"
}

# Generate a self-signed certificate for testing
resource "tls_private_key" "test_key_maximal" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test_cert_maximal" {
  private_key_pem = tls_private_key.test_key_maximal.private_key_pem

  subject {
    common_name  = "acc-test-certificate-maximal-${random_string.test_id_maximal.result}"
    organization = "Terraform Provider Test"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

resource "microsoft365_graph_beta_applications_application_certificate_credential" "test_maximal" {
  application_id = microsoft365_graph_beta_applications_application.test_app_maximal.id
  display_name   = "acc-test-certificate-maximal-${random_string.test_id_maximal.result}"

  key      = tls_self_signed_cert.test_cert_maximal.cert_pem
  encoding = "pem"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"

  start_date_time = "2027-01-01T00:00:00Z"
  end_date_time   = "2029-01-01T00:00:00Z"

  depends_on = [time_sleep.wait_for_app_maximal]
}
