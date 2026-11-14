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

# LW002: Mover workflow triggered by a department change
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw002_mover" {
  category     = "mover"
  display_name = "acc-test-lw002: mover department change ${random_string.suffix.result}"
  description  = "Notify the manager when a user changes department"

  execution_conditions = {
    scope = {
      rule = "(department eq 'acc-test-none')"
    }
    trigger = {
      attribute_change = {
        trigger_attribute = "department"
      }
    }
  }

  tasks = [
    {
      task_definition_id = "aab41899-9972-422a-9d97-f626014578b7"
      display_name       = "Send email to notify manager of user move"
    }
  ]
}
