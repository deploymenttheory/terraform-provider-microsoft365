# LW005: Maximal mover workflow with a group scope, membership change trigger and every task attribute
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw005_maximal" {
  category              = "mover"
  display_name          = "unit-test-lw005: maximal group membership"
  description           = "Runs when a user is added to the scoped group"
  is_enabled            = true
  is_scheduling_enabled = false

  execution_conditions = {
    on_demand_execution_only = false
    scope = {
      group_ids = ["11111111-1111-1111-1111-111111111111"]
    }
    trigger = {
      membership_change = {
        change_type = "add"
      }
    }
  }

  tasks = [
    {
      task_definition_id = "22085229-5809-45e8-97fd-270d28d66910"
      display_name       = "Add user to groups"
      description        = "Add the user to the team group"
      is_enabled         = true
      continue_on_error  = false
      arguments = [
        {
          name  = "groupID"
          value = "22222222-2222-2222-2222-222222222222"
        }
      ]
    },
    {
      task_definition_id = "aab41899-9972-422a-9d97-f626014578b7"
      display_name       = "Send email to notify manager of user move"
      description        = "Tell the manager about the move"
      is_enabled         = false
      continue_on_error  = true
      arguments = [
        {
          name  = "locale"
          value = "en-us"
        },
        {
          name  = "customSubject"
          value = "A user has moved"
        }
      ]
    }
  ]
}
