# Comprehensive configuration for tenant-wide group settings
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "test" {
  # Group.Unified template ID
  template_id = "62375ab9-6b52-47ed-826b-58e47e0e304b"

  # Complete set of settings
  values = [
    # Group creation controls
    {
      name  = "EnableGroupCreation"
      value = "false" # Only specific groups can create Microsoft 365 groups
    },
    {
      name  = "GroupCreationAllowedGroupId"
      value = "12345678-1234-1234-1234-123456789012" # Security group ID
    },

    # Naming policies
    {
      name  = "PrefixSuffixNamingRequirement"
      value = "[Contoso]-[GroupName]" # Enforce naming convention
    },
    {
      name  = "CustomBlockedWordsList"
      value = "CEO,Legal,HR" # Block specific words in group names
    },
    {
      name  = "EnableMSStandardBlockedWords"
      value = "true" # Enable Microsoft's list of blocked words
    },

    # Classification settings
    {
      name  = "ClassificationList"
      value = "Public,Internal,Confidential" # Available classifications
    },
    {
      name  = "ClassificationDescriptions"
      value = "Public:Public data,Internal:Internal data,Confidential:Confidential data" # Descriptions
    },
    {
      name  = "DefaultClassification"
      value = "Internal" # Default classification
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
    {
      name  = "GuestUsageGuidelinesUrl"
      value = "https://contoso.com/guestpolicies" # Link to guest usage guidelines
    },

    # Other settings
    {
      name  = "UsageGuidelinesUrl"
      value = "https://contoso.com/groupguidelines" # Link to general usage guidelines
    },
    {
      name  = "EnableMIPLabels"
      value = "true" # Enable sensitivity labels
    },
    {
      name  = "NewUnifiedGroupWritebackDefault"
      value = "true" # Enable group writeback to on-premises AD
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