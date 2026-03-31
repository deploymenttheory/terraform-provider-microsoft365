data "microsoft365_graph_identity_and_access_subscribed_skus" "test" {
  # "Intune" partial match for Microsoft_Intune_Suite present in the test tenant
  sku_part_number = "Intune"
}
