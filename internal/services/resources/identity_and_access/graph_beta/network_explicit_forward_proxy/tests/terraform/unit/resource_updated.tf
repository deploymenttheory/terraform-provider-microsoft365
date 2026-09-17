resource "microsoft365_graph_beta_identity_and_access_network_explicit_forward_proxy" "test" {
  internet_access = {
    is_enabled                            = false
    is_source_ip_session_affinity_enabled = false
    source_ip_session_affinity_options    = "none"
  }
}
