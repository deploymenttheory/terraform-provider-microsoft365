resource "microsoft365_graph_beta_groups_group_app_role_assignment" "minimal" {
  target_group_id    = "00000000-0000-0000-0000-000000000002"
  resource_object_id = "00000000-0000-0000-0000-000000000010"
  app_role_id        = "00000000-0000-0000-0000-000000000000"
}

