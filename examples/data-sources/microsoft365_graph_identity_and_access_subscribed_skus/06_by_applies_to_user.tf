data "microsoft365_graph_identity_and_access_subscribed_skus" "user_assignable" {
  applies_to = "User"

  timeouts = {
    read = "30s"
  }
}

output "user_license_inventory" {
  value = [
    for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.user_assignable.items : {
      sku_part_number    = sku.sku_part_number
      sku_id             = sku.sku_id
      consumed_units     = sku.consumed_units
      enabled_units      = sku.prepaid_units.enabled
      available_units    = sku.prepaid_units.enabled - sku.consumed_units
      capability_status  = sku.capability_status
      service_plan_count = length(sku.service_plans)
    }
  ]
  description = "User-assignable SKUs with license availability"
}
