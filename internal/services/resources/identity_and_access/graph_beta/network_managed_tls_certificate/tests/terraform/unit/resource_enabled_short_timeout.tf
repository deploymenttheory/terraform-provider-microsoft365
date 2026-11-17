resource "microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate" "test" {
  name              = "M-TLSi-unit-slow"
  common_name       = "Terraform TLS Inspection Root CA"
  organization_name = "Example"
  enabled           = true

  timeouts = {
    create = "2s"
    update = "2s"
  }
}
