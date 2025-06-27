resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "minimal" {
  group_lifetime_in_days = 180
  managed_group_types    = "All"
} 