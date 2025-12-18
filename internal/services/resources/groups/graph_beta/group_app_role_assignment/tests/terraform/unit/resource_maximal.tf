resource "microsoft365_graph_beta_groups_group_app_role_assignment" "maximal" {
  target_group_id    = "00000000-0000-0000-0000-000000000003"
  resource_object_id = "00000000-0000-0000-0000-000000000011"
  app_role_id        = "00000000-0000-0000-0000-000000000000"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

