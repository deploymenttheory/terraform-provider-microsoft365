# Example: Windows Quality Update Expedite Policy - Maximal Configuration (No Assignments)
resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "maximal_example" {
  display_name = "Critical Security Update Expedite Policy"
  description  = "Expedited deployment for critical security updates - January 2025"
  role_scope_tag_ids = ["0", "1"]

  # Required: Expedited update settings
  # Defines which quality update to expedite and reboot behavior
  expedited_update_settings = {
    # Quality update release to expedite
    # Valid values: "2025-12-09T00:00:00Z", "2025-11-20T00:00:00Z"
    quality_update_release = "2025-12-09T00:00:00Z"

    # If a reboot is required, select the number of days before it's enforced
    # Valid values: 0, 1, or 2
    # 0 = Immediate reboot after installation
    # 1 = Allow 1 day for user-initiated reboot
    # 2 = Allow 2 days for user-initiated reboot
    days_until_forced_reboot = 1
  }

  # Optional: Custom timeouts for resource operations
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}

