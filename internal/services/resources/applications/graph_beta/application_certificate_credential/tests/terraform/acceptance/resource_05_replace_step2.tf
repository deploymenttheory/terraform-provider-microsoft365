# Acceptance test Step 2: Deploy 3rd certificate with replace=true
# This should REMOVE all previous certificates (test_replace_1 and test_replace_2)
# and leave ONLY this certificate

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

# Generate third certificate (the replacement)
resource "tls_private_key" "test_key_replace_3" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "test_cert_replace_3" {
  private_key_pem = tls_private_key.test_key_replace_3.private_key_pem

  subject {
    common_name  = "acc-test-certificate-replace-3-${random_string.test_id_replace.result}"
    organization = "Terraform Provider Test"
  }

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

# Third certificate with replace=true
# This should REMOVE ALL existing certificates and leave only this one
resource "microsoft365_graph_beta_applications_application_certificate_credential" "test_replace_3" {
  application_id = microsoft365_graph_beta_applications_application.test_app_replace.id
  display_name   = "acc-test-certificate-replace-3-${random_string.test_id_replace.result}"

  key      = tls_self_signed_cert.test_cert_replace_3.cert_pem
  encoding = "pem"
  type     = "AsymmetricX509Cert"
  usage    = "Verify"

  # CRITICAL: This should remove ALL existing certificates
  replace_existing_certificates = true

  depends_on = [time_sleep.wait_for_app_replace]
}

# Data source to read the application and validate certificate count
data "microsoft365_graph_beta_applications_application" "verify_replace" {
  object_id = microsoft365_graph_beta_applications_application.test_app_replace.id

  depends_on = [microsoft365_graph_beta_applications_application_certificate_credential.test_replace_3]
}
