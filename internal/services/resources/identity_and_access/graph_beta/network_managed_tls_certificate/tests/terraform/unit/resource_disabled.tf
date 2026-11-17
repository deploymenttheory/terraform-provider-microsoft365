resource "microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate" "test" {
  name              = "M-TLSi-unit-enabled"
  common_name       = "Terraform TLS Inspection Root CA"
  organization_name = "Example"
  validity_months   = 60
  enabled           = false

  timeouts = {
    update = "15s"
  }
}
