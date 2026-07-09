resource "microsoft365_graph_beta_identity_and_access_network_private_network" "test" {
  name = "unit-test-private-network-invalid-range"

  dns_resolution_identification = {
    dns_servers     = ["8.8.8.8"]
    fqdn_to_resolve = "example.com"

    expected_ip_resolutions = [
      {
        type  = "ip_range"
        value = "192.168.1.1"
      }
    ]
  }
}
