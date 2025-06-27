# Comprehensive configuration for group-specific settings
resource "microsoft365_graph_beta_groups_group_settings" "maximal" {
  # Test group ID
  group_id = "12345678-1234-1234-1234-123456789012"

  # Group.Unified template ID (using the full-featured template)
  template_id = "62375ab9-6b52-47ed-826b-58e47e0e304b"

  # Group-specific settings that override tenant-wide defaults
  values = [
    # Classification settings
    {
      name  = "ClassificationList"
      value = "Confidential,Secret,Top Secret" # Custom classifications for this group
    },
    {
      name  = "DefaultClassification"
      value = "Confidential" # Default classification for this group
    },

    # Guest access settings
    {
      name  = "AllowGuestsToBeGroupOwner"
      value = "false" # Don't allow guests to be group owners
    },
    {
      name  = "AllowGuestsToAccessGroups"
      value = "true" # Allow guests to access groups
    },
    {
      name  = "AllowToAddGuests"
      value = "true" # Allow adding guests to groups
    },

    # Other settings
    {
      name  = "UsageGuidelinesUrl"
      value = "https://contoso.com/marketing-group-guidelines" # Group-specific guidelines
    }
  ]

  # Timeouts
  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
} 