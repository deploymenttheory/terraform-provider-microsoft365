resource "microsoft365_graph_beta_groups_license_assignment" "maximal" {
  group_id = "00000000-0000-0000-0000-000000000003"
  sku_id   = "44444444-4444-4444-4444-444444444444"
  disabled_plans = [
    "55555555-5555-5555-5555-555555555555",
    "66666666-6666-6666-6666-666666666666"
  ]
}
