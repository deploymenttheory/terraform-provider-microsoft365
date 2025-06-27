resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "maximal" {
  group_lifetime_in_days        = 365
  managed_group_types           = "Selected"
  alternate_notification_emails = "admin@example.com;notifications@example.com"
} 