resource "microsoft365_graph_beta_users_user_license_assignment" "minimal" {
  user_id = "00000000-0000-0000-0000-000000000002"
  add_licenses = [{
    sku_id = "33333333-3333-3333-3333-333333333333"
  }]
} 