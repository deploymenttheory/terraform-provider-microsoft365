# Example configurations for Microsoft 365 Group Lifecycle Policy
# This resource manages group lifecycle policies that set expiration periods for Microsoft 365 groups

# Scenario 1: Default lifecycle policy for all Microsoft 365 groups
# This policy applies to all groups and sets a 180-day expiration period
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "default_policy" {
  # Number of days before a group expires and needs to be renewed
  group_lifetime_in_days = 180

  # Apply to all Microsoft 365 groups
  managed_group_types = "All"

  # Optional: List of email addresses to send notifications for groups without owners
  # Multiple email addresses can be defined by separating with semicolons
  alternate_notification_emails = "admin@example.com;notifications@example.com"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Scenario 2: Short-term project groups policy
# This policy is for temporary project groups that should expire quickly
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "project_groups_policy" {
  # Short expiration period for project groups (90 days)
  group_lifetime_in_days = 90

  # Apply to selected groups only (requires manual assignment)
  managed_group_types = "Selected"

  # Multiple notification emails for project management
  alternate_notification_emails = "pm@example.com;project-admin@example.com;it-support@example.com"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Scenario 3: Long-term department groups policy
# This policy is for permanent department groups with longer expiration
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "department_groups_policy" {
  # Longer expiration period for department groups (365 days)
  group_lifetime_in_days = 365

  # Apply to selected groups only
  managed_group_types = "Selected"

  # Single notification email for department heads
  alternate_notification_emails = "dept-heads@example.com"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Scenario 4: Disabled lifecycle policy
# This policy effectively disables group lifecycle management
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "disabled_policy" {
  # Set a very long expiration period (10 years) to effectively disable
  group_lifetime_in_days = 3650

  # Don't apply to any groups
  managed_group_types = "None"

  # No notification emails needed since policy is disabled
  # alternate_notification_emails is omitted (optional field)

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Scenario 5: Compliance-focused policy
# This policy ensures regular review of groups for compliance purposes
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "compliance_policy" {
  # 6-month expiration for compliance review cycles
  group_lifetime_in_days = 180

  # Apply to all groups for comprehensive compliance
  managed_group_types = "All"

  # Multiple stakeholders for compliance notifications
  alternate_notification_emails = "compliance@example.com;security@example.com;legal@example.com;admin@example.com"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
} 