# Basic lookup examples
# Look up by display name
data "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "by_name" {
  display_name = "Level1-Support"
}

# Look up by ID
data "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "by_id" {
  id = "00000000-0000-0000-0000-000000000001"
}

# Output showing role scope tag details
output "role_scope_tag_details" {
  value = {
    # Basic details
    id           = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.id
    display_name = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.display_name
    description  = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.description
    is_built_in  = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.is_built_in
  }
}

# Use Case 1: Create new tag based on existing one
resource "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "clone" {
  display_name = "Clone - ${data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.display_name}"
  description  = "Cloned from: ${data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.description}"
  
  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Use Case 2: Conditional tag creation based on built-in status
resource "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "conditional" {
  count = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.is_built_in ? 0 : 1
  
  display_name = "Custom - ${data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.display_name}"
  description  = "Custom version of non-built-in tag"
  
  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Use Case 3: Look up multiple tags and compare
data "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "level1" {
  display_name = "Level1-Support"
}

data "microsoft365_graph_beta_device_and_app_management_role_scope_tag" "level2" {
  display_name = "Level2-Support"
}

output "support_tags_comparison" {
  value = {
    level1 = {
      id          = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.level1.id
      display_name = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.level1.display_name
      is_built_in  = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.level1.is_built_in
    }
    level2 = {
      id          = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.level2.id
      display_name = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.level2.display_name
      is_built_in  = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.level2.is_built_in
    }
    comparison = {
      both_built_in = (
        data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.level1.is_built_in && 
        data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.level2.is_built_in
      )
    }
  }
}

# Use Case 4: Create dynamic outputs based on tag properties
output "tag_summary" {
  value = {
    tag_info = {
      name        = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.display_name
      type        = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.is_built_in ? "Built-in" : "Custom"
      has_description = data.microsoft365_graph_beta_device_and_app_management_role_scope_tag.by_name.description != ""
    }
  }
}