resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "none" {
  group_lifetime_in_days        = 365
  managed_group_types           = "None"
  alternate_notification_emails = "admin@deploymenttheory.com"
}

