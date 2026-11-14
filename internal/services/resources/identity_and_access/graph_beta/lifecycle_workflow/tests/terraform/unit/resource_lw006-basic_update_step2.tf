# LW006: Basic property update (PATCH) - step 2
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw006_basic_update" {
  category              = "joiner"
  display_name          = "unit-test-lw006: basic property update renamed"
  description           = "Updated description"
  is_enabled            = true
  is_scheduling_enabled = false

  execution_conditions = {
    on_demand_execution_only = true
  }

  tasks = [
    {
      task_definition_id = "6fc52c9d-398b-4305-9763-15f42c1676fc"
      display_name       = "Enable user account"
    }
  ]
}
