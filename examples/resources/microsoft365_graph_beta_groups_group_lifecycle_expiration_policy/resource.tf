resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "example" {
  group_lifetime_in_days        = 365
  managed_group_types           = "All" // All, Selected, None
  alternate_notification_emails = "admin@deploymenttheory.com;notifications@deploymenttheory.com"
  overwrite_existing_policy     = false // Optional: Overwrite existing policy, defaults to false
} 