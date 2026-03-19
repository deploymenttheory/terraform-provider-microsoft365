# Unit Test 01: Find Windows Update product by catalog ID

data "microsoft365_graph_beta_windows_updates_product" "test" {
  search_type  = "catalog_id"
  search_value = "test-catalog-id-123"
}
