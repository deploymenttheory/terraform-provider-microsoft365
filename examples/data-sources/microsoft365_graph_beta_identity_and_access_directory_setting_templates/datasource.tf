# Example: Query directory setting templates
# This data source retrieves information about available directory setting templates in Microsoft 365
# These templates define the settings that can be applied to Microsoft 365 groups and other resources

# Get all directory setting templates
data "microsoft365_graph_beta_identity_and_access_directory_setting_templates" "all" {
  filter_type = "all"

  timeouts = {
    read = "1m"
  }
}

# Get a specific directory setting template by ID
# This example uses the Group.Unified template ID
data "microsoft365_graph_beta_identity_and_access_directory_setting_templates" "group_unified" {
  filter_type  = "id"
  filter_value = "62375ab9-6b52-47ed-826b-58e47e0e304b" # Group.Unified template ID

  timeouts = {
    read = "1m"
  }
}

# Get directory setting templates by display name filter
data "microsoft365_graph_beta_identity_and_access_directory_setting_templates" "group_templates" {
  filter_type  = "display_name"
  filter_value = "Group.Unified"

  timeouts = {
    read = "1m"
  }
}

# Output all template IDs and display names
output "all_template_ids" {
  value = [for template in data.microsoft365_graph_beta_identity_and_access_directory_setting_templates.all.directory_setting_templates : {
    id           = template.id
    display_name = template.display_name
  }]
  description = "List of all directory setting template IDs and display names"
}

# Output the Group.Unified template details
output "group_unified_template" {
  value = {
    id           = data.microsoft365_graph_beta_identity_and_access_directory_setting_templates.group_unified.directory_setting_templates[0].id
    display_name = data.microsoft365_graph_beta_identity_and_access_directory_setting_templates.group_unified.directory_setting_templates[0].display_name
    description  = data.microsoft365_graph_beta_identity_and_access_directory_setting_templates.group_unified.directory_setting_templates[0].description
    settings = [for value in data.microsoft365_graph_beta_identity_and_access_directory_setting_templates.group_unified.directory_setting_templates[0].values : {
      name          = value.name
      type          = value.type
      default_value = value.default_value
      description   = value.description
    }]
  }
  description = "Details of the Group.Unified template including all available settings"
} 