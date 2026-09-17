# The acceptance test substitutes the backed-up current value, never enabling tagging.
resource "microsoft365_graph_beta_identity_and_access_network_cross_tenant_access_settings" "test" {
  network_packet_tagging_status = "{{STATUS}}"
}
