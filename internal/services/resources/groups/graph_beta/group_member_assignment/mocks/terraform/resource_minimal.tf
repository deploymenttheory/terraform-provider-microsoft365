resource "microsoft365_graph_beta_groups_group_member_assignment" "minimal" {
  group_id = "00000000-0000-0000-0000-000000000002"
  member_id = "00000000-0000-0000-0000-000000000004"
  member_object_type = "User"
} 