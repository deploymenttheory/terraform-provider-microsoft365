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

# LW007: Task and trigger change (createNewVersion) - step 1
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw007_task_change" {
  category     = "joiner"
  display_name = "acc-test-lw007: task change ${random_string.suffix.result}"

  execution_conditions = {
    scope = {
      rule = "(department eq 'acc-test-none')"
    }
    trigger = {
      time_based_attribute = {
        attribute      = "employeeHireDate"
        offset_in_days = 1
      }
    }
  }

  tasks = [
    {
      task_definition_id = "6fc52c9d-398b-4305-9763-15f42c1676fc"
      display_name       = "Enable user account"
    },
    {
      task_definition_id = "70b29d51-b59a-4773-9280-8841dfd3f2ea"
      display_name       = "Send welcome email"
    }
  ]
}
