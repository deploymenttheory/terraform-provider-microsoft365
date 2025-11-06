resource "microsoft365_graph_beta_windows_365_azure_network_connection" "maximal" {
  display_name        = "Test Maximal Connection"
  connection_type     = "hybridAzureADJoin"
  ad_domain_name      = "example.local"
  ad_domain_username  = "testuser"
  ad_domain_password  = "TestPassword123!"
  organizational_unit = "OU=CloudPCs,DC=example,DC=local"
  resource_group_id   = "/subscriptions/22222222-2222-2222-2222-222222222222/resourceGroups/test-rg-maximal"
  subnet_id           = "/subscriptions/22222222-2222-2222-2222-222222222222/resourceGroups/test-rg-maximal/providers/Microsoft.Network/virtualNetworks/test-vnet-maximal/subnets/test-subnet-maximal"
  subscription_id     = "22222222-2222-2222-2222-222222222222"
  virtual_network_id  = "/subscriptions/22222222-2222-2222-2222-222222222222/resourceGroups/test-rg-maximal/providers/Microsoft.Network/virtualNetworks/test-vnet-maximal"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}