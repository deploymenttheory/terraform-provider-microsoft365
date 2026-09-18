# Configure TLS inspection, licensing and access policies before enabling EFP.
# Set is_enabled = true before creating hosted custom PAC files.
resource "microsoft365_graph_beta_identity_and_access_network_explicit_forward_proxy" "example" {
  internet_access = {
    is_enabled                            = false
    is_source_ip_session_affinity_enabled = true
    source_ip_session_affinity_options    = "useSessionId"
  }
}
