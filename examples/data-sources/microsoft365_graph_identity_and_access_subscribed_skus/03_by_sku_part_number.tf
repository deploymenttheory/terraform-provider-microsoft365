data "microsoft365_graph_identity_and_access_subscribed_skus" "e5_by_part_number" {
  sku_part_number = "E5"

  timeouts = {
    read = "30s"
  }
}

output "e5_skus_summary" {
  value = [
    for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.e5_by_part_number.items : {
      sku_part_number    = sku.sku_part_number
      consumed_units     = sku.consumed_units
      enabled_units      = sku.prepaid_units.enabled
      available_licenses = sku.prepaid_units.enabled - sku.consumed_units
    }
  ]
  description = "All SKUs containing 'E5' in their part number"
}
