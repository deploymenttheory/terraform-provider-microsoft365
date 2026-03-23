data "microsoft365_graph_identity_and_access_subscribed_skus" "by_account" {
  account_id = "f97aeefc-af85-414d-8ae4-b457f90efc40" // your tenant id

  timeouts = {
    read = "30s"
  }
}

output "tenant_license_summary" {
  value = {
    total_skus = length(data.microsoft365_graph_identity_and_access_subscribed_skus.tenant_skus.items)
    skus = [
      for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.tenant_skus.items : {
        name               = sku.sku_part_number
        total_licenses     = sku.prepaid_units.enabled
        used_licenses      = sku.consumed_units
        available_licenses = sku.prepaid_units.enabled - sku.consumed_units
        status             = sku.capability_status
      }
    ]
  }
  description = "Complete license summary for the tenant (account_id matches tenant_id)"
}
