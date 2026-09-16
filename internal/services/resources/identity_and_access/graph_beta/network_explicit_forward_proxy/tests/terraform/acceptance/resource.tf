resource "microsoft365_graph_beta_identity_and_access_network_explicit_forward_proxy" "test" {
  internet_access = {
    is_enabled                            = var.proxy_enabled
    is_source_ip_session_affinity_enabled = var.affinity_enabled
    source_ip_session_affinity_options    = var.affinity_options
  }
}
