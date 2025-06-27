# Minimal configuration for group-specific settings
resource "microsoft365_graph_beta_groups_group_settings" "minimal" {
  # Test group ID
  group_id = "12345678-1234-1234-1234-123456789012"

  # Group.Unified.Guest template ID
  template_id = "08d542b9-071f-4e16-94b0-74abb372e3d9"

  # Only essential settings
  values = [
    {
      name  = "AllowToAddGuests"
      value = "false" # Disable guest access for this specific group
    }
  ]
} 