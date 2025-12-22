# Example: Windows Quality Update Expedite Policy - Maximal Configuration with Assignments
resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "maximal_with_assignments" {
  display_name       = "Production Critical Update Expedite Policy"
  description        = "Expedited deployment for critical security updates targeting production devices with exclusions for test environments"
  role_scope_tag_ids = ["0", "1"]

  # Required: Expedited update settings
  expedited_update_settings = {
    # Quality update release to expedite
    # Valid values: "2025-12-09T00:00:00Z", "2025-11-20T00:00:00Z"
    quality_update_release = "2025-11-20T00:00:00Z"

    # Force reboot after 2 days to ensure update completion
    days_until_forced_reboot = 2
  }

  # Assignments: Target specific groups and exclude others
  assignments = [
    # Primary production group assignment
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111" # Production Devices Group
    },
    # Secondary production group assignment
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222" # Executive Devices Group
    },
    # Exclude test environment group
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333" # Test Devices Group
    },
    # Exclude development environment group
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444" # Development Devices Group
    }
  ]

  # Optional: Custom timeouts for resource operations
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}