# Example: Configure group-specific settings for Microsoft 365 groups
# This resource manages directory settings for a specific Microsoft 365 group

# Get the Group.Unified.Guest template details to see available settings
data "microsoft365_graph_beta_identity_and_access_directory_setting_templates" "group_unified_guest" {
  filter_type  = "display_name"
  filter_value = "Group.Unified.Guest"
}

# Data source to get the group ID
data "microsoft365_graph_beta_groups_group" "example" {
  display_name = "Marketing Team"
}

# Example 1: Configure group-specific guest access settings
# This overrides tenant-wide guest settings for this specific group
resource "microsoft365_graph_beta_groups_group_settings" "guest_settings" {
  group_id    = data.microsoft365_graph_beta_groups_group.example.id
  template_id = data.microsoft365_graph_beta_identity_and_access_directory_setting_templates.group_unified_guest.directory_setting_templates[0].id

  values = [
    {
      name  = "AllowToAddGuests"
      value = "false" # Disable guest access for this specific group
    }
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Get the Group.Unified template details
data "microsoft365_graph_beta_identity_and_access_directory_setting_templates" "group_unified" {
  filter_type  = "display_name"
  filter_value = "Group.Unified"
}

# Example 2: Configure group-specific unified settings
# This shows how to override other group settings for a specific group
resource "microsoft365_graph_beta_groups_group_settings" "unified_settings" {
  group_id    = data.microsoft365_graph_beta_groups_group.example.id
  template_id = data.microsoft365_graph_beta_identity_and_access_directory_setting_templates.group_unified.directory_setting_templates[0].id

  values = [
    {
      name  = "ClassificationList"
      value = "Confidential,Secret,Top Secret" # Custom classifications for this group
    },
    {
      name  = "DefaultClassification"
      value = "Confidential" # Default classification for this group
    },
    {
      name  = "UsageGuidelinesUrl"
      value = "https://contoso.com/marketing-group-guidelines" # Group-specific guidelines
    }
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Output the created settings IDs
output "guest_settings_id" {
  value       = microsoft365_graph_beta_groups_group_settings.guest_settings.id
  description = "The ID of the created group-specific guest settings"
}

output "unified_settings_id" {
  value       = microsoft365_graph_beta_groups_group_settings.unified_settings.id
  description = "The ID of the created group-specific unified settings"
}

# NOTE: In a real environment, you would typically create only one setting per template per group.
# The two resources shown here are for demonstration purposes only. 