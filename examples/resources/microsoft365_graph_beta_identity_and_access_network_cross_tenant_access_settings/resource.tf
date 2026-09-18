# One instance per tenant. Initial apply updates the existing singleton.
# Review tenant restrictions policies and the access impact before enabling tagging.
resource "microsoft365_graph_beta_identity_and_access_network_cross_tenant_access_settings" "example" {
  network_packet_tagging_status = "disabled"
}

# Destroy only releases Terraform management; it does not change these settings.
