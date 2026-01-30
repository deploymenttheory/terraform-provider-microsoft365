# Acceptance test: Application Certificate Credential with DER encoding (base64)
# DER is binary format, passed as base64-encoded string (e.g., using filebase64())
# Full dependency chain: random_string -> application -> time_sleep -> certificate_credential

resource "random_string" "test_id_der" {
  length  = 8
  special = false
  upper   = false
}

# Minimal application for certificate testing
resource "microsoft365_graph_beta_applications_application" "test_app_der" {
  display_name = "acc-test-app-cert-der-${random_string.test_id_der.result}"
  hard_delete  = true

  lifecycle {
    ignore_changes = [key_credentials]
  }
}

# Wait for eventual consistency after application creation
resource "time_sleep" "wait_for_app_der" {
  depends_on = [microsoft365_graph_beta_applications_application.test_app_der]

  create_duration = "30s"
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

resource "microsoft365_graph_beta_applications_application_certificate_credential" "test_der" {
  application_id = microsoft365_graph_beta_applications_application.test_app_der.id
  display_name   = "acc-test-certificate-der-${random_string.test_id_der.result}"

  key      = local.cert_der_base64
  encoding = "base64"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"

  depends_on = [time_sleep.wait_for_app_der]
}
