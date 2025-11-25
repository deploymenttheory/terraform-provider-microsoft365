resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "all" {
  group_lifetime_in_days        = 180
  managed_group_types           = "All"
  alternate_notification_emails = "admin@deploymenttheory.com"
} 