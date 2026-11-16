# Mover workflow scoped to members of a group, triggered when a user is added to it
resource "microsoft365_graph_beta_groups_group" "sales_team" {
  display_name     = "Sales Team"
  mail_nickname    = "sales-team"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "group_membership" {
  category              = "mover"
  display_name          = "Sales team onboarding"
  description           = "Runs when a user is added to the Sales Team group"
  is_enabled            = true
  is_scheduling_enabled = false

  execution_conditions = {
    scope = {
      group_ids = [microsoft365_graph_beta_groups_group.sales_team.id]
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
      description        = "Add the user to the sales resources group"
      arguments = [
        {
          name  = "groupID"
          value = "00000000-0000-0000-0000-000000000000"
        }
      ]
    },
    {
      task_definition_id = "aab41899-9972-422a-9d97-f626014578b7"
      display_name       = "Send email to notify manager of user move"
      continue_on_error  = true
      arguments = [
        {
          name  = "locale"
          value = "en-us"
        }
      ]
    }
  ]

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}
