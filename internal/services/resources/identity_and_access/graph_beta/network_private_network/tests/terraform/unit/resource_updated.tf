resource "microsoft365_graph_beta_identity_and_access_network_private_network" "test" {
  name = "unit-test-private-network-updated"

  app_ids = [
    "81ff8e33-3181-475b-ade6-e27147234bb3"
  ]

  dns_resolution_identification = {
    dns_servers     = ["8.8.8.8"]
    fqdn_to_resolve = "example.com"

    expected_ip_resolutions = [
      {
        type  = "ip_address"
        value = "192.168.1.11"
      },
      {
        type  = "ip_subnet"
        value = "192.168.1.1/16"
      },
      {
        type          = "ip_range"
        begin_address = "192.168.1.1"
        end_address   = "192.168.1.2"
      }
    ]
  }
}
