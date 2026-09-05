resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy" "test" {
  description = ""
  name        = "tf-api-probe-policy-updated"

  default_action = "allow"
}
