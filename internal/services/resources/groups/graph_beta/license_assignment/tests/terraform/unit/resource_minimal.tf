resource "microsoft365_graph_beta_groups_license_assignment" "minimal" {
  group_id = "00000000-0000-0000-0000-000000000002"
  sku_id   = "f30db892-07e9-47e9-837c-80727f46fd3d" # FLOW_FREE
}
