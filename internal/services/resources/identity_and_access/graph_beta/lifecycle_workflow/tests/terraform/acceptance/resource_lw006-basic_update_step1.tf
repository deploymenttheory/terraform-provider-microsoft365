# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Lifecycle Workflow
# ==============================================================================

# LW006: Basic property update (PATCH) - step 1
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw006_basic_update" {
  category              = "joiner"
  display_name          = "acc-test-lw006: basic property update ${random_string.suffix.result}"
  description           = "Initial description"
  is_enabled            = false
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
