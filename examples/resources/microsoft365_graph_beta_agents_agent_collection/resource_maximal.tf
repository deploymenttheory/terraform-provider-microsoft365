# Maximal Agent Collection configuration
# Creates an agent collection with all available fields configured
resource "microsoft365_graph_beta_agents_agent_collection" "maximal" {
  display_name = "IT Automation Agent Collection"
  owner_ids = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]
  description       = "A collection of IT automation agents for managing infrastructure and support workflows"
  managed_by        = "00000000-0000-0000-0000-000000000003"
  originating_store = "Terraform"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
