resource "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy" "example" {
  name           = "example-cloud-firewall-policy"
  description    = "Managed by Terraform"
  default_action = "allow"
}
