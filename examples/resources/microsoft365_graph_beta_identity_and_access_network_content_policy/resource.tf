resource "microsoft365_graph_beta_identity_and_access_network_content_policy" "example" {
  name           = "Content Policy"
  description    = "Global Secure Access content policy managed by Terraform"
  default_action = "allow"
}
