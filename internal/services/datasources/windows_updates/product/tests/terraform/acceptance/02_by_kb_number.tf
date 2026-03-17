# Test 02: Find Windows Update product by KB number
# This test retrieves product information using a KB article number

data "microsoft365_graph_beta_windows_updates_product" "test" {
  search_type  = "kb_number"
  search_value = "5029332"
}

output "product_count" {
  value = length(data.microsoft365_graph_beta_windows_updates_product.test.products)
}

output "product_names" {
  value = [
    for product in data.microsoft365_graph_beta_windows_updates_product.test.products :
    product.name
  ]
}

output "known_issues_count" {
  value = sum([
    for product in data.microsoft365_graph_beta_windows_updates_product.test.products :
    length(product.known_issues)
  ])
}
