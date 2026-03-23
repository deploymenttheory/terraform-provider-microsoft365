data "microsoft365_graph_identity_and_access_subscribed_skus" "org_skus" {
  account_name = "DeploymentTheory" // your tenant name

  timeouts = {
    read = "30s"
  }
}

output "org_license_usage" {
  value = [
    for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.org_skus.items : {
      sku_part_number    = sku.sku_part_number
      account_name       = sku.account_name
      total_licenses     = sku.prepaid_units.enabled
      used_licenses      = sku.consumed_units
      available_licenses = sku.prepaid_units.enabled - sku.consumed_units
      utilization_pct    = sku.prepaid_units.enabled > 0 ? (sku.consumed_units / sku.prepaid_units.enabled) * 100 : 0
    }
  ]
  description = "License usage summary filtered by account name (partial match)"
}
