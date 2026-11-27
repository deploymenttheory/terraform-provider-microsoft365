# Provision the manifest through an agent instance first.
# This resource adopts and updates it. Destroy only removes Terraform state.
resource "microsoft365_graph_beta_agents_agent_card_manifest" "example" {
  agent_card_manifest_id = "existing-manifest-id"
  display_name           = "Employee assistant"
  description            = "Answers employee questions"
  protocol_version       = "1.0"
  version                = "1.0.0"
  default_input_modes    = ["text"]
  default_output_modes   = ["text"]
  agent_provider = {
    organization = "Contoso"
    url          = "https://contoso.com"
  }
  capabilities = {
    streaming = true
  }
}
