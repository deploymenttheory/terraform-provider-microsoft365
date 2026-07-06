resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy" "allow_by_default" {
  name           = "Web Content Filtering Policy"
  description    = "Global Secure Access web filtering policy managed by Terraform"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy" "block_by_default" {
  name           = "Block By Default Web Content Filtering Policy"
  description    = "Global Secure Access web filtering policy staged by Terraform"
  default_action = "block"
}
