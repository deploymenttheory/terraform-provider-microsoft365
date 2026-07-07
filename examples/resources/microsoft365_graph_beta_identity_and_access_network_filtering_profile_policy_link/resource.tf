resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "example" {
  name        = "Example Security Profile"
  description = "Global Secure Access security profile managed by Terraform"
  priority    = 100
  state       = "enabled"
}

resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy" "example" {
  name           = "Example Web Filtering Policy"
  description    = "Global Secure Access web filtering policy managed by Terraform"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile_policy_link" "web_filtering" {
  filtering_profile_id = microsoft365_graph_beta_identity_and_access_network_filtering_profile.example.id
  policy_id            = microsoft365_graph_beta_identity_and_access_network_web_filtering_policy.example.id
  policy_type          = "web_filtering"
  state                = "enabled"
}
