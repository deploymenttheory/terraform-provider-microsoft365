# Example: Configure tenant-wide Microsoft 365 group settings
# This resource manages organization-wide policies for Microsoft 365 groups

# First, get the Group.Unified template details to see available settings
data "microsoft365_graph_beta_identity_and_access_directory_setting_templates" "group_unified" {
  filter_type  = "display_name"
  filter_value = "Group.Unified"
}

# Configure tenant-wide group creation and guest access policies using the Group.Unified template
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "unified_settings" {
  # Use the template ID from the data source
  template_id = data.microsoft365_graph_beta_identity_and_access_directory_setting_templates.group_unified.directory_setting_templates[0].id

  # Example 1: Minimal configuration with only essential settings
  # This shows how to configure just the most commonly used settings
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

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Example 2: Comprehensive configuration with naming policies and guest restrictions
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "comprehensive_settings" {
  # Use the template ID from the data source
  template_id = data.microsoft365_graph_beta_identity_and_access_directory_setting_templates.group_unified.directory_setting_templates[0].id

  values = [
    # Group creation controls
    {
      name  = "EnableGroupCreation"
      value = "false" # Only specific groups can create Microsoft 365 groups
    },
    {
      name  = "GroupCreationAllowedGroupId"
      value = "00000000-0000-0000-0000-000000000000" # Replace with actual security group ID
    },

    # Naming policies
    {
      name  = "PrefixSuffixNamingRequirement"
      value = "[Marketing]-[GroupName]" # Enforce naming convention
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
      value = "Public,Internal,Confidential,Highly Confidential" # Available classifications
    },
    {
      name  = "ClassificationDescriptions"
      value = "Public:Public data,Internal:Internal data,Confidential:Confidential data,Highly Confidential:Highly Confidential data" # Descriptions
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

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Output the created settings ID
output "unified_settings_id" {
  value       = microsoft365_graph_beta_groups_tenant_wide_group_settings.unified_settings.id
  description = "The ID of the created tenant-wide group settings"
}

# NOTE: In a real environment, you would typically create only one tenant-wide setting per template.
# The two resources shown here are for demonstration purposes only. 