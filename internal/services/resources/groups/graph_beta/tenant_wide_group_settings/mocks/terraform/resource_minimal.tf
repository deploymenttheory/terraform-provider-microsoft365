# Minimal configuration for tenant-wide group settings
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "test" {
  # Group.Unified template ID
  template_id = "62375ab9-6b52-47ed-826b-58e47e0e304b"

  # Only essential settings
  values = [
    {
      name  = "EnableGroupCreation"
      value = "true" # Allow users to create Microsoft 365 groups
    },
    {
      name  = "AllowGuestsToAccessGroups"
      value = "true" # Allow guest access to groups
    },
    {
      name  = "AllowToAddGuests"
      value = "true" # Allow adding guests to groups
    }
  ]
} 