# Unit Test 02: Find Windows Update product by KB number

data "microsoft365_graph_beta_windows_updates_product" "test" {
  search_type  = "kb_number"
  search_value = "5029332"
}
