resource "microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate" "test" {
  name              = "M-TLSi-unit-enabled"
  common_name       = "Terraform TLS Inspection Root CA"
  organization_name = "Contoso"
  validity_months   = 60
  enabled           = true
}
