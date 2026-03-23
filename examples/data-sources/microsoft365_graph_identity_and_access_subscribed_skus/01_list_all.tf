data "microsoft365_graph_identity_and_access_subscribed_skus" "all" {
  list_all = true

  timeouts = {
    read = "30s"
  }
}

output "all_skus_count" {
  value       = length(data.microsoft365_graph_identity_and_access_subscribed_skus.all.items)
  description = "Total number of subscribed SKUs"
}
