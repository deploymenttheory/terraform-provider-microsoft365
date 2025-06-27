resource "microsoft365_graph_beta_groups_group_owner_assignment" "minimal" {
  group_id          = "00000000-0000-0000-0000-000000000002"
  owner_id          = "00000000-0000-0000-0000-000000000004"
  owner_object_type = "User"
} 