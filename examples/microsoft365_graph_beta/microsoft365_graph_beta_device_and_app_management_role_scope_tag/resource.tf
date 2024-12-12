# Example of creating a basic role scope tag
resource "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "helpdesk" {
  display_name = "Helpdesk Support Tag"
  description  = "Role scope tag for helpdesk support staff"

  timeouts ={
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example of creating multiple related role scope tags
resource "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "it_support" {
  display_name = "IT Support Tag"
  description  = "Role scope tag for IT support teams"
}

resource "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "device_management" {
  display_name = "Device Management Tag"
  description  = "Role scope tag for device management teams"
}

# Example showing data source usage to reference an existing role scope tag
data "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "existing" {
  display_name = "Existing Tag"
}

# Example of using variables with role scope tags
variable "support_teams" {
  type = list(object({
    name        = string
    description = string
  }))
  default = [
    {
      name        = "Level1-Support"
      description = "First level support team scope"
    },
    {
      name        = "Level2-Support"
      description = "Second level support team scope"
    }
  ]
}

# Creating multiple tags using variables
resource "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "support_teams" {
  for_each = { for team in var.support_teams : team.name => team }

  display_name = each.value.name
  description  = each.value.description
}

# Output examples
output "helpdesk_tag_id" {
  value = microsoft365_graph_beta_device_and_app_management_role_scope_tag.helpdesk.id
}

output "all_support_team_ids" {
  value = [for tag in microsoft365_graph_beta_device_and_app_management_role_scope_tag.support_teams : tag.id]
}