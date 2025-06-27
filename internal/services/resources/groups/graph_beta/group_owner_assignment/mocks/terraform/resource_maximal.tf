resource "microsoft365_graph_beta_groups_group_owner_assignment" "maximal" {
  group_id          = "00000000-0000-0000-0000-000000000003"
  owner_id          = "00000000-0000-0000-0000-000000000005"
  owner_object_type = "ServicePrincipal"
} 