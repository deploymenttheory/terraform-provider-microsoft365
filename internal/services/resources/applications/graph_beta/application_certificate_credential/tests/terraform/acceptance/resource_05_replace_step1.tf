# Acceptance test Step 1: Deploy 2 pre-existing certificates
# These will be replaced in step 2

resource "random_string" "test_id_replace" {
  length  = 8
  special = false
  upper   = false
}

# Minimal application for certificate testing
resource "microsoft365_graph_beta_applications_application" "test_app_replace" {
  display_name = "acc-test-app-cert-replace-${random_string.test_id_replace.result}"
  hard_delete  = true

  lifecycle {
    ignore_changes = [key_credentials]
  }
}

# Wait for eventual consistency after application creation
resource "time_sleep" "wait_for_app_replace" {
  depends_on = [microsoft365_graph_beta_applications_application.test_app_replace]

  create_duration = "30s"
}

# Generate first certificate (pre-existing)
resource "tls_private_key" "test_key_replace_1" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test_cert_replace_1" {
  private_key_pem = tls_private_key.test_key_replace_1.private_key_pem

  subject {
    common_name  = "acc-test-certificate-replace-1-${random_string.test_id_replace.result}"
    organization = "Terraform Provider Test"
  }

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

# First pre-existing certificate
resource "microsoft365_graph_beta_applications_application_certificate_credential" "test_replace_1" {
  application_id = microsoft365_graph_beta_applications_application.test_app_replace.id
  display_name   = "acc-test-certificate-replace-1-${random_string.test_id_replace.result}"

  key      = tls_self_signed_cert.test_cert_replace_1.cert_pem
  encoding = "pem"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"

  replace_existing_certificates = false

  depends_on = [time_sleep.wait_for_app_replace]
}

# Generate second certificate (pre-existing)
resource "tls_private_key" "test_key_replace_2" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test_cert_replace_2" {
  private_key_pem = tls_private_key.test_key_replace_2.private_key_pem

  subject {
    common_name  = "acc-test-certificate-replace-2-${random_string.test_id_replace.result}"
    organization = "Terraform Provider Test"
  }

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

# Second pre-existing certificate
resource "microsoft365_graph_beta_applications_application_certificate_credential" "test_replace_2" {
  application_id = microsoft365_graph_beta_applications_application.test_app_replace.id
  display_name   = "acc-test-certificate-replace-2-${random_string.test_id_replace.result}"

  key      = tls_self_signed_cert.test_cert_replace_2.cert_pem
  encoding = "pem"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"

  replace_existing_certificates = false

  depends_on = [microsoft365_graph_beta_applications_application_certificate_credential.test_replace_1]
}
