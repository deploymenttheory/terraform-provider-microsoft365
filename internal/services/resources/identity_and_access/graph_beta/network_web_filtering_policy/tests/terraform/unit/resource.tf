resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy" "test" {
  name           = "unit-test-web-filtering-policy"
  description    = "Unit test web filtering policy"
  default_action = "allow"
}
