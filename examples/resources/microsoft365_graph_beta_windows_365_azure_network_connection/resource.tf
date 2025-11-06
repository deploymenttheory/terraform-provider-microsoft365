resource "microsoft365_graph_beta_windows_365_azure_network_connection" "example" {
  display_name        = "example-azure-network-connection"
  connection_type     = "hybridAzureADJoin"              # or "azureADJoin"
  ad_domain_name      = "ad.example.com"                 # (optional)
  ad_domain_username  = "admin@example.com"              # (optional)
  ad_domain_password  = var.ad_domain_password           # (optional, sensitive)
  organizational_unit = "OU=Computers,DC=example,DC=com" # (optional)
  resource_group_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg"
  subnet_id           = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/example-vnet/subnets/example-subnet"
  subscription_id     = "00000000-0000-0000-0000-000000000000"
  virtual_network_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/virtualNetworks/example-vnet"

  # Optionally, you can set timeouts
  # timeouts {
  #   create = "30m"
  #   update = "30m"
  #   delete = "30m"
  #   read   = "5m"
  # }
}

variable "ad_domain_password" {
  description = "The password for the AD domain join account."
  type        = string
  sensitive   = true
} 