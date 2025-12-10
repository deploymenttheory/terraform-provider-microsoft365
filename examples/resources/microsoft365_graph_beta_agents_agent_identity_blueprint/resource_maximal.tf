# Maximal Agent Identity Blueprint configuration
# Creates an agent identity blueprint with all available fields configured
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "maximal" {
  display_name = "Production AI Agent Blueprint"
  description  = "Blueprint for AI agents used in production workloads with full governance controls"

  sponsor_user_ids = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002",
  ]

  owner_user_ids = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002",
  ]

  tags = [
    "production",
    "ai-agent",
    "managed-by-terraform"
  ]

  # When true, permanently deletes from Entra ID on destroy (cannot be restored)
  # When false, moves to deleted items (can be restored within 30 days)
  hard_delete = true

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

