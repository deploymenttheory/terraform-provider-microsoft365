# Example: Basic Agent Identity with Tags
#
# This example shows the minimal configuration for an agent identity
# with optional tags for categorization.
#
# Prerequisites:
# - An existing Agent Identity Blueprint with app_id
# - The Agent Identity Blueprint must have a service principal created
# - At least one user to assign as sponsor and owner

resource "microsoft365_graph_beta_agents_agent_identity" "basic" {
  display_name                = "My Agent Identity"
  agent_identity_blueprint_id = "00000000-0000-0000-0000-000000000000" # Replace with blueprint app_id
  account_enabled             = true
  sponsor_ids                 = ["00000000-0000-0000-0000-000000000001"] # Replace with user IDs
  owner_ids                   = ["00000000-0000-0000-0000-000000000001"] # Replace with user IDs
  tags                        = ["production", "customer-service", "ai-agent"]

  # When true, permanently deletes from Entra ID on destroy (cannot be restored)
  # When false, moves to deleted items (can be restored within 30 days)
  hard_delete = true
}

