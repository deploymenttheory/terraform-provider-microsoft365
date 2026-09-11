# Creating policies and rules does not assign them to traffic.
resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy" "example" {
  name        = "Example prompt policy"
  description = "Prompt protection for AI applications"
}
