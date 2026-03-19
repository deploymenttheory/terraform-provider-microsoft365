# Test 01: Find Windows Update product by catalog ID
# This test retrieves product information using a catalog ID

data "microsoft365_graph_beta_windows_updates_product" "test" {
  search_type  = "catalog_id"
  search_value = "10cb1ba292c5586e22c9991be3f12fbd39f2ebf231cb5d201c67f42fbaccc567"
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

output "revision_count" {
  value = sum([
    for product in data.microsoft365_graph_beta_windows_updates_product.test.products :
    length(product.revisions)
  ])
}
