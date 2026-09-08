resource "microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate" "test" {
  name              = "M-TLSi-acc04"
  common_name       = "Terraform Acceptance Updated TLS Inspection Root CA"
  organization_name = "Deployment Theory Acceptance"
  enabled           = true
}

data "microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate" "test" {
  id = microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate.test.id
}
