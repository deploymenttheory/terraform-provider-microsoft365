# Use the client-resolved destination IP for Microsoft 365 traffic.
# This DNS setting is independent of EFP enablement and PAC hosting.
# Initial apply updates existing settings; destroy only removes state.
resource "microsoft365_graph_beta_identity_and_access_network_forwarding_options" "example" {
  skip_dns_lookup_state = "enabled"
}
