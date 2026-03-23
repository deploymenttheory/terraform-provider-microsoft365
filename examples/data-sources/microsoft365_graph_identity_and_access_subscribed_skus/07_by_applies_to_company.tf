data "microsoft365_graph_identity_and_access_subscribed_skus" "company_level" {
  applies_to = "Company"

  timeouts = {
    read = "30s"
  }
}

output "company_licenses" {
  value = [
    for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.company_level.items : {
      sku_part_number   = sku.sku_part_number
      sku_id            = sku.sku_id
      account_name      = sku.account_name
      capability_status = sku.capability_status
      consumed_units    = sku.consumed_units
    }
  ]
  description = "Company-level SKUs (not assignable to individual users)"
}
