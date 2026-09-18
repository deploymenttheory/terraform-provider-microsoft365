# Complete the EFP prerequisites before applying; this enables tenant-wide EFP.
# Reference an existing EFP resource instead if the singleton is already managed.
resource "microsoft365_graph_beta_identity_and_access_network_explicit_forward_proxy" "example" {
  internet_access = {
    is_enabled                            = true
    is_source_ip_session_affinity_enabled = true
    source_ip_session_affinity_options    = "useSessionId"
  }
}

# proxy.pac can use the literal ${GSAEFP} service placeholder.
resource "microsoft365_graph_beta_identity_and_access_network_proxy_auto_configuration" "example" {
  name       = "branch-office"
  is_enabled = false
  content    = file("${path.module}/proxy.pac")

  # EFP must be enabled before uploading, even when this PAC is disabled.
  depends_on = [microsoft365_graph_beta_identity_and_access_network_explicit_forward_proxy.example]
}
