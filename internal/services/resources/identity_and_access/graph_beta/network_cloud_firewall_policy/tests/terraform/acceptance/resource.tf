resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy" "test" {
  name           = "tf-api-probe-cfw-01a07bf8-${random_string.suffix.result}"
  description    = "Managed by Terraform"
  default_action = "allow"
}
