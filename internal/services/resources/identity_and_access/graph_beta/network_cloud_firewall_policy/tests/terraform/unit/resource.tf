resource "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy" "test" {
  name           = "tf-api-probe-cfw-01a07bf8-cloud-firewall-policy"
  description    = "Managed by Terraform"
  default_action = "allow"
}
