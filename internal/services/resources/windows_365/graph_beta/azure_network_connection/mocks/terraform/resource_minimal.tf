resource "microsoft365_graph_beta_windows_365_azure_network_connection" "minimal" {
  display_name       = "Test Minimal Connection"
  connection_type    = "hybridAzureADJoin"
  ad_domain_name     = "example.local"
  ad_domain_username = "testuser"
  ad_domain_password = "TestPassword123!"
  resource_group_id  = "/subscriptions/11111111-1111-1111-1111-111111111111/resourcegroups/test-rg"
  subnet_id          = "/subscriptions/11111111-1111-1111-1111-111111111111/resourcegroups/test-rg/providers/microsoft.network/virtualnetworks/test-vnet/subnets/test-subnet"
  subscription_id    = "11111111-1111-1111-1111-111111111111"
  virtual_network_id = "/subscriptions/11111111-1111-1111-1111-111111111111/resourcegroups/test-rg/providers/microsoft.network/virtualnetworks/test-vnet"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}