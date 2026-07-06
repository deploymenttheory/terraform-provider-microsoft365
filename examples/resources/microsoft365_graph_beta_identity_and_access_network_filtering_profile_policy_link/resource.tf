resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile_policy_link" "web_filtering" {
  filtering_profile_id = microsoft365_graph_beta_identity_and_access_network_filtering_profile.example.id
  policy_id            = microsoft365_graph_beta_identity_and_access_network_web_filtering_policy.example.id
  policy_type          = "web_filtering_policy"
  state                = "enabled"
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile_policy_link" "legacy_filtering" {
  filtering_profile_id = microsoft365_graph_beta_identity_and_access_network_filtering_profile.example.id
  policy_id            = microsoft365_graph_beta_identity_and_access_network_filtering_policy.example.id
  policy_type          = "filtering_policy"
  state                = "enabled"
  priority             = 100
  logging_state        = "enabled"
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_profile_policy_link" "tls_inspection" {
  filtering_profile_id = microsoft365_graph_beta_identity_and_access_network_filtering_profile.example.id
  policy_id            = "00000000-0000-0000-0000-000000000000"
  policy_type          = "tls_inspection_policy"
  state                = "enabled"
}
