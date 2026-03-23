data "microsoft365_graph_identity_and_access_subscribed_skus" "e5_sku" {
  sku_id = "2fd6bb84-ad40-4ec5-9369-a215b25c9952_06ebc4ee-1bb5-47dd-8120-11324bc54e06" // tenant id and sku id

  timeouts = {
    read = "30s"
  }
}

output "e5_sku_details" {
  value = length(data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items) > 0 ? {
    sku_part_number   = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].sku_part_number
    consumed_units    = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].consumed_units
    enabled_units     = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].prepaid_units.enabled
    capability_status = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].capability_status
    service_plans     = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].service_plans
  } : null
  description = "Microsoft 365 E5 SKU details including service plans"
}
