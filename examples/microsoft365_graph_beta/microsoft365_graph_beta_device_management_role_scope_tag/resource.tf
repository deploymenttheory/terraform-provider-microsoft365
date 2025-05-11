# Example of creating a basic role scope tag with a group assignment
resource "microsoft365_graph_beta_device_management_role_scope_tag" "helpdesk" {
  display_name = "Helpdesk Support Tag"
  description  = "Role scope tag for helpdesk support staff"

  assignments = ["00000000-0000-0000-0000-000000000001"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example of creating multiple related role scope tags with assignments
resource "microsoft365_graph_beta_device_management_role_scope_tag" "it_support" {
  display_name = "IT Support Tag"
  description  = "Role scope tag for IT support teams"

  assignments = ["00000000-0000-0000-0000-000000000002"]
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "device_management" {
  display_name = "Device Management Tag"
  description  = "Role scope tag for device management teams"

  assignments = [
    "00000000-0000-0000-0000-000000000003",
    "00000000-0000-0000-0000-000000000004"
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example showing data source usage to reference an existing role scope tag
data "microsoft365_graph_beta_device_management_role_scope_tag" "existing" {
  display_name = "Existing Tag"
}

# Example of using variables with role scope tags including assignments
variable "support_teams" {
  type = list(object({
    name        = string
    description = string
    group_ids   = list(string)
  }))
  default = [
    {
      name        = "Level1-Support"
      description = "First level support team scope"
      group_ids   = ["00000000-0000-0000-0000-000000000005"]
    },
    {
      name        = "Level2-Support"
      description = "Second level support team scope"
      group_ids   = ["00000000-0000-0000-0000-000000000006", "00000000-0000-0000-0000-000000000007"]
    }
  ]
}

# Creating multiple tags using variables
resource "microsoft365_graph_beta_device_management_role_scope_tag" "support_teams" {
  for_each = { for team in var.support_teams : team.name => team }

  display_name = each.value.name
  description  = each.value.description
  assignments  = each.value.group_ids

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Output examples
output "helpdesk_tag_id" {
  value = microsoft365_graph_beta_device_management_role_scope_tag.helpdesk.id
}

output "all_support_team_ids" {
  value = [for tag in microsoft365_graph_beta_device_management_role_scope_tag.support_teams : tag.id]
}

# Example of a role scope tag with conditional assignments based on environment
variable "environment" {
  type    = string
  default = "production"
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "environment_specific" {
  display_name = "Environment-Specific Support Tag"
  description  = "Role scope tag for ${var.environment} environment"

  assignments = (var.environment == "production"
    ? ["00000000-0000-0000-0000-000000000008"]
    : ["00000000-0000-0000-0000-000000000009"]
  )

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}