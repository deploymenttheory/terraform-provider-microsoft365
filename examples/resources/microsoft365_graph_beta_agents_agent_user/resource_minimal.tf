# Minimal Agent User Example
# This example shows only the required fields for creating an agent user

resource "microsoft365_graph_beta_agents_agent_user" "example" {
  display_name        = "Example Agent User"
  agent_identity_id   = "00000000-0000-0000-0000-000000000000" # ID of parent agent identity
  account_enabled     = true
  user_principal_name = "agent-user@contoso.com" # Must match your tenant's verified domain
  mail_nickname       = "agent-user"
  sponsor_ids         = ["11111111-1111-1111-1111-111111111111"] # User ID of sponsor
  hard_delete         = true
}

