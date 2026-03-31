data "microsoft365_graph_identity_and_access_subscribed_skus" "test" {
  # Microsoft Intune Suite SKU present in the test tenant
  sku_id = "2fd6bb84-ad40-4ec5-9369-a215b25c9952_a929cd4d-8672-47c9-8664-159c1f322ba8"
}
