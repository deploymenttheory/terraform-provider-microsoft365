resource "microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate" "test" {
  enabled = true
}

data "microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate" "test" {
  id = microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate.test.id
}
