resource "microsoft365_graph_beta_users_user_license_assignment" "minimal" {
  user_id = "00000000-0000-0000-0000-000000000002"
  add_licenses = [
    {
      sku_id = "a403ebcc-fae0-4ca2-8c8c-7a907fd6c235" # POWER_BI_STANDARD
    },
    {
      sku_id = "f30db892-07e9-47e9-837c-80727f46fd3d" # FLOW_FREE
    }
  ]
} 