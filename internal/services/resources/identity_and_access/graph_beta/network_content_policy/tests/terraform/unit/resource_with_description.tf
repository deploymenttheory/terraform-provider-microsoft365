resource "microsoft365_graph_beta_identity_and_access_network_content_policy" "test" {
  name           = "unit-test-content-policy-with-description"
  description    = "managed by Terraform"
  default_action = "allow"
}
