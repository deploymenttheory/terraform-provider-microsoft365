resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "enabled_profile" {
  name        = "Enabled Security Profile"
  description = "Global Secure Access security profile managed by Terraform"
  priority    = 100
  state       = "enabled"
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile" "disabled_profile" {
  name        = "Disabled Security Profile"
  description = "Global Secure Access security profile staged by Terraform"
  priority    = 200
  state       = "disabled"
}
